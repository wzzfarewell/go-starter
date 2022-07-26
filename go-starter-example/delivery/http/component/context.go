package component

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
