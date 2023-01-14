package middleware

import (
	"net/http"
	"strings"

	"github.com/Similadayo/backend/db"
	"github.com/Similadayo/backend/utils"
	"github.com/gin-gonic/gin"
)

func Authorization(c *gin.Context) {
	// Get the token from the Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found in request headers"})
		c.Abort()
		return
	}
	// Remove the "Bearer " prefix from the token
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Verify and decode the token
	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Check if the token is blacklisted
	db, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Message(false, "Error connecting to database"))
		c.Abort()
		return
	}
	defer db.Close()

	blacklisted, err := utils.IsTokenBlacklisted(db, tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking token blacklist"})
		c.Abort()
		return
	}
	if blacklisted {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
		c.Abort()
		return
	}

	// Add the claims to the request's context
	c.Set("claims", claims)
}
