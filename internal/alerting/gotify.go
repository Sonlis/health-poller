package alerting

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func AlertServiceUnealthy(serviceName, message, gotifyToken, gotifyHost string) error {
	var client http.Client
	requestBody := formatRequestBody(serviceName, message)
	r, err := http.NewRequest("POST", "http://"+gotifyHost+"//message", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+gotifyToken)

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return err
}

func formatRequestBody(serviceName, message string) []byte {
	return []byte(fmt.Sprintf(`{
        "title": "%s unealthy",
        "message": "%s"
    }`, serviceName, message))
}
