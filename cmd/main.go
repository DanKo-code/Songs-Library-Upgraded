package main

import (
	"SongsLibrary/internal/server"
	logrusCustom "SongsLibrary/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logrusCustom.InitLogger()

	err := godotenv.Load()
	if err != nil {
		logrusCustom.Logger.Fatalf("Error loading .env file")
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, "Successfully loaded environment variables")

	app := server.NewApp()

	if err := app.Run(os.Getenv("APP_PORT")); err != nil {
		logrusCustom.Logger.Fatalf("Error when running server: %s", err.Error())
	}
}
