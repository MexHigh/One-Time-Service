package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var hassApiToken string

func init() {
	hassApiToken = os.Getenv("SUPERVISOR_TOKEN")
	if hassApiToken == "" {
		panic("SUPERVISOR_TOKEN is empty or does not exist")
	}
}

func CallHomeAssistantService(sc *ServiceCall) error {
	split := strings.Split(sc.Service, ".")
	if len(split) != 2 {
		return errors.New("service string does not include exactly one dot")
	}

	url := fmt.Sprintf("%s/services/%s/%s", *hassApiUrl, split[0], split[1])

	jsonData, err := json.Marshal(sc.Data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+hassApiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		return nil
	} else {
		return errors.New("response status code is not 200, got " + strconv.Itoa(resp.StatusCode))
	}
}
