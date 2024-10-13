package song

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
	"context"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	GetSongs(context.Context, *dtos.GetSongsDTO) ([]models.Song, error)
	DeleteSong(ctx context.Context, id uuid.UUID) (*models.Song, error)
	UpdateSong(context.Context, *models.Song) (*models.Song, error)
	CreateSong(ctx context.Context, releaseDate time.Time, group string, songName string, lyrics string, link string) (*models.Song, error)
	GetSong(context.Context, uuid.UUID) (*models.Song, error)
	GetAuthorByName(ctx context.Context, authorName string) (*models.Author, error)
}
