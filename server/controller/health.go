package controller

import (
	"net/http"

	"worker-sample/config"
	"worker-sample/server/model"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	ServiceContext *config.ServiceContext
}

func NewHealthController(ctx *config.ServiceContext) *HealthController {
	return &HealthController{
		ServiceContext: ctx,
	}
}

func (h *HealthController) Status(c *gin.Context) model.Response {
	return model.Response{
		Code: http.StatusOK,
		Msg:  "OK!",
	}
}
