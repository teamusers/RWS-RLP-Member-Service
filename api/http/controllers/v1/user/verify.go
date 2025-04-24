package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"rlp-member-service/api/http/responses"
	"rlp-member-service/codes"
	"rlp-member-service/model"
	"rlp-member-service/system"
)

func VerifyRegister(c *gin.Context) {
	type RequestBody struct {
		Email string     `json:"email"`
	}

	var body RequestBody
	c.ShouldBindJSON(&body)
	email := body.Email
	if email == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Error: "Email is required"})
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
			Code: 	codes.FOUND,
			
		})
		return
	}
		c.JSON(404, responses.ApiResponse[any]{
			Message: "user not found",
			Data:    nil,
			Code: 	codes.NOT_FOUND,
			
		})

}

