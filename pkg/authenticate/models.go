package authenticate

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Account struct {
	Username string
	Password string
}

type Credentials struct {
	Username string
	Password string
}

type claims struct {
	jwt.MapClaims
	Username string `json:"username"`
}

type Token struct {
	TokenString string
	Expiration  time.Time
}
