package db

import (
	"context"
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
}

func (mdb *MongoDB) CreateNewQuiz(ctx *context.Context, newQuiz *Quiz) (interface{}, error) {
	//	Proceed to create new quiz
	ID, err := mdb.Collections.
		QuizCollection.Collection.InsertOne(*ctx, newQuiz)
	if err != nil {
		return nil, err
	}
	//	Quiz created successfully
	return ID, nil
}

func (mdb *MongoDB) UpdateScoreByID(ctx *context.Context,
	quizId primitive.ObjectID, answers []string) error {
	//	Check existence of quiz with given id
	findQuiz := mdb.Collections.
		QuestionCollection.Collection.
		FindOne(*ctx, bson.M{"_id": quizId})

	//	Handle error of retrieving the question with given username
	if findQuiz.Err() != nil {
		return findQuiz.Err()
	}

	//	Decode existed question to the declared variable existedQuiz
	existedQuiz := &Quiz{}
	if err := findQuiz.Decode(existedQuiz); err != nil {
		return err
	}

	//var checkScore float64
	//for idx, value := range answers {
	//	question, err := mdb.GetQuestionByID(ctx, existedQuiz.Questions[idx - 1])
	//}

	return nil

}
