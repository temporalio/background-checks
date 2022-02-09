package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/temporalio/background-checks/temporal/dataconverter"
	"go.temporal.io/sdk/converter"
)

func main() {
	frontend := os.Args[1]

	if frontend == "" {
		log.Fatal("frontend URL argument is required")
	}

	encoder := dataconverter.Encoder{KeyID: os.Getenv("DATACONVERTER_ENCRYPTION_KEY_ID")}
	handler := converter.NewPayloadEncoderHTTPHandler(&encoder)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontend)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		handler.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Addr: "0.0.0.0:8081",
	}

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
		srv.Close()
	case err := <-errCh:
		log.Fatalf("error: %v", err)
	}
}
