package profile

import (
	"Trivia_Quiz/internal/trivialQuiz/api/rest/server"
	"Trivia_Quiz/internal/trivialQuiz/db"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	db     *db.MongoDB
	logger *logrus.Logger
	ctx    *context.Context
}

func NewUserHandler(db *db.MongoDB, logger *logrus.Logger,
	ctx *context.Context) (server.Module, error) {
	return &UserHandler{
		db:     db,
		logger: logger,
		ctx:    ctx,
	}, nil
}

func (u *UserHandler) GetRoutes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPut,
			Path:    "/user",
			Handler: u.UpdateUserHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/user",
			Handler: u.RetrieveUserHandler,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/user",
			Handler: u.DeleteUserHandler,
		},
	}
}

func (u *UserHandler) RetrieveUserHandler(c *gin.Context) {
	loggedInUser, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "no user has been logged in, login is required",
		})
		return
	}

	//	Retrieve the logged-in username from database
	user, err := u.db.GetUserByUsername(u.ctx, loggedInUser.(string))
	if err != nil {
		u.logger.WithError(err).Warn("cannot retrieve user from database")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "cannot retrieve this user from database",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, userRetrievalResponse{
		Username: user.Username,
		Email:    user.Email,
		Birthday: user.Birthday,
	})
	u.logger.Infof("user %s retreived their data!", user.Username)
}

func (u *UserHandler) UpdateUserHandler(c *gin.Context) {
	loggedInUser, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "no user has been logged in, login is required",
		})
		return
	}

	//	Bind request body to updateRequestBody
	updatedUser := updateRequestBody{}
	err := c.ShouldBindJSON(&updatedUser)
	if err != nil {
		u.logger.WithError(err).
			Warnf("can not un marshal the updated user (%s) request body",
				loggedInUser.(string))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "can not unmarshal request body to its struct",
			"error":   err.Error(),
		})
		return
	}

	//	Update profile of the user which is logged in
	err = u.db.UpdateUserByUsername(u.ctx,
		&db.User{
			Password: updatedUser.Password,
			Email:    updatedUser.Email,
			Birthday: updatedUser.Birthday,
		},
		loggedInUser.(string))
	if err != nil {
		u.logger.WithError(err).Warn("cannot update logged-in user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "cannot update logged-in user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user has been updated successfully!",
	})
	u.logger.Infof("user %s has been updated successfully!", loggedInUser.(string))

}

func (u *UserHandler) DeleteUserHandler(c *gin.Context) {
	loggedInUser, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "no user has been logged in, login is required",
		})
		return
	}

	//	Delete the logged-in username from database
	err := u.db.DeleteUserByUsername(u.ctx, loggedInUser.(string))
	if err != nil {
		u.logger.WithError(err).Warn("cannot delete logged-in user from database")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "cannot delete logged-in user from database",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user was deleted successfully!",
	})
	u.logger.Infof("user %s was deleted successfully!", loggedInUser.(string))

}
