package routes

import (
	"go-template/controllers"
	"github.com/gin-gonic/gin"
	"go-template/middleware"
)

func RoutesAuth(router *gin.Engine) {
	// Rutas públicas (sin middleware)
	router.POST("/auth/register", controllers.CreateUser)
	router.POST("/auth/login", controllers.LoginUser)

	// Rutas protegidas (con middleware de autenticación)
	router.GET("/user/me", middleware.AuthMiddleware(), controllers.UserMe)


}
