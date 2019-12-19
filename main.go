package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-po/prometheus_operator"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: prometheus_operator.Provider})
}
