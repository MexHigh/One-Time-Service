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

	if details.UsesLeft <= 0 {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: "usage limit of token exceeded",
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

	if err := db.DecrementUseCountForToken(tokenParam); err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	ts := time.Now()
	tsString := ts.Local().Format("02.01.2006, 15:04:05 MST")

	msg := fmt.Sprintf("Token: `%s`\nExecuted macro: `%s`\nSubmitter IP: `%s`\nSubmission time: `%s`",
		tokenParam,
		details.MacroName,
		c.ClientIP(),
		tsString,
	)
	if details.Comment != nil {
		// add comment field if one was set or the token
		msg += fmt.Sprintf("\nComment: `%s`", *details.Comment)
	}

	if err := CreateHomeassistantNotification(
		"One Time Service Token submitted",
		msg,
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
