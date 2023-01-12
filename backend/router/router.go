package router

import (
	"github.com/Similadayo/backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/", controllers.HomeHandler)
	r.POST("/register", controllers.CreateUser)
	r.POST("/login", controllers.Login)
}
