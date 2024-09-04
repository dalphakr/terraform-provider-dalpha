package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Main function, calling the provider
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}
