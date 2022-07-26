package transport

const ErrorHandlerMiddlewareTmpl = `package middleware

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

`

const ContextComponentTmpl = `package component

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HandlerWrapper func(c *ContextWrapper) error

type ContextWrapper struct {
	*gin.Context
}

func (c *ContextWrapper) PathID() (int, error) {
	return strconv.Atoi(c.Param("id"))
}

func (c *ContextWrapper) Success(body any) {
	c.JSON(http.StatusOK, body)
}

func Wrap(f HandlerWrapper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := f(&ContextWrapper{Context: ctx}); err != nil {
			_ = ctx.Error(err)
		}
	}
}

`
