package routes

import (
	"task-management/controllers"
	"task-management/middleware"

	"github.com/gin-gonic/gin"
)

var task controllers.TaskInterface = &controllers.Task{}

func TaskRoutes(r *gin.Engine) {
	protected := r.Group("/users/:user_id/tasks")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("", task.GetAllTasks)            // /users/:user_id/tasks
		protected.GET("/:task_id", task.GetTaskByID)   // /users/:user_id/tasks/:task_id
		protected.POST("", task.CreateTask)            // /users/:user_id/tasks
		protected.PUT("/:task_id", task.UpdateTask)    // /users/:user_id/tasks/:task_id
		protected.DELETE("/:task_id", task.DeleteTask) // /users/:user_id/tasks/:task_id
	}
}
