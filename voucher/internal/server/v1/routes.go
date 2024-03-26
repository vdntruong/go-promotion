package v1

import (
	"net/http"

	"voucher/internal/usecase"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(handler *gin.Engine,  v usecase.VoucherUsecase) {
	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)
	// Prometheus metrics
	// handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		h.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
		newVoucherRoutes(h, v)
	}
}