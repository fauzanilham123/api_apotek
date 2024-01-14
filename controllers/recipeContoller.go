package controllers

import (
	"api_apotek/models"
	"math"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

	type recipeInput struct {
		No      			int       		`gorm:"unique" json:"no"`
		Tanggal 			string 			`json:"tanggal"`
		Nama_pasien 		string 			`json:"nama_pasien"`
		Nama_dokter 		string 			`json:"nama_dokter"`
		Nama_obat 			string 			`json:"obat_resep"`
		Jumlah_obat_resep 	int 			`json:"jumlah_obat_resep"`
		Flag 				int 			`json:"flag"`
	}


func GetAllRecipe(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var recipe []models.Recipe

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
		searchColumns := []string{"no", "tanggal", "nama_pasien","nama_dokter","nama_obat","jumlah_obat_resep"}

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
	query.Model(&recipe).Where("flag = 1").Count(&totalCount)

	// Calculate the total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

	// Calculate the offset for pagination
	offset := (pagination.Page - 1) * pagination.PerPage

	// Apply pagination and sorting
	err := query.Order(orderBy + " " + sortOrder).Offset(offset).Limit(pagination.PerPage).Find(&recipe).Error
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
		"data":         recipe,
		"current_page": pagination.Page,
		"last_page":    lastPage,
		"per_page":     pagination.PerPage,
		"nextPage":     nextPage,
		"prevPage":     prevPage,
		"totalPages":   totalPages,
		"totalCount":   totalCount,
	}

	checkAndLogActivity(c, "Get all recipe", response)
}

func CreateRecipe(c *gin.Context) {
	// Validate input
	var input recipeInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	// Create
	recipe := models.Recipe{No: input.No, Tanggal: input.Tanggal, Nama_pasien: input.Nama_pasien, Nama_dokter: input.Nama_dokter, Nama_obat: input.Nama_obat, Jumlah_obat_resep: input.Jumlah_obat_resep,Flag: 1, CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&recipe)

	SendResponse(c, recipe, "success")
	activityMessage := "Create recipe: " + input.Nama_obat
	activitylog(c, activityMessage)
}

func GetRecipeById(c *gin.Context) { // Get model if exist
	var recipe models.Recipe

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&recipe).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	checkAndLogActivity(c, "Get recipe by id "+c.Param("id"), recipe)
}

func UpdateRecipe(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var recipe models.Recipe
	if err := db.Where("id = ?", c.Param("id")).First(&recipe).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input recipeInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	oldName := recipe.Nama_obat

	var updatedInput models.Recipe
	updatedInput.No = input.No
	updatedInput.Tanggal = input.Tanggal
	updatedInput.Nama_dokter = input.Nama_dokter
	updatedInput.Nama_obat = input.Nama_obat
	updatedInput.Nama_pasien = input.Nama_pasien
	updatedInput.Jumlah_obat_resep = input.Jumlah_obat_resep
	updatedInput.UpdatedAt = time.Now()

	db.Model(&recipe).Updates(updatedInput)

	SendResponse(c, recipe, "success")
	activityMessage := "Update recipe:'" + oldName + "' to '" + input.Nama_obat + "'"
	activitylog(c, activityMessage)
}

func DeleteRecipe(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var recipe models.Recipe
	if err := db.Where("id = ?", c.Param("id")).First(&recipe).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Set the flag to 0
	if err := db.Model(&recipe).Update("flag", 0).Error; err != nil {
		SendError(c, "Failed to delete", err.Error())
		return
	}

	// Return success response
	SendResponse(c, recipe, "success")
	activityMessage := "Delete recipe: " + recipe.Nama_obat
	activitylog(c, activityMessage)
}