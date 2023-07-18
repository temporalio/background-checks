package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/temporalio/background-checks/temporal/dataconverter"

	"go.temporal.io/sdk/converter"
)

var portFlag int
var uiFlag string

func init() {
	flag.IntVar(&portFlag, "port", 8081, "Port to listen on")
	flag.StringVar(&uiFlag, "ui", "", "Temporal UI URL. Enables CORS which is required for access from Temporal UI")
}

func newCORSHTTPHandler(web string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", web)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Namespace")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()

	if uiFlag == "" {
		log.Fatalf("Please specific the -ui flag to enable UI support.\n")
	}

	handler := converter.NewPayloadCodecHTTPHandler(&dataconverter.Codec{KeyID: os.Getenv("DATACONVERTER_ENCRYPTION_KEY_ID")})
	handler = newCORSHTTPHandler(uiFlag, handler)

	srv := &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(portFlag),
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
		_ = srv.Close()
	case err := <-errCh:
		log.Fatal(err)
	}
}
