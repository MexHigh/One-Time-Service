package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var ( // flags
	dbPath               *string = flag.String("db", "./db.json", "Path to database JSON file")
	internalFrontendPath *string = flag.String("internal-frontend-path", "../frontend-internal", "Base path to static internal frontend files")
	publicFrontendPath   *string = flag.String("public-frontend-path", "../frontend-public", "Base path to static public frontend files")
	hassApiUrl           *string = flag.String("hass-api-url", "http://supervisor/core/api", "Custom base URL for Hass API")
	corsAllowDebug       *bool   = flag.Bool("cors-allow-all", false, "Allows all CORS request (for testing only!)")
	mockOptionsJson      *bool   = flag.Bool("mock-options-json", false, "Does not read from /data/options.json (for testing only!)")
)

var ( // global vars
	db           *DB
	baseTokenURL string // e.g. https://example.org/?token=
)

func main() {
	flag.Parse()
	if err := initBaseTokenURL(); err != nil {
		panic(err)
	}
	db = NewDB(*dbPath)

	if *corsAllowDebug || *mockOptionsJson {
		log.Println("[WARNING] One of -cors-allow-debug or -mock-options-json was set! This prevents Gin from using Release mode!")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	/// INTERNAL ROUTER ///
	internalRouter := gin.Default()
	internalRouter.SetTrustedProxies([]string{"172.30.32.2"}) // ingress IP

	// frontend route
	internalRouter.Use(static.Serve("/", static.LocalFile(*internalFrontendPath, true)))

	// api routes
	internalRouterApi := internalRouter.Group("/api/internal")
	if *corsAllowDebug {
		internalRouterApi.Use(corsAllowAll())
	}
	internalRouterApi.GET("/ping", handlePing)
	internalRouterApi.GET("/service-calls", handleGetServiceCallNames)
	internalRouterApi.GET("/service-call/details", handleGetServiceCall)
	internalRouterApi.POST("/service-call", handleCreateServiceCall)
	internalRouterApi.DELETE("/service-call", handleDeleteServiceCall)
	internalRouterApi.GET("/tokens", handleGetTokenNames)
	internalRouterApi.GET("/tokens/details", handleGetTokensWithDetails)
	internalRouterApi.GET("/token/details", getTokenDetails) // generic route implementation
	internalRouterApi.POST("/token", handleCreateToken)
	internalRouterApi.DELETE("/token", handleDeleteToken)

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

func initBaseTokenURL() error {
	options, err := getAddonOptions()
	if err != nil {
		return err
	}

	urlNoQuery, err := url.JoinPath(options.BaseURL, "/")
	if err != nil {
		return err
	}

	baseTokenURL = urlNoQuery + "?token="
	return nil
}
