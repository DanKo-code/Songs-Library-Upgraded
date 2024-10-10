package main

import (
	"SongsLibrary/internal/server"
	logrus "SongsLibrary/pkg/logger"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logrus.InitLogger()

	app := server.NewApp()

	if err := app.Run(os.Getenv("APP_PORT")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
