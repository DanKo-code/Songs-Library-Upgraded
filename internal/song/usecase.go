package song

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
	"github.com/google/uuid"
)

type UseCase interface {
	GetSongs(*dtos.GetSongsDTO) ([]models.Song, error)
	DeleteSong(id uuid.UUID) (*models.Song, error)
	UpdateSong(*models.Song) (*models.Song, error)
	CreateSong(group, song string) (*models.Song, error)
}

type MusixmatchUseCase interface {
	GetSongIP(groupName, song string) (string, string, string, error)
	GetLyrics(ip string) (string, error)
}
