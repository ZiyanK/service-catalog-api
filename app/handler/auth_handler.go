package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/ZiyanK/service-catalog-api/app/logger"
	"github.com/ZiyanK/service-catalog-api/app/middleware"
	"github.com/ZiyanK/service-catalog-api/app/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	log = logger.CreateLogger()
)

// AuthInput is a struct used to get the email and the password from the user
type AuthInput struct {
	Email    string `json:"email" validate:"required,email,max=50"`
	Password string `json:"password" validate:"required,min=7,max=15"`
}

// HandlerSignUp is a handler that signs in the user
func HandlerSignUp(c *gin.Context) {
	var body AuthInput

	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Error("Error while reading request body for signup", zap.Error(err))
		c.Status(http.StatusUnprocessableEntity)
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

	// hash password
	password, err := middleware.HashValue(body.Password)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	createUserObj := &model.User{
		UserUUID: uuid.New(),
		Email:    body.Email,
		Password: password,
	}

	err = createUserObj.CreateUser(context.TODO())
	if err != nil {
		if err.Error() == "mail exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Try using a different email.",
			})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	claims := AuthTokenClaims{
		UserUUID: createUserObj.UserUUID,
	}

	jwtSecret := viper.Get("jwt_secret").(string)
	token, err := generateAuthToken(jwtSecret, claims)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	respData := struct {
		Email       string    `json:"email"`
		CreatedAt   time.Time `json:"created_at"`
		AccessToken string    `json:"access_token"`
	}{
		Email:       createUserObj.Email,
		CreatedAt:   createUserObj.CreatedAt,
		AccessToken: token,
	}

	c.JSON(http.StatusCreated, gin.H{
		"msg":  "User created successfully",
		"data": respData,
	})
}

// HandlerLogin is a handler that logs in the user
func HandlerLogin(c *gin.Context) {
	var body AuthInput

	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Info("Error while reading request body for login", zap.Error(err))
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	err = validator.New().Struct(body)
	if err != nil {
		log.Info("validator error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid body",
		})
		return
	}

	user, err := model.GetUserByEmail(context.TODO(), body.Email)
	if err != nil {
		log.Error("Error when fetching user by email", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}
	// check if user found with email
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Invalid email or password. Please try again.",
		})
		return

	}

	// check if password matches
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		log.Info("Password does not match")
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Invalid email or password. Please try again.",
		})
		return
	}

	claims := AuthTokenClaims{
		UserUUID: user.UserUUID,
	}

	jwtSecret := viper.Get("jwt_secret").(string)
	token, err := generateAuthToken(jwtSecret, claims)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "User logged in successfully",
		"data": gin.H{
			"access_token": token,
		},
	})
}

type AuthTokenClaims struct {
	UserUUID uuid.UUID `json:"user_id"`
}

// GenerateAuthToken generates auth token for a given secret
func generateAuthToken(secret string, claimsToAdd AuthTokenClaims) (string, error) {
	var secretKeyInBytes = []byte(secret)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = claimsToAdd.UserUUID

	tokenString, err := token.SignedString(secretKeyInBytes)
	if err != nil {
		log.Error("Error while signing string for access token", zap.Error(err))
		return "", err
	}

	return tokenString, nil
}
