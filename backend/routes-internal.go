package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func handleGetMacros(c *gin.Context) {
	macros, err := db.GetMacroNames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: macros,
	})
}

func handleGetMacro(c *gin.Context) {
	tokenParam := c.Query("name")
	if tokenParam == "" {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: "name parameter is empty",
		})
		return
	}

	macro, err := db.GetMacro(tokenParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: macro,
	})
}

func handleCreateMacro(c *gin.Context) {
	var body struct {
		Name                  string `json:"name" binding:"required"`
		ServicePayloadYAMLb64 string `json:"service_payload_yaml_base64" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	yamlBytes, err := base64.StdEncoding.DecodeString(body.ServicePayloadYAMLb64)
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	var scCompat struct {
		Service string `json:"service"`
		// There are two supported ways to define an entity ID in YAML
		Target *map[string]interface{} `json:"target"` // (Prio 1) target -> entity_id [-> list element]
		Data   *map[string]interface{} `json:"data"`   // (Prio 2) data -> entity_id
		// we will move all entries in target into data, as the REST API does not support target contents
	}
	if err := yaml.Unmarshal(yamlBytes, &scCompat); err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	if scCompat.Service == "" {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: "service key cannot be empty",
		})
		return
	}

	if len(strings.Split(scCompat.Service, ".")) != 2 {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: "service key must follow home assistant service syntax",
		})
		return
	}

	if scCompat.Target != nil {
		var tempNewData map[string]interface{}
		if scCompat.Data == nil {
			tempNewData = make(map[string]interface{})
		} else {
			tempNewData = *scCompat.Data
		}
		for key, value := range *scCompat.Target {
			tempNewData[key] = value
		}
		scCompat.Data = &tempNewData
	}

	// TODO make comments

	var sd interface{}
	if scCompat.Data == nil {
		sd = make(map[string]interface{})
	} else {
		sd = *scCompat.Data
	}

	sc := &ServiceCall{
		Service: scCompat.Service,
		Data:    sd,
	}

	if err := db.AddMacro(body.Name, sc); err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, GenericResponse{
		Response: "Created",
	})
}

func handleDeleteMacro(c *gin.Context) {
	macroParam := c.Query("name")
	if macroParam == "" {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: "name parameter is empty",
		})
		return
	}

	// first, delete all tokens using this macro
	tokenNames, err := db.GetTokensByMacroName(macroParam)
	if err != nil {
		c.JSON(http.StatusNotFound, GenericResponse{
			Error: err.Error(),
		})
		return
	}
	for _, tokenName := range tokenNames {
		if err := db.DeleteToken(tokenName); err != nil {
			c.JSON(http.StatusInternalServerError, GenericResponse{
				Error: err.Error(),
			})
			return
		}
	}

	// then, delete the macro itself
	if err := db.DeleteMacro(macroParam); err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: "Deleted",
	})
}

func handleGetTokens(c *gin.Context) {
	tokens, err := db.GetTokenNames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: tokens,
	})
}

func handleCreateToken(c *gin.Context) {
	var body struct {
		MacroName string     `json:"macro_name" binding:"required"`
		Expires   *time.Time `json:"expires,omitempty"`
		Comment   *string    `json:"comment,omitempty"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	// check if macro exists
	if _, err := db.GetMacro(body.MacroName); err != nil {
		c.JSON(http.StatusNotFound, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	now := time.Now()

	var expires *time.Time
	if body.Expires != nil {
		expires = body.Expires
	} else {
		expires = nil
	}

	sc := &TokenDetails{
		MacroName: body.MacroName,
		Created:   &now,
		Expires:   expires,      // might be nil, but thats intended
		Comment:   body.Comment, // might be nil, but thats intended
	}

	randString := generateRandomString(32)

	if err := db.AddToken(randString, sc); err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, GenericResponse{
		Response: "Created token " + randString,
	})
}

func handleDeleteToken(c *gin.Context) {
	tokenParam := c.Query("token")
	if tokenParam == "" {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: "token parameter is empty",
		})
		return
	}

	if err := db.DeleteToken(tokenParam); err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: "Deleted",
	})
}

func handleGetShareUrl(c *gin.Context) {
	tokenParam := c.Query("token")
	if tokenParam == "" {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: "token parameter is empty",
		})
		return
	}

	_, err := db.GetTokenDetails(tokenParam)
	if err != nil {
		c.JSON(http.StatusNotFound, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	if *mockOptionsJson {
		c.JSON(http.StatusOK, GenericResponse{
			Response: "http://test-url/" + "?token=" + tokenParam,
		})
		return
	}

	file, err := os.Open("/data/options.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	var options struct {
		BaseURL string `json:"public_token_base_url"`
	}
	if err := json.Unmarshal(bytes, &options); err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: options.BaseURL + "?token=" + tokenParam,
	})
}
