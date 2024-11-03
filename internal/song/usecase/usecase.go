package usecase

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/dtos"
	logrusCustom "SongsLibrary/pkg/logger"
	"context"
	"fmt"
	songv1pb "github.com/DanKo-code/Protobuf-For-Songs-Library-Upgraded/protos/gen/go/song"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"strings"
	"time"
)

type SongUseCase struct {
	songRepo          song.Repository
	musixMatchUseCase song.MusixmatchUseCase
	gRPCClient        *grpc.ClientConn
}

func NewSongUseCase(songRepo song.Repository, musixMatchUseCase song.MusixmatchUseCase, gRPCClient *grpc.ClientConn) *SongUseCase {
	return &SongUseCase{songRepo: songRepo, musixMatchUseCase: musixMatchUseCase, gRPCClient: gRPCClient}
}

func (suc *SongUseCase) GetSongs(ctx context.Context, gsdto *dtos.GetSongsDTO) ([]models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetSongs UseCase with parameters: %+v", gsdto))

	songs, err := suc.songRepo.GetSongs(ctx, gsdto)
	if err != nil {
		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting GetSongs UseCase with songs: %+v", songs))

	return songs, nil
}

func (suc *SongUseCase) DeleteSong(ctx context.Context, id uuid.UUID) (*models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered DeleteSong UseCase with parameter: %+v", id.String()))

	deletedSong, err := suc.songRepo.DeleteSong(ctx, id)
	if err != nil {
		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting DeleteSongs UseCase with deleted song: %+v", deletedSong))

	return deletedSong, nil
}

func (suc *SongUseCase) UpdateSong(ctx context.Context, fieldsToUpdate *models.Song) (*models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered UpdateSongs UseCase with parameters: %+v", fieldsToUpdate))

	updatedSong, err := suc.songRepo.UpdateSong(ctx, fieldsToUpdate)
	if err != nil {
		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting UpdateSongs UseCase with updated song: %+v", updatedSong))

	return updatedSong, nil
}

func (suc *SongUseCase) CreateSong(ctx context.Context, group, songName string) (*models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered CreateSongs UseCase with parameters: group:%s, song:%s", group, songName))

	author, err := suc.songRepo.GetAuthorByName(ctx, songName)
	if err != nil && err.Error() != song.AuthorNotFound.Error() {
		return nil, err
	}
	if author != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, song.AuthorAlreadyExists.Error())

		return nil, song.AuthorAlreadyExists
	}

	//monolith
	/*ip, link, releaseDate, trackName, artistName, err := suc.musixMatchUseCase.GetSongData(ctx, group, songName)
	if err != nil {
		return nil, err
	}*/

	//micro services
	songv1pbClient := songv1pb.NewSongDataClient(suc.gRPCClient)

	getSongDataResponse, err := songv1pbClient.GetSongData(ctx, &songv1pb.GetSongDataRequest{Group: group, SongName: songName})
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())
		return nil, err
	}

	ip := getSongDataResponse.GetIp()
	link := getSongDataResponse.GetLink()
	releaseDate := getSongDataResponse.GetReleaseDate()
	trackName := getSongDataResponse.GetTrackName()
	artistName := getSongDataResponse.GetArtistName()

	lyrics, err := suc.musixMatchUseCase.GetLyrics(ctx, ip)
	if err != nil {
		return nil, err
	}

	releaseDateCasted, err := time.Parse("2006-01-02 15:04:05 -0700 MST", releaseDate)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return nil, err
	}

	createdSong, err := suc.songRepo.CreateSong(ctx, releaseDateCasted, artistName, trackName, lyrics, link)
	if err != nil {
		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting CreateSongs UseCase with created song: %+v", createdSong))

	return createdSong, nil
}

func (suc *SongUseCase) GetSongLyrics(ctx context.Context, gsldto *dtos.GetSongLyricsDTO) ([]string, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetSongLyrics UseCase with parameters: %+v", gsldto))

	existingSong, err := suc.songRepo.GetSong(ctx, gsldto.Id)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return nil, err
	}

	if existingSong.Text == "" {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, song.ErrorGetSongLyrics.Error())

		return nil, song.ErrorGetSongLyrics
	}

	verses := strings.Split(existingSong.Text, "\n\n")

	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Exiting GetSongLyrics UseCase with song lyrics verses: %+v", verses))

	offset := (gsldto.Page - 1) * gsldto.PageSize
	if offset > len(verses) {
		return []string{}, nil
	}

	end := offset + gsldto.PageSize
	if end > len(verses) {
		end = len(verses)
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting GetSongLyrics UseCase with song lyrics verses: %+v", verses[offset:end]))

	return verses[offset:end], nil
}
