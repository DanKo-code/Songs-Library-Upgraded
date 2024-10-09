package usecase

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/dtos"
	"github.com/google/uuid"
)

type SongUseCase struct {
	songRepo          song.Repository
	musixMatchUseCase song.MusixmatchUseCase
}

func NewSongUseCase(songRepo song.Repository, musixMatchUseCase song.MusixmatchUseCase) *SongUseCase {
	return &SongUseCase{songRepo: songRepo, musixMatchUseCase: musixMatchUseCase}
}

func (suc *SongUseCase) GetSongs(gsdto *dtos.GetSongsDTO) ([]models.Song, error) {

	songs, err := suc.songRepo.GetSongs(gsdto)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (suc *SongUseCase) DeleteSong(id uuid.UUID) (*models.Song, error) {
	deletedSong, err := suc.songRepo.DeleteSong(id)
	if err != nil {
		return nil, err
	}

	return deletedSong, nil
}

func (suc *SongUseCase) UpdateSong(fieldsToUpdate *models.Song) (*models.Song, error) {
	updatedSong, err := suc.songRepo.UpdateSong(fieldsToUpdate)
	if err != nil {
		return nil, err
	}

	return updatedSong, nil
}

func (suc *SongUseCase) CreateSong(group, song string) (*models.Song, error) {

	//send req for enrichment
	ip, link, releaseDate, err := suc.musixMatchUseCase.GetSongIP(group, song)
	if err != nil {
		return nil, err
	}

	lyrics, err := suc.musixMatchUseCase.GetLyrics(ip)
	if err != nil {
		return nil, err
	}

	createdSong, err := suc.songRepo.CreateSong(group, song, lyrics, link, releaseDate)
	if err != nil {
		return nil, err
	}

	return createdSong, nil
}
