package controllers

import (
	"LogtoUserProfile/globals"
	"LogtoUserProfile/storage"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/logto-io/go/client"
	"net/http"
)

var (
	ContentTypeHtml = "text/html; charset=utf-8"
)

func Home(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(getLogtoConfig(), &storage.SessionStorage{Session: session})

	authState := "You are not logged in to this website. :("
	if logtoClient.IsAuthenticated() {
		authState = "You are logged in to this website! :)"
	}

	homePage := `<h1>Hello Logto</h1>` +
		"<div>" + authState + "</div>" +
		`<div><a href="/sign-in">Sign In</a></div>` +
		`<div><a href="/sign-out">Sign Out</a></div>` +
		`<div><a href="/user-id-token-claims">ID Token Claims</a></div>` +
		`<div><a href="/protected">Protected Resource</a></div>`

	ctx.Data(http.StatusOK, ContentTypeHtml, []byte(homePage))
}

func SignIn(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(getLogtoConfig(), &storage.SessionStorage{Session: session})
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

func UserIdTokenClaims(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(getLogtoConfig(), &storage.SessionStorage{Session: session})

	idTokenClaims, err := logtoClient.GetIdTokenClaims()
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, idTokenClaims)
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

func Protected(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(getLogtoConfig(), &storage.SessionStorage{Session: session})

	if logtoClient.IsAuthenticated() {
		protectedPage := `
		<h1>Authenticated</h1>
		<div>Protected content</div>
		<div><a href="/">Home</a></div>
		`
		ctx.Data(http.StatusOK, ContentTypeHtml, []byte(protectedPage))
		return
	}

	unauthorizedPage := `
	<h1>Unauthorized</h1>
	<div>You cannot visit the protected content</div>
	<div><a href="/">Home</a></div>
	`
	ctx.Data(http.StatusOK, ContentTypeHtml, []byte(unauthorizedPage))
}

func getLogtoConfig() *client.LogtoConfig {
	return &client.LogtoConfig{
		Endpoint:  globals.Configuration.Logto.Endpoint,
		AppId:     globals.Configuration.Logto.AppId,
		AppSecret: globals.Configuration.Logto.AppSecret,
	}
}
