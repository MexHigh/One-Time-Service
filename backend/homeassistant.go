package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

type AddonOptions struct {
	BaseURL                 string `json:"public_token_base_url"`
	NotifyOnTokenSubmittion bool   `json:"notify_on_token_submission"`
	NotificationTarget      string `json:"notification_target"`
}

func getAddonOptions() (*AddonOptions, error) {
	if *mockOptionsJson {
		return &AddonOptions{
			BaseURL:                 "http://localhost:1337",
			NotifyOnTokenSubmittion: true,
			NotificationTarget:      "persistent_notification.create",
		}, nil
	}

	file, err := os.Open("/data/options.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var options AddonOptions
	if err := json.Unmarshal(bytes, &options); err != nil {
		return nil, err
	}

	return &options, nil
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

func CreateHomeassistantNotification(title, message string) error {
	addonOptions, err := getAddonOptions()
	if err != nil {
		return err
	}

	if !addonOptions.NotifyOnTokenSubmittion {
		return nil // do nothing
	}

	// Check if notification_target has the right format
	// TODO maybe check this at an init function?
	split := strings.Split(addonOptions.NotificationTarget, ".")
	if len(split) != 2 {
		return errors.New("service string does not include exactly one dot")
	}
	url := fmt.Sprintf("%s/services/%s/%s", *hassApiUrl, split[0], split[1])

	body := struct {
		Title   string `json:"title"`
		Message string `json:"message"`
	}{title, message}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	fmt.Println(url)
	fmt.Println(body)

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
