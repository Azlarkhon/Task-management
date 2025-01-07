package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"task-management/helper"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, helper.NewErrorResponse("No token provided"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, helper.NewErrorResponse("Invalid token"))
			c.Abort()
			return
		}
		tokenString := parts[1]

		_, claims, err := helper.VerifyJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, helper.NewErrorResponse("Invalid token"))
			c.Abort()
			return
		}

		c.Set("user_id", fmt.Sprintf("%v", claims["user_id"]))
		c.Next()
	}
}
