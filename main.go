package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/logto-io/go/client"
	"net/http"
)

var (
	ContentTypeHtml = "text/html; charset=utf-8"
	BASE_URL        = "http://localhost:5000"
)

func main() {
	logtoConfig := &client.LogtoConfig{
		// see .env.example for more details and examples
		Endpoint:  "",
		AppId:     "",
		AppSecret: "",
	}

	store := memstore.NewStore([]byte(""))

	router := gin.Default()

	router.Use(sessions.Sessions("logto-session", store))

	router.GET("/", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(logtoConfig, &SessionStorage{session: session})

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
	})

	router.GET("/sign-in", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(logtoConfig, &SessionStorage{session: session})
		signInUri, err := logtoClient.SignIn(BASE_URL + "/sign-in-callback")
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Redirect(http.StatusTemporaryRedirect, signInUri)
	})

	router.GET("/sign-in-callback", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(logtoConfig, &SessionStorage{session: session})

		err := logtoClient.HandleSignInCallback(ctx.Request)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, "/")
	})

	router.GET("/user-id-token-claims", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(logtoConfig, &SessionStorage{session: session})

		idTokenClaims, err := logtoClient.GetIdTokenClaims()

		if err != nil {
			ctx.String(http.StatusOK, err.Error())
		}

		ctx.JSON(http.StatusOK, idTokenClaims)
	})

	router.GET("/sign-out", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(logtoConfig, &SessionStorage{session: session})

		signOutUri, signOutErr := logtoClient.SignOut(BASE_URL)

		if signOutErr != nil {
			ctx.String(http.StatusOK, signOutErr.Error())
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, signOutUri)
	})

	router.GET("/protected", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(logtoConfig, &SessionStorage{session: session})

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
	})

	router.Run("0.0.0.0:5000")
}
