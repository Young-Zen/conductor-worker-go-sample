package server

import (
	"net/http"
	"path"

	"worker-sample/server/controller"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s HttpServer) setUpRouter() {
	healthController := controller.NewHealthController(s.ServiceContext)
	s.Engine.GET(s.joinBasePath("/health"), controller.Wrap(healthController.Status))

	s.Engine.GET(s.joinBasePath("metrics"), promHandler(promhttp.Handler()))

	userGroup := s.Engine.Group(s.joinBasePath("/user"))
	{
		userController := controller.NewUserController(s.ServiceContext)
		userGroup.POST("/", controller.Wrap(userController.AddUser))
		userGroup.GET("/:id", controller.Wrap(userController.GetUserById))
	}
}

func promHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func (s HttpServer) joinBasePath(relativePath string) string {
	basePath := s.ServiceContext.Config.Server.BasePath
	return path.Join(basePath, relativePath)
}
