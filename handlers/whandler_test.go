package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleWeatherRequest(t *testing.T) {
	cases := []struct {
		name            string
		method          string
		payload         interface{}
		expectedStatus  int
		expectedMessage string
	}{
		//state acronized
		{
			name:            "Success Case",
			method:          "POST",
			payload:         map[string]string{"city": "Seattle", "state": "WA"},
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},

		//state spelled out uppercase
		{
			name:            "Success Case",
			method:          "POST",
			payload:         map[string]string{"city": "Seattle", "state": "Washington"},
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},
		
		//lowercase city and state
		{
			name:            "Success Case",
			method:          "POST",
			payload:         map[string]string{"city": "seattle", "state": "washington"},
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},

		//lowercase city and and uppercase state
		{
			name:            "Success Case",
			method:          "POST",
			payload:         map[string]string{"city": "seattle", "state": "Washington"},
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},

		//uppercase city and and lowercase state
		{
			name:            "Success Case",
			method:          "POST",
			payload:         map[string]string{"city": "Seattle", "state": "washington"},
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},

		//invalid city and abbeviated state
		{
			name:            "Invalid Method",
			method:          "GET",
			payload:         map[string]string{"city": "somewhere", "state": "WA"},
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedMessage: "Method not allowed",
		},

		// empty fields
		{
			name:            "Invalid Payload",
			method:          "POST",
			payload:         map[string]string{"city": "", "state": ""},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid input. No entries for city or state.",
		},

		// empty city
		{
			name:            "Invalid Payload",
			method:          "POST",
			payload:         map[string]string{"city": "", "state": "Washington"},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid input. No entries for city or state.",
		},
		
		// empty state
		{
			name:            "Invalid Payload",
			method:          "POST",
			payload:         map[string]string{"city": "somewhere", "state": ""},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid input. No entries for city or state.",
		},
		

		//invalid city and state inputs	
		{
			name:            "Invalid Payload",
			method:          "POST",
			payload:         map[string]string{"city": "Sadfsdg", "state": "Words"},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid input. Please enter a valid city and state.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tc.payload)
			req, err := http.NewRequest(tc.method, "/weather", bytes.NewReader(payloadBytes))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandleWeatherRequest)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("expected status code %v but got %v", tc.expectedStatus, status)
			}

			if tc.expectedMessage != "" {
				if rr.Body.String() != tc.expectedMessage {
					t.Errorf("expected message %v but got %v", tc.expectedMessage, rr.Body.String())
				}
			}
		})
	}
}
