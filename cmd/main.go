package main

import (
	"Trivia_Quiz/config"
	"Trivia_Quiz/internal/trivialQuiz/api/rest/auth"
	"Trivia_Quiz/internal/trivialQuiz/api/rest/server"
	"Trivia_Quiz/internal/trivialQuiz/db"
	"Trivia_Quiz/pkg/authenticate"
	"context"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	cfg    *config.Config
	logger *logrus.Logger
	ctx    context.Context
)

func main() {

	logger = logrus.New()
	ctx = context.Background()
	cfg = &config.Config{}

	// Get fields in .env file
	err := godotenv.Load()
	if err != nil {
		logger.WithError(err).Error("can not load .env file")
		return
	}
	logger.Info("load .env file successfully!")

	//	Parse env to cfg
	if err := cleanenv.ReadEnv(cfg); err != nil {
		logger.WithError(err).Error("can not Parse env variables to config")
		return
	}
	logger.Info("env variables has been parsed to config successfully!")

	//	Start connecting to MongoDB
	mongodb, err := db.ConnectDB(cfg, logger)
	if err != nil {
		logger.WithError(err).Error("can not connect to database")
		return
	}
	logger.Info("connected to mongodb successfully!")

	//	Create collections which are in database if they are view
	err, collections := mongodb.CreateCollections()
	if err != nil {
		logger.WithError(err).Error("can not create collections")
		return
	}
	mongodb.Collections = collections
	logger.Info("collections has been created or they were already been in the database!")

	//	Setup authentication
	authenticator, err := authenticate.NewAuth(mongodb, 30, logger)
	if err != nil {
		logger.WithError(err).Fatal("error in authenticator middleware setup")
	}

	//	Create the authentication modules
	authModule, err := auth.NewAuthHandlers(mongodb, logger, &ctx, authenticator)
	if err != nil {
		logger.WithError(err).Fatal("error in creating the auth module")
	}

	//	Create Gin server
	ginServer := server.NewGinServer([]server.Module{
		authModule,
	})

	err = ginServer.HttpServer.ListenAndServe()
	if err != nil {
		logger.WithError(err).Fatal("error in listening and serve")
	}
}
