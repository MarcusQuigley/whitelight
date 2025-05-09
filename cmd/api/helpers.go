package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	e := dec.Decode(dst)
	if e != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {

		case errors.As(e, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)",
				syntaxError.Offset)

		case errors.Is(e, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(e, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)",
				unmarshalTypeError.Offset)

		case errors.Is(e, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(e, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes",
				maxBytesError.Limit)

		case errors.As(e, &invalidUnmarshalError):
			panic(e)
		default:
			return e
		}
	}
	e = dec.Decode(&struct{}{})
	if !errors.Is(e, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
