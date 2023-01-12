package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Similadayo/backend/db"
	"github.com/Similadayo/backend/models"
	"github.com/Similadayo/backend/utils"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	// Parse and validate user input
	var user models.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Message(false, "Invalid request body"))
		return
	}

	if user.Role == "" {
		user.Role = "customer"
	}
	if err := utils.Validate(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.Message(false, err.Error()))
		return
	}

	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer db.Close()

	// Check if email exists
	if err := db.Where("email = ?", user.Email).First(&models.User{}).Error; err == nil {
		c.JSON(http.StatusConflict, utils.Message(false, "Email already exists"))
		return
	}

	// Hash the user's password
	hashedPassword, err := utils.HashedPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error hashing password"))
		return
	}
	user.Password = string(hashedPassword)

	// Insert the user into the database

	// Save user to DB
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Generate a JSON web token
	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error generating JWT"))
		return
	}

	// Send a successful response to the client
	c.JSON(http.StatusOK, utils.Respond{
		Success: true,
		Data: gin.H{
			"token": token,
		},
	})
}

func Login(c *gin.Context) {
	// Parse and validate user input
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.Message(false, err.Error()))
		return
	}
	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, utils.Message(false, "Email and Password are required"))
		return
	}

	// Get user from database
	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error connecting to database"))
		return
	}
	defer db.Close()

	var foundUser models.User
	if err := db.Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, utils.Message(false, "Invalid email or password"))
		return
	}

	// Compare passwords
	if !utils.ComparePasswords(foundUser.Password, user.Password) {
		c.JSON(http.StatusBadRequest, utils.Message(false, "Invalid email or password"))
		return
	}

	// Generate JSON web token
	token, err := utils.GenerateToken(foundUser.Username, foundUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error generating JWT"))
		return
	}
	c.JSON(http.StatusOK, utils.Message(true, token))
}

func HomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to first closet")
}
