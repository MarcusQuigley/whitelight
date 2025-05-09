package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}
	e := app.writeJSON(w, http.StatusOK, env, nil)

	if e != nil {
		app.serverErrorResponse(w, r, e)
	}

}
