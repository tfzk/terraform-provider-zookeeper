package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

// New creates a new ZooKeeper Provider.
func New() (*schema.Provider, error) {
	// Hold a Clients Pool, so connection to the same ZooKeeper can be shared between resources.
	// The client is safe for concurrent usage.
	clientPool := client.NewPool()

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"servers": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperServer, nil),
				Description: "A comma separated list of 'host:port' pairs, pointing at ZooKeeper Server(s). " +
					"Can be set via `ZOOKEEPER_SERVERS` environment variable.",
			},
			"session_timeout": {
				Type:      schema.TypeInt,
				Optional:  true,
				Sensitive: false,
				DefaultFunc: schema.EnvDefaultFunc(
					client.EnvZooKeeperSessionSec,
					client.DefaultZooKeeperSessionSec,
				),
				Description: "How many seconds a session is considered valid after losing connectivity. " +
					"More information about ZooKeeper sessions can be found [here](#zookeeper-sessions). " +
					"Can be set via `ZOOKEEPER_SESSION` environment variable.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperUsername, nil),
				Description: "Username for digest authentication. " +
					"Can be set via `ZOOKEEPER_USERNAME` environment variable.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperPassword, nil),
				Description: "Password for digest authentication. " +
					"Can be set via `ZOOKEEPER_PASSWORD` environment variable.",
			},
			"tls_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperTLSEnable, nil),
				Description: "Use secure TLS connection when connecting to the server(s). " +
					"Can be set via `ZOOKEEPER_TLS_ENABLE` environment variable.",
			},
			"tls_skip_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperTLSSkipVerify, nil),
				Description: "Skip verification of server's certificate chain and host name. " +
					"Can be set via `ZOOKEEPER_TLS_SKIP_VERIFY` environment variable.",
			},
			"tls_root_ca_cert_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperTLSRootCertPath, nil),
				Description: "File path to the root CA certificate to use when connecting to the server(s) using TLS. " +
					"Can be set via `ZOOKEEPER_TLS_ROOT_CA_CERT_PATH` environment variable.",
			},
			"tls_cert_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperTLSCertPath, nil),
				Description: "File path to the certificate to use when connecting to the server(s) using TLS. " +
					"Can be set via `ZOOKEEPER_TLS_CERT_PATH` environment variable.",
			},
			"tls_key_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc(client.EnvZooKeeperTLSCertPath, nil),
				Description: "File path to the key to use when connecting to the server(s) using TLS. " +
					"Can be set via `ZOOKEEPER_TLS_KEY_PATH` environment variable.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zookeeper_znode":            resourceZNode(),
			"zookeeper_sequential_znode": resourceSeqZNode(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"zookeeper_znode": datasourceZNode(),
		},
		ConfigureContextFunc: func(_ context.Context, rscData *schema.ResourceData) (interface{}, diag.Diagnostics) {
			// Retrieve the given configuration
			servers := rscData.Get("servers").(string)
			sessionTimeout := rscData.Get("session_timeout").(int)
			username := rscData.Get("username").(string)
			password := rscData.Get("password").(string)

			tlsConfig := &client.TLSConfig{
				Enable:       rscData.Get("tls_enable").(bool),
				SkipVerify:   rscData.Get("tls_skip_verify").(bool),
				RootCertPath: rscData.Get("tls_root_ca_cert_path").(string),
				CertPath:     rscData.Get("tls_cert_path").(string),
				KeyPath:      rscData.Get("tls_key_path").(string),
			}

			if servers != "" {
				// NOTE: Client Pool above is in a closure here
				// because we don't have a way to add fields to the Provider.
				c, err := clientPool.GetOrCreateClient(servers, sessionTimeout, username, password, tlsConfig)
				if err != nil {
					// Report inability to connect internal Client
					return nil, diag.Errorf(
						"Unable creating ZooKeeper client against '%s': %v",
						servers,
						err,
					)
				}

				return c, diag.Diagnostics{}
			}

			// Report missing mandatory arguments
			return nil, diag.Errorf("Provider requires at least the '%s' argument", "servers")
		},
	}, nil
}
