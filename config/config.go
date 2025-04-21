package config

import (
	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/worker"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	Server Server
	MySQL  MySQL
	Worker Worker
}

type Server struct {
	Port     int
	BasePath string
}

type MySQL struct {
	Host     string
	Port     string
	Schemas  string
	Username string
	Password string
}

type Worker struct {
	BaseUrl         string
	Domain          string
	PollingInterval time.Duration
	BatchSize       int
	Username        string
	Password        string
}

type ServiceContext struct {
	Config *Config
	DB     *gorm.DB
	Worker *ConductorWorker
}

type ConductorWorker struct {
	APIClient      *client.APIClient
	WorkflowClient *client.WorkflowResourceApiService
	TaskRunner     *worker.TaskRunner
}

func NewServiceContext(configPath string, configName string) *ServiceContext {
	c := initConfig(configPath, configName)
	return &ServiceContext{
		Config: &c,
	}
}

func initConfig(configPath string, configName string) Config {
	v := viper.New()
	v.AutomaticEnv()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Read config failed: %+v", err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("Init config failed: %+v", err)
	}
	parseConfig(&c)
	return c
}

func parseConfig(c *Config) {
	c.MySQL.Host = evaluateConfig(c.MySQL.Host)
	c.MySQL.Port = evaluateConfig(c.MySQL.Port)
	c.MySQL.Username = evaluateConfig(c.MySQL.Username)
	c.MySQL.Password = evaluateConfig(c.MySQL.Password)
	c.Worker.BaseUrl = evaluateConfig(c.Worker.BaseUrl)
	c.Worker.Username = evaluateConfig(c.Worker.Username)
	c.Worker.Password = evaluateConfig(c.Worker.Password)
}

var envVarReg = regexp.MustCompile(`\$\{([^}]+)\}`)

func evaluateConfig(val string) string {
	return envVarReg.ReplaceAllStringFunc(val, func(s string) string {
		parts := strings.Split(s[2:len(s)-1], ":")
		envName := parts[0]
		defaultValue := ""
		if len(parts) > 1 {
			defaultValue = parts[1]
		}
		envValue := os.Getenv(envName)
		if envValue == "" {
			return defaultValue
		}
		return envValue
	})
}
