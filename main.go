package main

import (
	"task-management/config"
	"task-management/database"
	"task-management/middleware"
	"task-management/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDatabase()

	router := gin.Default()

	routes.UserRoutes(router)
	routes.TagRoutes(router)
	routes.TaskRoutes(router)
	routes.CommentRoutes(router)

	router.Use(middleware.CorsMiddleware)

	router.Run(":" + config.Init.Port)
}
