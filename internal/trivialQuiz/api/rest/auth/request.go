package auth

import "time"

type signupRequestBody struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	Birthday time.Time `json:"birthday"`
	Email    string    `json:"email"`
}

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
