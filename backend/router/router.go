package router

import (
	"github.com/Similadayo/backend/controllers"
	"github.com/Similadayo/backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	//User and Admin functions
	r.GET("/", controllers.HomeHandler)
	r.POST("/register", controllers.CreateUser)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)
	r.PUT("/update/:id", middleware.Authorization, controllers.UpdateUser)
	r.GET("/user/:id", middleware.Authorization, controllers.GetUser)
	r.GET("/users", middleware.Authorization, controllers.GetUsers)
	r.DELETE("/delete/:id", middleware.Authorization, controllers.DeleteUser)
	r.PUT("/suspend/:id", middleware.Authorization, controllers.SuspendUser)
	r.POST("/forgotpassword", controllers.ForgotPassword)
	r.POST("/resetpassword", controllers.ResetPassword)

	// Category Routes
	r.POST("/createcategory", controllers.CreateCategory)
	r.GET("/categories", controllers.GetCategories)
	r.GET("/category/:id", controllers.GetCategory)
	r.PUT("/updatecategory/:id", controllers.UpdateCategory)
	r.DELETE("/deletecategory/:id", controllers.DeleteCategory)

	// product function
	r.POST("/createproduct", controllers.CreateProduct)

}
