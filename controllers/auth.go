package controllers

import (
	"bolalar-akademiyasi/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Auth(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	validUsername := os.Getenv("AUTH_USERNAME")
	validPassword := os.Getenv("AUTH_PASSWORD")

	if credentials.Username == validUsername && credentials.Password == validPassword {
		token, err := utils.GenerateToken(credentials.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
	}
}
