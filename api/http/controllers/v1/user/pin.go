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
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	var user model.User
	err := db.Where("email = ?", burn_pin.Email).First(&user).Error
	userExists := (err == nil)

	if userExists {
		user.BurnPin = burn_pin.BurnPin
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to update user"})
			return
		}
	
		resp := responses.ApiResponse[any]{
			Message: "user updated",
			Data: nil,
			Code: codes.SUCCESSFUL,
		}
		c.JSON(http.StatusOK, resp)
	}
	c.JSON(http.StatusNotFound, responses.ApiResponse[any]{
		Message: "user not found",
		Data:  nil,
		Code: codes.NOT_FOUND,
	})
}