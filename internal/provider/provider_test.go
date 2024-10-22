package provider_test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	testifyAssert "github.com/stretchr/testify/assert"
	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
	"github.com/tfzk/terraform-provider-zookeeper/internal/provider"
)

func TestProvider(t *testing.T) {
	assert := testifyAssert.New(t)

	provider, err := provider.New()
	assert.NoError(err)

	assert.NoError(provider.InternalValidate())
}

//nolint:gochecknoglobals
var providersMapMutex sync.Mutex

//nolint:gochecknoglobals
var providersMap = make(map[string]*schema.Provider)

// providerFactoriesMap associates to each Provider factory instance, a name.
//
// WARN: This is important as this will be the name the provider will be expected
// to have when executing the acceptance tests.
// Fail to match the provider expected name will mean that the underlying binary
// terraform, used during acceptance tests, will error complaining it can't find
// the provider and `terraform init` should be executed.
func providerFactoriesMap() map[string]func() (*schema.Provider, error) {
	const providerName = "zookeeper"

	return map[string]func() (*schema.Provider, error){
		providerName: func() (*schema.Provider, error) {
			providersMapMutex.Lock()
			defer providersMapMutex.Unlock()

			// Return previously-initialized copy of the provider
			if zkProv, found := providersMap[providerName]; found {
				return zkProv, nil
			}

			// Create new copy of the provider and store it, before returning it
			zkProv, err := provider.New()
			if err != nil {
				return nil, fmt.Errorf("could not create provider %s: %w", providerName, err)
			}
			providersMap[providerName] = zkProv
			return zkProv, nil
		},
	}
}

// checkPreconditions should be used with the field `PreCheck` of resource.TestCase.
func checkPreconditions(t *testing.T) {
	if v := os.Getenv(client.EnvZooKeeperServer); v == "" {
		t.Fatalf("Environment variable '%s' must be set for acceptance tests", client.EnvZooKeeperServer)
	}
}

// getTestZKClient can be used during test to procure a client.Client.

// confirmAllZNodeDestroyed should be used with the field `CheckDestroy` of resource.TestCase.
func confirmAllZNodeDestroyed(s *terraform.State) error {
	fmt.Println("[DEBUG] Confirming all ZNodes have been removed")
	zkClient, err := client.NewClientFromEnv()
	if err != nil {
		return fmt.Errorf("failed to create new Client: %w", err)
	}
	defer zkClient.Close()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zookeeper_znode" && rs.Type != "zookeeper_sequential_znode" {
			continue
		}

		// Confirm ZNode has been destroyed
		if exists, _ := zkClient.Exists(rs.Primary.ID); exists {
			return fmt.Errorf("ZNode '%s' still exists", rs.Primary.ID)
		}
	}

	return nil
}
