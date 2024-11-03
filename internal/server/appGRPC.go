package server

import (
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/delivery/grpc/songGRPC"
	songpostgres "SongsLibrary/internal/song/repository/postgres"
	songusecase "SongsLibrary/internal/song/usecase"
	"SongsLibrary/internal/validators"
	logrusCustom "SongsLibrary/pkg/logger"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"os"
	"os/signal"
)

type AppGRPC struct {
	gRPCServer *grpc.Server
	gRPCClient *grpc.ClientConn
	songUC     song.UseCase
}

func NewAppGRPC() (*AppGRPC, error) {

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

	//gRPCClient
	conn, err := grpc.NewClient("localhost:3025", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &AppGRPC{
		songUC: songusecase.NewSongUseCase(songRepo, musixMatchUseCase, conn),
	}, nil
}

func (app *AppGRPC) Run(port string) error {

	validate := validator.New()
	err := validate.RegisterValidation("DateValidation", validators.DateValidation)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return err
	}

	app.gRPCServer = grpc.NewServer()

	songGRPC.Register(app.gRPCServer, validate, app.songUC)

	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Starting gRPC server on port %s", port))

	go func() {
		if err := app.gRPCServer.Serve(listen); err != nil {
			logrusCustom.Logger.Fatalf("Failed to serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("stopping gRPC server %s", port))
	app.gRPCServer.GracefulStop()

	if app.gRPCClient != nil {
		if err := app.gRPCClient.Close(); err != nil {
			logrusCustom.LogWithLocation(logrus.ErrorLevel, fmt.Sprintf("Failed to close gRPC connection: %v", err))
		} else {
			logrusCustom.LogWithLocation(logrus.InfoLevel, "gRPC connection closed successfully")
		}
	}

	return nil
}
