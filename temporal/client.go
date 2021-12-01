package temporal

import (
	"os"

	"go.temporal.io/sdk/client"
)

func NewClient(options client.Options) (client.Client, error) {
	if options.HostPort == "" {
		options.HostPort = os.Getenv("TEMPORAL_GRPC_ENDPOINT")
	}

	return client.NewClient(options)
}
