package home

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rlp-member-service/api/http/responses"
	"rlp-member-service/api/interceptor"
	model "rlp-member-service/models"
	"rlp-member-service/system"
)

// GetUsers handles GET /users
// If a user with the provided email already exists, it returns an error that the email not exists.
// If no user is found, it continues to generate an OTP.
func Login(c *gin.Context) {
	// Bind the JSON payload to LoginRequest struct.
	appID := c.GetHeader("AppID")
	if appID == "" {
		resp := responses.ErrorResponse{
			Error: "APPId not found",
		}
		c.JSON(http.StatusMethodNotAllowed, resp)
		return
	}

	email := c.Param("email")
	if email == "" {
		resp := responses.ErrorResponse{
			Error: "Valid email is required as query parameter",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

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
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	token, err := interceptor.GenerateToken(appID)
	if err != nil {
		resp := responses.ErrorResponse{
			Error: "Failed to generate token",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	expiration := 365 * 24 * time.Hour
	expiresAt := time.Now().Add(expiration).Unix()

	user.SessionToken = token
	user.SessionExpiry = expiresAt

	/*
		if err := db.Save(&user).Error; err != nil {
			resp := responses.ErrorResponse{
				Error: "Failed to update user with session token",
			}
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
	*/

	c.JSON(http.StatusOK, gin.H{
		"message": "email found",
		"data": gin.H{
			"loginSessionToken": token,
			"login_expireIn":    expiresAt,
		},
	})

}
