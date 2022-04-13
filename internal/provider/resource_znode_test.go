package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceZNodeBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      testParentAndChildZNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: parentAndChildZNodeConfig("/"+acctest.RandString(10), "parent data", "child data"),
				Check:  testParentAndChildZNode,
			},
		},
	})
}

func parentAndChildZNodeConfig(parentPath, parentData, childData string) string {
	return fmt.Sprintf(`
	resource "zookeeper_znode" "parent" {
		path = "%s"
		data = "%s"
	}
	resource "zookeeper_znode" "child" {
		path = "${zookeeper_znode.parent.path}/child"
		data = "%s"
	}`, parentPath, parentData, childData)
}

func testParentAndChildZNode(s *terraform.State) error {
	// Check ZNode paths hierarchy
	parentResourceState := s.Modules[0].Resources["zookeeper_znode.parent"]
	if parentResourceState == nil {
		return fmt.Errorf("'zookeeper_znode.parent' not found in state")
	}
	parentInstanceState := parentResourceState.Primary

	childResourceState := s.Modules[0].Resources["zookeeper_znode.child"]
	if childResourceState == nil {
		return fmt.Errorf("'zookeeper_znode.child' not found in state")
	}
	childInstanceState := childResourceState.Primary

	// Check IDs and Paths
	if !strings.HasPrefix(childInstanceState.ID, parentInstanceState.ID) {
		return fmt.Errorf("'zookeeper_znode.child' is not a child ZNode of 'zookeeper_znode.parent'")
	}

	if parentInstanceState.Attributes["path"] != parentInstanceState.ID {
		return fmt.Errorf("'zookeeper_znode.parent.path' does not match 'zookeeper_znode.parent.id'")
	}

	if childInstanceState.Attributes["path"] != childInstanceState.ID {
		return fmt.Errorf("'zookeeper_znode.child.path' does not match 'zookeeper_znode.child.id'")
	}

	zkClient := getTestAccClient()

	// Check data (parent)
	parentZNode, err := zkClient.Read(parentInstanceState.Attributes["path"])
	if err != nil {
		return err
	}
	if parentInstanceState.Attributes["data"] != string(parentZNode.Data) {
		return fmt.Errorf("'zookeeper_znode.parent.data' does not match underlying ZNode")
	}
	if parentInstanceState.Attributes["data"] != "parent data" {
		return fmt.Errorf("'zookeeper_znode.parent.data' does not match expected value")
	}

	// Check data (child)
	childZNode, err := zkClient.Read(childInstanceState.Attributes["path"])
	if err != nil {
		return err
	}
	if childInstanceState.Attributes["data"] != string(childZNode.Data) {
		return fmt.Errorf("'zookeeper_znode.child.data' does not match underlying ZNode")
	}
	if childInstanceState.Attributes["data"] != "child data" {
		return fmt.Errorf("'zookeeper_znode.child.data' does not match expected value")
	}

	return nil
}

func testParentAndChildZNodeDestroy(s *terraform.State) error {
	zkClient := getTestAccClient()

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
