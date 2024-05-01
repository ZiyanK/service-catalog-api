package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ZiyanK/service-catalog-api/app/middleware"
	"github.com/ZiyanK/service-catalog-api/app/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// ServiceInput is a struct used to take the name and description of the service
type ServiceInput struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description" validate:"required,min=20"`
}

// HandlerCreateService creates a new service for the user
func HandlerCreateService(c *gin.Context) {
	userUUID, err := middleware.GetUserUUID(c)
	if err != nil {
		log.Error("Error getting user_uuid", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	var body ServiceInput

	err = c.ShouldBindJSON(&body)
	if err != nil {
		log.Error("Error while reading request body for signup", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid body.",
		})
		return
	}

	err = validator.New().Struct(body)
	if err != nil {
		log.Info("validator error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid body.",
		})
		return
	}

	service := &model.Service{
		Name:        body.Name,
		Description: body.Description,
		UserUUID:    userUUID,
	}

	err = service.CreateService(context.TODO())
	if err != nil {
		if err.Error() == "service exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Service with same name exists.",
			})
			return
		}

		log.Error("Error while creating service", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

// HandlerGetServices fetches all the services created of the user along with filtering, pagination and sorting
func HandlerGetServices(c *gin.Context) {
	userUUID, err := middleware.GetUserUUID(c)
	if err != nil {
		log.Error("Error getting user_uuid", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	// Get the limit query parameter from the URL
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid limit value"})
		return
	}

	// Get the offset query parameter from the URL
	offsetStr := c.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid offset value"})
		return
	}

	// Get the name and orderBy query parameter from the URL
	name := c.Query("name")
	orderBy := c.Query("orderBy")

	services, err := model.GetServices(context.TODO(), userUUID, limit, offset, name, orderBy)
	if err != nil {
		log.Error("Error while fetching services", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(services) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"msg": "No services found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "Services fetched successfully.",
		"data": services,
	})
}

// HandlerGetService fetches as service and all the versions available for the service
func HandlerGetService(c *gin.Context) {
	userUUID, err := middleware.GetUserUUID(c)
	if err != nil {
		log.Error("Error getting user_uuid", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	serviceIDstring := c.Param("id")
	serviceID, err := strconv.Atoi(serviceIDstring)
	if err != nil {
		log.Info("invalid service id")
		c.Status(http.StatusNotFound)
		return
	}

	service, err := model.GetService(context.TODO(), serviceID, userUUID)
	if err != nil {
		log.Error("Error fetching service", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(service) == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": service,
		"msg":  "Service fetched successfully.",
	})
}

// HandlerUpdateService updates the name and the description of the given service
func HandlerUpdateService(c *gin.Context) {
	userUUID, err := middleware.GetUserUUID(c)
	if err != nil {
		log.Error("Error getting user_uuid", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	serviceIDstring := c.Param("id")
	serviceID, err := strconv.Atoi(serviceIDstring)
	if err != nil {
		log.Info("invalid service id")
		c.Status(http.StatusNotFound)
		return
	}

	var body ServiceInput

	err = c.ShouldBindJSON(&body)
	if err != nil {
		log.Error("Error while reading request body for signup", zap.Error(err))
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	err = validator.New().Struct(body)
	if err != nil {
		log.Info("validator error", zap.Error(err))
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	service := model.Service{
		ServiceID:   serviceID,
		UserUUID:    userUUID,
		Name:        body.Name,
		Description: body.Description,
	}

	err = service.UpdateService(context.TODO())
	if err != nil {
		if err.Error() == "service does not exist" {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// HandlerDeleteService deletes a service and all the versions associated to the service
func HandlerDeleteService(c *gin.Context) {
	userUUID, err := middleware.GetUserUUID(c)
	if err != nil {
		log.Error("Error getting user_uuid", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	serviceIDstring := c.Param("id")
	serviceID, err := strconv.Atoi(serviceIDstring)
	if err != nil {
		log.Info("invalid service id")
		c.Status(http.StatusNotFound)
		return
	}

	service := model.Service{
		ServiceID: serviceID,
		UserUUID:  userUUID,
	}

	err = service.DeleteService(context.TODO())
	if err != nil {
		if err.Error() == "service does not exist" {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
