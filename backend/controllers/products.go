package controllers

import (
	"net/http"

	"github.com/Similadayo/backend/db"
	"github.com/Similadayo/backend/models"
	"github.com/Similadayo/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can create products"})
		return
	}
	// Bind the request body to the product struct
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		validateError := err.(validator.ValidationErrors)
		message := validateError[0].Field() + " is " + validateError[0].Tag()
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	// Get the category from the request body
	var category models.Category
	if err := db.DB.Where("id = ?", product.CategoryID).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category id"})
		return
	}

	// Append  sizes and colors to the product
	if len(product.Sizes) > 0 {
		for _, size := range product.Sizes {
			var dbSize models.Size
			if err := db.DB.Where("name = ?", size.Name).First(&dbSize).Error; err != nil {
				if err := db.DB.Create(&size).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating new color"})
					return
				}
			}
			db.DB.Model(&product).Association("Sizes").Append(dbSize)
		}
	}
	// Add new colors to the database and associate them with the product
	for _, color := range product.Colors {
		var dbColor models.Color
		if err := db.DB.Where("name = ?", color.Name).First(&dbColor).Error; err != nil {
			if err := db.DB.Create(&color).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating new color"})
				return
			}
			dbColor = color
		}
		db.DB.Model(&product).Association("Colors").Append(dbColor)
	}

	// Save the product to the database
	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}
	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}
