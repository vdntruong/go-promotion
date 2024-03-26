package http

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"voucher/internal/server/health"

	"github.com/gin-gonic/gin"
)

type healthHandlers struct {
	db          *sql.DB
}

func NewHealthHandlers(db *sql.DB) health.Handlers {
	return &healthHandlers{db: db}
}

func (h *healthHandlers) Liveness() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	}
}

func (h *healthHandlers) Readiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		dbPingStatus := h.db.PingContext(ctx) == nil

		c.JSON(http.StatusOK, gin.H{
			"database": dbPingStatus,
		})
	}
}
