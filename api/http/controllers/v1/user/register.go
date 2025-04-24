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

func CreateUser(c *gin.Context) {
	appID := c.GetHeader("AppID")

	var body requests.Register
	c.ShouldBindJSON(&body)

	db := system.GetDb()
	newUser := body.User
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	token, err := interceptor.GenerateToken(appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}
	expiration := 365 * 24 * time.Hour
	newUser.SessionToken = token
	newUser.SessionExpiry = time.Now().Add(expiration).Unix()

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusCreated, responses.ApiResponse[any]{
		Message: "user created",
		Data:    newUser,
		Code:    codes.SUCCESSFUL,
	})
}

func VerifyUserExistence(c *gin.Context) {

	var body requests.RegisterVerification
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
		c.JSON(http.StatusOK, responses.ApiResponse[any]{
			Message: "user found",
			Data:    nil,
			Code:    codes.FOUND,
		})
		return
	}
	c.JSON(404, responses.ApiResponse[any]{
		Message: "user not found",
		Data:    nil,
		Code:    codes.NOT_FOUND,
	})

}
