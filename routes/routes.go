package routes

import (
	"LogtoUserProfile/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", controllers.Home)
	router.GET("/sign-in", controllers.SignIn)
	router.GET("/sign-in-callback", controllers.SignInCallback)
	router.GET("/user-profile", controllers.UserProfile)
	router.GET("/sign-out", controllers.SignOut)

	apiV1 := router.Group("/api/v1")

	userprofile := apiV1.Group("userprofile")
	userprofile.GET("/", controllers.UserProfile)
	userprofile.PATCH("/", controllers.UpdateUserProfile)

	apiV1.PATCH("change-password", controllers.ChangePassword)
}
