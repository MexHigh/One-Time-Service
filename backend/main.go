package main

import (
	"flag"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	dbPath               *string = flag.String("db", "./db.json", "Path to database JSON file")
	internalFrontendPath *string = flag.String("internal-frontend-path", "../frontend-internal", "Base path to static internal frontend files")
	publicFrontendPath   *string = flag.String("public-frontend-path", "../frontend-public", "Base path to static public frontend files")
	hassApiUrl           *string = flag.String("hass-api-url", "http://supervisor/core/api", "Custom base URL for Hass API")

	db *DB
)

func main() {
	flag.Parse()
	db = NewDB(*dbPath)

	/// INTERNAL ROUTER ///
	internalRouter := gin.Default()

	// frontend route
	internalRouter.Use(static.Serve("/", static.LocalFile(*internalFrontendPath, true)))

	// api routes
	internalRouterApi := internalRouter.Group("/api/internal")
	internalRouterApi.GET("/ping", handlePing)
	internalRouterApi.GET("/token/details", getTokenDetails)
	// TODO

	go internalRouter.Run(":8099")

	/// PUBLIC ROUTER ///
	publicRouter := gin.Default()

	// serve public frontend
	publicRouter.Use(static.Serve("/", static.LocalFile(*publicFrontendPath, true)))

	// public routes
	publicRouterApi := publicRouter.Group("/api/public")
	publicRouterApi.GET("/ping", handlePing)
	publicRouterApi.GET("/token/details", getTokenDetails)
	publicRouterApi.POST("/token/submit", handleCodeSubmit)

	go publicRouter.Run(":1337")

	// block
	select {}
}
