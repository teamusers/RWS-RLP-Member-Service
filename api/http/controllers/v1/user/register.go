package home

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rlp-member-service/api/http/responses"
	"rlp-member-service/codes"
	model "rlp-member-service/models"
	"rlp-member-service/system"
)

func GetUser(c *gin.Context) {
	email := c.Param("email")
	signUpType := c.Param("sign_up_type")

	if email == "" || signUpType == "" {
		resp := responses.ErrorResponse{
			Error: "Valid email and sign_up_type are required as query parameters",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	db := system.GetDb()

	var user model.User
	err := db.Where("email = ?", email).First(&user).Error
	if err == nil {

		resp := responses.APIResponse{
			Message: "email registered",
			Data:    user,
		}
		c.JSON(codes.CODE_EMAIL_REGISTERED, resp)
		return
	}
	if err != gorm.ErrRecordNotFound {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(codes.CODE_EMAIL_REGISTERED, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "email not registered",
	}
	c.JSON(http.StatusOK, resp)
}

func CreateUser(c *gin.Context) {
	db := system.GetDb()
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	var existingUser model.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		resp := responses.ErrorResponse{
			Error: "Email already exists",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	} else if err != gorm.ErrRecordNotFound {
		// Some other error occurred while querying.
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := db.Create(&user).Error; err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "user created",
		Data:    user,
	}
	c.JSON(http.StatusCreated, resp)

}
