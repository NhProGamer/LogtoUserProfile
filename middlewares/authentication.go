package middlewares

import (
	"LogtoUserProfile/globals"
	"LogtoUserProfile/storage"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/logto-io/go/client"
	"net/http"
)

func LogtoAuth(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(&globals.LogtoConfig, &storage.SessionStorage{Session: session})

	if logtoClient.IsAuthenticated() {
		ctx.Next()
	} else {
		ctx.Redirect(http.StatusTemporaryRedirect, "/sign-in")
		ctx.Abort()
	}
}
