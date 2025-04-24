package controller

import (
	"net/http"
	"strconv"
	"worker-sample/config"
	"worker-sample/server/dao"
	"worker-sample/server/form"
	"worker-sample/server/model"
	"worker-sample/server/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	ServiceContext *config.ServiceContext
	UserService    service.UserService
}

func NewUserController(ctx *config.ServiceContext) *UserController {
	userDao := dao.UserRepo{
		ServiceContext: ctx,
	}
	userService := service.UserServiceImpl{
		ServiceContext: ctx,
		UserDao:        &userDao,
	}
	return &UserController{
		ServiceContext: ctx,
		UserService:    &userService,
	}
}

func (c *UserController) AddUser(ctx *gin.Context) model.Response {
	var req form.AddUserReq
	if err := ctx.ShouldBind(&req); err != nil {
		return model.Response{Code: http.StatusBadRequest, Msg: "Request param error", Err: err}
	}

	user, err := c.UserService.AddUser(req)
	if err != nil {
		return model.Response{Code: http.StatusInternalServerError, Msg: "Add user failed", Err: err}
	}
	return model.Response{Code: http.StatusOK, Data: user}
}

func (c *UserController) GetUserById(ctx *gin.Context) model.Response {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return model.Response{Code: http.StatusBadRequest, Msg: "Request param error", Err: err}
	}

	user, err := c.UserService.GetUserById(id)
	if err != nil {
		return model.Response{Code: http.StatusInternalServerError, Msg: "Get user failed", Err: err}
	}
	return model.Response{Code: http.StatusOK, Data: user}
}
