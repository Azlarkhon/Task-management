package routes

import (
	"task-management/controllers"
	"task-management/middleware"

	"github.com/gin-gonic/gin"
)

var tag controllers.TagInterface = &controllers.Tag{}

func TagRoutes(r *gin.Engine) {
	protected := r.Group("/users/:user_id/tags")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("", tag.GetAllTags)                      // /users/:user_id/tags
		protected.GET("/:tag_id", tag.GetTagById)              // /users/:user_id/tags/:tag_id
		protected.POST("", tag.CreateTag)                      // /users/:user_id/tags
		protected.PUT("/:tag_id", tag.UpdateTagById)           // /users/:user_id/tags/:tag_id
		protected.DELETE("/:tag_id", tag.DeleteTagById)        // /users/:user_id/tags/:tag_id
		protected.POST("/add_tag", tag.AddTagToTask)           // /users/:user_id/tags/add_tag
		protected.DELETE("/delete_tag", tag.DeleteTagFromTask) // /users/:user_id/tags/delete_tag
		protected.GET("/task_tags", tag.GetAllTagsOfTask)      // /users/:user_id/tags/task_tags
	}
}
