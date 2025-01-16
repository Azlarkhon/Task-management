package app

import (
	"task-management/config"
	"task-management/database"
	"task-management/middleware"
	"task-management/routes"

	"github.com/gin-gonic/gin"
)

func Run() error {
	database.ConnectDatabase()

	router := gin.Default()

	router.Use(middleware.CorsMiddleware)

	routes.UserRoutes(router)
	routes.TagRoutes(router)
	routes.TaskRoutes(router)
	routes.CommentRoutes(router)
	routes.GoogleLoginRoutes(router)

	return router.Run(":" + config.Init.Port)
}
