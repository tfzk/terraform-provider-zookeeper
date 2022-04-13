package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

func datasourceZNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZNodeRead,
		Schema: map[string]*schema.Schema{
			fieldPath: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			fieldData: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			fieldStat: &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceZNodeRead(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Get(fieldPath).(string)

	znode, err := zkClient.Read(znodePath)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable read ZNode from '%s': %v", znodePath, err),
		})
	}

	// Terraform will use the ZNode.Path as unique identifier for this Data Source
	rscData.SetId(znode.Path)

	return setResourceDataFromZNode(rscData, &znode, diags)
}
