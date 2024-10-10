package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	Logger.SetLevel(logrus.DebugLevel)

	Logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.SetOutput(file)
	} else {
		Logger.SetOutput(os.Stdout)
		return
	}

	multiWriter := io.MultiWriter(file, os.Stdout)
	Logger.SetOutput(multiWriter)
}

func LogWithLocation(level logrus.Level, msg string) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		Logger.Log(level, msg)
		return
	}

	rootDir, err := os.Getwd()
	if err != nil {
		Logger.Log(level, msg)
		return
	}

	relPath, err := filepath.Rel(rootDir, file)
	if err != nil {
		Logger.Log(level, msg)
		return
	}

	relPath = filepath.ToSlash(relPath)

	Logger.WithFields(
		logrus.Fields{
			"location": fmt.Sprintf("%s:%d", relPath, line),
		}).Log(level, msg)
}
