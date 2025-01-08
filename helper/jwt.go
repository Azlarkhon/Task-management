package helper

import (
	"fmt"
	"net/http"
	"task-management/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,                                     // Custom claim: User ID
		"exp":     time.Now().Add(24 * 30 * time.Hour).Unix(), // Expiration: 1 month
		"iat":     time.Now().Unix(),                          // Issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(config.Init.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyJWT(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.Init.JWTSecret), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if float64(time.Now().Unix()) > exp {
				return nil, nil, fmt.Errorf("token is expired")
			}
		}
		return token, claims, nil
	}

	return nil, nil, jwt.ErrSignatureInvalid
}

func CheckAuthenticationAndAuthorization(c *gin.Context, userID string) (string, bool) {
	loggedInUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("User is not authenticated"))
		return "", false
	}

	loggedInUserIDStr, ok := loggedInUserID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Invalid user ID type"))
		return "", false
	}

	if loggedInUserIDStr != userID {
		c.JSON(http.StatusForbidden, NewErrorResponse("You are not authorized to view this user"))
		return "", false
	}

	return loggedInUserIDStr, true
}
