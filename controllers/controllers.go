package controllers

import (
	"LogtoUserProfile/globals"
	"LogtoUserProfile/storage"
	"LogtoUserProfile/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/logto-io/go/client"
	"github.com/logto-io/go/core"
	"net/http"
)

var (
	ContentTypeHtml = "text/html; charset=utf-8"
)

func Home(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(getLogtoConfig(), &storage.SessionStorage{Session: session})

	if logtoClient.IsAuthenticated() {
		userInfos, err := utils.FetchUserInfos(logtoClient, getLogtoConfig())
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

func SignIn(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(getLogtoConfig(), &storage.SessionStorage{Session: session})
	/*signInUri, err := logtoClient.SignIn(&client.SignInOptions{
		RedirectUri: globals.Configuration.Server.ServerURL + "/sign-in-callback",
	})*/
	signInUri, err := logtoClient.SignIn(globals.Configuration.Server.ServerURL + "/sign-in-callback")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, signInUri)
}

func SignInCallback(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(getLogtoConfig(), &storage.SessionStorage{Session: session})

	err := logtoClient.HandleSignInCallback(ctx.Request)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func UserProfile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(getLogtoConfig(), &storage.SessionStorage{Session: session})

	userInfos, err := utils.FetchUserInfos(logtoClient, getLogtoConfig())
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, userInfos)
}

func SignOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(getLogtoConfig(), &storage.SessionStorage{Session: session})

	signOutUri, signOutErr := logtoClient.SignOut(globals.Configuration.Server.ServerURL)
	if signOutErr != nil {
		ctx.String(http.StatusOK, signOutErr.Error())
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, signOutUri)
}

func getLogtoConfig() *client.LogtoConfig {
	return &client.LogtoConfig{
		Endpoint:  globals.Configuration.Logto.Endpoint,
		AppId:     globals.Configuration.Logto.AppId,
		AppSecret: globals.Configuration.Logto.AppSecret,
		Scopes:    []string{core.UserScopeProfile, core.UserScopeCustomData, core.UserScopeEmail, core.ReservedScopeOpenId, "family_name", "familyName"},
	}
}
