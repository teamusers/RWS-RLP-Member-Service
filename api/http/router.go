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
		usersGroup.GET("/login/:email", user.Login)
		usersGroup.GET("/register/:email/:sign_up_type", user.GetUser)
		usersGroup.POST("/register", user.CreateUser)

		//GET - LBE-5 - /api/v1/user/pin - burn PIN update
	}
}
