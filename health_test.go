package health

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_ServeHTTP(t *testing.T) {
	type Scenario struct {
		Name                 string
		useJSON              bool
		status               Status
		expectedResponseBody string
		expectedResponseCode int
	}
	scenarios := []Scenario{
		{
			Name:                 "text-up",
			useJSON:              false,
			status:               Up,
			expectedResponseBody: "UP",
			expectedResponseCode: 200,
		},
		{
			Name:                 "text-down",
			useJSON:              false,
			status:               Down,
			expectedResponseBody: "DOWN",
			expectedResponseCode: 500,
		},
		{
			Name:                 "json-up",
			useJSON:              true,
			status:               Up,
			expectedResponseBody: `{"status":"UP"}`,
			expectedResponseCode: 200,
		},
		{
			Name:                 "json-down",
			useJSON:              true,
			status:               Down,
			expectedResponseBody: `{"status":"DOWN"}`,
			expectedResponseCode: 500,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			handler := Handler().WithJSON(scenario.useJSON)
			SetStatus(scenario.status)

			request, _ := http.NewRequest("GET", "/health", nil)
			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)
			if responseRecorder.Code != scenario.expectedResponseCode {
				t.Errorf("expected GET /health to return status code %d, got %d", scenario.expectedResponseCode, responseRecorder.Code)
			}
			body, _ := ioutil.ReadAll(responseRecorder.Body)
			if string(body) != scenario.expectedResponseBody {
				t.Errorf("expected GET /health to return %s, got %s", scenario.expectedResponseBody, string(body))
			}
		})
	}

}
