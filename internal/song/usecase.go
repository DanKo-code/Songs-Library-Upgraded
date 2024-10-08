package song

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
)

type UseCase interface {
	GetSongs(*dtos.GetSongsDTO) ([]models.Song, error)
	/*GetSongText(ctx context.Context) (string, error)
	DeleteSong(ctx context.Context) (*models.Song, error)
	UpdateSong(ctx context.Context) (*models.Song, error)
	CreateSong(ctx context.Context) (*models.Song, error)*/
}
