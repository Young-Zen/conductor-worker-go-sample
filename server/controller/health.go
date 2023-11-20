package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"worker-sample/config"
	"worker-sample/server/model"
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
