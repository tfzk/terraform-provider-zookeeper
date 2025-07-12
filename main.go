// Package main is the entry point for the Terraform provider.
package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/tfzk/terraform-provider-zookeeper/internal/provider"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate echo "\n*** tfplugindocs: generating documentation... ***"
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate
//go:generate echo "*** tfplugindocs: generated! ***"
//
// Validate the documentation generated (above) by `tfplugindocs`
//go:generate echo "\n*** tfplugindocs: validation Documentation ***"
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs validate
//go:generate echo "*** tfplugindocs: validated! ***"

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
