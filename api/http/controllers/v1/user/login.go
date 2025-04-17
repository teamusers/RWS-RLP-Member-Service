package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rlp-member-service/api/http/responses"
	"rlp-member-service/api/interceptor"
	"rlp-member-service/codes"
	"rlp-member-service/model"
	"rlp-member-service/system"
)

func Login(c *gin.Context) {
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

	db := system.GetDb()
	var user model.User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			resp := responses.APIResponse{
				Message: "email not found",
				Data: responses.LoginResponse{
					LoginSessionToken: "",
					LoginExpireIn:     0,
				},
			}
			c.JSON(codes.CODE_EMAIL_NOTFOUND, resp)
			return
		}
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

	if err := db.Save(&user).Error; err != nil {
		resp := responses.ErrorResponse{
			Error: "Failed to update user with session token",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "email found",
		Data: responses.LoginResponse{
			LoginSessionToken: token,
			LoginExpireIn:     expiresAt,
		},
	}
	c.JSON(http.StatusOK, resp)

}
