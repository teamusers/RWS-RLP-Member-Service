package user

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"rlp-member-service/api/http/requests"
	"rlp-member-service/api/http/responses"
	"rlp-member-service/codes"
	"rlp-member-service/model"
	"rlp-member-service/system"
)

func CreateUser(c *gin.Context) {

	var body requests.Register
	c.ShouldBindJSON(&body)

	db := system.GetDb()
	newUser := body.User
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	if err := db.Create(&newUser).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "UNIQUE") {
			c.JSON(http.StatusConflict, responses.InternalErrorResponse())
			return
		}

		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, responses.ApiResponse[any]{
		Message: "user created",
		Data:    newUser,
		Code:    codes.SUCCESSFUL,
	})
}

func VerifyUserExistenceByField(c *gin.Context) {
	field := c.Param("field")

	allowed := map[string]bool{"email": true, "gr_id": true}
	if !allowed[field] {
		c.JSON(http.StatusBadRequest, responses.InvalidQueryParametersErrorResponse())
		return
	}

	var body map[string]string
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	value, ok := body[field]
	if !ok || value == "" {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	db := system.GetDb()
	var user model.User

	err := db.Where(fmt.Sprintf("%v = ?", field), value).First(&user).Error
	userExists := (err == nil)

	if userExists {
		c.JSON(http.StatusConflict, responses.ApiResponse[any]{
			Code:    codes.FOUND,
			Message: "user found",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, responses.ApiResponse[any]{
		Code:    codes.NOT_FOUND,
		Message: "user not found",
		Data:    nil,
	})

}
