package temporal

import (
	"os"

	"github.com/temporalio/background-checks/temporal/dataconverter"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

func NewClient(options client.Options) (client.Client, error) {
	if options.HostPort == "" {
		options.HostPort = os.Getenv("TEMPORAL_GRPC_ENDPOINT")
	}

	options.DataConverter = dataconverter.NewEncryptionDataConverter(
		converter.GetDefaultDataConverter(),
		dataconverter.DataConverterOptions{},
	)
	options.ContextPropagators = []workflow.ContextPropagator{dataconverter.NewContextPropagator()}

	return client.NewClient(options)
}
