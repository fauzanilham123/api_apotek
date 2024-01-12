package controllers

import (
	"api_apotek/models"
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []models.User

	sort := c.DefaultQuery("sort", "asc")
	// Default to ascending if not provided
	sortOrder := "ASC"
	if sort == "desc" {
		sortOrder = "DESC"
	}

	pagination := ExtractPagination(c)

	// Get all query parameters and loop through them
	queryParams := c.Request.URL.Query()
	// Remove 'page' and 'perPage' keys from queryParams
	delete(queryParams, "page")
	delete(queryParams, "perPage")
	delete(queryParams, "sort")

	// Loop through query parameters for filtering
	for column, values := range queryParams {
		value := values[0] // In case there are multiple values, we take the first one

		// Apply filtering condition if the value is not empty
		if value != "" {
			db = db.Where(column+" LIKE ?", "%"+value+"%")
		}
	}

	// Count the total number of records
	var totalCount int64
	db.Model(&users).Count(&totalCount)

	// Calculate the total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

	// Calculate the offset for pagination
	offset := (pagination.Page - 1) * pagination.PerPage

	// Apply pagination and sorting
	err := db.Order("id " + sortOrder).Offset(offset).Limit(pagination.PerPage).Preload("User").Find(&users).Error
	if err != nil {
		SendError(c, "internal server error", err.Error())
		return
	}

	// Calculate "last_page" based on total pages
	lastPage := totalPages

	// Calculate "nextPage" and "prevPage"
	nextPage := pagination.Page + 1
	if nextPage > totalPages {
		nextPage = 1
	}

	prevPage := pagination.Page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	response := map[string]interface{}{
		"data":         users,
		"current_page": pagination.Page,
		"last_page":    lastPage,
		"per_page":     pagination.PerPage,
		"nextPage":     nextPage,
		"prevPage":     prevPage,
		"totalPages":   totalPages,
		"totalCount":   totalCount,
	}

	SendResponse(c, response, "success")
}
