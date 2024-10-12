package http

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/dtos"
	logrusCustom "SongsLibrary/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	useCase  song.UseCase
	validate *validator.Validate
}

func NewHandler(useCase song.UseCase, validate *validator.Validate) *Handler {
	return &Handler{
		useCase:  useCase,
		validate: validate,
	}
}

// GetSongs
// @Summary Retrieve a list of songs
// @Description Fetch a list of songs from the library with filtering options such as name, group name, release date, text, link, and pagination. Each song can be filtered based on the available query parameters.
// @Tags Songs
// @Produce  json
// @Param name query string false "Name of the song" maxlength(100)
// @Param group_name query string false "Name of the group" maxlength(100)
// @Param release_date query string false "Release date of the song" format(date)
// @Param text query string false "Lyrics of the song" maxlength(10000)
// @Param link query string false "Link to the song" format(url)
// @Param page query int false "Page number for pagination" minimum(1)
// @Param page_size query int false "Number of songs per page" minimum(1) maximum(100)
// @Success 200 {array} models.Song "List of songs"
// @Failure 400 {object} string "Invalid input data"
// @Failure 404 {object} string "Songs not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/songs [get]
func (h *Handler) GetSongs(c *gin.Context) {
	var gsdto dtos.GetSongsDTO

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetSongs Hanlder with parameters: %+v", gsdto))

	if err := c.ShouldBindQuery(&gsdto); err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidInputData.Error()})
		return
	}
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Successfully binded song parameters: %+v", gsdto))

	err := h.validate.Struct(gsdto)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidInputData.Error()})
		return
	}
	logrusCustom.LogWithLocation(logrus.InfoLevel, "Successfully validated parameters")

	gsdto.SetDefaults()
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Setted default parameters: %+v", gsdto))

	songs, err := h.useCase.GetSongs(&gsdto)
	if err != nil {

		if err.Error() == song.SongsNotFound.Error() {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": song.SongsNotFound.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": song.SongsNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

// DeleteSong
// @Summary Delete a song by its ID
// @Description Remove a song from the library using its UUID. The song ID should be in UUID format.
// @Tags Songs
// @Produce json
// @Param id path string true "UUID of the song to delete" format(uuid)
// @Success 200 {object} models.Song "Deleted song details"
// @Failure 400 {object} string "Invalid song ID format"
// @Failure 404 {object} string "Song not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/songs/{id} [delete]
func (h *Handler) DeleteSong(c *gin.Context) {
	id := c.Param("id")

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered DeleteSong Hanlder with parameter: %s", id))

	convertedId, err := uuid.Parse(id)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidSongIdFormat.Error()})
		return
	}

	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Successfully converted songId to uuid format: %s", convertedId.String()))

	deleteSong, err := h.useCase.DeleteSong(convertedId)
	if err != nil {

		if err.Error() == song.SongsNotFound.Error() {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": song.SongsNotFound.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleteSong": deleteSong})
}

// UpdateSong
// @Summary Update a song by its ID
// @Description Update the details of a song in the library using its UUID. The song ID should be in UUID format. The request body should contain the fields to be updated.
// @Tags Songs
// @Produce json
// @Param id path string true "UUID of the song to update" format(uuid)
// @Param fieldsToUpdate body dtos.UpdateSongsDTO true "Fields to update"
// @Success 200 {object} models.Song "Updated song details"
// @Failure 400 {object} string "Invalid input data"
// @Failure 404 {object} string "Song not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/songs/{id} [put]
func (h *Handler) UpdateSong(c *gin.Context) {

	var fieldsToUpdate dtos.UpdateSongsDTO
	id := c.Param("id")
	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered UpdateSong Hanlder with parameter: id: , %+v", fieldsToUpdate))

	convertedId, err := uuid.Parse(id)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidSongIdFormat.Error()})
		return
	}
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Successfully converted songId to uuid format: %s", convertedId.String()))

	if err := c.ShouldBindJSON(&fieldsToUpdate); err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidInputData.Error()})
		return
	}
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Successfully binded song parameters: %+v", fieldsToUpdate))

	err = h.validate.Struct(fieldsToUpdate)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidInputData.Error()})
		return
	}
	logrusCustom.LogWithLocation(logrus.InfoLevel, "Successfully validated parameters")

	var releaseDateCasted time.Time
	if fieldsToUpdate.ReleaseDate != "" {
		releaseDateCasted, err = time.Parse(time.RFC3339, fieldsToUpdate.ReleaseDate)
		if err != nil {
			logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidInputData.Error()})
			return
		}
	}

	var songToUpdate models.Song = models.Song{
		ID:          convertedId,
		Name:        fieldsToUpdate.Name,
		GroupName:   fieldsToUpdate.GroupName,
		Text:        fieldsToUpdate.Text,
		Link:        fieldsToUpdate.Link,
		ReleaseDate: releaseDateCasted}

	updateSong, err := h.useCase.UpdateSong(&songToUpdate)
	if err != nil {

		if err.Error() == song.SongsNotFound.Error() {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": song.SongsNotFound.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Updated Song": updateSong})
}

// CreateSong
// @Summary Create a new song
// @Description Create a new song in the library by providing song details in the request body. The group and song name will be converted to lowercase before saving.
// @Tags Songs
// @Produce json
// @Param createSongDTO body dtos.CreateSongDTO true "Details of the song to create"
// @Success 200 {object} models.Song "Created song details"
// @Failure 400 {object} string "Invalid input data"
// @Failure 404 {object} string "Data not found"
// @Failure 409 {object} string "Song already exists"
// @Failure 500 {object} string "Internal server error"
// @Router /api/songs [post]
func (h *Handler) CreateSong(c *gin.Context) {
	var createSongDTO dtos.CreateSongDTO

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered CreateSongs Hanlder with parameters: %+v", createSongDTO))

	if err := c.ShouldBindJSON(&createSongDTO); err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidInputData.Error()})
		return
	}
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Successfully binded song parameters: %+v", createSongDTO))

	err := h.validate.Struct(createSongDTO)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidInputData.Error()})
		return
	}
	logrusCustom.LogWithLocation(logrus.InfoLevel, "Successfully validated parameters")

	createSong, err := h.useCase.CreateSong(strings.ToLower(createSongDTO.Group), strings.ToLower(createSongDTO.Song))
	if err != nil {

		if err.Error() == song.ErrorGetSongData.Error() {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": song.ErrorGetSongData.Error()})
			return
		}

		if err.Error() == song.ErrorGetSongLyrics.Error() {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": song.ErrorGetSongLyrics.Error()})
			return
		}

		if err.Error() == song.SongAlreadyExists.Error() {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": song.SongAlreadyExists.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Created Song": createSong})
}

// GetSongLyrics
// @Summary Retrieve lyrics of a song
// @Description Fetch the lyrics of a specific song identified by its ID. Optional query parameters can be used to filter the results further.
// @Tags Songs
// @Produce json
// @Param id path string true "ID of the song"
// @Param page query int false "Page number for pagination" minimum(1)
// @Param page_size query int false "Number of songs per page" minimum(1) maximum(100)
// @Success 200 {object} string "Lyrics of the song"
// @Failure 400 {object} string "Invalid input data"
// @Failure 404 {object} string "Song not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/songs/{id}/lyrics [get]
func (h *Handler) GetSongLyrics(c *gin.Context) {

	var gsldtp dtos.GetSongLyricsDTO
	id := c.Param("id")
	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetSongLyrics Hanlder with parameter: id: , %s", id))

	convertedId, err := uuid.Parse(id)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidSongIdFormat.Error()})
		return
	}
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Successfully converted songId to uuid format: %s", convertedId.String()))

	if err := c.ShouldBindQuery(&gsldtp); err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidInputData.Error()})
		return
	}
	gsldtp.Id = convertedId

	err = h.validate.Struct(gsldtp)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": song.InvalidInputData.Error()})
		return
	}
	logrusCustom.LogWithLocation(logrus.InfoLevel, "Successfully validated parameters")

	gsldtp.SetDefaults()

	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Setted default parameters: %+v", gsldtp))

	lyrics, err := h.useCase.GetSongLyrics(&gsldtp)
	if err != nil {

		if err.Error() == song.ErrorGetSongData.Error() {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": song.ErrorGetSongData.Error()})
			return
		}

		if err.Error() == song.ErrorGetSongLyrics.Error() {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": song.ErrorGetSongData.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": song.SongsNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lyrics": lyrics})
}
