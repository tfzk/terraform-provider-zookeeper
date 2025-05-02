package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceZNode(t *testing.T) {
	parentPath := "/" + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { checkPreconditions(t) },
		ProviderFactories: providerFactoriesMap(),
		CheckDestroy:      confirmAllZNodeDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "zookeeper_znode" "parent" {
						path = "%s"
						data = "parent data"
					}
					resource "zookeeper_znode" "child" {
						path = "${zookeeper_znode.parent.path}/child"
						data = "child data"
					}`, parentPath,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Parent checks
					resource.TestCheckResourceAttr("zookeeper_znode.parent", "path", parentPath),
					resource.TestCheckResourceAttrPair(
						"zookeeper_znode.parent",
						"path",
						"zookeeper_znode.parent",
						"id",
					),
					resource.TestCheckResourceAttr("zookeeper_znode.parent", "data", "parent data"),
					resource.TestCheckResourceAttr(
						"zookeeper_znode.parent",
						"data_base64",
						"cGFyZW50IGRhdGE=",
					),
					// Child checks
					resource.TestCheckResourceAttr(
						"zookeeper_znode.child",
						"path",
						parentPath+"/child",
					),
					resource.TestCheckResourceAttrPair(
						"zookeeper_znode.child",
						"path",
						"zookeeper_znode.child",
						"id",
					),
					resource.TestCheckResourceAttr("zookeeper_znode.child", "data", "child data"),
					resource.TestCheckResourceAttr(
						"zookeeper_znode.child",
						"data_base64",
						"Y2hpbGQgZGF0YQ==",
					),
				),
			},
			{
				ResourceName:      "zookeeper_znode.parent",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "zookeeper_znode.child",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceZNode_MultipleSharedPath(t *testing.T) {
	sharedPath := "/" + acctest.RandString(
		5,
	) + "/" + acctest.RandString(
		5,
	) + "/" + acctest.RandString(
		5,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { checkPreconditions(t) },
		ProviderFactories: providerFactoriesMap(),
		CheckDestroy:      confirmAllZNodeDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "zookeeper_znode" "Z01" {
						path = "%[1]s/01"
					}
					resource "zookeeper_znode" "Z02" {
						path = "%[1]s/02"
					}
					resource "zookeeper_znode" "Z03" {
						path = "%[1]s/03"
					}
					resource "zookeeper_znode" "Z04" {
						path = "%[1]s/04"
					}
					resource "zookeeper_znode" "Z05" {
						path = "%[1]s/05"
					}
					`, sharedPath,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zookeeper_znode.Z01", "path", sharedPath+"/01"),
					resource.TestCheckResourceAttr("zookeeper_znode.Z02", "path", sharedPath+"/02"),
					resource.TestCheckResourceAttr("zookeeper_znode.Z03", "path", sharedPath+"/03"),
					resource.TestCheckResourceAttr("zookeeper_znode.Z04", "path", sharedPath+"/04"),
					resource.TestCheckResourceAttr("zookeeper_znode.Z05", "path", sharedPath+"/05"),
				),
			},
		},
	})
}

func TestAccResourceZNode_Base64(t *testing.T) {
	sharedPath := "/" + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { checkPreconditions(t) },
		ProviderFactories: providerFactoriesMap(),
		CheckDestroy:      confirmAllZNodeDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "zookeeper_znode" "src" {
						path = "%[1]s/src"
						data = "Forza Napoli!"
					}
					resource "zookeeper_znode" "dst" {
						path = "%[1]s/dst"
						data_base64 = zookeeper_znode.src.data_base64
					}`, sharedPath,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"zookeeper_znode.src",
						"data",
						"zookeeper_znode.dst",
						"data",
					),
					resource.TestCheckResourceAttr("zookeeper_znode.dst", "data", "Forza Napoli!"),
					resource.TestCheckResourceAttrPair(
						"zookeeper_znode.src",
						"data_base64",
						"zookeeper_znode.dst",
						"data_base64",
					),
					resource.TestCheckResourceAttr(
						"zookeeper_znode.dst",
						"data_base64",
						"Rm9yemEgTmFwb2xpIQ==",
					),
				),
			},
			{
				ResourceName:      "zookeeper_znode.src",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "zookeeper_znode.dst",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceZNode_DefaultACL(t *testing.T) {
	path := "/" + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { checkPreconditions(t) },
		ProviderFactories: providerFactoriesMap(),
		CheckDestroy:      confirmAllZNodeDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "zookeeper_znode" "test_default_acl" {
						path = "%s"
						data = "Default ACL Test"
					}`, path),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"zookeeper_znode.test_default_acl",
						"path",
						path,
					),
					resource.TestCheckResourceAttr(
						"zookeeper_znode.test_default_acl",
						"data",
						"Default ACL Test",
					),
					resource.TestCheckResourceAttr(
						"zookeeper_znode.test_default_acl",
						"acl.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"zookeeper_znode.test_default_acl",
						"acl.0.scheme",
						"world",
					),
					resource.TestCheckResourceAttr(
						"zookeeper_znode.test_default_acl",
						"acl.0.id",
						"anyone",
					),
					resource.TestCheckResourceAttr(
						"zookeeper_znode.test_default_acl",
						"acl.0.permissions",
						"31",
					),
				),
			},
		},
	})
}

func TestAccResourceZNode_WithACL(t *testing.T) {
	path := "/" + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { checkPreconditions(t) },
		ProviderFactories: providerFactoriesMap(),
		CheckDestroy:      confirmAllZNodeDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "zookeeper_znode" "test_acl" {
						path = "%s"
						data = "ACL Test"
						acl {
							scheme      = "world"
							id          = "anyone"
							permissions = 31
						}
					}`, path),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zookeeper_znode.test_acl", "path", path),
					resource.TestCheckResourceAttr("zookeeper_znode.test_acl", "data", "ACL Test"),
					resource.TestCheckResourceAttr("zookeeper_znode.test_acl", "acl.#", "1"),
					resource.TestCheckResourceAttr(
						"zookeeper_znode.test_acl",
						"acl.0.scheme",
						"world",
					),
					resource.TestCheckResourceAttr(
						"zookeeper_znode.test_acl",
						"acl.0.id",
						"anyone",
					),
					resource.TestCheckResourceAttr(
						"zookeeper_znode.test_acl",
						"acl.0.permissions",
						"31",
					),
				),
			},
		},
	})
}
