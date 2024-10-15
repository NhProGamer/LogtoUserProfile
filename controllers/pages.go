package controllers

import (
	"LogtoUserProfile/globals"
	"LogtoUserProfile/storage"
	"LogtoUserProfile/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/logto-io/go/client"
	"net/http"
)

func Home(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(&globals.LogtoConfig, &storage.SessionStorage{Session: session})

	if logtoClient.IsAuthenticated() {
		userInfos, err := utils.FetchUserInfos(logtoClient, &globals.LogtoConfig)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"name":           userInfos.Name,
			"username":       userInfos.Username,
			"email":          userInfos.Email,
			"profilePicture": userInfos.Picture,
			"givenName":      userInfos.GivenName,
			"familyName":     userInfos.FamilyName,
		})
	} else {
		ctx.Redirect(http.StatusTemporaryRedirect, "/sign-in")
	}
}
