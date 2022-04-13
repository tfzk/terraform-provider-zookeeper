package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	testifyAssert "github.com/stretchr/testify/assert"

	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

func TestProvider(t *testing.T) {
	assert := testifyAssert.New(t)

	provider, err := New()
	assert.NoError(err)

	assert.NoError(provider.InternalValidate())
}

// testAccProviderFactories associates to each Provider factory instance, a name.
//
// WARN: This is important as this will be the name the provider will be expected
// to have when executing the acceptance tests.
// Fail to match the provider expected name will mean that the underlying binary
// terraform, used during acceptance tests, will error complaining it can't find
// the provider and `terraform init` should be executed.
func testAccProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"zookeeper": func() (*schema.Provider, error) {
			return New()
		},
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ZOOKEEPER_SERVERS"); v == "" {
		t.Fatal("ZOOKEEPER_SERVERS must be set for acceptance tests")
	}
}

func getTestAccClient() *client.Client {
	zkClient, _ := client.NewClient(os.Getenv("ZOOKEEPER_SERVERS"), 10)
	return zkClient
}
