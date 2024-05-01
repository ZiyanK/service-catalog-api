package model

import (
	"context"
	"strings"
	"testing"

	database "github.com/ZiyanK/service-catalog-api/app/db"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	userUUID string
)

type configuration struct {
	DSN       string `mapstructure:"dsn"`
	Port      string `mapstructure:"port"`
	JWTSecret string `mapstructure:"jwt_secret"`
	Mode      string `mapstructure:"mode"`
}

var (
	config   configuration
	UserUUID uuid.UUID
)

func setupTest() {
	viper.AddConfigPath("../..")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file", zap.String("err", err.Error()))
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("unable to decode into struct", zap.String("err", err.Error()))
	}

	if err := database.InitConn(config.DSN); err != nil {
		log.Fatal("Error connecting to db: ", zap.Error(err))
	}
	userUUID = "d90f9b49-dcd9-4feb-8250-d013098e45ee"
}

func appendExplain(query string) string {
	return "EXPLAIN " + query
}

func TestGetUserByEmail(t *testing.T) {
	setupTest()

	rows, err := db.NamedExplainQuery(context.Background(), appendExplain(queryCheckUserExist), map[string]interface{}{
		"email": "jd@gmail.com",
	})
	if err != nil {
		t.Fatal("Failed to execute query:", err)
	}
	defer rows.Close()

	// Analyze the query execution plan
	var plan string
	var indexUsed bool
	for rows.Next() {
		if err := rows.Scan(&plan); err != nil {
			t.Fatal("Failed to scan row:", err)
		}

		// Check if the index is being used
		if strings.Contains(plan, "Index Only") {
			log.Info("Index scan being used")
			indexUsed = true
			break
		}
	}

	if !indexUsed {
		t.Error("Expected index scan but index is not being used")
	}

	if err := rows.Err(); err != nil {
		t.Fatal("Error iterating over rows:", err)
	}
}
