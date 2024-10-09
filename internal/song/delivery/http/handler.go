package http

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type Handler struct {
	useCase song.UseCase
}

func NewHandler(useCase song.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) GetSongs(c *gin.Context) {
	var gsdto dtos.GetSongsDTO

	if err := c.ShouldBindQuery(&gsdto); err != nil {
		log.Println(err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if gsdto.Page == 0 {
		gsdto.Page = song.DefaultPage
	}
	if gsdto.PageSize == 0 {
		gsdto.PageSize = song.DefaultPageSize
	}

	songs, err := h.useCase.GetSongs(&gsdto)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

func (h *Handler) DeleteSong(c *gin.Context) {
	id := c.Param("id")

	convertedId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Song ID format"})
		return
	}

	deleteSong, err := h.useCase.DeleteSong(convertedId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleteSong": deleteSong})
}

func (h *Handler) UpdateSong(c *gin.Context) {
	var fieldsToUpdate models.Song
	if err := c.ShouldBindJSON(&fieldsToUpdate); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateSong, err := h.useCase.UpdateSong(&fieldsToUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"Updated Song": updateSong})
}

func (h *Handler) CreateSong(c *gin.Context) {
	var createSongDTO dtos.CreateSongDTO
	if err := c.ShouldBindJSON(&createSongDTO); err != nil {
		log.Println(err)
		return
	}

	createSong, err := h.useCase.CreateSong(createSongDTO.Group, createSongDTO.Song)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Created Song": createSong})
}
