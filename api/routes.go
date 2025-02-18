package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/songs", GetSongsHandler)
		api.GET("/songs/:id/verses", GetSongVersesHandler)
		api.POST("/songs", AddSongHandler)
		api.PUT("/songs/:id", UpdateSongHandler)
		api.DELETE("/songs/:id", DeleteSongHandler)
	}
}
