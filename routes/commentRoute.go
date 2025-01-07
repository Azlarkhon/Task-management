package routes

import (
	"task-management/controllers"
	"task-management/middleware"

	"github.com/gin-gonic/gin"
)

var comment controllers.CommentInterface = &controllers.Comment{}

func CommentRoutes(r *gin.Engine) {
	protected := r.Group("/users/:user_id/comments")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("", comment.GetAllCommentsOfTask)         // /users/:user_id/comments
		protected.POST("", comment.CreateComment)               // /users/:user_id/comments
		protected.DELETE("/:comment_id", comment.DeleteComment) // /users/:user_id/comments/:comment_id
	}
}
