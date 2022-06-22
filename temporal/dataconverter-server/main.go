package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/temporalio/background-checks/temporal/dataconverter"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/log/tag"
)

func newClaimMapper(providerURL string, logger log.Logger) authorization.ClaimMapper {
	authConfig := config.Authorization{
		JWTKeyProvider: config.JWTKeyProvider{
			KeySourceURIs: []string{providerURL},
		},
		ClaimMapper: "default",
	}

	provider := authorization.NewDefaultTokenKeyProvider(
		&authConfig,
		logger,
	)

	return authorization.NewDefaultJWTClaimMapper(provider, &authConfig, logger)
}

func newOauthPayloadEncoderHTTPHandler(frontendURL string, providerURL string, encoders map[string]dataconverter.Encoder, logger log.Logger) (http.HandlerFunc, error) {
	if frontendURL == "" {
		return nil, fmt.Errorf("frontend URL is required")
	}
	if providerURL == "" {
		return nil, fmt.Errorf("oauth provider URL is required")
	}
	if len(encoders) == 0 {
		return nil, fmt.Errorf("a namespace to encoder mapping is required")
	}

	handlers := make(map[string]http.Handler, len(encoders))
	for namespace, encoder := range encoders {
		handlers[namespace] = converter.NewPayloadEncoderHTTPHandler(&encoder)
	}

	mapper := newClaimMapper(providerURL, logger)

	oauthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Namespace")

		if r.Method == "OPTIONS" {
			return
		}

		namespace := r.Header.Get("X-Namespace")
		if namespace == "" {
			http.Error(w, "X-Namespace header must be set", http.StatusBadRequest)
			return
		}

		handler, ok := handlers[namespace]
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		authInfo := authorization.AuthInfo{
			AuthToken: r.Header.Get("Authorization"),
			Audience:  frontendURL,
		}

		claims, err := mapper.GetClaims(&authInfo)
		if err != nil {
			logger.Warn("unable to parse claims")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// If they have no role in this namespace they will get RoleUndefined
		role := claims.Namespaces[namespace]

		authorized := false

		switch {
		case strings.HasSuffix(r.URL.Path, "/decode"):
			if role >= authorization.RoleReader {
				authorized = true
			}
		case strings.HasSuffix(r.URL.Path, "/encode"):
			if role >= authorization.RoleWriter {
				authorized = true
			}
		}

		if authorized {
			handler.ServeHTTP(w, r)
			return
		}

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}

	return oauthHandler, nil
}

func main() {
	logger := log.NewCLILogger()

	if len(os.Args) < 2 {
		logger.Fatal("frontend URL and provider URL arguments are required")
	}

	frontend := os.Args[1]
	provider := os.Args[2]

	// When supporting multiple namespaces, add them here.
	// Per-namespace KeyID and so on can be set at this level.
	// Only handle encoding for the default namespace in this example.
	encoders := map[string]dataconverter.Encoder{
		"default": {KeyID: os.Getenv("DATACONVERTER_ENCRYPTION_KEY_ID")},
	}

	oauthHandler, err := newOauthPayloadEncoderHTTPHandler(frontend, provider, encoders, logger)
	if err != nil {
		logger.Fatal("unable to create oauth handler", tag.Error(err))
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: oauthHandler,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
		srv.Close()
	case err := <-errCh:
		logger.Fatal("error", tag.NewErrorTag(err))
	}
}
