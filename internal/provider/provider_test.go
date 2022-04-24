package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	testifyAssert "github.com/stretchr/testify/assert"

	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

func TestProvider(t *testing.T) {
	assert := testifyAssert.New(t)

	provider, err := New()
	assert.NoError(err)

	assert.NoError(provider.InternalValidate())
}

// providerFactoriesMap associates to each Provider factory instance, a name.
//
// WARN: This is important as this will be the name the provider will be expected
// to have when executing the acceptance tests.
// Fail to match the provider expected name will mean that the underlying binary
// terraform, used during acceptance tests, will error complaining it can't find
// the provider and `terraform init` should be executed.
func providerFactoriesMap() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"zookeeper": func() (*schema.Provider, error) {
			return New()
		},
	}
}

// checkPreconditions should be used with the field `PreCheck` of resource.TestCase.
func checkPreconditions(t *testing.T) {
	if v := os.Getenv("ZOOKEEPER_SERVERS"); v == "" {
		t.Fatal("ZOOKEEPER_SERVERS must be set for acceptance tests")
	}
}

// getTestZKClient can be used during test to procure a client.Client
func getTestZKClient() *client.Client {
	zkClient, _ := client.NewClient(os.Getenv("ZOOKEEPER_SERVERS"), 10)
	return zkClient
}

// confirmAllZNodeDestroyed should be used with the field `CheckDestroy` of resource.TestCase.
func confirmAllZNodeDestroyed(s *terraform.State) error {
	zkClient := getTestZKClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zookeeper_znode" {
			continue
		}

		// Confirm ZNode has been destroyed
		if exists, _ := zkClient.Exists(rs.Primary.ID); exists {
			return fmt.Errorf("ZNode '%s' still exists", rs.Primary.ID)
		}
	}

	return nil
}
