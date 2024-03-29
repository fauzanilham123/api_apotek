package controllers

import (
	"api_apotek/models"
	"api_apotek/utils/token"
	"math"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var activityLogMap = make(map[string]time.Time)

func GetAllLogActivity(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var logActivity []models.LogActivity

	sort := c.DefaultQuery("sort", "asc")
	// Default to ascending if not provided
	sortOrder := "ASC"
	if sort == "desc" {
		sortOrder = "DESC"
	}
	orderBy := c.DefaultQuery("orderBy", "id")

	searchQuery := c.Query("search")

	// If search query is provided, apply the search condition
	if searchQuery != "" {
		// Specify the columns you want to search in
		searchColumns := []string{"user_id", "activity"}

		// Create a dynamic OR query for each searchable column
		var orQueries []string
		var args []interface{}
		for _, column := range searchColumns {
			orQueries = append(orQueries, column+" LIKE ?")
			args = append(args, "%"+searchQuery+"%")
		}

		// Combine OR queries with AND condition
		db = db.Where(strings.Join(orQueries, " OR "), args...)
	}


	pagination := ExtractPagination(c)

	// Get all query parameters and loop through them
	queryParams := c.Request.URL.Query()
	// Remove 'page' and 'perPage' keys from queryParams
	delete(queryParams, "page")
	delete(queryParams, "perPage")
	delete(queryParams, "sort")
	delete(queryParams, "orderBy")
	delete(queryParams, "search")

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
	db.Model(&logActivity).Count(&totalCount)

	// Calculate the total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

	// Calculate the offset for pagination
	offset := (pagination.Page - 1) * pagination.PerPage

	// Apply pagination and sorting
	err := db.Order(orderBy + " " + sortOrder).Preload("User").Offset(offset).Limit(pagination.PerPage).Find(&logActivity).Error

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
		"data":         logActivity,
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

func activitylog(c *gin.Context, activity string) error {
	// Dapatkan UserID dari token otentikasi
	userID, errs := token.ExtractTokenID(c)
	if errs != nil {
		SendError(c, "error", "Error extracting user ID from token")
	}
	// Lakukan logging aktivitas di sini, contoh: simpan ke dalam database
	logActivity := models.LogActivity{
		UserID:   userID,
		Time:     time.Now(),
		Activity: activity,
		Method:   c.Request.Method,
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&logActivity).Error; err != nil {
		return err
	}
	return nil
}

func checkAndLogActivity(c *gin.Context, activityMessage string, response interface{}) {
	currentTime := time.Now()

	// Periksa apakah aktivitas sudah ada dalam map
	lastActivityTime, exists := activityLogMap[activityMessage]

	if exists && currentTime.Sub(lastActivityTime).Minutes() < 3 {
		// Jika kurang dari 3 menit, jangan masukkan ke logactivity
		SendResponse(c, response, "success")
		return
	}
	SendResponse(c, response, "success")

	// Jika waktu lebih dari 3 menit atau tidak ada aktivitas sebelumnya
	// Lakukan proses logactivity seperti biasa
	activitylogShow(c, activityMessage)
	activityLogMap[activityMessage] = currentTime // Perbarui waktu aktivitas terakhir
}

func activitylogShow(c *gin.Context, activity string) error {
	// Lakukan logging aktivitas di sini, contoh: simpan ke dalam database
	logActivity := models.LogActivity{
		UserID:   1,
		Time:     time.Now(),
		Activity: activity,
		Method:   c.Request.Method,
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&logActivity).Error; err != nil {
		return err
	}
	return nil
}

func activitylogin(c *gin.Context, activity string, userID uint) error {
	// Lakukan logging aktivitas di sini, contoh: simpan ke dalam database
	logActivity := models.LogActivity{
		UserID:   userID,
		Time:     time.Now(),
		Activity: activity,
		Method:   c.Request.Method,
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&logActivity).Error; err != nil {
		return err
	}
	return nil
}