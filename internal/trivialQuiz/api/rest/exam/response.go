package exam

import "go.mongodb.org/mongo-driver/bson/primitive"

type createQuizResponse struct {
	ID        interface{}             `json:"_id"`
	Questions []questionsQuizResponse `json:"questions"`
}
type questionsQuizResponse struct {
	QId     primitive.ObjectID `json:"QId"`
	Text    string             `json:"text"`
	Option1 string             `json:"option1"`
	Option2 string             `json:"option2"`
	Option3 string             `json:"option3"`
	Option4 string             `json:"option4"`
}
