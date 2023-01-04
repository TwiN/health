package health_test

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TwiN/health"
)

func BenchmarkHealthHandler_ServeHTTP(b *testing.B) {
	h := health.Handler().WithJSON(true)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			request, _ := http.NewRequest("GET", "/health", http.NoBody)
			responseRecorder := httptest.NewRecorder()
			h.ServeHTTP(responseRecorder, request)
			if n := rand.Intn(100); n < 1 {
				health.SetStatus(health.Down)
			} else if n < 5 {
				health.SetStatus(health.Up)
			}
		}
	})
	b.ReportAllocs()
}
