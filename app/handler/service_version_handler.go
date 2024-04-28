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

// ServiceVersionInput is a struct used to take the version and changelog of a new version for a given service
type ServiceVersionInput struct {
	Version   string `json:"version" validate:"required,min=2"`
	Changelog string `json:"changelog"`
}

// HandlerCreateServiceVersion created a new version for a given service
func HandlerCreateServiceVersion(c *gin.Context) {
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

	var body ServiceVersionInput

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

	serviceVersion := model.ServiceVersion{
		Version:   body.Version,
		Changelog: body.Changelog,
		ServiceID: serviceID,
	}

	err = serviceVersion.CreateServiceVersion(context.TODO(), userUUID)
	if err != nil {
		switch err.Error() {
		case "service does not exist":
			c.Status(http.StatusNotFound)
			return
		case "service version exists":
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Service with same version exists.",
			})
			return
		}

		log.Error("Error while creating service version", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

// HandlerDeleteServiceVersion deletes an existing version for a given service
func HandlerDeleteServiceVersion(c *gin.Context) {
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

	svIDstring := c.Param("vid")
	svID, err := strconv.Atoi(svIDstring)
	if err != nil {
		log.Info("invalid service id")
		c.Status(http.StatusNotFound)
		return
	}

	err = model.DeleteServiceVersion(context.TODO(), userUUID, serviceID, svID)
	if err != nil {
		if err.Error() == "service version does not exist" {
			c.Status(http.StatusNotFound)
			return
		}
		log.Error("Error while deleting service version", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
