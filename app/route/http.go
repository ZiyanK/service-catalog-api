package route

import (
	"github.com/ZiyanK/service-catalog-api/app/handler"
	"github.com/ZiyanK/service-catalog-api/app/middleware"
	"github.com/gin-gonic/gin"
)

const (
	pathPing = "/ping"

	pathSignup = "/signup"
	pathLogin  = "/login"

	pathUser      = "/user"
	pathServices  = "/services"
	pathService   = "/service"
	pathServiceID = "/service/:id"

	pathServiceIDVersion   = "/service/:id/version"
	pathServiceIDVersionID = "/service/:id/version/:vid"
)

func AddRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET(pathPing, func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	// Auth routes
	router.POST(pathSignup, handler.HandlerSignUp)
	router.POST(pathLogin, handler.HandlerLogin)

	// Protected routes
	router.Use(middleware.VerifyAuthToken)

	// User routes
	router.GET(pathUser, handler.HandlerGetUser)
	router.PUT(pathUser, handler.HandlerUpdateUser)

	// Service routes
	router.GET(pathServices, handler.HandlerGetServices)
	router.POST(pathService, handler.HandlerCreateService)
	router.GET(pathServiceID, handler.HandlerGetService)
	router.PUT(pathServiceID, handler.HandlerUpdateService)
	router.DELETE(pathServiceID, handler.HandlerDeleteService)

	// Service version routes
	router.POST(pathServiceIDVersion, handler.HandlerCreateServiceVersion)
	router.DELETE(pathServiceIDVersionID, handler.HandlerDeleteServiceVersion)

	return router
}
