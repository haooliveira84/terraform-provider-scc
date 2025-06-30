//go:generate go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
//go:generate tfplugindocs generate --provider-name "scc" --rendered-provider-name "SAP Cloud Connector"
package main

import (
	"context"
	"flag"
	"log"

	"github.com/SAP/terraform-provider-scc/scc/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address:         "registry.terraform.io/sap/scc",
		Debug:           debug,
		ProtocolVersion: 6,
	}

	err := providerserver.Serve(context.Background(), provider.New, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
