package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/tfzk/terraform-provider-zookeeper/internal/provider"
)

func main() {
	p, err := provider.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize provider: %v\n", err)
		os.Exit(1)
	}

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return p
		},
	})
}
