package usecase

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/dtos"
	logrusCustom "SongsLibrary/pkg/logger"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type SongUseCase struct {
	songRepo          song.Repository
	musixMatchUseCase song.MusixmatchUseCase
}

func NewSongUseCase(songRepo song.Repository, musixMatchUseCase song.MusixmatchUseCase) *SongUseCase {
	return &SongUseCase{songRepo: songRepo, musixMatchUseCase: musixMatchUseCase}
}

func (suc *SongUseCase) GetSongs(gsdto *dtos.GetSongsDTO) ([]models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetSongs UseCase with parameters: %+v", gsdto))

	songs, err := suc.songRepo.GetSongs(gsdto)
	if err != nil {
		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting GetSongs UseCase with songs: %+v", songs))

	return songs, nil
}

func (suc *SongUseCase) DeleteSong(id uuid.UUID) (*models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered DeleteSong UseCase with parameter: %+v", id.String()))

	deletedSong, err := suc.songRepo.DeleteSong(id)
	if err != nil {
		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting DeleteSongs UseCase with deleted song: %+v", deletedSong))

	return deletedSong, nil
}

func (suc *SongUseCase) UpdateSong(fieldsToUpdate *models.Song) (*models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered UpdateSongs UseCase with parameters: %+v", fieldsToUpdate))

	updatedSong, err := suc.songRepo.UpdateSong(fieldsToUpdate)
	if err != nil {
		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting UpdateSongs UseCase with updated song: %+v", updatedSong))

	return updatedSong, nil
}

func (suc *SongUseCase) CreateSong(group, songName string) (*models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered CreateSongs UseCase with parameters: group:%s, song:%s", group, songName))

	createdSong, err := suc.songRepo.GetSongByName(songName)
	if err != nil && err.Error() != song.SongsNotFound.Error() {
		return nil, err
	}
	if createdSong != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, song.SongAlreadyExists.Error())

		return nil, song.SongAlreadyExists
	}

	ip, link, releaseDate, err := suc.musixMatchUseCase.GetSongIP(group, songName)
	if err != nil {
		return nil, err
	}

	lyrics, err := suc.musixMatchUseCase.GetLyrics(ip)
	if err != nil {
		return nil, err
	}

	releaseDateCasted, err := time.Parse(time.RFC3339, releaseDate)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return nil, err
	}

	createdSong, err = suc.songRepo.CreateSong(releaseDateCasted, group, songName, lyrics, link)
	if err != nil {
		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting CreateSongs UseCase with created song: %+v", createdSong))

	return createdSong, nil
}

func (suc *SongUseCase) GetSongLyrics(gsldto *dtos.GetSongLyricsDTO) ([]string, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetSongLyrics UseCase with parameters: %+v", gsldto))

	existingSong, err := suc.songRepo.GetSong(gsldto.Id)
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
