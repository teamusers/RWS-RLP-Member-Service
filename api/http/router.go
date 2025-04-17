package http

import (
	v1 "rlp-member-service/api/http/controllers/v1"
	user "rlp-member-service/api/http/controllers/v1/user"
	"rlp-member-service/api/interceptor"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {

	v1Group := e.Group("/v1")
	v1Group.POST("/auth", v1.AuthHandler)
	usersGroup := v1Group.Group("/user", interceptor.HttpInterceptor())
	{
		// The endpoints below will all require a valid access token.
		usersGroup.POST("/", user.VerifyOrRegisterOrLogin)
		usersGroup.PUT("/pin", user.UpdateBurnPin)
	}
}
