package main

import (
	"fmt"
	"net/http"
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
	fmt.Fprintf(w, "show movie id: %d", id)
}
