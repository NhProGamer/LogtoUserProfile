package controllers

import (
	"LogtoUserProfile/globals"
	"LogtoUserProfile/logto"
	"LogtoUserProfile/storage"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/logto-io/go/client"
	"net/http"
)

func ChangePassword(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(&globals.LogtoConfig, &storage.SessionStorage{Session: session})
	tokenClaims, err := logtoClient.GetIdTokenClaims()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if logtoClient.IsAuthenticated() {
		oldPassword := ctx.DefaultQuery("oldPassword", "")
		newPassword := ctx.DefaultQuery("newPassword", "")
		if oldPassword == "" || newPassword == "" {
			ctx.String(http.StatusBadRequest, "")
			return
		}
		state, err := logto.VerifyUserPassword(tokenClaims.Sub, oldPassword)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		if !state {
			ctx.String(http.StatusUnprocessableEntity, "Bad password!")
			return
		}

		err = logto.PatchUserPassword(tokenClaims.Sub, newPassword)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, "")
	} else {
		ctx.String(http.StatusForbidden, "")
	}
}
