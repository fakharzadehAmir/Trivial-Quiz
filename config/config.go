package config

type Config struct {
	Database struct {
		Host         string `env:"DATABASE_HOST" env-default:"localhost"`
		Port         string `env:"DATABASE_PORT" env-default:"27017"`
		DatabaseName string `env:"DATABASE_NAME" env-default:"TriviaQuiz"`
		Username     string `env:"DATABASE_USERNAME" env-default:"admin"`
		Password     string `env:"DATABASE_PASSWORD" env-default:"admin"`
	}
}
