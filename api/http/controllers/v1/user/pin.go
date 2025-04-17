package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"rlp-member-service/api/http/requests"
	"rlp-member-service/api/http/responses"
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
	db.Where("email = ?", burn_pin.Email).First(&user)

	user.BurnPin = burn_pin.BurnPin
	if err := db.Create(&user).Error; err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to update user"})
		return
	}

	resp := responses.APIResponse{
		Message: "user updated",
	}
	c.JSON(http.StatusCreated, resp)

}