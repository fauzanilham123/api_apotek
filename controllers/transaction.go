package controllers

import (
	"api_apotek/models"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type inputTransaction struct{
	No 			int 		`json:"no"`
	Tanggal 	string 	`json:"tanggal"`
	Nama_kasir 	string 		`json:"nama_kasir"`
	Total_bayar int			`json:"total_bayar"`
	UserID   	uint      	`gorm:"column:user_id" json:"id_user"`
	DrugID   	uint      	`gorm:"column:drug_id" json:"id_drug"`
	RecipeID   	uint      	`gorm:"column:recipe_id" json:"id_recipe"`
}

func GetAllTransaction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var transaction []models.Transaction

	sort := c.DefaultQuery("sortBy", "asc")
	// Default to ascending if not provided
	sortOrder := "ASC"
	if sort == "desc" {
		sortOrder = "DESC"
	}
	orderBy := c.DefaultQuery("orderBy", "id")

	pagination := ExtractPagination(c)
	query := db.Where("flag = 1")

	// Get all query parameters and loop through them
	queryParams := c.Request.URL.Query()
	// Remove 'page' and 'perPage' keys from queryParams
	delete(queryParams, "page")
	delete(queryParams, "perPage")
	delete(queryParams, "sort")
	delete(queryParams, "orderBy")
	for column, values := range queryParams {
		value := values[0] // In case there are multiple values, we take the first one

		// Apply filtering condition if the value is not empty
		if value != "" {
			query = query.Where(column+" LIKE ?", "%"+value+"%")
		}
	}

	var totalCount int64
	query.Model(&transaction).Where("flag = 1").Count(&totalCount)

	// Calculate the total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

	// Calculate the offset for pagination
	offset := (pagination.Page - 1) * pagination.PerPage

	// Apply pagination and sorting
	err := query.Order(orderBy + " " + sortOrder).Preload("User").Preload("Drug").Preload("Recipe").Offset(offset).Limit(pagination.PerPage).Find(&transaction).Error
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
		"data":         transaction,
		"current_page": pagination.Page,
		"last_page":    lastPage,
		"per_page":     pagination.PerPage,
		"nextPage":     nextPage,
		"prevPage":     prevPage,
		"totalPages":   totalPages,
		"totalCount":   totalCount,
	}

	checkAndLogActivity(c, "Get all transaction", response)
}

func CreateTransaction(c *gin.Context) {
	// Validate input
	var input inputTransaction
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	// Create
	transaction := models.Transaction{No: input.No, Tanggal: input.Tanggal, Nama_kasir: input.Nama_kasir, Total_bayar: input.Total_bayar, UserID: input.UserID, DrugID: input.DrugID,RecipeID: input.RecipeID, Flag: 1, CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&transaction)

	strNumber := strconv.Itoa(input.No)

	SendResponse(c, transaction, "success")
	activityMessage := "Create transaction: " + strNumber
	activitylog(c, activityMessage)
}

func GetTransactionById(c *gin.Context) { // Get model if exist
	var transaction models.Transaction

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	checkAndLogActivity(c, "Get transaction by id "+c.Param("id"), transaction)
}

func UpdateTransaction(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var transaction models.Transaction
	if err := db.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input inputTransaction
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	oldName := transaction.No

	var updatedInput models.Transaction
	updatedInput.No = input.No
	updatedInput.Tanggal = input.Tanggal
	updatedInput.Total_bayar = input.Total_bayar
	updatedInput.Nama_kasir = input.Nama_kasir
	updatedInput.UserID = input.UserID
	updatedInput.RecipeID = input.RecipeID
	updatedInput.DrugID = input.DrugID
	updatedInput.UpdatedAt = time.Now()

	db.Model(&transaction).Updates(updatedInput)

	strNumber := strconv.Itoa(input.No)
	strNumber2 := strconv.Itoa(oldName)


	SendResponse(c, transaction, "success")
	activityMessage := "Update transaction:'" + strNumber2 + "' to '" + strNumber + "'"
	activitylog(c, activityMessage)
}

func DeleteTransaction(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var transaction models.Transaction
	if err := db.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Set the flag to 0
	if err := db.Model(&transaction).Update("flag", 0).Error; err != nil {
		SendError(c, "Failed to delete", err.Error())
		return
	}

	strNumber := strconv.Itoa(transaction.No)


	// Return success response
	SendResponse(c, transaction, "success")
	activityMessage := "Delete transaction: " + strNumber
	activitylog(c, activityMessage)
}