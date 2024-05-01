package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZiyanK/service-catalog-api/app/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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

func SetupTest() *gin.Engine {
	viper.AddConfigPath("../..")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file", zap.String("err", err.Error()))
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("unable to decode into struct", zap.String("err", err.Error()))
	}

	r := gin.New()
	if err := db.InitConn(config.DSN); err != nil {
		log.Fatal("Error connecting to db: ", zap.Error(err))
	}
	gin.SetMode(gin.ReleaseMode)
	return r
}

func TestSignUp(t *testing.T) {
	router := SetupTest()

	email := "jd@gmail.com"
	route := "/signup"
	router.POST(route, HandlerSignUp)
	body := AuthInput{
		Email:    email,
		Password: "johndoe123",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	type UserData struct {
		Email       string `json:"email"`
		AccessToken string `json:"access_token"`
	}

	type Response struct {
		Data UserData `json:"data"`
		Msg  string   `json:"msg"`
	}

	var responseBody Response

	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		log.Info("failed to unmarshal body")
		t.Fail()
	}

	assert.Equal(t, email, responseBody.Data.Email)

	// Case fail: User already present
	req, _ = http.NewRequest(http.MethodPost, route, bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin(t *testing.T) {
	router := SetupTest()

	route := "/login"
	router.POST(route, HandlerLogin)
	body := AuthInput{
		Email:    "jd@gmail.com",
		Password: "johndoe123",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
