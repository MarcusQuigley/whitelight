package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, e := strconv.ParseInt(params.ByName("id"), 10, 64)
	if e != nil || id < 1 {

		return 0, errors.New("invalid id param")
	}
	return id, nil
}
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	//js, e := json.Marshal(data)
	js, e := json.MarshalIndent(data, "", "\t")
	if e != nil {
		return e
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
