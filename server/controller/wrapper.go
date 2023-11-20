package controller

import (
	"github.com/gin-gonic/gin"
	"worker-sample/server/model"
)

type Controller func(ctx *gin.Context) model.Response

func Wrap(c Controller) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := c(ctx)
		ctx.JSON(resp.Code, resp)
	}
}
