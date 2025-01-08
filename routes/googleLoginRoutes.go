package routes

import (
	"task-management/controllers"

	"github.com/gin-gonic/gin"
)

var google controllers.GoogleInterface = &controllers.Google{}

func GoogleLoginRoutes(r *gin.Engine) {
	r.GET("/auth/google/login", google.GoogleLogin)       // /auth/google/login
	r.GET("/auth/google/callback", google.GoogleCallback) // /auth/google/callback
}
