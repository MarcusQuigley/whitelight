package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	app := application{}
	t.Run("gets healtcheck", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/v1/healthceck", nil)
		response := httptest.NewRecorder()
		handler := http.HandlerFunc(app.healthCheckHandler)

		handler.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("Content-Type"), "application/json")
		var healthDetails struct {
			Status      string `json:"status"`
			Environment string `json:"environment"`
			Version     string `json:"version"`
		}

		if err := json.NewDecoder(response.Body).Decode(&healthDetails); err != nil {
			t.Fatal(err)
		}

		expectedStatus := "available"
		if healthDetails.Status != "available" {
			t.Errorf("handler returned wrong status: got %v want %v",
				healthDetails.Status, expectedStatus)
		}
	})
}

func assertStatus(t testing.TB, got, want int) {
	if got != want {
		t.Errorf("wrong status code: got %d want %d", got, want)
	}
}

func assertContentType(t testing.TB, got, want string) {
	if got != want {
		t.Errorf("wrong content type: got %v want %v", got, want)
	}
}
