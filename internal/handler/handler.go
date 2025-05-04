package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func ProtectedEndpoint(c *gin.Context) {
	userID := c.GetString("userID")
	userEmail := c.GetString("userEmail")

	c.JSON(http.StatusOK, gin.H{
		"message": "You are authorized",
		"userID":  userID,
		"email":   userEmail,
	})
}
