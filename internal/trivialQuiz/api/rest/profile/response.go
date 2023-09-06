package profile

import "time"

type userRetrievalResponse struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Birthday time.Time `json:"birthday"`
}
