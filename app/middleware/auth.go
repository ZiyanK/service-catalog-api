package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/ZiyanK/service-catalog-api/app/logger"
	model "github.com/ZiyanK/service-catalog-api/app/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	log = logger.CreateLogger()
)

// VerifyAuthToken is used to verify any given auth token with the secret provided in the env
func VerifyAuthToken(c *gin.Context) {
	authorizationToken := c.Request.Header.Get("Authorization")
	if authorizationToken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokensArray := strings.Split(authorizationToken, " ")
	if len(tokensArray) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := tokensArray[1]

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		jwtSecret := viper.Get("jwt_secret").(string)
		var secretKeyInBytes = []byte(jwtSecret)
		return secretKeyInBytes, nil

	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !parsedToken.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userUUIDstring, ok := claims["user_id"].(string)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userUUID, err := uuid.Parse(userUUIDstring)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err := model.GetUserByID(c, userUUID)
	if err == sql.ErrNoRows {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Set("user_uuid", user.UserUUID)
	c.Next()
}

// HashValue returns the hashed string of the value given
func HashValue(value string) (string, error) {
	byteValue := []byte(value)

	// generate hash of the from byte slice
	hash, err := bcrypt.GenerateFromPassword(byteValue, bcrypt.MinCost)
	if err != nil {
		log.Error("Error while generating hashed value: ", zap.Error(err))
		return "", err
	}

	// convert the bytes to a string
	return string(hash), nil
}

// GetUserUUID returns the uuid of the user making the request
func GetUserUUID(c *gin.Context) (uuid.UUID, error) {
	var userUUID uuid.UUID
	userID, exists := c.Get("user_uuid")
	if !exists {
		err := errors.New("user_uuid does not exists in gin.Context")
		return uuid.Nil, err
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		err := errors.New("error while changing type of userID from string to userUUID")
		return uuid.Nil, err
	}

	return userUUID, nil

}
