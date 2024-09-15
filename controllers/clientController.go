package controllers

import (
	"bolalar-akademiyasi/database"
	"bolalar-akademiyasi/models"
	"net/http"
	"strconv"
	"strings"
	"math"

	"github.com/gin-gonic/gin"
)

func GetClients(c *gin.Context) {
	var clients []models.Client
	var totalClients int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "created_at:desc")
	search := c.DefaultQuery("search", "")

	offset := (page - 1) * limit

	sortField := strings.Split(sort, ":")[0]
	sortOrder := strings.Split(sort, ":")[1]

	// Convert sortField to the correct database column name if necessary
	if sortField == "CreatedAt" {
		sortField = "created_at"
	}

	query := database.DB.Model(&models.Client{})

	// Apply search if provided
	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("name ILIKE ? OR phone_number ILIKE ?", searchQuery, searchQuery)
	}

	// Count total clients (after applying search)
	query.Count(&totalClients)

	// Apply sorting
	if sortOrder == "desc" {
		query = query.Order(sortField + " DESC")
	} else {
		query = query.Order(sortField)
	}

	// Apply pagination
	query.Offset(offset).Limit(limit).Find(&clients)

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalClients) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":       clients,
		"page":       page,
		"limit":      limit,
		"sort":       sort,
		"search":     search,
		"totalPages": totalPages,
		"totalItems": totalClients,
	})
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
