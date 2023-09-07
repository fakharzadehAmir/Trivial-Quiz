package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Text       string             `bson:"text"`
	Options    []string           `bson:"options"`
	Correct    string             `bson:"correct"`
	Difficulty string             `bson:"difficulty"`
	Category   string             `bson:"category"`
}

// String returns a string representing the fields of User
func (q *Question) String() string {
	return fmt.Sprintf("Question{%+v}", *q)
}

//	CreateNewQuestion returns error if it can't store new question that you wanted to create
func (mdb *MongoDB) CreateNewQuestion(ctx *context.Context, newQuestion *Question) error {
	//	Proceed to create new question
	_, err := mdb.Collections.
		QuestionCollection.Collection.InsertOne(*ctx, newQuestion)
	if err != nil {
		return err
	}
	//	Question created successfully
	return nil
}

//	GetQuestionByID returns the Question and error for the time you want to create a new quiz
func (mdb *MongoDB) GetQuestionByID(ctx *context.Context,
	qId primitive.ObjectID) (*Question, error) {
	//	Check existence of question with given id
	findQuestion := mdb.Collections.
		QuestionCollection.Collection.
		FindOne(*ctx, bson.M{"_id": qId})

	//	Handle error of retrieving the question with given username
	if findQuestion.Err() != nil {
		return nil, findQuestion.Err()
	}

	//	Decode existed question to the declared variable existedQuestion
	existedQuestion := &Question{}
	if err := findQuestion.Decode(existedQuestion); err != nil {
		return nil, err
	}

	return existedQuestion, nil
}

//	GetQuestionsByCategoryDifficulty returns array of questions with different category and difficulty
func (mdb *MongoDB) GetQuestionsByCategoryDifficulty(ctx context.Context,
	category string, difficulty string) ([]*Question, error) {
	//	Create a slice to store the result
	var questions []*Question

	//	Find documents that match the given category and difficulty
	cursor, err := mdb.Collections.
		QuestionCollection.Collection.
		Find(ctx, bson.M{
			"category":   category,
			"difficulty": difficulty,
		})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	//	Iterate through the cursor and decode documents into the questions slice
	for cursor.Next(ctx) {
		var question Question
		if err := cursor.Decode(&question); err != nil {
			return nil, err
		}
		questions = append(questions, &question)
	}

	//	Check for errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

