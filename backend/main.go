package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	/// INTERNAL ROUTER ///
	internalRouter := gin.Default()

	// serve internal frontend
	internalRouter.Use(static.Serve("/", static.LocalFile("../frontend-internal", true)))

	// internal routes
	internalRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong (admin)",
		})
	})

	go internalRouter.Run(":8099")

	/// PUBLIC ROUTER ///
	publicRouter := gin.Default()

	// serve public frontend
	publicRouter.Use(static.Serve("/", static.LocalFile("../frontend-public", true)))

	// public routes
	publicRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	go publicRouter.Run(":1337")

	// block
	select {}
}
