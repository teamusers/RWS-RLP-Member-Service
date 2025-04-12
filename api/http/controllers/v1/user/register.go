package home

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rlp-member-service/api/http/requests"
	model "rlp-member-service/models"
	"rlp-member-service/system"
)

// GetUsers handles GET /users
// If a user with the provided email already exists, it returns an error that the email already exists.
// If no user is found, it continues to generate an OTP.
func GetUser(c *gin.Context) {
	var req requests.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid email and sign_up_type are required in the request body"})
		return
	}
	email := req.Email

	db := system.GetDb()

	var user model.User
	err := db.Where("email = ?", email).First(&user).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "email registered",
			"data":    user,
		})
		return
	}
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "email not registered",
	})
}

func CreateUser(c *gin.Context) {
	db := system.GetDb()
	var user model.User
	// Bind the incoming JSON payload to the user struct.
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a user with the same email already exists.
	var existingUser model.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		// Record found - email already exists.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		// Some other error occurred while querying.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set timestamps for the new record.
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Create the user along with any associated phone numbers.
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"data":    user,
	})
}

// DeleteUser handles DELETE /users/:id - delete a user and cascade delete phone numbers.
func DeleteUser(c *gin.Context) {
	db := system.GetDb()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// The foreign key constraint (with ON DELETE CASCADE) in the database will handle the deletion of associated phone numbers.
	if err := db.Delete(&model.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
