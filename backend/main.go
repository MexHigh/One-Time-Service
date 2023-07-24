package main

import (
	"flag"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	dbPath *string = flag.String("db", "./db.json", "Path to database JSON file")

	db *DB
)

func main() {
	flag.Parse()

	db = NewDB(*dbPath)

	/// INTERNAL ROUTER ///
	internalRouter := gin.Default()

	// serve internal frontend
	internalRouter.Use(static.Serve("/", static.LocalFile("../frontend-internal", true)))

	// internal routes
	internalRouter.GET("/")

	go internalRouter.Run(":8099")

	/// PUBLIC ROUTER ///
	publicRouter := gin.Default()

	// serve public frontend
	publicRouter.Use(static.Serve("/", static.LocalFile("../frontend-public", true)))

	// public routes
	publicRouter.POST("/submit", handleCodeSubmit)

	go publicRouter.Run(":1337")

	// block
	select {}
}
