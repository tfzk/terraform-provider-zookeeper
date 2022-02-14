package zookeeper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tfzk/client"
)

const (
	typeZNode    = "zookeeper_znode"
	typeSeqZNode = "zookeeper_sequential_znode"

	fieldPath       = "path"
	fieldPathPrefix = "path_prefix"
	fieldData       = "data"
	fieldStat       = "stat"
)

func setResourceDataFromZNode(rscData *schema.ResourceData, znode *client.ZNode, diags diag.Diagnostics) diag.Diagnostics {
	if err := rscData.Set(fieldPath, znode.Path); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := rscData.Set(fieldData, string(znode.Data)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := rscData.Set(fieldStat, znode.StatAsMap()); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

func getFieldDataFromResourceData(rscData *schema.ResourceData) []byte {
	znodeDataRaw, exists := rscData.GetOk(fieldData)
	if exists {
		znodeDataStr := znodeDataRaw.(string)
		return []byte(znodeDataStr)
	} else {
		return nil
	}
}
