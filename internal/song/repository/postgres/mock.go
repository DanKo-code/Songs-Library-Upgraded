package postgres

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetSongs(ctx context.Context, gsdto *dtos.GetSongsDTO) ([]models.Song, error) {
	args := m.Called(ctx, gsdto)
	if songs, ok := args.Get(0).([]models.Song); ok {
		return songs, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) DeleteSong(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	args := m.Called(ctx, id)
	if song, ok := args.Get(0).(*models.Song); ok {
		return song, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) UpdateSong(ctx context.Context, song *models.Song) (*models.Song, error) {
	args := m.Called(ctx, song)
	if updatedSong, ok := args.Get(0).(*models.Song); ok {
		return updatedSong, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) CreateSong(ctx context.Context, releaseDate time.Time, group string, songName string, lyrics string, link string) (*models.Song, error) {
	args := m.Called(ctx, releaseDate, group, songName, lyrics, link)
	if newSong, ok := args.Get(0).(*models.Song); ok {
		return newSong, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) GetSong(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	args := m.Called(ctx, id)
	if song, ok := args.Get(0).(*models.Song); ok {
		return song, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) GetAuthorByName(ctx context.Context, authorName string) (*models.Author, error) {
	args := m.Called(ctx, authorName)
	if author, ok := args.Get(0).(*models.Author); ok {
		return author, args.Error(1)
	}
	return nil, args.Error(1)
}
