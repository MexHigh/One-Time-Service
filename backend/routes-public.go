package main

import (
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

	if details.Expires != nil && details.Expires.After(time.Now()) {
		c.JSON(http.StatusUnauthorized, GenericResponse{
			Error: "token is expired",
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

	c.JSON(http.StatusOK, GenericResponse{
		Response: "Ok",
	})
}
