package http

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"ekyc/internal/server/health"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type healthHandlers struct {
	db          *sql.DB
	redisClient *redis.Client
}

func NewHealthHandlers(db *sql.DB, redisClient *redis.Client) health.Handlers {
	return &healthHandlers{db: db, redisClient: redisClient}
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

		redisPingStatus := h.redisClient.Ping(ctx).Err() == nil
		dbPingStatus := h.db.PingContext(ctx) == nil

		status := redisPingStatus && dbPingStatus
		c.JSON(http.StatusOK, gin.H{
			"success":  status,
			"redis":    redisPingStatus,
			"database": dbPingStatus,
		})
	}
}
