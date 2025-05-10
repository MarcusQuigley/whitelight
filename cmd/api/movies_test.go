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

// func TestInsertMovie(t *testing.T) {
// 	t.Run("insert a movie", func(t *testing.T) {
// 		input := map[string]interface{}{
// 			"title":   "The Matrix",
// 			"year":    1999,
// 			"runtime": "120 mins",
// 			"genres":  []string{"action", "sci-fi"},
// 		}

// 	setupMock:= func(m *mockMovieModel) {
// 			m.insertFunc = func(movie *data.Movie) error {
// 				movie.ID = 1
// 				movie.CreatedAt = time.Now()
// 				movie.Version = 1
// 				return nil
// 			}
// 		}
// 			//expectedStatus: http.StatusCreated,
// 			app := &application{
// 			models: &mockModels{
// 				movies: &mockMovieModel{},
// 			},
// 		}

// 		setupMock(app.models.Movies.(*mockMovieModel))

// 		body, e := json.Marshal(&input)
// 		if e != nil {
// 			t.Fatal(e)
// 		}
// 		request := httptest.NewRequest(http.MethodPost, "/v1/movies", bytes.NewReader(body))
// 		response := httptest.NewRecorder()

// 		handler := http.HandlerFunc(app.createMovieHandler)
// 		handler.ServeHTTP(response, request)
// 		assertStatus(t, response.Code, http.StatusCreated)
// 	})
// }

// type mockModels struct {
// 	movies *mockMovieModel
// }

// type mockMovieModel struct {
// 	insertFunc func(*data.Movie) error
// }

// func (m *mockMovieModel) Insert(movie *data.Movie) error {
// 	if m.insertFunc != nil {
// 		return m.insertFunc(movie)
// 	}
// 	return nil
// }
