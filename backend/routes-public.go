package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleCodeSubmit(c *gin.Context) {
	tokenParam := c.Query("token")
	if tokenParam == "" {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: "token parameter is empty",
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: fmt.Sprintf("Token is: %s", tokenParam),
	})
}
