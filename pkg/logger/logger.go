package logger

import (
	"go-auth-backend-api/internal/config/env"

	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {

	switch env.AppEnv.APP_ENV {
	case "production":
		Log, _ = zap.NewProduction()

	case "development":
		Log, _ = zap.NewDevelopment()
	default:
		panic("invalid APP_ENV: " + env.AppEnv.APP_ENV)
	}
}
