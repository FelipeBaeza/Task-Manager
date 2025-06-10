package routes

import (
	"go-template/controllers"
	"go-template/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// routes for tasks
	router.POST("/tasks", middleware.AuthMiddleware(), controllers.CreateTask)
	router.GET("/tasks", middleware.AuthMiddleware(), controllers.GetTasks)
	router.GET("/tasks/:id", middleware.AuthMiddleware(), controllers.GetTaskByID)
	router.PUT("/tasks/:id", middleware.AuthMiddleware(), controllers.UpdateTaskStatus)
	router.DELETE("/tasks/:id", middleware.AuthMiddleware(), controllers.DeleteTask)

	// routes for comments 
	router.POST("/tasks/:id/comments", middleware.AuthMiddleware(), controllers.CreateComment)
	router.GET("/tasks/:id/comments", middleware.AuthMiddleware(), controllers.GetCommentsByTask)
	router.DELETE("/comments/:commentId", middleware.AuthMiddleware(), controllers.DeleteComment)
}
