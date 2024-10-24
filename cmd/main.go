package main

import (
	"SongsLibrary/internal/constants"
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

	var protocol constants.Protocol = constants.Protocol(os.Getenv("PROTOCOL"))

	switch protocol {
	case constants.HTTP:
		startHTTPServer()
	case constants.GRPC:
		startGRPCServer()
	}
}

func startHTTPServer() {
	app := server.NewApp()

	if err := app.Run(os.Getenv("APP_PORT")); err != nil {
		logrusCustom.Logger.Fatalf("Error when running server: %s", err.Error())
	}
}

func startGRPCServer() {
	app := server.NewAppGRPC()

	err := app.Run(os.Getenv("APP_PORT"))
	if err != nil {
		logrusCustom.Logger.Fatalf("Error when running server: %s", err.Error())
	}
}
