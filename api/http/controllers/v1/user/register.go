package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"rlp-member-service/api/http/requests"
	"rlp-member-service/api/http/responses"
	"rlp-member-service/api/interceptor"
	"rlp-member-service/codes"
	"rlp-member-service/system"
)

func Register(c *gin.Context) {
	appID := c.GetHeader("AppID")


	var body requests.Register
	c.ShouldBindJSON(&body)

	db := system.GetDb()
	
	newUser := body.User
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

		token, err := interceptor.GenerateToken(appID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to generate token"})
			return
		}
		expiration := 365 * 24 * time.Hour
		newUser.SessionToken = token
		newUser.SessionExpiry = time.Now().Add(expiration).Unix()

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, responses.ApiResponse[any]{
		Message: "user created",
		Data:    newUser,
		Code: codes.SUCCESSFUL,
	})
}