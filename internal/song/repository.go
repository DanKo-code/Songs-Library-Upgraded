package song

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	GetSongs(*dtos.GetSongsDTO) ([]models.Song, error)
	DeleteSong(id uuid.UUID) (*models.Song, error)
	UpdateSong(*models.Song) (*models.Song, error)
	CreateSong(releaseDate time.Time, group string, songName string, lyrics string, link string) (*models.Song, error)
	GetSong(id uuid.UUID) (*models.Song, error)
	GetSongByName(name string) (*models.Song, error)
}
