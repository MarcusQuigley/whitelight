package main

import (
	"context"
	"encoding/json"
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

func TestWriteJSON(t *testing.T) {
	// Create a new application instance

	// Test cases
	tests := []struct {
		name            string
		status          int
		data            any
		headers         http.Header
		expectedStatus  int
		expectedBody    string
		expectedHeaders map[string]string
	}{
		{
			name:   "Simple JSON object",
			status: http.StatusOK,
			data: map[string]string{
				"message": "test",
			},
			headers: http.Header{
				"X-Custom-Header": []string{"test-value"},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\n\t\"message\": \"test\"\n}\n",
			expectedHeaders: map[string]string{
				"Content-Type":    "application/json",
				"X-Custom-Header": "test-value",
			},
		},
		{
			name:           "Empty object",
			status:         http.StatusCreated,
			data:           map[string]interface{}{},
			headers:        http.Header{},
			expectedStatus: http.StatusCreated,
			expectedBody:   "{}\n",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:   "Nested JSON object",
			status: http.StatusOK,
			data: map[string]interface{}{
				"user": map[string]string{
					"name": "John",
					"role": "admin",
				},
			},
			headers:        http.Header{},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\n\t\"user\": {\n\t\t\"name\": \"John\",\n\t\t\"role\": \"admin\"\n\t}\n}\n",
			expectedHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	app := &application{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new ResponseRecorder
			rr := httptest.NewRecorder()

			// Call writeJSON
			err := app.writeJSON(rr, tt.status, tt.data, tt.headers)
			if err != nil {
				t.Errorf("writeJSON() error = %v", err)
				return
			}

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("writeJSON() status = %v, want %v", rr.Code, tt.expectedStatus)
			}

			// Check response body
			if rr.Body.String() != tt.expectedBody {
				t.Errorf("writeJSON() body = %v, want %v", rr.Body.String(), tt.expectedBody)
			}

			// Check headers
			for key, expectedValue := range tt.expectedHeaders {
				if got := rr.Header().Get(key); got != expectedValue {
					t.Errorf("writeJSON() header %s = %v, want %v", key, got, expectedValue)
				}
			}

			// Verify the response is valid JSON
			var jsonCheck interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &jsonCheck); err != nil {
				t.Errorf("writeJSON() produced invalid JSON: %v", err)
			}
		})
	}
}
