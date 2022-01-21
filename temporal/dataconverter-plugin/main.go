package main

import (
	"github.com/hashicorp/go-plugin"
	"go.temporal.io/sdk/converter"
	cliplugin "go.temporal.io/server/tools/cli/plugin"

	"github.com/temporalio/background-checks/temporal/dataconverter"
)

func main() {
	var pluginMap = map[string]plugin.Plugin{
		cliplugin.DataConverterPluginType: &cliplugin.DataConverterPlugin{
			Impl: dataconverter.NewEncryptionDataConverter(
				converter.GetDefaultDataConverter(),
				dataconverter.DataConverterOptions{},
			),
		},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: cliplugin.PluginHandshakeConfig,
		Plugins:         pluginMap,
	})
}
