package controllers

import (
	"bolalar-akademiyasi/database"
	"bolalar-akademiyasi/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetClients(c *gin.Context) {
	var clients []models.Client
	database.DB.Find(&clients)
	c.JSON(http.StatusOK, clients)
}

func GetClient(c *gin.Context) {
	id := c.Param("id")
	var client models.Client
	if err := database.DB.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	c.JSON(http.StatusOK, client)
}

func CreateClient(c *gin.Context) {
	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&client)
	c.JSON(http.StatusOK, client)
}

func UpdateClient(c *gin.Context) {
	id := c.Param("id")
	var client models.Client
	if err := database.DB.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&client)
	c.JSON(http.StatusOK, client)
}

func DeleteClient(c *gin.Context) {
	id := c.Param("id")
	var client models.Client
	if err := database.DB.Delete(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Client deleted"})
}
