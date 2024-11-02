package controllers

import (
	"bolalar-akademiyasi/database"
	"bolalar-akademiyasi/models"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var stringToStatus = map[string]models.Status{
	"active":   models.Active,
	"inactive": models.Inactive,
	"pending":  models.Pending,
	"agree":    models.Agree,
}

func GetClients(c *gin.Context) {
	var clients []models.Client
	var totalClients int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "id:desc") // Default to sorting by id
	search := c.DefaultQuery("search", "")
	statusFilter := c.DefaultQuery("status", "")

	offset := (page - 1) * limit

	// Validate and split the sort parameter
	sortParts := strings.Split(sort, ":")
	if len(sortParts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort parameter"})
		return
	}
	sortField := sortParts[0]
	sortOrder := sortParts[1]

	// Convert sortField to the correct database column name if necessary
	switch sortField {
	case "CreatedAt":
		sortField = "created_at"
	case "Status":
		sortField = "status"
	}

	query := database.DB.Model(&models.Client{})

	// Apply search if provided
	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("name ILIKE ? OR phone_number ILIKE ?", searchQuery, searchQuery)
	}

	// Apply status filter if provided
	if statusFilter != "" {
		status, ok := stringToStatus[statusFilter]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status filter"})
			return
		}
		query = query.Where("status = ?", status.String())
	}

	// Count total clients (after applying search and status filter)
	if err := query.Count(&totalClients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count clients"})
		return
	}

	// Apply sorting
	if sortOrder == "desc" {
		query = query.Order(sortField + " DESC")
	} else {
		query = query.Order(sortField + " ASC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&clients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clients"})
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalClients) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":       clients,
		"page":       page,
		"limit":      limit,
		"sort":       sort,
		"search":     search,
		"status":     statusFilter,
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
	if client.ChatID != 0 {
		client.Source = models.Telegram
	}
	client.Status = models.Active

	if err := database.DB.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Client created successfully",
		"client":  client,
	})
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
	if err := database.DB.Save(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Client updated successfully",
		"client":  client,
	})
}

func DeleteClient(c *gin.Context) {
	id := c.Param("id")
	var client models.Client
	if err := database.DB.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	if err := database.DB.Delete(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Client deleted successfully",
		"id":      id,
	})
}
