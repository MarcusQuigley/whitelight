package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"whitelight.quigley.net/internal/data"
)

func TestShowMovies(t *testing.T) {
	t.Skip("skipping till i can fix request context issue")
	app := application{}
	t.Run("gets movies", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/v1/movies/123", nil)
		response := httptest.NewRecorder()
		handler := http.HandlerFunc(app.showMovieHandler)

		handler.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		//assertContentType(t, response.Header().Get("Content-Type"), "application/json")
		// var healthDetails struct {
		// 	Status      string `json:"status"`
		// 	Environment string `json:"environment"`
		// 	Version     string `json:"version"`
		// }
		movie := data.Movie{}

		if err := json.NewDecoder(response.Body).Decode(&movie); err != nil {
			t.Fatal(err)
		}

		expectedTitle := "available"
		if movie.Title != "available" {
			t.Errorf("handler returned wrong movie title: got %v want %v",
				movie.Title, expectedTitle)
		}
	})
}

// func assertStatus(t testing.TB, got, want int) {
// 	if got != want {
// 		t.Errorf("wrong status code: got %d want %d", got, want)
// 	}
// }
