package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

func New() (*schema.Provider, error) {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"servers": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc("ZOOKEEPER_SERVERS", nil),
				Description: "A string containing a comma separated list of 'host:port' pairs",
			},
			"session_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc("ZOOKEEPER_SESSION", 10),
				Description: "How many seconds a session is considered valid after losing connectivity",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zookeeper_znode":            resourceZNode(),
			"zookeeper_sequential_znode": resourceSeqZNode(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"zookeeper_znode": datasourceZNode(),
		},
		ConfigureContextFunc: configureProviderContext,
	}, nil
}

func configureProviderContext(ctx context.Context, rscData *schema.ResourceData) (interface{}, diag.Diagnostics) {
	servers := rscData.Get("servers").(string)
	sessionTimeout := rscData.Get("session_timeout").(int)

	if servers != "" {
		c, err := client.NewClient(servers, sessionTimeout)

		if err != nil {
			// Report inability to connect internal Client
			return nil, diag.Errorf("Unable creating ZooKeeper client against '%s': %v", servers, err)
		}

		return c, diag.Diagnostics{}
	}

	// Report missing mandatory arguments
	return nil, diag.Errorf("Provider requires at least the '%s' argument", "servers")
}
