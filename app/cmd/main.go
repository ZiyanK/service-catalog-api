package main

import (
	"fmt"
	"net/http"

	"github.com/ZiyanK/service-catalog-api/app/db"
	"github.com/ZiyanK/service-catalog-api/app/logger"
	"github.com/ZiyanK/service-catalog-api/app/route"
	"github.com/gin-contrib/cors"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "*"}, // TODO: restrict origins
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.RemoveExtraSlash = true

	log.Info("Server up and running")
	err := router.Run(fmt.Sprintf(":%v", config.Port))
	if err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}

// configration is a struct used to get the environment variable from the config.yaml file
type configuration struct {
	DSN       string `mapstructure:"DSN"`
	Port      string `mapstructure:"PORT"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

var (
	config configuration
)

// LoadConfig is used to load the configuration
func LoadConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		viper.AutomaticEnv()

		dsn := viper.GetString("DSN")
		jwtSecret := viper.GetString("JWT_SECRET")
		port := viper.GetString("PORT")

		log.Info("config", zap.Any("DSN", dsn))

		config.DSN = dsn
		config.JWTSecret = jwtSecret
		config.Port = port
	} else {
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatal("unable to decode into struct", zap.String("err", err.Error()))
			return err
		}
	}

	return nil
}
