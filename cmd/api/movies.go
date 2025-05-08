package main

import (
	"fmt"
	"net/http"
	"time"

	"whitelight.quigley.net/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create movei")
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
