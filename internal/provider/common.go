package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

func setResourceDataFromZNode(rscData *schema.ResourceData, znode *client.ZNode, diags diag.Diagnostics) diag.Diagnostics {
	if err := rscData.Set("path", znode.Path); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := rscData.Set("data", string(znode.Data)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := rscData.Set("stat", znode.StatAsMap()); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

func getFieldDataFromResourceData(rscData *schema.ResourceData) []byte {
	znodeDataRaw, exists := rscData.GetOk("data")
	if exists {
		znodeDataStr := znodeDataRaw.(string)
		return []byte(znodeDataStr)
	} else {
		return nil
	}
}
