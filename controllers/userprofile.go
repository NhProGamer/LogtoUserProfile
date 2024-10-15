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

func UserProfile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(&globals.LogtoConfig, &storage.SessionStorage{Session: session})

	if logtoClient.IsAuthenticated() {
		userInfos, err := utils.FetchUserInfos(logtoClient, &globals.LogtoConfig)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, userInfos)
	} else {
		ctx.String(http.StatusForbidden, "")
	}
}

func UpdateUserProfile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(&globals.LogtoConfig, &storage.SessionStorage{Session: session})

	if logtoClient.IsAuthenticated() {

		ctx.JSON(http.StatusOK, "")
	} else {
		ctx.String(http.StatusForbidden, "")
	}
}
