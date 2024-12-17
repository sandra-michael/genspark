package handlers

import (
	"Assignment1/models/mockmodels"
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHealthCheck(t *testing.T) {
	//name
	//input anything response writer request
	//output "I am working and healthy"

	tt := [...]struct {
		name             string
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "OK",
			expectedStatus:   http.StatusOK,
			expectedResponse: "I am working and healthy",
		},
	}

	ctrl := gomock.NewController(t)

	// NewMockService would give us the implementation of the
	// interface that we can set in handlers struct
	mockDb := mockmodels.NewMockService(ctrl)

	// Creating the handler with the mocked service and validator
	h := Handler{
		c:        mockDb,          // Passing the mocked service
		validate: validator.New(), // Initializing the validator for input validation
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Create a response recorder to capture the response
			// it is an implementation of ResponseWriter
			rec := httptest.NewRecorder()

			// constructing the request
			req := httptest.NewRequest(http.MethodGet, "/v1/tasks/health", nil)

			h.healthCheck(rec, req)

			// checking if expected output matches or not
			require.Equal(t, tc.expectedStatus, rec.Code)

			body := rec.Body.String()
			body = strings.TrimSpace(body)

			require.Equal(t, tc.expectedResponse, body)

		})
	}

}

func TestCreateTask(t *testing.T) {
	//name
	//input body which is of json models.NewTask
	//output id created

	// mockTask := models.Task{
	// 	ID:          1,
	// 	Name:        "Implement API Endpoint",
	// 	Description: "Create a new API endpoint to handle task creation requests.",
	// 	Status:      "NEW",
	// }

	tt := [...]struct {
		name             string
		body             []byte // Body to send to the request
		expectedStatus   int
		expectedResponse string
		MockStore        func(m *mockmodels.MockService)
	}{
		{
			name: "OK",
			body: []byte(`{
					"name": "Implement API Endpoint",
					"description": "Create a new API endpoint to handle task creation requests."
					}`),
			expectedStatus:   http.StatusOK,
			expectedResponse: "1",
			MockStore: func(m *mockmodels.MockService) {
				// setting the expectations for the mock call
				m.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(1, nil).Times(1)
			},
		},
		{
			name:             "Error unmarshal",
			body:             []byte("dsfew"),
			expectedStatus:   http.StatusExpectationFailed,
			expectedResponse: "Error while unmarshal",
			MockStore: func(m *mockmodels.MockService) {
				// setting the expectations for the mock call
				m.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(1, nil).Times(0)
			},
		},
		{
			name: "Error validate Desc",
			body: []byte(`{
				"name": "Impdsf ",
				"description": "C"
				}`),
			expectedStatus:   http.StatusExpectationFailed,
			expectedResponse: "Error while validationKey: 'NewTask.Description' Error:Field validation for 'Description' failed on the 'min' tag",
			MockStore: func(m *mockmodels.MockService) {
				// setting the expectations for the mock call
				m.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(1, nil).Times(0)
			},
		},
		{
			name: "Error validate Name and Desc",
			body: []byte(`{
				"name": "I",
				"description": "C"
				}`),
			expectedStatus:   http.StatusExpectationFailed,
			expectedResponse: "Error while validationKey: 'NewTask.Name' Error:Field validation for 'Name' failed on the 'min' tag +Key: 'NewTask.Description' Error:Field validation for 'Description' failed on the 'min' tag",
			MockStore: func(m *mockmodels.MockService) {
				// setting the expectations for the mock call
				m.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(1, nil).Times(0)
			},
		},
	}

	ctrl := gomock.NewController(t)

	// NewMockService would give us the implementation of the
	// interface that we can set in handlers struct
	mockDb := mockmodels.NewMockService(ctrl)

	// Creating the handler with the mocked service and validator
	h := Handler{
		c:        mockDb,          // Passing the mocked service
		validate: validator.New(), // Initializing the validator for input validation
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.MockStore(mockDb)
			// Create a response recorder to capture the response
			// it is an implementation of ResponseWriter
			rec := httptest.NewRecorder()

			// constructing the request
			req := httptest.NewRequest(http.MethodPost, "/v1/tasks/", bytes.NewReader(tc.body))

			h.createTask(rec, req)

			// checking if expected output matches or not
			require.Equal(t, tc.expectedStatus, rec.Code)

			body := rec.Body.String()
			body = strings.TrimSpace(body)

			require.Equal(t, tc.expectedResponse, body)

		})
	}

}
