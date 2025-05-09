package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"whitelight.quigley.net/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}
	e := json.NewDecoder(r.Body).Decode(&input)
	if e != nil {
		app.errorResponse(w, r, http.StatusBadRequest, e.Error())
	}
	fmt.Fprintf(w, "%+v\n", input)

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, e := app.readIDParam(r)

	if e != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	// Encode the struct to JSON and send it as the HTTP response.
	e = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if e != nil {
		app.serverErrorResponse(w, r, e)
	}
}
