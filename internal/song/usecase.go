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
}
