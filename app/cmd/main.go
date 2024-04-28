package main

import (
	"fmt"

	"github.com/ZiyanK/service-catalog-api/app/db"
	"github.com/ZiyanK/service-catalog-api/app/logger"
	"github.com/ZiyanK/service-catalog-api/app/route"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	log = logger.CreateLogger()
)

func main() {
	// load config file
	if err := LoadConfig(); err != nil {
		panic(err)
	}

	// Databases Init
	if err := db.InitConn(config.DSN); err != nil {
		log.Fatal("Failed to conenct to the database", zap.Error(err))
	}

	// HTTP API
	router := route.AddRouter()

	log.Info("Server up and running")
	err := router.Run(fmt.Sprintf(":%v", config.Port))
	if err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}

// configration is a struct used to get the environment variable from the config.yaml file
type configuration struct {
	DSN       string `mapstructure:"dsn"`
	Port      string `mapstructure:"port"`
	JWTSecret string `mapstructure:"jwt_secret"`
	Mode      string `mapstructure:"mode"`
}

var (
	config configuration
)

// LoadConfig is used to load the configuration
func LoadConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file", zap.String("err", err.Error()))
		return err
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("unable to decode into struct", zap.String("err", err.Error()))
		return err
	}

	return nil
}
