package home

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rlp-member-service/api/http/services"
	model "rlp-member-service/models"
	"rlp-member-service/system"
)

// SignUpRequest represents the expected JSON structure for the request body.
type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// GetUsers handles GET /users
// If a user with the provided email already exists, it returns an error that the email not exists.
// If no user is found, it continues to generate an OTP.
func Login(c *gin.Context) {
	var req LoginRequest
	// Bind the JSON payload to LoginRequest struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid email is required in the request body"})
		return
	}
	email := req.Email

	// Get a database handle.
	db := system.GetDb()

	// Attempt to find a user by email.
	var user model.User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		// If no user is found, return an error.
		if err == gorm.ErrRecordNotFound {
			c.JSON(201, gin.H{
				"message": "email not found",
				"data": gin.H{
					"loginSessionToken": nil,
					"login_expireIn":    nil,
				},
			})
			return
		}
		// For any other errors, return an internal server error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the response with the custom JSON format.
	c.JSON(http.StatusOK, gin.H{
		"message": "email found",
		"data": gin.H{
			"loginSessionToken": user.SessionToken,
			"login_expireIn":    user.SessionExpiry,
		},
	})

}
func CreateUserSessionToken(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid email is required in the request body"})
		return
	}
	email := req.Email

	db := system.GetDb()

	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	otpService := services.NewOTPService()
	ctx := context.Background()
	otpResp, err := otpService.GenerateOTP(ctx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		return
	}

	user.SessionToken = otpResp.OTP
	user.SessionExpiry = otpResp.ExpiresAt

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user with session token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session token generated",
		"data": gin.H{
			"loginSessionToken": user.SessionToken,
			"login_expireIn":    user.SessionExpiry,
		},
	})
}

