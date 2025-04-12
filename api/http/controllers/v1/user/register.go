package home

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rlp-member-service/api/http/requests"
	"rlp-member-service/api/http/responses"
	model "rlp-member-service/models"
	"rlp-member-service/system"
)

// GetUsers handles GET /users
// If a user with the provided email already exists, it returns an error that the email already exists.
// If no user is found, it continues to generate an OTP.
func GetUser(c *gin.Context) {
	var req requests.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := responses.ErrorResponse{
			Error: "Valid email and sign_up_type are required in the request body",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	email := req.Email

	db := system.GetDb()

	var user model.User
	err := db.Where("email = ?", email).First(&user).Error
	if err == nil {

		resp := responses.APIResponse{
			Message: "email registered",
			Data:    user,
		}
		c.JSON(http.StatusConflict, resp)
		return
	}
	if err != gorm.ErrRecordNotFound {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "email not registered",
	}
	c.JSON(http.StatusConflict, resp)
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
		resp := responses.ErrorResponse{
			Error: "Email already exists",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	} else if err != gorm.ErrRecordNotFound {
		// Some other error occurred while querying.
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
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

	resp := responses.APIResponse{
		Message: "user created",
		Data:    user,
	}
	c.JSON(http.StatusCreated, resp)

}
