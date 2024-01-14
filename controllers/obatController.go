package controllers

import (
	"api_apotek/models"
	"math"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type obatInput struct {
	Kode        	string     		`gorm:"unique" json:"kode"`
	Name        	string  		`json:"name"`
	ExpiredDate     string      	`json:"expired_date"`
	Jumlah        	int      		`json:"jumlah"`
	HargaPerUnit    int      		`json:"harga_per_unit"`
	CreatedAt 		time.Time 		`json:"created_at"`
	UpdatedAt 		time.Time 		`json:"updated_at"`
}

func GetAllObat(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var obat []models.Drug

	sort := c.DefaultQuery("sortBy", "asc")
	// Default to ascending if not provided
	sortOrder := "ASC"
	if sort == "desc" {
		sortOrder = "DESC"
	}
	orderBy := c.DefaultQuery("orderBy", "id")

	pagination := ExtractPagination(c)
	query := db.Where("flag = 1")

	searchQuery := c.Query("search")

	// If search query is provided, apply the search condition
	if searchQuery != "" {
		// Specify the columns you want to search in
		searchColumns := []string{"kode", "obat", "name","expired_date","jumlah","harga_per_unit"}

		// Create a dynamic OR query for each searchable column
		var orQueries []string
		var args []interface{}
		for _, column := range searchColumns {
			orQueries = append(orQueries, column+" LIKE ?")
			args = append(args, "%"+searchQuery+"%")
		}

		// Combine OR queries with AND condition
		query = query.Where(strings.Join(orQueries, " OR "), args...)
	}


	// Get all query parameters and loop through them
	queryParams := c.Request.URL.Query()
	// Remove 'page' and 'perPage' keys from queryParams
	delete(queryParams, "page")
	delete(queryParams, "perPage")
	delete(queryParams, "sort")
	delete(queryParams, "orderBy")
	delete(queryParams, "search")
	for column, values := range queryParams {
		value := values[0] // In case there are multiple values, we take the first one

		// Apply filtering condition if the value is not empty
		if value != "" {
			query = query.Where(column+" LIKE ?", "%"+value+"%")
		}
	}

	var totalCount int64
	query.Model(&obat).Where("flag = 1").Count(&totalCount)

	// Calculate the total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

	// Calculate the offset for pagination
	offset := (pagination.Page - 1) * pagination.PerPage

	// Apply pagination and sorting
	err := query.Order(orderBy + " " + sortOrder).Offset(offset).Limit(pagination.PerPage).Find(&obat).Error
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
		"data":         obat,
		"current_page": pagination.Page,
		"last_page":    lastPage,
		"per_page":     pagination.PerPage,
		"nextPage":     nextPage,
		"prevPage":     prevPage,
		"totalPages":   totalPages,
		"totalCount":   totalCount,
	}

	checkAndLogActivity(c, "Get all obat", response)
}

func CreateObat(c *gin.Context) {
	// Validate input
	var input obatInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}


	// Create
	obat := models.Drug{Kode: input.Kode, Name: input.Name, ExpiredDate: input.ExpiredDate, Jumlah: input.Jumlah, HargaPerUnit: input.HargaPerUnit, Flag: 1, CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&obat)

	SendResponse(c, obat, "success")
	activityMessage := "Create obat: " + input.Name
	activitylog(c, activityMessage)
}

func GetObatById(c *gin.Context) { // Get model if exist
	var obat models.Drug

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&obat).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	checkAndLogActivity(c, "Get obat by id "+c.Param("id"), obat)
}

func UpdateObat(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var obat models.Drug
	if err := db.Where("id = ?", c.Param("id")).First(&obat).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input obatInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}



	oldName := obat.Name

	var updatedInput models.Drug
	updatedInput.Kode = input.Kode
	updatedInput.Name = input.Name
	updatedInput.ExpiredDate = input.ExpiredDate
	updatedInput.Jumlah = input.Jumlah
	updatedInput.HargaPerUnit = input.HargaPerUnit
	updatedInput.UpdatedAt = time.Now()

	db.Model(&obat).Updates(updatedInput)

	SendResponse(c, obat, "success")
	activityMessage := "Update obat:'" + oldName + "' to '" + input.Name + "'"
	activitylog(c, activityMessage)
}

func DeleteObat(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var obat models.Drug
	if err := db.Where("id = ?", c.Param("id")).First(&obat).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Set the flag to 0
	if err := db.Model(&obat).Update("flag", 0).Error; err != nil {
		SendError(c, "Failed to delete", err.Error())
		return
	}

	// Return success response
	SendResponse(c, obat, "success")
	activityMessage := "Delete obat: " + obat.Name
	activitylog(c, activityMessage)
}