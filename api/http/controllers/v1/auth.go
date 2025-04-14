package v1

import (
	"crypto/hmac"
	"fmt"
	"net/http"

	"rlp-member-service/codes"
	"rlp-member-service/model"
	"rlp-member-service/system"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rlp-member-service/api/http/requests"
	"rlp-member-service/api/http/responses"
	"rlp-member-service/api/http/services"
	"rlp-member-service/api/interceptor"
)

func getSecretKey(db *gorm.DB, appID string) (string, error) {
	var channel model.SysChannel
	if err := db.Where("app_id = ?", appID).First(&channel).Error; err != nil {
		return "", fmt.Errorf("failed to get secret key for appID %s: %w", appID, err)
	}
	return channel.AppKey, nil
}
func AuthHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		resp := responses.ErrorResponse{
			Error: "Method Not Allowed",
		}
		c.JSON(http.StatusMethodNotAllowed, resp)
		return
	}
	if c.GetHeader("Content-Type") != "application/json" {
		resp := responses.ErrorResponse{
			Error: "Content-Type must be application/json",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	appID := c.GetHeader("AppID")
	if appID == "" {
		resp := responses.APIResponse{
			Message: "invalid appid",
			Data: responses.AuthResponse{
				AccessToken: "",
			},
		}
		c.JSON(codes.CODE_INVALID_APPID, resp)
		return
	}
	var req requests.AuthRequest
	if err := c.BindJSON(&req); err != nil {
		resp := responses.ErrorResponse{
			Error: "Invalid JSON body",
		}
		c.JSON(http.StatusMethodNotAllowed, resp)
		return
	}

	db := system.GetDb()
	secretKey, err := getSecretKey(db, appID)

	if err != nil || secretKey == "" {
		resp := responses.APIResponse{
			Message: "invalid appid",
			Data: responses.AuthResponse{
				AccessToken: "",
			},
		}
		c.JSON(codes.CODE_INVALID_APPID, resp)
		return
	}

	authReq, err := services.GenerateSignatureWithParams(appID, req.Nonce, req.Timestamp, secretKey)
	if err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	expectedSignature := authReq.Signature
	if !hmac.Equal([]byte(expectedSignature), []byte(req.Signature)) {
		resp := responses.APIResponse{
			Message: "invalid signature",
			Data: responses.AuthResponse{
				AccessToken: "",
			},
		}
		c.JSON(codes.CODE_INVALID_SIGNATURE, resp)
		return
	}

	token, err := interceptor.GenerateToken(appID)
	if err != nil {
		resp := responses.ErrorResponse{
			Error: "Failed to generate token",
		}
		c.JSON(http.StatusMethodNotAllowed, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "token successfully generated",
		Data: responses.AuthResponse{
			AccessToken: token,
		},
	}
	c.JSON(http.StatusOK, resp)
}
