package http

import (
	"context"
	"net/http"
	"time"

	"ekyc/config"
	"ekyc/internal/server/health"

	"github.com/gin-gonic/gin"
)

type healthHandlers struct {
	cfg *config.Config
	// 3-party services
}

func NewHealthHandlers(cfg *config.Config) health.Handlers {
	return &healthHandlers{cfg: cfg}
}

func (h *healthHandlers) Liveness() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	}
}

func (h *healthHandlers) Readiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		// check 3-party service

		c.JSON(http.StatusOK, gin.H{
			"status":       "OK",
			"3party":       "OK",
			"second3party": "OK",
		})
	}
}
