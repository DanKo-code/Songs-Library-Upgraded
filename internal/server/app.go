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
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
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
	songUC     song.UseCase
}

func NewApp() (*App, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered NewApp function"))

	db := initDB()

	songRepo := songpostgres.NewSongRepository(db)
	musixMatchUseCase := songusecase.CreateMusixMatchUseCase(
		os.Getenv("MMLAPI_BASE_URL"),
		os.Getenv("MMLAPI_GET_SONG_IP_PATH"),
		os.Getenv("MMLAPI_GET_LYRICS_PATH"),
		os.Getenv("MMLAPI_API_KEY"),

		os.Getenv("GAPI_BASE_URL"),
		os.Getenv("GAPI_GET_SONG_RELEASE_DATE"),
		os.Getenv("GAPI_AUTHORIZATION"),
		&http.Client{},
	)

	//conn, err := grpc.NewClient("localhost:3025", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("127.0.0.1:3024", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &App{
		songUC: songusecase.NewSongUseCase(songRepo, musixMatchUseCase, conn),
	}, nil
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

	initData(db)

	return db
}

func initData(db *gorm.DB) {
	var count int64
	db.Model(&models.Author{}).Count(&count)
	if count == 0 {
		logrusCustom.LogWithLocation(logrus.InfoLevel, "Adding initial data to the database")

		authors := []models.Author{
			{
				ID:        uuid.New(),
				GroupName: "creedence clearwater revival",
			},
			{
				ID:        uuid.New(),
				GroupName: "михаил круг",
			},
		}

		releaseDateCastedFirstSong, _ := time.Parse("2006-01-02", "1969-11-02")
		releaseDateCastedSecondSong, _ := time.Parse("2006-01-02", "1970-12-07")
		releaseDateCastedThirdSong, _ := time.Parse("2006-01-02", "1993-11-30")

		songs := []models.Song{
			{
				ID:          uuid.New(),
				Name:        "fortunate son",
				AuthorId:    authors[0].ID,
				ReleaseDate: releaseDateCastedFirstSong,
				Text:        "some folks are born made to wave the flag\\nthey're red, white and blue\\nand when the band plays \\\"hail to the chief\\\"\\nthey point the cannon at you, lord\\n\\nit ain't me, it ain't me\\ni ain't no senator's son, son\\nit ain't me, it ain't me\\ni ain't no fortunate one\\n\\nsome folks are born, silver spoon in hand",
				Link:        "https://www.musixmatch.com/lyrics/Creedence-Clearwater-Revival/Fortunate-Son?utm_source=application&utm_campaign=api&utm_medium=DanKoKode%3A1409625027081",
			},
			{
				ID:          uuid.New(),
				Name:        "фраер",
				AuthorId:    authors[1].ID,
				ReleaseDate: releaseDateCastedThirdSong,
				Text:        "что ж ты, фраер, сдал назад\\nне по масти я тебе\\nты смотри в мои глаза\\nбрось трепаться о судьбе\\n\\nведь с тобой мой мусорок\\nя попутала рамсы\\nзавязала узелок\\nкак тугие две косы\\n\\nпомню как ты подошел\\nкак поскрипывал паркет\\nкак поставил на мой стол\\nчайных роз большой букет\\n\\nя решила ты - скокарь\\nили вор-авторитет\\nоказалось просто тварь",
				Link:        "https://www.musixmatch.com/lyrics/%D0%9C%D0%B8%D1%85%D0%B0%D0%B8%D0%BB-%D0%9A%D1%80%D1%83%D0%B3/%D0%A4%D1%80%D0%B0%D0%B5%D1%80?utm_source=application&utm_campaign=api&utm_medium=DanKoKode%3A1409625027081",
			},
			{
				ID:          uuid.New(),
				Name:        "девочка - пай",
				AuthorId:    authors[1].ID,
				ReleaseDate: releaseDateCastedSecondSong,
				Text:        "в тебе было столько желанья\\nи месяц над нами светил\\nкогда по маляве, придя на свиданье\\nя розы тебе подарил\\n\\nкакой ты казалась серьёзной\\nкачала в ответ головой\\nкогда я сказал, что отнял эти розы\\nв киоске на первой ямской\\n\\nкак было тепло, что нас с тобой вместе свело\\nдевочка-пай, рядом жиган и хулиган\\nв нашей твери нету таких даже среди шкур центровых\\nдевочка-пай, ты не грусти и не скучай\\n\\nпонты просадил я чуть позже\\nв делах узелки затянул",
				Link:        "https://www.musixmatch.com/lyrics/%D0%9C%D0%B8%D1%85%D0%B0%D0%B8%D0%BB-%D0%9A%D1%80%D1%83%D0%B3/%D0%94%D0%B5%D0%B2%D0%BE%D1%87%D0%BA%D0%B0-%D0%9F%D0%B0%D0%B9?utm_source=application&utm_campaign=api&utm_medium=DanKoKode%3A1409625027081",
			},
		}

		db.Create(&authors)
		db.Create(&songs)

		logrusCustom.LogWithLocation(logrus.InfoLevel, "Initial data added successfully")
	} else {
		logrusCustom.LogWithLocation(logrus.InfoLevel, "Data already exists in the database, skipping initialization")
	}
}
