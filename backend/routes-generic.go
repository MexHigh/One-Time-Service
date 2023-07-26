package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, GenericResponse{
		Response: "pong",
	})
}

func getTokenDetails(c *gin.Context) {
	tokenParam := c.Query("token")
	if tokenParam == "" {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: "token parameter is empty",
		})
		return
	}

	details, err := db.GetTokenDetails(tokenParam)
	if err != nil {
		c.JSON(http.StatusNotFound, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: details,
	})
}
