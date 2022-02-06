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
	dc := dataconverter.Encoder{}

	srv := &http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: converter.NewPayloadEncoderHTTPHandler(&dc, os.Args[1]),
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
