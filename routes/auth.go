package routes

import (
	"github.com/gin-gonic/gin"
	"go-test/controllers"
	"go-test/enums"
	"go-test/middlewares"
)

func RegisterAuthRoutes(router *gin.Engine) {
	router.POST("/auth/login", controllers.Login)
	router.POST("/auth/logout", middlewares.AuthMiddleware(), controllers.Logout)

	router.POST("/register", controllers.Register)

	router.GET("/profile", middlewares.AuthMiddleware(), controllers.Profile)
	router.GET("/users", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.Admin()), controllers.UsersList)
	router.POST("/admin", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.SuperAdmin()), controllers.CreateAdmin)
	router.POST("/delete/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Roles{}.Admin()), controllers.DeleteUser)
}
