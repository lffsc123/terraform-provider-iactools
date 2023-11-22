package main

import (
	"context"
	"flag"
	"log"

	"terraform-provider-firewall/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like ")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/lffsc123/firewall",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
