package main

import (
	"LogtoUserProfile/config"
	"LogtoUserProfile/globals"
	"LogtoUserProfile/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func main() {
	var err error
	globals.Configuration, err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	store := memstore.NewStore([]byte(globals.Configuration.Server.Secret))
	router := gin.Default()
	router.Use(sessions.Sessions("logto-session", store))

	routes.RegisterRoutes(router)

	router.Run(globals.Configuration.Server.Host + ":" + strconv.Itoa(globals.Configuration.Server.Port))
}
