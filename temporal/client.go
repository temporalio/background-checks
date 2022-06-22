package temporal

import (
	"crypto/tls"
	"os"

	"github.com/temporalio/background-checks/temporal/dataconverter"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

var (
	// mTLS Certificates - copied from project /secrets to container /run/secrets during Docker build
	ClientCertFile       = os.Getenv("TEMPORAL_CLI_TLS_CERT")
	ClientCertPrivateKey = os.Getenv("TEMPORAL_CLI_TLS_KEY")

	// gRPC Endpoint for Temporal Cloud - Example: "accounting-production.f45a2.tmprl.cloud:7233"
	// gRPC Endpoint for Local Docker Compose: host.docker.internal:7233"
	gRPCEndpoint = os.Getenv("TEMPORAL_CLI_ADDRESS")

	// e.g. "accounting-production.f45a2.tmprl.cloud"
	ServerName = os.Getenv("TEMPORAL_CLI_TLS_SERVERNAME")

	// e.g. "accounting-production.f45a2"
	NamespaceID = os.Getenv("TEMPORAL_CLI_NAMESPACE")
)

// NewClient returns a new Temporal Client
func NewClient(options client.Options) (client.Client, error) {
	if options.HostPort == "" {
		options.HostPort = gRPCEndpoint
	}

	options.DataConverter = dataconverter.NewEncryptionDataConverter(
		converter.GetDefaultDataConverter(),
		dataconverter.DataConverterOptions{KeyID: os.Getenv("DATACONVERTER_ENCRYPTION_KEY_ID")},
	)

	if ServerName != "" {
		clientCert, err := tls.LoadX509KeyPair(ClientCertFile, ClientCertPrivateKey)
		if err != nil {
			return nil, err
		}

		if NamespaceID != "" {
			options.Namespace = NamespaceID
		} else {
			options.Namespace = "default"
		}

		options.ConnectionOptions = client.ConnectionOptions{
			TLS: &tls.Config{
				Certificates: []tls.Certificate{clientCert},
				ServerName:   ServerName,
			},
		}
	}

	return client.NewClient(options)
}

/*
	- TEMPORAL_CLI_ADDRESS=smitty.temporal-dev.tmprl.cloud:7233
	- TEMPORAL_CLI_NAMESPACE=smitty.temporal-dev
	- TEMPORAL_CLI_TLS_SERVERNAME=smitty.temporal-dev.tmprl.cloud
	- TEMPORAL_CLI_TLS_CERT=/run/secrets/mTLS_Cert
	- TEMPORAL_CLI_TLS_KEY=/run/secrets/mTLS_Private_Key
	- DATACONVERTER_ENCRYPTION_KEY_ID=secret
*/
