package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	e := app.writeJSON(w, http.StatusOK, data, nil)

	if e != nil {
		app.logger.Error(e.Error())
		http.Error(w, "The server encountered a problem and could not process your request",
			http.StatusInternalServerError)
	}

}
