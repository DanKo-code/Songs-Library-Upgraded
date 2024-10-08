package song

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
	"github.com/google/uuid"
)

type UseCase interface {
	GetSongs(*dtos.GetSongsDTO) ([]models.Song, error)
	DeleteSong(id uuid.UUID) (*models.Song, error)
	/*GetSongText(ctx context.Context) (string, error)
	UpdateSong(ctx context.Context) (*models.Song, error)
	CreateSong(ctx context.Context) (*models.Song, error)*/
}
