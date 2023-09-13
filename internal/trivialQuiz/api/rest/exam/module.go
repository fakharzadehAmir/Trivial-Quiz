package exam

import (
	"Trivia_Quiz/internal/trivialQuiz/api/rest/server"
	"Trivia_Quiz/internal/trivialQuiz/db"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type QuizHandler struct {
	db     *db.MongoDB
	logger *logrus.Logger
	ctx    *context.Context
}

func NewQuizHandler(db *db.MongoDB, logger *logrus.Logger,
	ctx *context.Context) (server.Module, error) {
	return &QuizHandler{
		db:     db,
		logger: logger,
		ctx:    ctx,
	}, nil
}

func (q *QuizHandler) GetRoutes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/create-exam",
			Handler: q.CreateExamHandler,
		},
		{
			Method:  http.MethodPut,
			Path:    "/submit-exam/:quizId",
			Handler: q.SubmitExamHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/check-exam/:quizId/answers",
			Handler: q.CheckExamHandler,
		},
	}
}

func (q *QuizHandler) CreateExamHandler(c *gin.Context) {
	//	Get logged-in username
	loggedInUser, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "no user has been logged in, login is required",
		})
		return
	}
	//	Declare a struct to define the query parameters
	var queryParams struct {
		Category   string `form:"category"`
		Difficulty string `form:"difficulty"`
	}

	//	Bind queries in path to its struct and handle it error
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		q.logger.WithError(err).Warn("can not bind queries")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not bind given queries",
			"err":     err.Error(),
		})
		return
	}

	//	Get at most 10 question of given category and difficulty
	questions, err := q.db.GetQuestionsByCategoryDifficulty(q.ctx, queryParams.Category, queryParams.Difficulty)
	if err != nil {
		q.logger.WithError(err).Warn("can not return questions with given category and difficulty")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can not return questions with given category and difficulty",
			"err":     err.Error(),
		})
		return
	}

	//	Create a new quiz with given category and difficulty
	quizId, err := q.db.CreateNewQuiz(q.ctx, &db.Quiz{
		CreatedAt:    time.Now(),
		Category:     queryParams.Category,
		Difficulty:   queryParams.Difficulty,
		UserAnswerer: loggedInUser.(string),
		Questions:    questions,
	})
	if err != nil {
		q.logger.WithError(err).Warn("can not create quiz with given category and difficulty")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can not create quiz with given category and difficulty",
			"err":     err.Error(),
		})
		return
	}
	q.logger.Info("new quiz created successfully!")

	//	retrieve quiz to answer it for user
	var createdQuiz = &createQuizResponse{
		ID: quizId,
	}
	for _, question := range questions {
		var questionResponse = questionsQuizResponse{
			QId:     question.ID,
			Text:    question.Text,
			Option1: question.Options[0],
			Option2: question.Options[1],
			Option3: question.Options[2],
			Option4: question.Options[3],
		}
		createdQuiz.Questions = append(createdQuiz.Questions, questionResponse)
	}
	c.JSON(http.StatusCreated, *createdQuiz)
}

func (q *QuizHandler) SubmitExamHandler(c *gin.Context) {
	//	Get qid params in url and cast it to ObjectID
	qid := c.Param("quizId")
	quizId, err := primitive.ObjectIDFromHex(qid)
	if err != nil {
		q.logger.WithError(err).Warn("can not cast quiz id to ObjectID")
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can not cast quiz id to ObjectID",
			"error":   err.Error(),
		})
		return
	}

	//	Check logged-in username with user of the quiz with given id
	loggedInUser, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "no user has been logged in, login is required",
		})
		return
	}
	correctUser, err := q.db.GetUserOfQuizByID(q.ctx, quizId, loggedInUser.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can not retrieve the quiz with given id",
		})
		return
	}
	if !correctUser {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "the user of given quiz is not matched with logged-in user",
		})
		return
	}

	//	Parse request body for submitting quiz
	answerExam := submitExamRequest{}
	err = c.ShouldBindJSON(&answerExam)
	if err != nil {
		q.logger.WithError(err).Warn("can not unmarshal the exam submission request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not unmarshal request body of the exam submission",
			"error":   err.Error(),
		})
		return
	}

	//	Evaluation of user's answer
	var evalExam = 0
	for _, answer := range answerExam.MyExam {
		correctAnswer, err := q.db.GetQuestionAnswerByID(q.ctx, answer.QuestionID)
		if err != nil {
			q.logger.WithError(err).Warnf("can not get the answer of the question with this id %s", answer.QuestionID)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "can not get the answer of the question",
				"error":   err.Error(),
			})
			return
		}
		if correctAnswer == answer.MyAnswer {
			evalExam++
		}
	}

	//	Updating score in the database
	totalScore := float64(evalExam) / float64(len(answerExam.MyExam)) * 10
	err = q.db.UpdateScoreByID(q.ctx, quizId, totalScore)
	if err != nil {
		q.logger.WithError(err).Warn("can not submit your answers")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can not submit your answers",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, submitExamResponse{
		YourScore: totalScore,
	})
}

func (q *QuizHandler) CheckExamHandler(c *gin.Context) {
	//	Get qid params in url and cast it to ObjectID
	qid := c.Param("quizId")
	quizId, err := primitive.ObjectIDFromHex(qid)
	if err != nil {
		q.logger.WithError(err).Warn("can not cast quiz id to ObjectID")
		c.JSON(http.StatusNotFound, gin.H{
			"message": "can not cast quiz id to ObjectID",
			"error":   err.Error(),
		})
		return
	}

	//	Check logged-in username with user of the quiz with given id
	loggedInUser, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "no user has been logged in, login is required",
		})
		return
	}
	correctUser, err := q.db.GetUserOfQuizByID(q.ctx, quizId, loggedInUser.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can not retrieve the quiz with given id",
		})
		return
	}
	if !correctUser {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "the user of given quiz is not matched with logged-in user",
		})
		return
	}

	//	Get Question of the quiz with the given id if it's been answered
	questions, err := q.db.GetQuestionsByID(q.ctx, quizId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "this quiz has not been answered",
		})
		return
	}

	//	Initialize retrieved questions in its struct response body
	resBody := checkExamResponse{}
	for _, question := range questions {
		resBody.Answers = append(resBody.Answers,
			checkQuestionResponse{
				Text:          question.Text,
				CorrectAnswer: question.Correct,
				Option1:       question.Options[0],
				Option2:       question.Options[1],
				Option3:       question.Options[2],
				Option4:       question.Options[3],
			})
	}
	c.JSON(http.StatusOK, resBody)
}
