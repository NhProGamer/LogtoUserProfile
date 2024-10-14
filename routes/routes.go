package routes

import (
	"LogtoUserProfile/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", controllers.Home)
	router.GET("/sign-in", controllers.SignIn)
	router.GET("/sign-in-callback", controllers.SignInCallback)
	router.GET("/user-id-token-claims", controllers.UserIdTokenClaims)
	router.GET("/sign-out", controllers.SignOut)
	router.GET("/protected", controllers.Protected)
}
