package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"rlp-member-service/api/http/requests"
	"rlp-member-service/api/http/responses"
	"rlp-member-service/api/interceptor"
	"rlp-member-service/codes"
	"rlp-member-service/model"
	"rlp-member-service/system"
)

func Login(c *gin.Context) {
	appID := c.GetHeader("AppID")


	var body requests.LoginRequest
	c.ShouldBindJSON(&body)
	email := body.Email
	if email == "" {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	db := system.GetDb()
	var user model.User

	err := db.Where("email = ?", email).First(&user).Error
	userExists := (err == nil)

	if userExists {
			token, err := interceptor.GenerateToken(appID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
				return
			}
			expiration := 365 * 24 * time.Hour
			user.SessionToken = token
			user.SessionExpiry = time.Now().Add(expiration).Unix()

			if err := db.Save(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
				return
			}
				
			c.JSON(http.StatusOK, responses.ApiResponse[responses.LoginResponse]{
				Message: "user updated",
				Data:    responses.LoginResponse {LoginSessionToken: user.SessionToken, LoginExpireIn: user.SessionExpiry},
				Code: 	codes.FOUND,
			})
			return
	

	}
	c.JSON(http.StatusOK, responses.ApiResponse[responses.LoginResponse]{
		Message: "user not found",
		Data:    responses.LoginResponse {},
		Code: 	codes.NOT_FOUND,
	})
}