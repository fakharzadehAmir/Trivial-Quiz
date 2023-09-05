package auth

import (
	"Trivia_Quiz/internal/trivialQuiz/api/rest/server"
	"Trivia_Quiz/internal/trivialQuiz/db"
	"Trivia_Quiz/pkg/authenticate"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AuthenticateHandler struct {
	db            *db.MongoDB
	logger        *logrus.Logger
	ctx           *context.Context
	authenticator *authenticate.Auth
}

func NewAuthHandlers(db *db.MongoDB, logger *logrus.Logger,
	ctx *context.Context, authenticator *authenticate.Auth) (server.Module, error) {
	return &AuthenticateHandler{
		db:            db,
		logger:        logger,
		ctx:           ctx,
		authenticator: authenticator,
	}, nil
}
func (ah *AuthenticateHandler) GetRoutes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/signup",
			Handler: ah.SignupHandler,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: ah.LoginHandler,
		},
	}
}
func (ah *AuthenticateHandler) SignupHandler(c *gin.Context) {
	//	Parse request body for new user
	inputUser := signupRequestBody{}
	err := c.ShouldBindJSON(&inputUser)
	if err != nil {
		ah.logger.WithError(err).Warn("can not un marshal the new user request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not unmarshal request body",
			"error":   err.Error(),
		})
		return
	}
	//	Add the new user to the database
	err = ah.db.CreateNewUser(ah.ctx, &db.User{
		Username: inputUser.Username,
		Password: inputUser.Password,
		Birthday: inputUser.Birthday,
		Email:    inputUser.Email,
	})
	if err != nil {
		ah.logger.WithError(err).Warn("can not create a new user")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not create a new user",
			"error":   err.Error(),
		})
		return
	}

	//	the user has been added successfully
	c.JSON(http.StatusCreated, gin.H{
		"message": "new user has been created successfully!",
	})
	ah.logger.Info("new user has been created successfully!")
}

func (ah *AuthenticateHandler) LoginHandler(c *gin.Context) {
	//	Parse request body for logged-in user
	inputUser := loginRequestBody{}
	err := c.ShouldBindJSON(&inputUser)
	if err != nil {
		ah.logger.WithError(err).Warn("can not un marshal the new user request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not unmarshal request body",
			"error":   err.Error(),
		})
		return
	}

	//	check login middleware
	accessToken, err := ah.authenticator.Login(ah.ctx, &authenticate.Credentials{
		Username: inputUser.Username,
		Password: inputUser.Password,
	})
	if err != nil {
		ah.logger.WithError(err).Warn("cannot login and authorize")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "cannot authorize",
			"error":   err.Error(),
		})
		return
	}

	//	the user logged in successfully
	c.JSON(http.StatusOK, authenticateResponseBody{
		AccessToken: accessToken.TokenString,
	})
	ah.logger.Info("user logged in successfully!")
}
