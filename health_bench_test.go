package health

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkHealthHandler_ServeHTTP(b *testing.B) {
	h := Handler().WithJSON(true)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			request, _ := http.NewRequest("GET", "/health", http.NoBody)
			responseRecorder := httptest.NewRecorder()
			h.ServeHTTP(responseRecorder, request)
			if n := rand.Intn(100); n < 1 {
				SetStatus(Down)
			} else if n < 5 {
				SetStatus(Up)
			}
		}
	})
	b.ReportAllocs()
}
