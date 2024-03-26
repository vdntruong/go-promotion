package health

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	Liveness() gin.HandlerFunc
	Readiness() gin.HandlerFunc
}
