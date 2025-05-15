package main

import (
	"errors"
	"fmt"
	"net/http"

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
	e = app.models.Movies.Insert(movie)
	if e != nil {
		app.serverErrorResponse(w, r, e)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	e = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if e != nil {
		app.serverErrorResponse(w, r, e)
	}

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, e := app.readIDParam(r)

	if e != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie, e := app.models.Movies.Get(id)
	if e != nil {
		switch {
		case errors.Is(e, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, e)
		}
	}

	e = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if e != nil {
		app.serverErrorResponse(w, r, e)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, e := app.readIDParam(r)
	if e != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, e := app.models.Movies.Get(id)
	if e != nil {
		switch {
		case errors.Is(e, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, e)
		}
	}

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	e = app.readJSON(w, r, &input)
	if e != nil {
		app.badRequestResponse(w, r, e)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	e = app.models.Movies.Update(movie)
	if e != nil {
		app.serverErrorResponse(w, r, e)
		return
	}
	e = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if e != nil {
		app.serverErrorResponse(w, r, e)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, e := app.readIDParam(r)
	if e != nil {
		app.notFoundResponse(w, r)
		return
	}
	e = app.models.Movies.Delete(id)
	if e != nil {
		switch {
		case errors.Is(e, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, e)
		}
		return
	}
	e = app.writeJSON(w, http.StatusNoContent, nil, nil)
	if e != nil {
		app.serverErrorResponse(w, r, e)
	}
}
