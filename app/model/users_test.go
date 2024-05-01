package model

import (
	"context"
	"strings"
	"testing"

	database "github.com/ZiyanK/service-catalog-api/app/db"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type configuration struct {
	DSN       string `mapstructure:"dsn"`
	Port      string `mapstructure:"port"`
	JWTSecret string `mapstructure:"jwt_secret"`
	Mode      string `mapstructure:"mode"`
}

var (
	config      configuration
	userUUID    uuid.UUID
	userUUIDStr string
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
	userUUIDStr = "d90f9b49-dcd9-4feb-8250-d013098e45ee"
	userUUID, _ = uuid.Parse(userUUIDStr)
}

func appendExplain(query string) string {
	return "EXPLAIN " + query
}

func TestCreateUser(t *testing.T) {
	setupTest()

	password := "12345678"
	byteValue := []byte(password)

	// generate hash of the from byte slice
	hash, err := bcrypt.GenerateFromPassword(byteValue, bcrypt.MinCost)
	if err != nil {
		log.Error("Error while generating hashed value: ", zap.Error(err))
		t.Error("Error while generating hashed value")
	}

	password = string(hash)

	user := User{
		UserUUID: uuid.New(),
		Email:    "jd1@gmail.com",
		Password: password,
	}

	err = user.CreateUser(context.Background())
	if err != nil {
		t.Error("Error while creating user")
	}
}

func TestGetUserByEmail(t *testing.T) {
	setupTest()

	user, err := GetUserByEmail(context.Background(), "test@gmail.com")
	if err != nil {
		t.Error("Error while fetching user by email")
	}

	assert.Equal(t, "test@gmail.com", user.Email)
	assert.Equal(t, userUUID, user.UserUUID)
}

func TestGetUserByID(t *testing.T) {
	setupTest()

	user, err := GetUserByID(context.Background(), userUUID)
	if err != nil {
		t.Error("Error while fetching user by email")
	}

	assert.Equal(t, "test@gmail.com", user.Email)
	assert.Equal(t, userUUID, user.UserUUID)
}

func TestUpdateUser(t *testing.T) {
	setupTest()

	err := UpdateUser(context.Background(), "test1@gmail.com", userUUID)
	if err != nil {
		t.Error("Error while creating user")
	}
}

// Test query indexes

// TestQueryGetUserByEmail is used to test whether the index is used to fetch user by email
func TestQueryGetUserByEmail(t *testing.T) {
	setupTest()

	rows, err := db.NamedExplainQuery(context.Background(), appendExplain(queryCheckUserExist), map[string]interface{}{
		"email": "test@gmail.com",
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
		t.Log("error in index used")
		t.Error("Expected index scan but index is not being used")
	}

	if err := rows.Err(); err != nil {
		t.Log("error in scanning rows")
		t.Fatal("Error iterating over rows:", err)
	}
}
