package server

import (
	"fmt"
	"worker-sample/config"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type HttpServer struct {
	Engine         *gin.Engine
	ServiceContext *config.ServiceContext
}

func InitHttpServer(ctx *config.ServiceContext) HttpServer {
	server := newHttpServer(ctx)
	server.setUpRouter()
	server.run()
	return server
}

func newHttpServer(ctx *config.ServiceContext) HttpServer {
	engine := gin.New()
	engine.Use(config.GinLogger(), gin.Recovery())
	return HttpServer{
		Engine:         engine,
		ServiceContext: ctx,
	}
}

func (s HttpServer) run() {
	err := s.Engine.Run(fmt.Sprintf(":%d", s.ServiceContext.Config.Server.Port))
	if err != nil {
		log.Fatalf("Init Http server failed: %+v", err)
	}
}
