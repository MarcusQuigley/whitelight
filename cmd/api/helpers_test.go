package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestReadIDParam(t *testing.T) {
	t.Skip("skipping till i can fix mac debug issue")
	app := &application{}

	tests := []struct {
		name       string
		paramValue string
		wantID     int64
		wantErr    error
	}{
		{
			name:       "valid id",
			paramValue: "123",
			wantID:     123,
			wantErr:    nil,
		},
		// {
		// 	name:       "invalid id - non-numeric",
		// 	paramValue: "abc",
		// 	wantID:     0,
		// 	wantErr:    errors.New("invalid id param"),
		// },
		// {
		// 	name:       "invalid id - negative",
		// 	paramValue: "-5",
		// 	wantID:     0,
		// 	wantErr:    errors.New("invalid id param"),
		// },
		// {
		// 	name:       "invalid id - zero",
		// 	paramValue: "0",
		// 	wantID:     0,
		// 	wantErr:    errors.New("invalid id param"),
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)

			params := httprouter.Params{
				{Key: "id", Value: tc.paramValue},
			}
			//req = req.WithContext(httprouter.WithParams(req.Context(), params))
			req = addParamsToRequest(req, params)

			id, err := app.readIDParam(req)

			if id != tc.wantID {
				t.Errorf("got ID %d, want %d", id, tc.wantID)
			}

			if (err != nil && tc.wantErr == nil) || (err == nil && tc.wantErr != nil) || (err != nil && err.Error() != tc.wantErr.Error()) {
				t.Errorf("got error %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func addParamsToRequest(r *http.Request, params httprouter.Params) *http.Request {
	// this is a workaround to emulate what httprouter does internally
	const paramsKey = "ParamsKey" // shadow the unexported key
	type contextKey string
	ctx := context.WithValue(r.Context(), contextKey(paramsKey), params)
	return r.WithContext(ctx)
}
