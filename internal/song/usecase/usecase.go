package usecase

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/dtos"
	"github.com/google/uuid"
)

type SongUseCase struct {
	songRepo song.Repository
}

func NewSongUseCase(songRepo song.Repository) *SongUseCase {
	return &SongUseCase{songRepo: songRepo}
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
