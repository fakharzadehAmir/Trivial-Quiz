package exam

import (
	"Trivia_Quiz/internal/trivialQuiz/api/rest/server"
	"Trivia_Quiz/internal/trivialQuiz/db"
	"context"
	"github.com/sirupsen/logrus"
)

type QuizHandler struct {
	db     *db.MongoDB
	logger *logrus.Logger
	ctx    *context.Context
}

func NewQuizHandler(db *db.MongoDB, logger *logrus.Logger,
	ctx *context.Context) (server.Module, error) {
	return &QuizHandler{
		db:     db,
		logger: logger,
		ctx:    ctx,
	}, nil
}

func (q *QuizHandler) GetRoutes() []server.Route {
	return []server.Route{}
}
