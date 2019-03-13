package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/rayene/terraform-provider-statuspageio/statuspageio"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: statuspageio.Provider,
	})
}
