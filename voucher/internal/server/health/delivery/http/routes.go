package http

import (
	"net/http"

	"voucher/internal/server/health"

	"github.com/gin-gonic/gin"
)

func MapHealthRoutes(g *gin.Engine, healthHandlers health.Handlers) {
	g.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	g.GET("/liveness", healthHandlers.Liveness())
	g.GET("/readiness", healthHandlers.Readiness())
}