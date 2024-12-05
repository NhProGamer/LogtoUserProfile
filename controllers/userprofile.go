package controllers

import (
	"LogtoUserProfile/globals"
	"LogtoUserProfile/logto"
	"LogtoUserProfile/storage"
	"LogtoUserProfile/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/logto-io/go/client"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
	tokenClaims, err := logtoClient.GetIdTokenClaims()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if logtoClient.IsAuthenticated() {
		var payload interface{}
		avatar := ctx.DefaultQuery("avatar", "")
		name := ctx.DefaultQuery("name", "")
		givenName := ctx.DefaultQuery("given_name", "")
		familyName := ctx.DefaultQuery("family_name", "")

		if givenName == "" && familyName == "" {
			payload = logto.PatchProfilePayloadLite{
				Avatar: avatar,
				Name:   name,
			}
		} else {
			payload = logto.PatchProfilePayload{
				Avatar: avatar,
				Name:   name,
				Profile: logto.ProfilePayload{
					GivenName:  givenName,
					FamilyName: familyName,
				},
			}
		}
		err := logto.PatchUserProfile(tokenClaims.Sub, payload)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.String(http.StatusOK, "")
	} else {
		ctx.String(http.StatusForbidden, "")
	}
}

func GetPfp(ctx *gin.Context) {
	userId := ctx.Query("userId")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user query parameter is required"})
		return
	}
	pattern := `^[a-zA-Z_][a-zA-Z0-9_]*$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(userId) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "illegal user id"})
		return
	}

	filePath := filepath.Join("pfp", userId+".gif")
	if _, err := filepath.Abs(filePath); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "profile picture not found"})
		return
	}
	ctx.Header("Content-Type", "image/webp")
	ctx.File(filePath)
}
func UploadPfp(ctx *gin.Context) {
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(&globals.LogtoConfig, &storage.SessionStorage{Session: session})
	tokenClaims, err := logtoClient.GetIdTokenClaims()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if logtoClient.IsAuthenticated() {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file"})
			return
		}

		src, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file: " + err.Error()})
			return
		}

		picture, err := utils.ConvertToGIF(src, filepath.Ext(file.Filename))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot convert file: " + err.Error()})
			return
		}

		savePath := filepath.Join("pfp", tokenClaims.Username)
		if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create directory"})
			return
		}

		pictureFile, err := os.Create(savePath + ".gif")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed while creating file on disk"})
			return
		}

		_, err = io.Copy(pictureFile, picture)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed while saving image"})
			return
		}

		src.Close()
		pictureFile.Close()

		ctx.JSON(http.StatusOK, gin.H{"message": "profile picture uploaded successfully"})
	}
}
