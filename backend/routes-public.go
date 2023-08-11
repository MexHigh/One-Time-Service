package main

import (
	"fmt"
	"net/http"
	"time"

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

	details, err := db.GetTokenDetails(tokenParam)
	if err != nil {
		c.JSON(http.StatusNotFound, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	if details.Expires != nil && time.Now().After(*details.Expires) {
		c.JSON(http.StatusUnauthorized, GenericResponse{
			Error: "token has expired",
		})
		return
	}

	sc, err := db.GetMacro(details.MacroName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	if err := CallHomeAssistantService(sc); err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	// delete ("invalidate") token on success
	if err := db.DeleteToken(tokenParam); err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	if err := CreateHomeassistantNotification(
		"One Time Service Token submitted",
		fmt.Sprintf("Token: `%s`\nExecuted macro: `%s`\nSubmitter IP: `%s`",
			tokenParam,
			details.MacroName,
			c.RemoteIP(),
		),
	); err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: "Ok",
	})
}
