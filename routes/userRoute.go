package routes

import (
	"task-management/controllers"
	"task-management/middleware"

	"github.com/gin-gonic/gin"
)

var user controllers.UserInterface = &controllers.User{}

func UserRoutes(r *gin.Engine) {
	r.POST("/users/signup", user.SignUp) // /users/signup -> Sign up
	r.POST("/users/login", user.Login)   // /users/login -> Login

	protected := r.Group("/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/:user_id", user.GetUserByID)   // /users/:user_id -> Get user by ID
		protected.PUT("/:user_id", user.UpdateUser)    // /users/:user_id -> Update user by ID
		protected.DELETE("/:user_id", user.DeleteUser) // /users/:user_id -> Delete user by ID
	}
}
