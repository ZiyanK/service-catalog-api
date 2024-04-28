package handler

import (
	"context"
	"net/http"

	"github.com/ZiyanK/service-catalog-api/app/middleware"
	"github.com/ZiyanK/service-catalog-api/app/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HandlerGetUser fetches the user details (email)
func HandlerGetUser(c *gin.Context) {
	userUUID, err := middleware.GetUserUUID(c)
	if err != nil {
		log.Error("Error getting user_uuid", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	user, err := model.GetUserByID(context.TODO(), userUUID)
	if err != nil {
		log.Error("Error fetching user info", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Error fetching user info.",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "Fetched user info successfully",
		"data": user,
	})
}

// HandlerUpdateUser updates the user email
func HandlerUpdateUser(c *gin.Context) {
	userUUID, err := middleware.GetUserUUID(c)
	if err != nil {
		log.Error("Error getting user_uuid", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	type updateUserReq struct {
		Email string `json:"email" validate:"required,email,max=50"`
	}

	var body updateUserReq
	err = c.ShouldBindJSON(&body)
	if err != nil {
		log.Error("Error while reading request body for signup", zap.Error(err))
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	// check for existing email and update if not present
	err = model.UpdateUser(context.TODO(), body.Email, userUUID)
	if err != nil {
		if err.Error() == "mail exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "User with this mail already exists.",
			})
			return
		}

		log.Error("Error while updating user email", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
