package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Quiz struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Difficulty   string             `bson:"difficulty"`
	Questions    []Question         `bson:"questions"`
	UserAnswerer string             `bson:"user_answerer"`
	//	The Score is out of 10
	Score      float64   `bson:"score"`
	Category   string    `bson:"category"`
	AnsweredAt time.Time `bson:"answered_at"`
	CreatedAt  time.Time `bson:"created_at"`
	IsAnswered bool      `bson:"is_answered"`
}

func (mdb *MongoDB) CreateNewQuiz(ctx *context.Context, newQuiz *Quiz) (interface{}, error) {
	//	Proceed to create new quiz
	ID, err := mdb.Collections.
		QuizCollection.Collection.InsertOne(*ctx, newQuiz)
	if err != nil {
		return nil, err
	}
	//	Quiz created successfully
	return ID.InsertedID, nil
}

func (mdb *MongoDB) UpdateScoreByID(ctx *context.Context,
	quizId primitive.ObjectID, userScore float64) error {
	// Retrieve the quiz document
	var quiz Quiz
	err := mdb.Collections.QuizCollection.Collection.
		FindOne(*ctx, bson.M{"_id": quizId}).Decode(&quiz)
	if err != nil {
		return err
	}
	//	Update if it has not been answered
	if !quiz.IsAnswered {
		//	Eval answered time
		answeredTime := time.Now()
		//	Check time limit of exam
		if answeredTime.Sub(quiz.CreatedAt) < 10*time.Minute {
			//	Update user's score and answered time
			_, err := mdb.Collections.QuizCollection.Collection.
				UpdateOne(*ctx, bson.M{"_id": quizId},
					bson.M{
						"$set": bson.M{
							"score":       userScore,
							"answered_at": answeredTime,
							"is_answered": true,
						},
					})
			return err
		}
		return errors.New("time limit exceeded, cannot update score")
	}
	return errors.New("quiz has already been answered, score cannot be updated")
}

func (mdb *MongoDB) GetUserOfQuizByID(ctx *context.Context,
	quizId primitive.ObjectID, loggedInUser string) (bool, error) {
	// Retrieve the quiz document
	var quiz Quiz
	err := mdb.Collections.QuizCollection.Collection.
		FindOne(*ctx, bson.M{"_id": quizId}).Decode(&quiz)
	if err != nil {
		return false, err
	}
	//	quiz has been retrieved successfully
	return quiz.UserAnswerer == loggedInUser, nil
}
