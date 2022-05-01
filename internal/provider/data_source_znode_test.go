package provider

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
						data = "source znode data"
					}
					data "zookeeper_znode" "dst" {
						path = zookeeper_znode.src.path

						# It's necessary to add 'depends_on' on data sources in Terraform 0.13
						# because of this issue: https://github.com/hashicorp/terraform/issues/25961
						depends_on = [
							zookeeper_znode.src,
						]
					}`, srcPath,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "path", "zookeeper_znode.src", "path"),
					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "data", "zookeeper_znode.src", "data"),
					resource.TestCheckResourceAttrPair("data.zookeeper_znode.dst", "stat", "zookeeper_znode.src", "stat"),
				),
			},
		},
	})
}
