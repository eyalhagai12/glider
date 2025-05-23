package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFuncR[Req any, Res any] func(*gin.Context, Req) (Res, error)
type HandlerFunc[Res any] func(*gin.Context) (Res, error)

func HandlerFromFunc[Req any, Res any](f HandlerFuncR[Req, Res], successCode int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Req
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := f(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(successCode, res)
	}
}
