package postgres

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestGetSongsRepository_Success(t *testing.T) {
	mockRepo := new(MockRepository)

	// Настройка тестовых данных
	gsdto := &dtos.GetSongsDTO{
		Name:      "testsong",
		GroupName: "testgroup",
		Page:      1,
		PageSize:  10,
	}

	// Подготовка мок-данных
	mockSongs := []models.Song{
		{
			ID:          uuid.New(),
			Name:        "testsong",
			AuthorId:    uuid.New(),
			Author:      models.Author{GroupName: "testgroup"},
			ReleaseDate: time.Now(),
			Text:        "Some lyrics",
			Link:        "http://example.com",
		},
	}

	// Настройка мока для возврата данных
	mockRepo.On("GetSongs", mock.Anything, gsdto).Return(mockSongs, nil)

	ctx := context.Background()
	songs, err := mockRepo.GetSongs(ctx, gsdto)

	assert.NoError(t, err)
	assert.Equal(t, mockSongs, songs)
	mockRepo.AssertExpectations(t)
}
