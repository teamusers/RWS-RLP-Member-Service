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
		usersGroup.POST("/register", user.CreateUser)
		usersGroup.PUT("/login", user.Login)
		usersGroup.PUT("/pin", user.UpdateBurnPin)
		usersGroup.POST("/verify", user.VerifyUserExistence)
	}
}
