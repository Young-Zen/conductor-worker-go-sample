package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	Server Server
	MySQL  MySQL
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

type ServiceContext struct {
	Config *Config
	DB     *gorm.DB
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
}

var envVarReg = regexp.MustCompile("\\$\\{([^}]+)\\}")

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
