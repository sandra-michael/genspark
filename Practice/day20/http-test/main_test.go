package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoubleHandler(t *testing.T) {
	// Name of the test
	// Input
	// output
	tt := [...]struct {
		name           string // Name of the test
		queryParam     string // Query parameter to pass
		expectedStatus int    // Expected HTTP status code
		expectedBody   string // Expected response body
	}{
		{
			name:           "OK",
			queryParam:     "10",
			expectedStatus: http.StatusOK,
			expectedBody:   "20",
		},
		{
			name:           "Missing Value",
			queryParam:     "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing value",
		},
		{
			name:           "Not a number ",
			queryParam:     "val",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "not a number: " + "val",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			// Create a response recorder to capture the response
			// it is an implementation of ResponseWriter
			rec := httptest.NewRecorder()

			// constructing the request
			req := httptest.NewRequest(http.MethodGet, "/double?v="+tc.queryParam, nil)
			// no error should happen while constructing the request
			//require.NoError(t, err)

			// calling the actual handler function
			doubleHandler(rec, req)

			// checking if expected output matches or not
			require.Equal(t, tc.expectedStatus, rec.Code)

			body := rec.Body.String()
			body = strings.TrimSpace(body)

			require.Equal(t, tc.expectedBody, body)

		})

	}
}
