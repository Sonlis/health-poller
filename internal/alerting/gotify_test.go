package alerting

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAlertServiceUnealthy(t *testing.T) {
	server := CreateMockServer(t)
	err := AlertServiceUnealthy("bk", "fastfood is unealthy", "verySecureToken", server.URL)
	if err != nil {
		t.Errorf("Error testing sending alerts to gotify: %v", err)
	}

}

func TestFormatRequestBody(t *testing.T) {
	wanted := []byte(`{
        "title": "fastfood unealthy",
        "message": "But everybody knew that"
    }`,
	)
	got := formatRequestBody("fastfood", "But everybody knew that")
	if !bytes.Equal(wanted, got) {
		t.Errorf("Expected %v, got %v", string(wanted), string(got))
	}
}

func CreateMockServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/message" {
			t.Errorf("Expected to request '/message', got: %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Alerting Gotify is done through a POST request, the header setting the content type must be present.")
		}

		if r.Header.Get("Authorization") == "" {
			t.Errorf("An authorization header holding the gotify token is necessary.")
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
            {
                "id": 1,
                "appId": 1,
                "message": "fastfood is unealthy",
                "Title": "bk"
            }
            `))
		}
	}))
}
