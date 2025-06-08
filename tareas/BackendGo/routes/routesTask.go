package routes

import (
	"go-template/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	// routes for tasks
	router.POST("/tasks", controllers.CreateTask)
	router.GET("/tasks", controllers.GetTasks)
	router.PUT("/tasks/:id", controllers.UpdateTaskStatus)
	router.GET("/tasks/:id", controllers.GetTaskByID)
	router.DELETE("/tasks/:id", controllers.DeleteTask)
}
