package http

import (
	"SongsLibrary/internal/song"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc song.UseCase) {
	h := NewHandler(uc)

	authEndPoints := router.Group("/api")
	{
		authEndPoints.GET("/songs", h.GetSongs)
		authEndPoints.DELETE("/songs/:id", h.DeleteSong)
		authEndPoints.PUT("/songs", h.UpdateSong)

	}

}
