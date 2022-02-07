package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"

	"github.com/temporalio/background-checks/temporal/dataconverter"
	"go.temporal.io/sdk/converter"
)

//go:embed iframe.go.html
var iframeHTML string
var iframeHTMLTemplate = template.Must(template.New("iframe").Parse(iframeHTML))

func main() {
	frontend := os.Args[1]

	if frontend == "" {
		log.Fatal("frontend URL argument is required")
	}

	http.HandleFunc("/js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		err := iframeHTMLTemplate.Execute(w, map[string]string{
			"BasePath": r.URL.Path[0 : len(r.URL.Path)-3],
			"Frontend": frontend,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	encoder := dataconverter.Encoder{}
	http.Handle("/", converter.NewPayloadEncoderHTTPHandler(&encoder))

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
