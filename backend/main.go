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
	corsAllowDebug       *bool   = flag.Bool("cors-allow-all", false, "Allows all CORS request (for testing only!)")
	mockOptionsJson      *bool   = flag.Bool("mock-options-json", false, "Does not read from /data/options.json (for testing only!)")

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
	if *corsAllowDebug {
		internalRouterApi.Use(corsAllowAll())
	}
	internalRouterApi.GET("/ping", handlePing)
	internalRouterApi.GET("/macros", handleGetMacros)
	internalRouterApi.GET("/macro/details", handleGetMacro)
	internalRouterApi.POST("/macro", handleCreateMacro)
	internalRouterApi.DELETE("/macro", handleDeleteMacro)
	internalRouterApi.GET("/tokens", handleGetTokens)
	internalRouterApi.GET("/token/details", getTokenDetails) // generic route implementation
	internalRouterApi.POST("/token", handleCreateToken)
	internalRouterApi.DELETE("/token", handleDeleteToken)
	internalRouterApi.GET("/token/share-url", handleGetShareUrl)

	go internalRouter.Run(":8099")

	/// PUBLIC ROUTER ///
	publicRouter := gin.Default()

	// serve public frontend
	publicRouter.Use(static.Serve("/", static.LocalFile(*publicFrontendPath, true)))

	// public routes
	publicRouterApi := publicRouter.Group("/api/public")
	if *corsAllowDebug {
		publicRouterApi.Use(corsAllowAll())
	}
	publicRouterApi.GET("/ping", handlePing)
	publicRouterApi.GET("/token/details", getTokenDetails) // generic route implementation
	publicRouterApi.POST("/token/submit", handleCodeSubmit)

	go publicRouter.Run(":1337")

	// block
	select {}
}
