package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

func datasourceZNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZNodeRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stat": {
				Type:     schema.TypeMap,
			"data_base64": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Content of the ZNode, encoded in Base64. " +
					"Use this if content is binary (i.e. sequence of bytes).",
			},
		},
	}
}

func dataSourceZNodeRead(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Get("path").(string)

	znode, err := zkClient.Read(znodePath)
	if err != nil {
		return diag.Errorf("Unable read ZNode from '%s': %v", znodePath, err)
	}

	// Terraform will use the ZNode.Path as unique identifier for this Data Source
	rscData.SetId(znode.Path)

	return setAttributesFromZNode(rscData, znode, diag.Diagnostics{})
}
