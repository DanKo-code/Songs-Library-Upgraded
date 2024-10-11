package server

import (
	_ "SongsLibrary/docs"
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	songhttp "SongsLibrary/internal/song/delivery/http"
	songpostgres "SongsLibrary/internal/song/repository/postgres"
	songusecase "SongsLibrary/internal/song/usecase"
	"SongsLibrary/internal/validators"
	logrusCustom "SongsLibrary/pkg/logger"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer *http.Server

	songUC song.UseCase
}

func NewApp() *App {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered NewApp function"))

	db := initDB()

	songRepo := songpostgres.NewSongRepository(db)
	musixMatchUseCase := songusecase.CreateMusixMatchUseCase(
		os.Getenv("MMLAPI_BASE_URL"),
		os.Getenv("MMLAPI_GET_SONG_IP_PATH"),
		os.Getenv("MMLAPI_GET_LYRICS_PATH"),
		os.Getenv("MMLAPI_API_KEY"),
	)

	return &App{
		songUC: songusecase.NewSongUseCase(songRepo, musixMatchUseCase),
	}
}

func (a *App) Run(port string) error {
	router := gin.Default()

	validate := validator.New()
	err := validate.RegisterValidation("DateValidation", validators.DateValidation)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return err
	}

	songhttp.RegisterHTTPEndpoints(router, a.songUC, validate)

	router.GET(os.Getenv("SWAGGER_PATH"), ginSwagger.WrapHandler(swaggerFiles.Handler))

	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Starting server on port %s", port))

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			logrusCustom.Logger.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit
	logrusCustom.LogWithLocation(logrus.InfoLevel, "Gracefully shutting down server...")

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB() *gorm.DB {

	logrusCustom.LogWithLocation(logrus.InfoLevel, "Entered initDB function")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SLLMODE"),
	)

	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Loaded dsn: %s", dsn))

	logrusCustom.LogWithLocation(logrus.InfoLevel, "Connecting to db")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logrusCustom.Logger.Fatalf("Failed to connect database")
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, "Successfully connected to db")

	logrusCustom.LogWithLocation(logrus.InfoLevel, "Starting migrating db")

	err = db.AutoMigrate(&models.Song{})
	if err != nil {
		logrusCustom.Logger.Fatalf("Failed migrate db")
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, "Successfully migrated to db")

	return db
}
