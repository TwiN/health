package health

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_ServeHTTP(t *testing.T) {
	type Scenario struct {
		Name                 string
		useJSON              bool
		status               Status
		reason               string
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
			Name:                 "text-up-reason",
			useJSON:              false,
			status:               Up,
			reason:               "reason",
			expectedResponseBody: "UP: reason",
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
			Name:                 "text-down-reason",
			useJSON:              false,
			status:               Down,
			reason:               "reason",
			expectedResponseBody: "DOWN: reason",
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
			Name:                 "json-up-reason",
			useJSON:              true,
			status:               Up,
			reason:               "Error",
			expectedResponseBody: `{"status":"UP","reason":"Error"}`,
			expectedResponseCode: 200,
		},
		{
			Name:                 "json-down",
			useJSON:              true,
			status:               Down,
			expectedResponseBody: `{"status":"DOWN"}`,
			expectedResponseCode: 500,
		},
		{
			Name:                 "json-down-reason",
			useJSON:              true,
			status:               Down,
			reason:               "Error",
			expectedResponseBody: `{"status":"DOWN","reason":"Error"}`,
			expectedResponseCode: 500,
		},
		{
			Name:                 "json-down-reason-with-quotes",
			useJSON:              true,
			status:               Down,
			reason:               `error "with" quotes`,
			expectedResponseBody: `{"status":"DOWN","reason":"error \"with\" quotes"}`,
			expectedResponseCode: 500,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			handler := Handler().WithJSON(scenario.useJSON)
			SetStatus(scenario.status)
			SetReason(scenario.reason)

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

func TestSetStatus(t *testing.T) {
	SetStatus(Up)
	if GetStatus() != Up {
		t.Error("expected status to be 'Up', got", GetStatus())
	}
	SetStatus(Down)
	if GetStatus() != Down {
		t.Error("expected status to be 'Down', got", GetStatus())
	}
	SetStatus(Up)
}

func TestSetReason(t *testing.T) {
	SetReason("hello")
	if GetReason() != "hello" {
		t.Error("expected reason to be 'hello', got", GetReason())
	}
	SetReason("world")
	if GetReason() != "world" {
		t.Error("expected reason to be 'world', got", GetReason())
	}
	SetReason("")
	if GetReason() != "" {
		t.Error("expected reason to be '', got", GetReason())
	}
}

func TestSetStatusAndReason(t *testing.T) {
	SetStatusAndReason(Down, "for what")
	if GetStatus() != Down {
		t.Error("expected status to be 'Down', got", GetStatus())
	}
	if GetReason() != "for what" {
		t.Error("expected reason to be 'hello', got", GetReason())
	}
}

func TestSetStatusAndResetReason(t *testing.T) {
	SetStatusAndReason(Down, "for what")
	if GetStatus() != Down {
		t.Error("expected status to be 'Down', got", GetStatus())
	}
	if GetReason() != "for what" {
		t.Error("expected reason to be 'for what', got", GetReason())
	}
	SetStatusAndResetReason(Up)
	if GetStatus() != Up {
		t.Error("expected status to be 'Up', got", GetStatus())
	}
	if GetReason() != "" {
		t.Error("expected reason to be '', got", GetReason())
	}
}

func TestSetHealthy(t *testing.T) {
	SetStatusAndReason(Down, "for what")
	if GetStatus() != Down {
		t.Error("expected status to be 'Down', got", GetStatus())
	}
	if GetReason() != "for what" {
		t.Error("expected reason to be 'for what', got", GetReason())
	}
	SetHealthy()
	if GetStatus() != Up {
		t.Error("expected status to be 'Up', got", GetStatus())
	}
	if GetReason() != "" {
		t.Error("expected reason to be '', got", GetReason())
	}
}

func TestSetUnhealthy(t *testing.T) {
	SetStatusAndReason(Up, "")
	if GetStatus() != Up {
		t.Error("expected status to be '', got", GetStatus())
	}
	if GetReason() != "" {
		t.Error("expected reason to be '', got", GetReason())
	}
	SetUnhealthy("for what")
	if GetStatus() != Down {
		t.Error("expected status to be 'Down', got", GetStatus())
	}
	if GetReason() != "for what" {
		t.Error("expected reason to be 'for what', got", GetReason())
	}
}
