package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHealth(t *testing.T) {
	server := CreateMockServer(t)
	expectedJSON := map[string]string{
		"status": "ok",
	}
	health, err := GetHealth(server.URL, "/healthy", expectedJSON)
	if err != nil {
		t.Errorf("Error querying test health server: %v", err)
	}
	if health != "" {
		t.Errorf("Health should be an empty string, got %s", health)
	}

	expectedHealthReport := "status returned ew, while ok was expected"
	health, err = GetHealth(server.URL, "/unealthy", expectedJSON)
	if err != nil {
		t.Errorf("Error querying test health server: %v", err)
	}
	if health != expectedHealthReport {
		t.Errorf("The unealthy report should be %s, got %s", expectedHealthReport, health)
	}

	expectedJSON = map[string]string{
		"mah?": "sjdf",
	}

	expectedHealthReport = "status field was not expected in JSON response"
	health, err = GetHealth(server.URL, "/unealthy", expectedJSON)
	if err != nil {
		t.Errorf("Error querying test health server: %v", err)
	}
	if health != expectedHealthReport {
		t.Errorf("The unealthy report should be %s, got %s", expectedHealthReport, health)
	}
}

func TestCheckHealth(t *testing.T) {
	healthStatus := map[string]string{
		"body": "alive",
		"mind": "whoknows",
	}
	expectedStatus := map[string]string{
		"body": "alive",
		"mind": "whoknows",
	}

	checkedHealth := checkHealth(healthStatus, expectedStatus)
	if checkedHealth != "" {
		t.Errorf("The health status is the one expected, but the check reports it as unealthy: %s", checkedHealth)
	}

	expectedStatus = map[string]string{
		"body": "alive",
		"mind": "optimistic",
	}

	checkedHealth = checkHealth(healthStatus, expectedStatus)
	want := "mind returned whoknows, while optimistic was expected"
	if checkedHealth != want {
		t.Errorf("Expected health status report: %s, got %s", want, checkedHealth)
	}

	expectedStatus = map[string]string{
		"body": "alive",
		"soul": "inspired",
	}
	checkedHealth = checkHealth(healthStatus, expectedStatus)
	want = "mind field was not expected in JSON response"
	if checkedHealth != want {
		t.Errorf("Expected health status report: %s, got %s", want, checkedHealth)
	}
}

func CreateMockServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthy" && r.URL.Path != "/unealthy" {
			t.Errorf("Expected to request '/healthy' or '/unealthy', got: %s", r.URL.Path)
		} else {
			if r.URL.Path == "/healthy" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
            {
                "status": "ok"
            }
            `))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
            {
                "status": "ew"
            }
            `))
			}
		}
	}))
}
