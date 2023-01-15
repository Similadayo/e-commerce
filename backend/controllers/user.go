package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Similadayo/backend/db"
	"github.com/Similadayo/backend/mailer"
	"github.com/Similadayo/backend/models"
	"github.com/Similadayo/backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found in cookie"})
		return
	}
	// Verify the user's token
	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Get the user ID from the request parameters
	userID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the user is an admin
	if claims.Username != claims.Username {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only users can update other users"})
		return
	}

	// Bind the updated user data from the request body
	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	// Update the user in the database
	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to the database"})
		return
	}
	defer db.Close()

	if err := db.Model(&models.User{}).Where("id = ?", userID).Updates(updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating the user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func GetUser(c *gin.Context) {
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

	// Get the user ID from the request parameters
	userID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the user is an admin
	if claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can access this resource"})
		return
	}

	// Get the user from the database
	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error connecting to database"))
		return
	}
	defer db.Close()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   user,
	})
}

func GetUsers(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can access this route"})
		return
	}

	// Get all users from the database
	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}
	defer db.Close()

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	// Return the users with hashed passwords
	c.JSON(http.StatusOK, gin.H{"users": sanitizeUsers(users)})
}

func sanitizeUsers(users []models.User) []models.User {
	for i, user := range users {
		user.Password = ""
		users[i] = user
	}
	return users
}

func DeleteUser(c *gin.Context) {
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

	// Get the user ID from the request parameters
	userID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the user from the database
	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error connecting to database"))
		return
	}
	defer db.Close()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the user trying to delete their own account
	if claims.Username != user.Username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own account"})
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully deleted user",
	})
}

func SuspendUser(c *gin.Context) {
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

	// Get the user ID from the request parameters
	userID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the user from the database
	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error connecting to database"))
		return
	}
	defer db.Close()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Create new Suspension struct and associate it with the user
	var suspension models.Suspension
	suspension.UserID = user.ID
	suspension.StartTime = time.Now()
	suspension.EndTime = time.Now().Add(time.Duration(24) * time.Hour)
	suspension.Reason = "Violation of terms of service"
	db.Create(&suspension)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "User suspended successfully",
	})
}

func ForgotPassword(c *gin.Context) {
	// Get the email address from the request body
	var request struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get the user from the database
	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error connecting to database"))
		return
	}
	defer db.Close()

	var user models.User
	if err := db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user from database"})
		return
	}

	// Generate a password reset token
	token, err := utils.GeneratePasswordResetToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating password reset token"})
		return
	}

	// Send an email with the password reset link
	link := fmt.Sprintf("http://example.com/reset-password?token=%s", token)
	if err := mailer.SendPasswordResetEmail(user.Email, link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending password reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func ResetPassword(c *gin.Context) {
	// Get the token from the request body
	token := c.PostForm("token")

	// Verify the token
	claims, err := utils.VerifyPasswordResetToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	// Get the new password from the request body
	password := c.PostForm("password")

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Get the user from the database
	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error connecting to database"))
		return
	}
	defer db.Close()

	var user models.User
	if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Update the user's password
	user.Password = string(hashedPassword)
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Password reset successfully",
	})
}

func HomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to first closet")
}
