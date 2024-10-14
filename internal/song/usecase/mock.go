package usecase

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSongUseCase struct {
	mock.Mock
}

func (m *MockSongUseCase) GetSongs(ctx context.Context, gsdto *dtos.GetSongsDTO) ([]models.Song, error) {
	args := m.Called(ctx, gsdto)
	if songs, ok := args.Get(0).([]models.Song); ok {
		return songs, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSongUseCase) DeleteSong(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	args := m.Called(ctx, id)
	if song, ok := args.Get(0).(*models.Song); ok {
		return song, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSongUseCase) UpdateSong(ctx context.Context, song *models.Song) (*models.Song, error) {
	args := m.Called(ctx, song)
	if updatedSong, ok := args.Get(0).(*models.Song); ok {
		return updatedSong, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSongUseCase) CreateSong(ctx context.Context, group, song string) (*models.Song, error) {
	args := m.Called(ctx, group, song)
	if newSong, ok := args.Get(0).(*models.Song); ok {
		return newSong, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSongUseCase) GetSongLyrics(ctx context.Context, dto *dtos.GetSongLyricsDTO) ([]string, error) {
	args := m.Called(ctx, dto)
	if lyrics, ok := args.Get(0).([]string); ok {
		return lyrics, args.Error(1)
	}
	return nil, args.Error(1)
}
