package authenticate

import (
	"context"
	"crypto/rand"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Auth struct {
	mongoDB MongoDB
	//	jwtSecretKey is the JWT secret key.
	//	Each time the server starts, new key is generated.
	jwtSecretKey          []byte
	jwtExpirationDuration time.Duration
	logger                *logrus.Logger
}

// NewAuth creates new instance of Auth for authenticating user accounts.
func NewAuth(authMongoDB MongoDB, jwtExpirationMinutes int64,
	authLogger *logrus.Logger) (*Auth, error) {
	secretKey, err := generateRandomKey()
	if err != nil {
		return nil, err
	}

	//	Check authMongoDB
	if authMongoDB == nil {
		return nil, errors.New("the authenticate database is essential")
	}

	//	Check the logger
	if authLogger == nil {
		authLogger = logrus.New()
	}

	return &Auth{
		mongoDB:               authMongoDB,
		logger:                authLogger,
		jwtSecretKey:          secretKey,
		jwtExpirationDuration: time.Duration(int64(time.Minute) * jwtExpirationMinutes),
	}, nil
}

// Login Check login input and if everything was ok, it creates JWT token.
func (a *Auth) Login(ctx *context.Context, cred *Credentials) (*Token, error) {

	//	Check existence of user
	account, err := a.mongoDB.GetAccountByUsername(ctx, cred.Username)
	if err != nil {
		return nil, err
	}

	//	Check Password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(cred.Password))
	if err != nil {
		return nil, err
	}

	return a.GenerateToken(cred)
}

func (a *Auth) GetAccountByToken(ctx context.Context, token string) (*Account, error) {
	//	Handle empty access token
	if token == "" {
		return nil, errors.New("access denied: the access token is empty")
	}

	//	Validate jwt token
	claim, err := a.checkToken(token)
	if err != nil {
		return nil, errors.New("access denied: the access token is invalid")
	}

	//	Get tue user entity
	account, err := a.mongoDB.GetAccountByUsername(&ctx, claim.Username)
	if err != nil {
		return nil, errors.New("access denied: cannot fetch such a user")
	}

	return account, nil
}
func (a *Auth) checkToken(tokenStr string) (*claims, error) {
	c := &claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, c, func(token *jwt.Token) (interface{}, error) {
		return a.jwtSecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token")
		}
		a.logger.WithError(err).Warn("cannot validate the token of the user")
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("unauthorized")
	}

	return c, err
}
func (a *Auth) GenerateToken(cred *Credentials) (*Token, error) {
	//	Create the JWT token
	expirationTime := time.Now().Add(a.jwtExpirationDuration)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		Username: cred.Username,
		MapClaims: jwt.MapClaims{
			"expired_at": expirationTime.Unix(),
		},
	})

	//	Calculate the signed string format of JWT key
	tokenString, err := tokenJWT.SignedString(a.jwtSecretKey)
	if err != nil {
		return nil, err
	}
	return &Token{
		TokenString: tokenString,
		Expiration:  expirationTime,
	}, nil
}

// generateRandomKey
// Each time Auth is initialized, this function is called to generate another key
func generateRandomKey() ([]byte, error) {
	jwtKey := make([]byte, 32)
	if _, err := rand.Read(jwtKey); err != nil {
		return nil, err
	}
	return jwtKey, nil
}
