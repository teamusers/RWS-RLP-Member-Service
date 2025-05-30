package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"rlp-member-service/api/http/requests"
	"rlp-member-service/api/http/responses"
	"rlp-member-service/codes"
	"rlp-member-service/model"
	"rlp-member-service/system"
)

func UpdateBurnPin(c *gin.Context) {
	db := system.GetDb()
	var burn_pin requests.UpdateBurnPin
	if err := c.ShouldBindJSON(&burn_pin); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}
	var user model.User
	err := db.Where("email = ?", burn_pin.Email).First(&user).Error
	userExists := (err == nil)

	if userExists {
		user.BurnPin = burn_pin.BurnPin
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
			return
		}

		resp := responses.ApiResponse[any]{
			Code:    codes.SUCCESSFUL,
			Message: "user updated",
			Data:    nil,
		}
		c.JSON(http.StatusOK, resp)
	}
	c.JSON(http.StatusConflict, responses.ApiResponse[any]{
		Code:    codes.NOT_FOUND,
		Message: "user not found",
		Data:    nil,
	})
}
