package http

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/dtos"
	"SongsLibrary/internal/song/usecase"
	"SongsLibrary/internal/validators"
	logrusCustom "SongsLibrary/pkg/logger"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setup() (*gin.Engine, *usecase.MockSongUseCase, *validator.Validate) {
	gin.SetMode(gin.TestMode)
	logrusCustom.InitLogger()

	mockUseCase := new(usecase.MockSongUseCase)

	validate := validator.New()
	err := validate.RegisterValidation("DateValidation", validators.DateValidation)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())
	}

	r := gin.Default()
	RegisterHTTPEndpoints(r, mockUseCase, validate)

	return r, mockUseCase, validate
}

func performRequest(r *gin.Engine, method, url string) (*httptest.ResponseRecorder, error) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	r.ServeHTTP(w, req)
	return w, nil
}

func TestGetSongsHandler_Success(t *testing.T) {
	r, mockUseCase, _ := setup()

	gsdto := dtos.GetSongsDTO{
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
			ReleaseDate: time.Date(1981, 9, 23, 0, 0, 0, 0, time.UTC),
			Text:        "Test text",
			Link:        "http://example.com",
		},
	}

	mockUseCase.On("GetSongs", mock.Anything, &gsdto).Return(mockSongs, nil)

	w, err := performRequest(r, http.MethodGet, "/api/songs?name=testsong&group_name=testgroup&page=1&page_size=10")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]models.Song
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, mockSongs, response["songs"])
}

func TestGetSongsHandler_FailureSongsNotFound(t *testing.T) {
	r, mockUseCase, _ := setup()

	gsdto := dtos.GetSongsDTO{
		Name:      "testsong",
		GroupName: "testgroup",
		Page:      1,
		PageSize:  10,
	}

	mockUseCase.On("GetSongs", mock.Anything, &gsdto).Return(nil, song.SongsNotFound)

	w, err := performRequest(r, http.MethodGet, "/api/songs?name=testsong&group_name=testgroup&page=1&page_size=10")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, song.SongsNotFound.Error(), response["error"])
}

func TestGetSongsHandler_BindQueryError(t *testing.T) {
	r, _, _ := setup()

	w, err := performRequest(r, http.MethodGet, "/api/songs?link=123")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, song.InvalidInputData.Error(), response["error"])
}

func TestGetSongsHandler_ValidateQueryError(t *testing.T) {
	r, _, _ := setup()

	w, err := performRequest(r, http.MethodGet, "/api/songs?release_date=2026-01-01")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, song.InvalidInputData.Error(), response["error"])
}

func TestGetSongsHandler_TimeoutFailure(t *testing.T) {
	r, mockUseCase, _ := setup()

	gsdto := dtos.GetSongsDTO{
		Name:      "testsong",
		GroupName: "testgroup",
		Page:      1,
		PageSize:  10,
	}

	mockUseCase.On("GetSongs", mock.Anything, &gsdto).Run(func(args mock.Arguments) {
		time.Sleep(6 * time.Second)
	}).Return(nil, context.DeadlineExceeded)

	w, err := performRequest(r, http.MethodGet, "/api/songs?name=testsong&group_name=testgroup&page=1&page_size=10")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
