package controllers

import (
	"net/http"

	"github.com/Similadayo/backend/db"
	"github.com/Similadayo/backend/models"
	"github.com/Similadayo/backend/utils"
	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	// Get the token from the Authorization header
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found in request headers"})
		return
	}

	// Verify the token
	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Check if the user is an admin
	if claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can access this resource"})
		return
	}

	// Bind the request body to a Product struct
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get a DB connection
	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error connecting to database"))
		return
	}
	defer db.Close()

	// Save the product to the database
	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error creating product"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Product created successfully",
		"data":    product,
	})
}
