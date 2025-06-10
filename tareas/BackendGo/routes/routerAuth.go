package routes

import (
	"go-template/controllers"
	"go-template/middleware"

	"github.com/gin-gonic/gin"
)

func RoutesAuth(router *gin.Engine) {

	router.POST("/auth/register", controllers.CreateUser)
	router.POST("/auth/login", controllers.LoginUser)

	// Protected routes
	router.GET("/user/me", middleware.AuthMiddleware(), controllers.UserMe)
	router.GET("/users", middleware.AuthMiddleware(), controllers.GetAllUsers)
}
