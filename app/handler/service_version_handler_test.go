package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZiyanK/service-catalog-api/app/middleware"
	"github.com/stretchr/testify/assert"
)

func TestHandlerCreateServiceVersion(t *testing.T) {
	router := SetupTest()

	route := "/service/:id/version"
	router.Use(middleware.VerifyAuthToken)
	router.POST(route, HandlerCreateServiceVersion)
	svBody := ServiceVersionInput{
		Version:   "v1.0.0",
		Changelog: "this is the changelog",
	}
	jsonValue, _ := json.Marshal(svBody)

	req, _ := http.NewRequest(http.MethodPost, "/service/2/version", bytes.NewBuffer(jsonValue))
	AddAuthorizationHeader(req)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Case fail: User already present
	req, _ = http.NewRequest(http.MethodPost, "/service/2/version", bytes.NewBuffer(jsonValue))
	AddAuthorizationHeader(req)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerDeleteServiceVersion(t *testing.T) {
	router := SetupTest()

	route := "/service/:id/version/:vid"
	router.Use(middleware.VerifyAuthToken)
	router.DELETE(route, HandlerDeleteService)
	req, _ := http.NewRequest(http.MethodDelete, "/service/2/version/1", nil)
	AddAuthorizationHeader(req)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
