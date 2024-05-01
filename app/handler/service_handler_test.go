package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZiyanK/service-catalog-api/app/middleware"
	"github.com/ZiyanK/service-catalog-api/app/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func AddAuthorizationHeader(req *http.Request) error {
	jwtSecret := "secret"
	userUUID, err := uuid.Parse("d90f9b49-dcd9-4feb-8250-d013098e45ee")
	if err != nil {
		log.Error("error parsing uuid", zap.Error(err))
	}

	claims := AuthTokenClaims{
		UserUUID: userUUID,
	}
	token, err := generateAuthToken(jwtSecret, claims)
	if err != nil {
		log.Error("error generating auth token", zap.Error(err))
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}

func TestHandlerCreateService(t *testing.T) {
	router := SetupTest()

	route := "/service"
	router.Use(middleware.VerifyAuthToken)
	router.POST(route, HandlerCreateService)
	body := ServiceInput{
		Name:        "backend",
		Description: "this service has the backend",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/service", bytes.NewBuffer(jsonValue))
	AddAuthorizationHeader(req)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Case fail: User already present
	req, _ = http.NewRequest(http.MethodPost, route, bytes.NewBuffer(jsonValue))
	AddAuthorizationHeader(req)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	body.Name = "frontend"
	body.Description = "this service is the frontend"

	jsonValue, _ = json.Marshal(body)

	req, _ = http.NewRequest(http.MethodPost, "/service", bytes.NewBuffer(jsonValue))
	AddAuthorizationHeader(req)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestHandlerGetServices(t *testing.T) {
	router := SetupTest()

	route := "/services"
	router.Use(middleware.VerifyAuthToken)
	router.GET(route, HandlerGetServices)
	req, _ := http.NewRequest(http.MethodGet, "/services?limit=10&offset=0", nil)
	AddAuthorizationHeader(req)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	type Response struct {
		Data []model.Service `json:"data"`
		Msg  string          `json:"msg"`
	}

	var responseBody Response
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		log.Info("failed to unmarshal body")
		t.Fail()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, responseBody.Data[0].Name, "backend")
	assert.Equal(t, responseBody.Data[0].Description, "this service has the backend")
	assert.Equal(t, responseBody.Data[1].Name, "frontend")
	assert.Equal(t, responseBody.Data[1].Description, "this service is the frontend")
}

func TestHandlerGetServicesOrderByDesc(t *testing.T) {
	router := SetupTest()

	route := "/services"
	router.Use(middleware.VerifyAuthToken)
	router.GET(route, HandlerGetServices)
	req, _ := http.NewRequest(http.MethodGet, "/services?limit=10&offset=0&orderBy=DESC", nil)
	AddAuthorizationHeader(req)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	type Response struct {
		Data []model.Service `json:"data"`
		Msg  string          `json:"msg"`
	}

	var responseBody Response
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		log.Info("failed to unmarshal body")
		t.Fail()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, responseBody.Data[0].Name, "frontend")
	assert.Equal(t, responseBody.Data[0].Description, "this service is the frontend")
	assert.Equal(t, responseBody.Data[1].Name, "backend")
	assert.Equal(t, responseBody.Data[1].Description, "this service has the backend")
}

func TestHandlerGetService(t *testing.T) {
	router := SetupTest()

	route := "/service/:id"
	router.Use(middleware.VerifyAuthToken)
	router.GET(route, HandlerGetService)
	req, _ := http.NewRequest(http.MethodGet, "/service/1", nil)
	AddAuthorizationHeader(req)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	type Response struct {
		Data []model.ServiceWithVersions `json:"data"`
		Msg  string                      `json:"msg"`
	}

	var responseBody Response
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		log.Info("failed to unmarshal body")
		t.Fail()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, responseBody.Data[0].Name, "backend")
	assert.Equal(t, responseBody.Data[0].Description, "this service has the backend")
}

func TestHandlerUpdateService(t *testing.T) {
	router := SetupTest()

	route := "/service/:id"
	router.Use(middleware.VerifyAuthToken)
	router.PUT(route, HandlerUpdateService)
	body := ServiceInput{
		Name:        "frontend",
		Description: "this service has the frontend",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPut, "/service/1", bytes.NewBuffer(jsonValue))
	AddAuthorizationHeader(req)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandlerDeleteService(t *testing.T) {
	router := SetupTest()

	route := "/service/:id"
	router.Use(middleware.VerifyAuthToken)
	router.DELETE(route, HandlerDeleteService)
	req, _ := http.NewRequest(http.MethodDelete, "/service/1", nil)
	AddAuthorizationHeader(req)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
