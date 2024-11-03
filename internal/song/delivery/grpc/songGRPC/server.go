package songGRPC

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/dtos"
	logrusCustom "SongsLibrary/pkg/logger"
	"context"
	"fmt"
	songv1 "github.com/DanKo-code/Protobuf-For-Songs-Library-Upgraded/protos/gen/go/song"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type serverGRPC struct {
	songv1.UnimplementedSongServer
	validate *validator.Validate
	usecase  song.UseCase
}

func Register(gRPC *grpc.Server, validator *validator.Validate, usecase song.UseCase) {
	songv1.RegisterSongServer(gRPC, &serverGRPC{
		validate: validator,
		usecase:  usecase,
	})
}

func (s *serverGRPC) GetSongs(ctx context.Context, req *songv1.GetSongsRequest) (*songv1.GetSongsResponseList, error) {

	var gsdto dtos.GetSongsDTO

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetSongs gRPC Hanlder with parameters: %+v", gsdto))

	gsdto = dtos.GetSongsDTO{
		Id:          req.GetId(),
		Name:        req.GetName(),
		GroupName:   req.GetGroupName(),
		ReleaseDate: req.GetReleaseDate(),
		Text:        req.GetText(),
		Link:        req.GetLink(),
		Page:        int(req.GetPage()),
		PageSize:    int(req.GetPageSize()),
	}

	err := s.validate.Struct(gsdto)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return nil, status.Error(codes.InvalidArgument, "")
	}

	gsdto.SetDefaults()
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Setted default parameters: %+v", gsdto))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	gsdto.Text = strings.ToLower(gsdto.Text)

	songs, err := s.usecase.GetSongs(ctx, &gsdto)
	if err != nil {

		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		if err.Error() == song.SongsNotFound.Error() {
			return nil, status.Error(codes.NotFound, "")
		}

		return nil, status.Error(codes.Internal, "")
	}

	songsResponseList := convertSongToSongsResponseList(songs)

	return songsResponseList, nil
}

func convertSongToSongsResponseList(songsList []models.Song) *songv1.GetSongsResponseList {

	var songsResponseList songv1.GetSongsResponseList

	for _, specificSong := range songsList {
		songsResponseList.Songs = append(songsResponseList.GetSongs(), &songv1.GetSongsResponse{
			Id:          specificSong.ID.String(),
			Name:        specificSong.Name,
			AuthorId:    specificSong.AuthorId.String(),
			AuthorName:  specificSong.Author.GroupName,
			ReleaseDate: specificSong.ReleaseDate.String(),
			Text:        specificSong.Text,
			Link:        specificSong.Link,
		})
	}

	return &songsResponseList
}

func (s *serverGRPC) DeleteSong(ctx context.Context, req *songv1.DeleteSongsRequest) (*songv1.DeleteSongsResponse, error) {
	id := req.GetId()

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered DeleteSong Hanlder with parameter: %s", id))

	convertedId, err := uuid.Parse(id)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return nil, status.Error(codes.InvalidArgument, "")
	}

	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Successfully converted songId to uuid format: %s", convertedId.String()))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	deletedSong, err := s.usecase.DeleteSong(ctx, convertedId)
	if err != nil {

		if err.Error() == song.SongsNotFound.Error() {

			return nil, status.Error(codes.NotFound, "")
		}

		return nil, status.Error(codes.Internal, "")
	}

	return &songv1.DeleteSongsResponse{
		Id:          deletedSong.ID.String(),
		Name:        deletedSong.Name,
		AuthorId:    deletedSong.AuthorId.String(),
		AuthorName:  deletedSong.Author.GroupName,
		ReleaseDate: deletedSong.ReleaseDate.String(),
		Text:        deletedSong.Text,
		Link:        deletedSong.Link,
	}, nil
}
