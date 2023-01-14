package router

import (
	"github.com/Similadayo/backend/controllers"
	"github.com/Similadayo/backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/", controllers.HomeHandler)
	r.POST("/register", controllers.CreateUser)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)
	r.PUT("/update/:id", middleware.Authorization, controllers.UpdateUser)
	r.GET("/user/:id", middleware.Authorization, controllers.GetUser)
	r.GET("/users", middleware.Authorization, controllers.GetUsers)
	r.DELETE("/delete/:id", middleware.Authorization, controllers.DeleteUser)
	r.PUT("/suspend/:id", middleware.Authorization, controllers.SuspendUser)
}
