package main

import (
	"LogtoUserProfile/config"
	"LogtoUserProfile/globals"
	"LogtoUserProfile/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/logto-io/go/client"
	"github.com/logto-io/go/core"
	"log"
	"strconv"
)

func main() {
	var err error
	globals.Configuration, err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	globals.LogtoConfig = client.LogtoConfig{
		Endpoint:  globals.Configuration.Logto.Endpoint,
		AppId:     globals.Configuration.Logto.AppId,
		AppSecret: globals.Configuration.Logto.AppSecret,
		Scopes:    []string{core.UserScopeProfile, core.UserScopeCustomData, core.UserScopeEmail, core.ReservedScopeOpenId, "family_name", "familyName"},
	}

	store := memstore.NewStore([]byte(globals.Configuration.Server.Secret))
	if globals.Configuration.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	err = router.SetTrustedProxies(globals.Configuration.Server.TrustedProxies)
	if err != nil {
		log.Fatal(err)
	}

	router.LoadHTMLGlob("templates/*")
	router.Use(sessions.Sessions("logto-session", store))

	routes.RegisterRoutes(router)

	err = router.Run(globals.Configuration.Server.Host + ":" + strconv.Itoa(globals.Configuration.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
}
