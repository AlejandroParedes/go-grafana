package sentry

import (
	"log"
	"time"

	"go-grafana/internal/config"

	"github.com/getsentry/sentry-go"
)

// InitSentry initializes the Sentry client
func InitSentry(cfg *config.Config) {
	if cfg.Sentry.DSN == "" {
		log.Println("Sentry DSN not provided, skipping initialization")
		return
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.Sentry.DSN,
		// Set tracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production.
		TracesSampleRate: 0.1,
		EnableTracing:    true,
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)

	log.Println("Sentry initialized successfully")
}
