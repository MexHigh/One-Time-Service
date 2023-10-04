package main

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type TokenDetailsWithShareURL struct {
	TokenDetails `json:",inline"`
	ShareURL     string `json:"share_url"`
}

func handleGetServiceCallNames(c *gin.Context) {
	serviceCalls, err := db.GetServiceCallNames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: serviceCalls,
	})
}

func handleGetServiceCall(c *gin.Context) {
	serviceCallParam := c.Query("name")
	if serviceCallParam == "" {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: "name parameter is empty",
		})
		return
	}

	serviceCall, err := db.GetServiceCall(serviceCallParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: serviceCall,
	})
}

func handleCreateServiceCall(c *gin.Context) {
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

	if err := db.AddServiceCall(body.Name, sc); err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, GenericResponse{
		Response: "Created",
	})
}

func handleDeleteServiceCall(c *gin.Context) {
	serviceCallParam := c.Query("name")
	if serviceCallParam == "" {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: "name parameter is empty",
		})
		return
	}

	// first, delete all tokens using this service call
	tokenNames, err := db.GetTokensByServiceCallName(serviceCallParam)
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

	// then, delete the service call itself
	if err := db.DeleteServiceCall(serviceCallParam); err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: "Deleted",
	})
}

func handleGetTokenNames(c *gin.Context) {
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

func handleGetTokensWithDetails(c *gin.Context) {
	tokens, err := db.GetAllTokenDetails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	// append share URLs to all tokens
	tokensWithShareURL := make(map[string]*TokenDetailsWithShareURL)
	for token, details := range tokens {
		tokensWithShareURL[token] = &TokenDetailsWithShareURL{
			*details,
			baseTokenURL + token,
		}
	}

	c.JSON(http.StatusOK, GenericResponse{
		Response: tokensWithShareURL,
	})
}

func handleCreateToken(c *gin.Context) {
	var body struct {
		ServiceCallName string     `json:"service_call_name" binding:"required"`
		Expires         *time.Time `json:"expires,omitempty"`
		UsesMax         int        `json:"uses_max"`
		Comment         *string    `json:"comment,omitempty"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	// check if service call exists
	if _, err := db.GetServiceCall(body.ServiceCallName); err != nil {
		c.JSON(http.StatusNotFound, GenericResponse{
			Error: err.Error(),
		})
		return
	}

	now := time.Now()

	// set expiry time
	var expires *time.Time
	if body.Expires != nil {
		expires = body.Expires
	} else {
		expires = nil
	}

	// set max usages default
	var usagesMax int
	if body.UsesMax == 0 { // default (nothing specified)
		usagesMax = 1 // when nothing was specified (or 0 explicitly) set to one usage
	} else {
		usagesMax = body.UsesMax
	}

	sc := &TokenDetails{
		ServiceCallName: body.ServiceCallName,
		Created:         &now,
		Expires:         expires, // might be nil, but thats intended
		UsesMax:         usagesMax,
		UsesLeft:        usagesMax,
		Comment:         body.Comment, // might be nil, but thats intended
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
