package utils

import (
	"github.com/gin-gonic/gin"
)

func ResponseHandler(c *gin.Context, httpStatus int, data interface{}, err error) {
	if err != nil {
		c.JSON(httpStatus, gin.H{
			"success": false,
			"error":   err.Error(),
			"data":    data,
		})
		return
	}
	c.JSON(httpStatus, gin.H{
		"success": true,
		"error":   nil,
		"data":    data,
	})
}
