package exam

import (
	"Trivia_Quiz/internal/trivialQuiz/api/rest/server"
	"Trivia_Quiz/internal/trivialQuiz/db"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
			Path:    "/give-exam/:qid",
			Handler: q.CreateExamHandler,
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
