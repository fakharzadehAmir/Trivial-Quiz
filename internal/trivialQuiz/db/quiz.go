package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	QuizName     string             `bson:"quiz_name"`
	Difficulty   string             `bson:"difficulty"`
	Questions    primitive.A        `bson:"questions"`
	UserAnswerer string             `bson:"user_answerer"`
	Score        float64            `bson:"score"`
	Category     string             `bson:"category"`
}

func (mdb *MongoDB) CreateNewQuiz(ctx *context.Context, newQuiz *Quiz) error {
	//	Proceed to create new quiz
	_, err := mdb.Collections.
		QuizCollection.Collection.InsertOne(*ctx, newQuiz)
	if err != nil {
		return err
	}
	//	Quiz created successfully
	return nil
}
