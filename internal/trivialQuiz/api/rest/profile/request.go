package profile

import "time"

type updateRequestBody struct {
	Password string    `json:"password"`
	Birthday time.Time `json:"birthday"`
	Email    string    `json:"email"`
}
