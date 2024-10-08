package http

import (
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/dtos"
	"github.com/gin-gonic/gin"
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
