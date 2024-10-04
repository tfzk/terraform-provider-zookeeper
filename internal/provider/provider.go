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
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperServer, nil),
				Description: "A comma separated list of 'host:port' pairs, pointing at ZooKeeper Server(s).",
			},
			"session_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperSessionSec, client.DefaultZooKeeperSessionSec),
				Description: "How many seconds a session is considered valid after losing connectivity. " +
					"More information about ZooKeeper sessions can be found [here](#zookeeper-sessions).",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperUsername, nil),
				Description: "Username for digest authentication. Can be set via `ZOOKEEPER_USERNAME` environment variable.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperPassword, nil),
				Description: "Password for digest authentication. Can be set via `ZOOKEEPER_PASSWORD` environment variable.",
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

func configureProviderContext(_ context.Context, rscData *schema.ResourceData) (interface{}, diag.Diagnostics) {
	servers := rscData.Get("servers").(string)
	sessionTimeout := rscData.Get("session_timeout").(int)
	username := rscData.Get("username").(string)
	password := rscData.Get("password").(string)

	if servers != "" {
		c, err := client.NewClient(servers, sessionTimeout, username, password)

		if err != nil {
			// Report inability to connect internal Client
			return nil, diag.Errorf("Unable creating ZooKeeper client against '%s': %v", servers, err)
		}

		return c, diag.Diagnostics{}
	}

	// Report missing mandatory arguments
	return nil, diag.Errorf("Provider requires at least the '%s' argument", "servers")
}
