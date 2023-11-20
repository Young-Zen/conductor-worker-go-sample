package config

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		log.Infof("| %3d | %13v | %15s | %s | %s |",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.ClientIP(),
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}
