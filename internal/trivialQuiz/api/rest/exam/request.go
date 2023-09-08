package exam

import "go.mongodb.org/mongo-driver/bson/primitive"

type giveExamRequest struct {
	MyExam []answerQuestionsRequest `json:"my_exam"`
}

type answerQuestionsRequest struct {
	QuestionID primitive.ObjectID `json:"question_id"`
	MyAnswer   string             `json:"my_answer"`
}
