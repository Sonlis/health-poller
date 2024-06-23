package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

const (
	gotifyUrl = "http://localhost:8080"
)

type GotifyApp struct {
	Token string `json:"token"`
	Id    int    `json:"id"`
}

type GotifyMessages struct {
	Messages []struct {
		ID       int    `json:"id"`
		Appid    int    `json:"appid"`
		Message  string `json:"message"`
		Title    string `json:"title"`
		Priority int    `json:"priority"`
	} `json:"messages"`
}

func TestMain(m *testing.M) {
	appBody := []byte(`{
        "name": "testapp"
    }`)

	gotifyApp, err := createGotifyApp(appBody)
	if err != nil {
		panic(fmt.Sprintf("Error creating the gotify app: %v", err))
	}
	os.Setenv("GOTIFY_TOKEN", gotifyApp.Token)
	main()
	os.Exit(m.Run())
}

func TestCheckMessages(t *testing.T) {
	messages, err := getApplicationMessages()
	if err != nil {
		t.Errorf("Error getting the messages amount: %v", err)
	}
	if len(messages.Messages) != 1 {
		t.Errorf("Expected 1 messages, got %d", len(messages.Messages))
	}
}

func getApplicationMessages() (GotifyMessages, error) {
	var gotifyResponse GotifyMessages
	var c http.Client

	r, err := http.NewRequest("GET", gotifyUrl+"/application/1/message", nil)
	if err != nil {
		return gotifyResponse, err
	}
	r.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4=") // Base 64 value of admin:admin, the default credentials.

	resp, err := c.Do(r)
	if err != nil {
		return gotifyResponse, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return gotifyResponse, err
	}
	err = json.Unmarshal(body, &gotifyResponse)
	return gotifyResponse, err
}

func createGotifyApp(requestBody []byte) (GotifyApp, error) {
	var gotifyResponse GotifyApp
	var c http.Client
	r, err := http.NewRequest("POST", gotifyUrl+"/application", bytes.NewBuffer(requestBody))
	if err != nil {
		return gotifyResponse, err
	}
	r.Header.Add("Content-type", "application/json")
	r.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4=") // Base 64 value of admin:admin, the default credentials.

	resp, err := c.Do(r)
	if err != nil {
		return gotifyResponse, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return gotifyResponse, err
	}
	err = json.Unmarshal(body, &gotifyResponse)
	return gotifyResponse, err
}
