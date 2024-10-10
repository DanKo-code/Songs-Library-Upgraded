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
