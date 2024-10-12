package song

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
	"context"
	"github.com/google/uuid"
)

type UseCase interface {
	GetSongs(context.Context, *dtos.GetSongsDTO) ([]models.Song, error)
	DeleteSong(ctx context.Context, id uuid.UUID) (*models.Song, error)
	UpdateSong(context.Context, *models.Song) (*models.Song, error)
	CreateSong(ctx context.Context, group, song string) (*models.Song, error)
	GetSongLyrics(ctx context.Context, dto *dtos.GetSongLyricsDTO) ([]string, error)
}

type MusixmatchUseCase interface {
	GetSongIP(ctx context.Context, groupName, song string) (string, string, string, error)
	GetLyrics(ctx context.Context, ip string) (string, error)
}
