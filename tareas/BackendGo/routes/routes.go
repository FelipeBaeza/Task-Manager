package routes

import (
	"go-template/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	// routes for tasks
	router.POST("/tasks", controllers.CreateTask)
	router.GET("/tasks", controllers.GetTasks)
	router.PUT("/tasks/:id/complete", controllers.CompletedTask)

}
