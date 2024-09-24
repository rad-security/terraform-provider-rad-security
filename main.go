package main

import (
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	radsecurity "github.com/rad-security/terraform-provider-rad-security/internal/rad-security"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	debugMode := flag.Bool("debug", false, "set to true to run the provider in debug mode")
	flag.Parse()
	plugin.Serve(&plugin.ServeOpts{
		Debug:        *debugMode,
		ProviderAddr: "registry.terraform.io/rad-security/rad-security",
		ProviderFunc: func() *schema.Provider {
			return radsecurity.Provider()
		},
	})
}
