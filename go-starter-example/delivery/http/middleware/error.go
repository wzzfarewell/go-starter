package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wzzfarewell/go-mod/infrastructure/logger"
	"go.uber.org/zap"
)

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(404, gin.H{
			"error": "resource not found",
		})
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			if err != nil {
				logger.Error("http error", zap.Error(err))
				c.AbortWithStatusJSON(400, gin.H{
					"error": err.Error(),
				})
			}
		}
	}
}
