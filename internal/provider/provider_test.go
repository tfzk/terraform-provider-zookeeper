package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	testifyAssert "github.com/stretchr/testify/assert"
	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
	"github.com/tfzk/terraform-provider-zookeeper/internal/provider"
)

func TestProvider(t *testing.T) {
	assert := testifyAssert.New(t)

	p, err := provider.New()
	assert.NoError(err)

	assert.NoError(p.InternalValidate())
}

//nolint:unparam
func providerFactoriesMap() map[string]func() (*schema.Provider, error) {
	// Instantiate the provider in advance...
	p, err := provider.New()
	if err != nil {
		panic(fmt.Errorf("failed to instantiate provider: %w", err))
	}

	// ... then return it within the factory method.
	// This avoids the tests creating a new provider (and new connections) every single time.
	return map[string]func() (*schema.Provider, error){
		"zookeeper": func() (*schema.Provider, error) {
			return p, nil
		},
	}
}

// checkPreconditions should be used with the field `PreCheck` of resource.TestCase.
func checkPreconditions(t *testing.T) {
	if v := os.Getenv(client.EnvZooKeeperServer); v == "" {
		t.Fatalf("Environment variable '%s' must be set for acceptance tests", client.EnvZooKeeperServer)
	}
}

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
