package controllers

import (
	"api_apotek/models"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

	type inputUser struct {
		Username    string        	`gorm:"not null;unique" json:"username"`
		Name 		string 		  	`gorm:"not null;unique" json:"name"`
		Email       string        	`gorm:"not null;unique" json:"email"`
		Password    string        	`gorm:"not null;" json:"password"`
		Role    	string        	`gorm:"not null" json:"role"`
	}

func GetAllUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []models.User

	sort := c.DefaultQuery("sort", "asc")
	// Default to ascending if not provided
	sortOrder := "ASC"
	if sort == "desc" {
		sortOrder = "DESC"
	}

	orderBy := c.DefaultQuery("orderBy", "id")


	pagination := ExtractPagination(c)

	// Get all query parameters and loop through them
	queryParams := c.Request.URL.Query()
	// Remove 'page' and 'perPage' keys from queryParams
	delete(queryParams, "page")
	delete(queryParams, "perPage")
	delete(queryParams, "sort")
	delete(queryParams, "orderBy")

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
	err := db.Order(orderBy + " " + sortOrder).Offset(offset).Limit(pagination.PerPage).Find(&users).Error
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

func CreateUser(c *gin.Context) {
	// Validate input
	db := c.MustGet("db").(*gorm.DB)
	var input inputUser
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}


	// Check if email already exists in the database
	var existingUser models.User
	if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		SendError(c, "Email already registered", "Email already registered")
		return
	}

	u := models.User{}
	u.Username = input.Username
	u.Name = input.Name
	u.Email = input.Email
	u.Password = input.Password
	u.Role = input.Role

	savedUser, err := u.SaveUser(db)

	if err != nil {
		SendError(c, "error", err.Error())
		return
	}

	user := map[string]string{
		"username": input.Username,
		"email":    input.Email,
		"name":    input.Name,
		"role":    input.Role,
	}
	SendResponse(c, user, "create user success")
	activityMessage := "Create user"
	activitylogin(c, activityMessage, savedUser.ID)
}

func GetUserById(c *gin.Context) { // Get model if exist
	var user models.User

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	checkAndLogActivity(c, "Get user by id "+c.Param("id"), user)
}

func UpdateUser(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input inputUser
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}



	oldName := user.Name

	var updatedInput models.User
	updatedInput.Username = input.Username
	updatedInput.Name = input.Name
	updatedInput.Password = input.Password
	updatedInput.Email = input.Email
	updatedInput.Role = input.Role
	updatedInput.UpdatedAt = time.Now()

	db.Model(&user).Updates(updatedInput)

	SendResponse(c, user, "success")
	activityMessage := "Update user:'" + oldName + "' to '" + input.Name + "'"
	activitylog(c, activityMessage)
}

func DeleteUser(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Hapus log aktivitas terkait
	if err := db.Where("user_id = ?", user.ID).Delete(&models.LogActivity{}).Error; err != nil {
		SendError(c, "Failed to delete user's log activities", err.Error())
		return
	}

	// Hapus pengguna
	if err := db.Delete(&user).Error; err != nil {
		SendError(c, "Failed to delete user", err.Error())
		return
	}
	// Return success response
	SendResponse(c, user, "success")
	activityMessage := "Delete user: " + user.Name
	activitylog(c, activityMessage)
}
