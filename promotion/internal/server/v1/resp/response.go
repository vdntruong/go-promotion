package resp

import (
	"net/http"
	"promotion/internal/model"

	"github.com/gin-gonic/gin"
)

func RespondData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
	return
}


func RespondDataWithCode(c *gin.Context, data interface{}, code int) {
	c.JSON(code, gin.H{
		"data": data,
	})
	return
}

func RespondError(c *gin.Context, err error) {
	if rerr, ok := err.(*model.ResponseError); ok {
		c.JSON(rerr.StatusCode, gin.H{
			"error": rerr.DescError(),
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
	return
}

func RespondErrorWithCode(c *gin.Context, err error, code int) {
	if rerr, ok := err.(*model.ResponseError); ok {
		c.JSON(code, gin.H{
			"error": rerr.DescError(),
		})
		return
	}

	c.JSON(code, gin.H{
		"error": err.Error(),
	})
	return
}

func Respond(c *gin.Context, data interface{}, err error) {
	if err != nil {
		RespondError(c, err)
		return
	}
	RespondData(c, data)
}
