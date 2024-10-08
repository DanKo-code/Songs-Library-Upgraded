package song

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
	"github.com/google/uuid"
)

type Repository interface {
	GetSongs(*dtos.GetSongsDTO) ([]models.Song, error)
	DeleteSong(id uuid.UUID) (*models.Song, error)
}
