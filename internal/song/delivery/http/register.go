package http

import (
	"SongsLibrary/internal/song"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc song.UseCase, validator *validator.Validate) {
	h := NewHandler(uc, validator)

	authEndPoints := router.Group("/api")
	{
		authEndPoints.GET("/songs", h.GetSongs)
		authEndPoints.DELETE("/songs/:id", h.DeleteSong)
		authEndPoints.PUT("/songs", h.UpdateSong)
		authEndPoints.POST("/songs", h.CreateSong)
		authEndPoints.GET("/songs/:id/lyrics", h.GetSongLyrics)
	}
}
