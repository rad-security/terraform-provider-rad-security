package main

import (
	"flag"

	radsecurity "github.com/rad-security/terraform-provider-rad-security/internal/rad-security"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

var (
	Version string = "dev"
	Commit  string = ""
)

func main() {
	debugMode := flag.Bool("debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        *debugMode,
		ProviderAddr: "registry.terraform.io/rad-security/rad-security",
		ProviderFunc: radsecurity.New(Version),
	}

	plugin.Serve(opts)
}
