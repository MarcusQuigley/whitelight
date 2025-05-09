package main

import (
	"fmt"
	"net/http"
	"time"

	"whitelight.quigley.net/internal/data"
	"whitelight.quigley.net/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	e := app.readJSON(w, r, &input) //json.NewDecoder(r.Body).Decode(&input)
	if e != nil {
		app.badRequestResponse(w, r, e)
		return
	}
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
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
