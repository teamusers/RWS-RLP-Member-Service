package user

import (
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"

	"rlp-member-service/api/http/responses"
	"rlp-member-service/api/interceptor"
	"rlp-member-service/model"
	"rlp-member-service/system"
)

func VerifyOrRegisterOrLogin(c *gin.Context) {
	email := c.Query("email")
	shouldUpdateToken := c.Query("updateSessionToken") == "true"
	appID := c.GetHeader("AppID")

	if email == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Error: "Email is required"})
		return
	}

	type RequestBody struct {
		User model.User `json:"user"`
	}
	var body RequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil && err.Error() != "EOF" {
		// Only return if it's not an empty body (EOF = empty but valid)
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Error: "Invalid request body"})
		return
	}

	db := system.GetDb()
	var user model.User

	err = db.Where("email = ?", email).First(&user).Error
	userExists := (err == nil)

	if userExists {
		if shouldUpdateToken {
			if appID == "" {
				c.JSON(http.StatusMethodNotAllowed, responses.ErrorResponse{Error: "AppID header required for session update"})
				return
			}
			token, err := interceptor.GenerateToken(appID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to generate token"})
				return
			}
			expiration := 365 * 24 * time.Hour
			user.SessionToken = token
			user.SessionExpiry = time.Now().Add(expiration).Unix()

			if err := db.Save(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to update user"})
				return
			}
		}

		c.JSON(http.StatusOK, responses.APIResponse{
			Message: "user found",
			Data:    user,
		})
		return
	}

	// If user is not found and no body was provided, don't create

	if reflect.DeepEqual(body.User, model.User{}) {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{Error: "User not found"})
		return
	}

	// Create new user
	newUser := body.User
	newUser.Email = email
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	if appID != "" {
		token, err := interceptor.GenerateToken(appID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to generate token"})
			return
		}
		expiration := 365 * 24 * time.Hour
		newUser.SessionToken = token
		newUser.SessionExpiry = time.Now().Add(expiration).Unix()
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Message: "user created",
		Data:    newUser,
	})
}

