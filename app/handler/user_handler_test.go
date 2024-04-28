package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZiyanK/service-catalog-api/app/middleware"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetUser(t *testing.T) {
	router := SetupTest()

	route := "/user"
	router.Use(middleware.VerifyAuthToken)
	router.GET(route, HandlerGetUser)
	req, _ := http.NewRequest(http.MethodGet, route, nil)
	AddAuthorizationHeader(req)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
