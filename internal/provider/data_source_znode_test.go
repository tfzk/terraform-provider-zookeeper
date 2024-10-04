package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceZNode(t *testing.T) {
	srcPath := "/" + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { checkPreconditions(t) },
		ProviderFactories: providerFactoriesMap(),
		CheckDestroy:      confirmAllZNodeDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "zookeeper_znode" "src" {
						path = "%s"
						data = "Forza Napoli!"
					}
					data "zookeeper_znode" "dst" {
						depends_on = [zookeeper_znode.src]
						path 	   = zookeeper_znode.src.path
					}`, srcPath,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "path", "zookeeper_znode.src", "path"),

					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "data", "zookeeper_znode.src", "data"),
					resource.TestCheckResourceAttr("data.zookeeper_znode.dst", "data", "Forza Napoli!"),

					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "data_base64", "zookeeper_znode.src", "data_base64"),
					resource.TestCheckResourceAttr("data.zookeeper_znode.dst", "data_base64", "Rm9yemEgTmFwb2xpIQ=="),

					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat", "zookeeper_znode.src", "stat"),

					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat.0.czxid", "zookeeper_znode.src", "stat.0.czxid"),
					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat.0.mzxid", "zookeeper_znode.src", "stat.0.mzxid"),
					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat.0.pzxid", "zookeeper_znode.src", "stat.0.pzxid"),

					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat.0.ctime", "zookeeper_znode.src", "stat.0.ctime"),
					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat.0.mtime", "zookeeper_znode.src", "stat.0.mtime"),

					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat.0.version", "zookeeper_znode.src", "stat.0.version"),
					resource.TestCheckResourceAttr("data.zookeeper_znode.dst", "stat.0.version", "0"),
					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat.0.cversion", "zookeeper_znode.src", "stat.0.cversion"),
					resource.TestCheckResourceAttr("data.zookeeper_znode.dst", "stat.0.cversion", "0"),
					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat.0.aversion", "zookeeper_znode.src", "stat.0.aversion"),
					resource.TestCheckResourceAttr("data.zookeeper_znode.dst", "stat.0.aversion", "0"),

					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat.0.ephemeral_owner", "zookeeper_znode.src", "stat.0.ephemeral_owner"),

					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat.0.data_length", "zookeeper_znode.src", "stat.0.data_length"),
					resource.TestCheckResourceAttr("data.zookeeper_znode.dst", "stat.0.data_length", "13"),

					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat.0.num_children", "zookeeper_znode.src", "stat.0.num_children"),
					resource.TestCheckResourceAttr("data.zookeeper_znode.dst", "stat.0.num_children", "0"),

					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "acl.0.scheme", "zookeeper_znode.src", "acl.0.scheme"),
					resource.TestCheckResourceAttr("data.zookeeper_znode.dst", "acl.0.scheme", "world"),
					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "acl.0.id", "zookeeper_znode.src", "acl.0.id"),
					resource.TestCheckResourceAttr("data.zookeeper_znode.dst", "acl.0.id", "anyone"),
					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "acl.0.permissions", "zookeeper_znode.src", "acl.0.permissions"),
					resource.TestCheckResourceAttr("data.zookeeper_znode.dst", "acl.0.permissions", "31"),
				),
			},
		},
	})
}
