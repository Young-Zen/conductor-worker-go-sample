package server

import (
	"path"
	"worker-sample/server/controller"
)

func (s HttpServer) setUpRouter() {
	healthController := controller.NewHealthController(s.ServiceContext)
	s.Engine.GET(s.joinBasePath("/health"), controller.Wrap(healthController.Status))

	userGroup := s.Engine.Group(s.joinBasePath("/user"))
	{
		userController := controller.NewUserController(s.ServiceContext)
		userGroup.POST("/", controller.Wrap(userController.AddUser))
		userGroup.GET("/:id", controller.Wrap(userController.GetUserById))
	}
}

func (s HttpServer) joinBasePath(relativePath string) string {
	basePath := s.ServiceContext.Config.Server.BasePath
	return path.Join(basePath, relativePath)
}
