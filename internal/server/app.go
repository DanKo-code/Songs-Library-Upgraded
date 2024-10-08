package server

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	songhttp "SongsLibrary/internal/song/delivery/http"
	songpostgres "SongsLibrary/internal/song/repository/postgres"
	songusecase "SongsLibrary/internal/song/usecase"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
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
	db := initDB()

	songRepo := songpostgres.NewSongRepository(db)

	return &App{
		songUC: songusecase.NewSongUseCase(songRepo),
	}
}

func (a *App) Run(port string) error {
	router := gin.Default()

	songhttp.RegisterHTTPEndpoints(router, a.songUC)

	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SLLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.Song{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}

	return db
}
