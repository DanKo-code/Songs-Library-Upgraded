package usecase

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/dtos"
	postgres "SongsLibrary/internal/song/repository/postgres"
	logrusCustom "SongsLibrary/pkg/logger"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"time"
)

func TestGetSongsUseCase_Success(t *testing.T) {
	mockRepo := new(postgres.MockRepository)
	logrusCustom.InitLogger()

	musixMatchUseCase := CreateMusixMatchUseCase(
		"https://api.musixmatch.com/ws/1.1/",
		"track.search?q_artist=%s&q_track=%s&apikey=%s",
		"track.lyrics.get?commontrack_id=%s&apikey=%s",
		"738351c40ea2c5b8ab58f24665ce3bc0",

		"https://api.genius.com/",
		"RELEASE_DATE='search?q=%s %s",
		"Bearer biWMNIIRWL62TeubvtPoP6JIeLHa1r7rDDG65ry0oXM0HagiD7YewmXChSCIMPy-",
		&http.Client{},
	)

	suc := NewSongUseCase(mockRepo, musixMatchUseCase)

	gsdto := &dtos.GetSongsDTO{
		Name:      "testsong",
		GroupName: "testgroup",
		Page:      1,
		PageSize:  10,
	}

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

	mockRepo.On("GetSongs", mock.Anything, gsdto).Return(mockSongs, nil)

	ctx := context.Background()
	songs, err := suc.GetSongs(ctx, gsdto)

	assert.NoError(t, err)
	assert.Equal(t, mockSongs, songs)

	mockRepo.AssertExpectations(t)
}

func TestGetSongsUseCase_Fail(t *testing.T) {
	mockRepo := new(postgres.MockRepository)
	logrusCustom.InitLogger()

	musixMatchUseCase := CreateMusixMatchUseCase(
		"https://api.musixmatch.com/ws/1.1/",
		"track.search?q_artist=%s&q_track=%s&apikey=%s",
		"track.lyrics.get?commontrack_id=%s&apikey=%s",
		"738351c40ea2c5b8ab58f24665ce3bc0",

		"https://api.genius.com/",
		"RELEASE_DATE='search?q=%s %s",
		"Bearer biWMNIIRWL62TeubvtPoP6JIeLHa1r7rDDG65ry0oXM0HagiD7YewmXChSCIMPy-",
		&http.Client{},
	)

	suc := NewSongUseCase(mockRepo, musixMatchUseCase)

	gsdto := &dtos.GetSongsDTO{
		Name:      "testsong",
		GroupName: "testgroup",
		Page:      1,
		PageSize:  10,
	}

	mockRepo.On("GetSongs", mock.Anything, gsdto).Return(nil, song.SongsNotFound)

	ctx := context.Background()
	songRes, err := suc.GetSongs(ctx, gsdto)

	assert.Nil(t, songRes)
	assert.Error(t, err, song.AuthorNotFound)
}
