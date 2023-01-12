package utils

import (
	"github.com/gin-gonic/gin"
)

type Respond struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Response(c *gin.Context, status int, data interface{}) {
	var respond Respond
	if data != nil {
		respond = Respond{
			Success: true,
			Data:    data,
		}
	} else {
		respond = Respond{
			Success: false,
			Error:   "Error processing request",
		}
	}
	c.JSON(status, respond)
}
