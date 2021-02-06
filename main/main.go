package main

import (
	"net/http"
	"time"

	"github.com/TwinProduction/health"
)

func main() {
	router := http.NewServeMux()
	router.Handle("/health", health.Handler().WithJSON(true))
	health.SetStatus(health.Up)
	health.SetStatus(health.Down)
	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	server.ListenAndServe()
}
