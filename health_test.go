package health_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TwiN/health"
)

func TestHealthHandler_ServeHTTP(t *testing.T) {
	defer health.SetHealthy()
	type Scenario struct {
		Name                 string
		useJSON              bool
		status               health.Status
		reason               string
		expectedResponseBody string
		expectedResponseCode int
	}
	scenarios := []Scenario{
		{
			Name:                 "text-up",
			useJSON:              false,
			status:               health.Up,
			expectedResponseBody: "UP",
			expectedResponseCode: 200,
		},
		{
			Name:                 "text-up-reason",
			useJSON:              false,
			status:               health.Up,
			reason:               "reason",
			expectedResponseBody: "UP: reason",
			expectedResponseCode: 200,
		},
		{
			Name:                 "text-down",
			useJSON:              false,
			status:               health.Down,
			expectedResponseBody: "DOWN",
			expectedResponseCode: 500,
		},
		{
			Name:                 "text-down-reason",
			useJSON:              false,
			status:               health.Down,
			reason:               "reason",
			expectedResponseBody: "DOWN: reason",
			expectedResponseCode: 500,
		},
		{
			Name:                 "json-up",
			useJSON:              true,
			status:               health.Up,
			expectedResponseBody: `{"status":"UP"}`,
			expectedResponseCode: 200,
		},
		{
			Name:                 "json-up-reason",
			useJSON:              true,
			status:               health.Up,
			reason:               "Error",
			expectedResponseBody: `{"status":"UP","reason":"Error"}`,
			expectedResponseCode: 200,
		},
		{
			Name:                 "json-down",
			useJSON:              true,
			status:               health.Down,
			expectedResponseBody: `{"status":"DOWN"}`,
			expectedResponseCode: 500,
		},
		{
			Name:                 "json-down-reason",
			useJSON:              true,
			status:               health.Down,
			reason:               "Error",
			expectedResponseBody: `{"status":"DOWN","reason":"Error"}`,
			expectedResponseCode: 500,
		},
		{
			Name:                 "json-down-reason-with-quotes",
			useJSON:              true,
			status:               health.Down,
			reason:               `error "with" quotes`,
			expectedResponseBody: `{"status":"DOWN","reason":"error \"with\" quotes"}`,
			expectedResponseCode: 500,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			handler := health.Handler().WithJSON(scenario.useJSON)
			health.SetStatus(scenario.status)
			health.SetReason(scenario.reason)

			request, _ := http.NewRequest("GET", "/health", http.NoBody)
			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)
			if responseRecorder.Code != scenario.expectedResponseCode {
				t.Errorf("expected GET /health to return status code %d, got %d", scenario.expectedResponseCode, responseRecorder.Code)
			}
			body, _ := io.ReadAll(responseRecorder.Body)
			if string(body) != scenario.expectedResponseBody {
				t.Errorf("expected GET /health to return %s, got %s", scenario.expectedResponseBody, string(body))
			}
		})
	}
}

func TestHealthHandler_GetResponseStatusCodeAndBody(t *testing.T) {
	defer health.SetHealthy()
	handler := health.Handler().WithJSON(true)
	health.SetStatus(health.Up)

	statusCode, body := handler.GetResponseStatusCodeAndBody()
	if statusCode != 200 {
		t.Error("expected status code to be 200, got", statusCode)
	}
	if string(body) != `{"status":"UP"}` {
		t.Error("expected body to be {\"status\":\"UP\"}, got", string(body))
	}
}

func TestSetStatus(t *testing.T) {
	defer health.SetHealthy()
	health.SetStatus(health.Up)
	if health.GetStatus() != health.Up {
		t.Error("expected status to be 'Up', got", health.GetStatus())
	}
	health.SetStatus(health.Down)
	if health.GetStatus() != health.Down {
		t.Error("expected status to be 'Down', got", health.GetStatus())
	}
	health.SetStatus(health.Up)
}

func TestSetReason(t *testing.T) {
	defer health.SetHealthy()
	health.SetReason("hello")
	if health.GetReason() != "hello" {
		t.Error("expected reason to be 'hello', got", health.GetReason())
	}
	health.SetReason("world")
	if health.GetReason() != "world" {
		t.Error("expected reason to be 'world', got", health.GetReason())
	}
	health.SetReason("")
	if health.GetReason() != "" {
		t.Error("expected reason to be '', got", health.GetReason())
	}
}

func TestSetStatusAndReason(t *testing.T) {
	defer health.SetHealthy()
	health.SetStatusAndReason(health.Down, "for what")
	if health.GetStatus() != health.Down {
		t.Error("expected status to be 'Down', got", health.GetStatus())
	}
	if health.GetReason() != "for what" {
		t.Error("expected reason to be 'hello', got", health.GetReason())
	}
}

func TestSetStatusAndResetReason(t *testing.T) {
	defer health.SetHealthy()
	health.SetStatusAndReason(health.Down, "for what")
	if health.GetStatus() != health.Down {
		t.Error("expected status to be 'Down', got", health.GetStatus())
	}
	if health.GetReason() != "for what" {
		t.Error("expected reason to be 'for what', got", health.GetReason())
	}
	health.SetStatusAndResetReason(health.Up)
	if health.GetStatus() != health.Up {
		t.Error("expected status to be 'Up', got", health.GetStatus())
	}
	if health.GetReason() != "" {
		t.Error("expected reason to be '', got", health.GetReason())
	}
}

func TestSetHealthy(t *testing.T) {
	defer health.SetHealthy()
	health.SetStatusAndReason(health.Down, "for what")
	if health.GetStatus() != health.Down {
		t.Error("expected status to be 'Down', got", health.GetStatus())
	}
	if health.GetReason() != "for what" {
		t.Error("expected reason to be 'for what', got", health.GetReason())
	}
	health.SetHealthy()
	if health.GetStatus() != health.Up {
		t.Error("expected status to be 'Up', got", health.GetStatus())
	}
	if health.GetReason() != "" {
		t.Error("expected reason to be '', got", health.GetReason())
	}
}

func TestSetUnhealthy(t *testing.T) {
	defer health.SetHealthy()
	health.SetStatusAndReason(health.Up, "")
	if health.GetStatus() != health.Up {
		t.Error("expected status to be '', got", health.GetStatus())
	}
	if health.GetReason() != "" {
		t.Error("expected reason to be '', got", health.GetReason())
	}
	health.SetUnhealthy("for what")
	if health.GetStatus() != health.Down {
		t.Error("expected status to be 'Down', got", health.GetStatus())
	}
	if health.GetReason() != "for what" {
		t.Error("expected reason to be 'for what', got", health.GetReason())
	}
}
