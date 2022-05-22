package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSeqZNode_FromDir(t *testing.T) {
	seqFromDir := "/" + acctest.RandString(10) + "/"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { checkPreconditions(t) },
		ProviderFactories: providerFactoriesMap(),
		CheckDestroy:      confirmAllZNodeDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "zookeeper_sequential_znode" "from_dir" {
						path_prefix = "%s"
						data = "sequential znode created by passing a dir"
					}`, seqFromDir,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("zookeeper_sequential_znode.from_dir", "path", regexp.MustCompile(`^`+seqFromDir+`\d{10}`)),
					resource.TestCheckResourceAttrPair("zookeeper_sequential_znode.from_dir", "path", "zookeeper_sequential_znode.from_dir", "id"),
					resource.TestCheckResourceAttr("zookeeper_sequential_znode.from_dir", "data", "sequential znode created by passing a dir"),
					resource.TestCheckResourceAttr("zookeeper_sequential_znode.from_dir", "data_base64", "c2VxdWVudGlhbCB6bm9kZSBjcmVhdGVkIGJ5IHBhc3NpbmcgYSBkaXI="),
				),
			},
			{
				ResourceName:      "zookeeper_sequential_znode.from_dir",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceSeqZNode_FromPrefix(t *testing.T) {
	seqFromPrefix := "/" + acctest.RandString(10) + "/prefix-"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { checkPreconditions(t) },
		ProviderFactories: providerFactoriesMap(),
		CheckDestroy:      confirmAllZNodeDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "zookeeper_sequential_znode" "from_prefix" {
						path_prefix = "%s"
						data = "sequential znode created by passing a prefix"
					}`, seqFromPrefix,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("zookeeper_sequential_znode.from_prefix", "path", regexp.MustCompile(`^`+seqFromPrefix+`\d{10}`)),
					resource.TestCheckResourceAttrPair("zookeeper_sequential_znode.from_prefix", "path", "zookeeper_sequential_znode.from_prefix", "id"),
					resource.TestCheckResourceAttr("zookeeper_sequential_znode.from_prefix", "data", "sequential znode created by passing a prefix"),
					resource.TestCheckResourceAttr("zookeeper_sequential_znode.from_prefix", "data_base64", "c2VxdWVudGlhbCB6bm9kZSBjcmVhdGVkIGJ5IHBhc3NpbmcgYSBwcmVmaXg="),
				),
			},
			{
				ResourceName:      "zookeeper_sequential_znode.from_prefix",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
