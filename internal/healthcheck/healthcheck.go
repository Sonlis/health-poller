package healthcheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetHealth(host, path string, expectedJSON map[string]string) (string, error) {
	var body []byte
	resp, err := http.Get(host + path)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode > 299 {
		return "", errors.New(fmt.Sprintf("Status code %d is not a succesfull status code.", resp.StatusCode))
	}

	// There's a possibility that we do not want to set an expected JSON to be returned by the service.
	// In that case, the service is considered healthy from the http response code.
	if expectedJSON == nil {
		return "", nil
	}

	var healthStatus map[string]string
	err = json.Unmarshal(body, &healthStatus)
	if err != nil {
		return "", err
	}

	return checkHealth(healthStatus, expectedJSON), nil
}

func checkHealth(healthStatus map[string]string, expectedBody map[string]string) string {
	for key, value := range healthStatus {
		if expectedValue, ok := expectedBody[key]; ok {
			if value == expectedValue {
				continue
			}
			return formatHealthError(key, value, expectedValue)
		}
		return formatKeyError(key)
	}
	return ""
}
