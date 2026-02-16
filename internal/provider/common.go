// Package provider contains the Terraform provider implementation.
package provider

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"

	"github.com/go-zookeeper/zk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

// ErrACLPermNotAnInt is returned when the permissions value of an ACL is not an integer.
//
// This is a special error that is used to distinguish between a user error (e.g. a user
// providing a string value for the permissions field) and a provider error (e.g. a provider
// providing an invalid ACL).
//
// This error is not returned by the provider itself, but by the provider's helper functions
// that parse ACLs from the Terraform Schema.
var ErrACLPermNotAnInt = errors.New("acl permissions value is not an integer")

const (
	zNodeLinkForDesc = "[ZooKeeper ZNode](https://zookeeper.apache.org/doc/current/zookeeperProgrammers.html#sc_zkDataModel_znodes)"
)

// setAttributesFromZNode takes a *client.ZNode and populates the *schema.ResourceData with its content.
func setAttributesFromZNode(
	rscData *schema.ResourceData,
	znode *client.ZNode,
	diags diag.Diagnostics,
) diag.Diagnostics {
	if err := rscData.Set("path", znode.Path); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := rscData.Set("data", string(znode.Data)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := rscData.Set(
		"data_base64",
		base64.StdEncoding.EncodeToString(znode.Data),
	); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := rscData.Set("stat", []interface{}{zNodeStatToMap(znode)}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Convert ACLs from []zk.ACL to []map[string]interface{}
	aclConfigs := make([]map[string]interface{}, 0, len(znode.ACL))
	for _, acl := range znode.ACL {
		aclConfig := map[string]interface{}{
			"scheme":      acl.Scheme,
			"id":          acl.ID,
			"permissions": acl.Perms,
		}
		aclConfigs = append(aclConfigs, aclConfig)
	}

	if err := rscData.Set("acl", aclConfigs); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// statSchema provides the *schema.Schema to represent the ZNode Stat Structure.
// For more info: https://zookeeper.apache.org/doc/r3.5.9/zookeeperProgrammers.html#sc_zkStatStructure.
func statSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"czxid": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The zxid of the change that caused this znode to be created.",
				},
				"mzxid": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The zxid of the change that last modified this znode.",
				},
				"pzxid": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The zxid of the change that last modified children of this znode.",
				},
				"ctime": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The time in milliseconds from epoch when this znode was created.",
				},
				"mtime": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The time in milliseconds from epoch when this znode was last modified.",
				},
				"version": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The number of changes to the data of this znode.",
				},
				"cversion": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The number of changes to the children of this znode.",
				},
				"aversion": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The number of changes to the ACL of this znode.",
				},
				"ephemeral_owner": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The session id of the owner of this znode if the znode is an ephemeral node. If it is not an ephemeral node, it will be zero.",
				},
				"data_length": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The length of the data field of this znode.",
				},
				"num_children": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The number of children of this znode.",
				},
			},
		},
		Description: "[ZooKeeper Stat Structure](https://zookeeper.apache.org/doc/current/zookeeperProgrammers.html#sc_zkStatStructure) of the ZNode. " +
			"More details about `stat` can be found [here](../../docs#the-stat-structure).",
	}
}

// zNodeStatToMap is a helper that returns the zk.Stat contained to in client.ZNode,
// in the form of Terraform Schema compliant map.
func zNodeStatToMap(z *client.ZNode) map[string]interface{} {
	return map[string]interface{}{
		"czxid":           z.Stat.Czxid,
		"mzxid":           z.Stat.Mzxid,
		"pzxid":           z.Stat.Pzxid,
		"ctime":           z.Stat.Ctime,
		"mtime":           z.Stat.Mtime,
		"version":         z.Stat.Version,
		"cversion":        z.Stat.Cversion,
		"aversion":        z.Stat.Aversion,
		"ephemeral_owner": z.Stat.EphemeralOwner,
		"data_length":     z.Stat.DataLength,
		"num_children":    z.Stat.NumChildren,
	}
}

// getDataBytesFromResourceData reads the `data` or `data_base64` fields from the given *schema.ResourceData.
//
// If both fields are not set, it returns `nil` bytes, meaning the ZNode related to this resource/data-source
// has no content.
func getDataBytesFromResourceData(rscData *schema.ResourceData) ([]byte, error) {
	if dataRaw, exists := rscData.GetOk("data"); exists {
		return []byte(dataRaw.(string)), nil
	}

	if dataRawBase64, exists := rscData.GetOk("data_base64"); exists {
		dataBytes, err := base64.StdEncoding.DecodeString(dataRawBase64.(string))
		if err != nil {
			return nil, fmt.Errorf("decoding 'data_base64' from Base64 failed: %w", err)
		}
		return dataBytes, nil
	}

	return nil, nil
}

func parseACLsFromResourceData(rscData *schema.ResourceData) ([]zk.ACL, error) {
	aclConfigs := rscData.Get("acl").([]interface{})
	acls := make([]zk.ACL, 0, len(aclConfigs))

	for _, aclConfig := range aclConfigs {
		aclMap := aclConfig.(map[string]interface{})
		scheme := aclMap["scheme"].(string)
		id := aclMap["id"].(string)
		permissionsValue, ok := aclMap["permissions"].(int)
		if !ok {
			return nil, ErrACLPermNotAnInt
		}
		if permissionsValue < math.MinInt32 || permissionsValue > math.MaxInt32 {
			return nil, NewACLPermissionsValueOutOfRangeError(permissionsValue)
		}
		permissions := int32(permissionsValue)

		acls = append(acls, zk.ACL{
			Scheme: scheme,
			ID:     id,
			Perms:  permissions,
		})
	}

	if len(acls) == 0 {
		acls = zk.WorldACL(zk.PermAll)
	}

	return acls, nil
}
