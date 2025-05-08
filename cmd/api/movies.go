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
		http.NotFound(w, r)
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
	e = app.writeJSON(w, http.StatusOK, movie, nil)
	if e != nil {
		app.logger.Error(e.Error())
		http.Error(w, "The server encountered a problem and could not process your request",
			http.StatusInternalServerError)
	}
}
