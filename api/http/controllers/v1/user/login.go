package home

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rlp-member-service/api/http/requests"
	"rlp-member-service/api/interceptor"
	model "rlp-member-service/models"
	"rlp-member-service/system"
)

// GetUsers handles GET /users
// If a user with the provided email already exists, it returns an error that the email not exists.
// If no user is found, it continues to generate an OTP.
func Login(c *gin.Context) {
	var req requests.LoginRequest
	// Bind the JSON payload to LoginRequest struct.
	appID := c.GetHeader("AppID")
	if appID == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "APPId not found",
		})
		return
	}

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

	token, err := interceptor.GenerateToken(appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	expiration := 365 * 24 * time.Hour
	expiresAt := time.Now().Add(expiration).Unix()

	user.SessionToken = token
	user.SessionExpiry = expiresAt

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user with session token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "email found",
		"data": gin.H{
			"loginSessionToken": token,
			"login_expireIn":    expiresAt,
		},
	})

}
