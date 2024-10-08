package song

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
)

type Repository interface {
	GetSongs(*dtos.GetSongsDTO) ([]models.Song, error)
}
