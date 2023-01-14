package controllers

import (
	"encoding/json"
	"net/http"
	"time"

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
	if foundUser.Role != "is_admin" {
		c.JSON(http.StatusForbidden, utils.Message(false, "You do not have permission to access this resource"))
		return
	}
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"token":   token,
	})
}

func Logout(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found in cookie"})
		return
	}

	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token in claims"})
		return
	}

	// Add the token to the blacklist
	expiresAt := time.Unix(claims.ExpiresAt, 0).UTC()
	if err := utils.AddToBlacklist(db.DB, tokenString, expiresAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add token to blacklist"})
		return
	}
	// Clear the token cookie
	c.SetCookie("token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func UpdateUser(c *gin.Context) {
	// Get the user ID from the request parameters
	userID := c.Param("id")

	// Get the authenticated user's ID from the JWT
	authUserID := c.MustGet("user_id").(string)

	// Verify that the authenticated user has permission to update the user
	if authUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this user"})
		return
	}

	// Get the updated user information from the request body
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user in the database
	if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func HomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to first closet")
}
