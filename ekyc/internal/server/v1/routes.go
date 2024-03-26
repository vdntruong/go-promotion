package v1

import (
	"net/http"

	"ekyc/internal/usecase"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
// @title           eKYC API
// @version         1.0
// @description     This is eKYC API Spec.
// @termsOfService  http://swagger.io/terms/
// @host			localhost:3001
// @BasePath		/v1
func NewRouter(handler *gin.Engine, c usecase.UserUsecase) {
	// Swagger
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Prometheus metrics
	// handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		h.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
		newUserRoutes(h, c)
	}
}
