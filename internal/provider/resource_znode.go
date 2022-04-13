package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

func resourceZNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceZNodeCreate,
		ReadContext:   resourceZNodeRead,
		UpdateContext: resourceZNodeUpdate,
		DeleteContext: resourceZNodeDelete,
		Schema: map[string]*schema.Schema{
			fieldPath: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			fieldData: {
				Type:     schema.TypeString,
				Optional: true,
			},
			fieldStat: {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceZNodeCreate(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}

	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Get(fieldPath).(string)

	znode, err := zkClient.Create(znodePath, getFieldDataFromResourceData(rscData))
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to create ZNode '%s': %v", znodePath, err),
		})
	}

	// Terraform will use the ZNode.Path as unique identifier for this Resource
	rscData.SetId(znode.Path)
	rscData.MarkNewResource()

	return setResourceDataFromZNode(rscData, &znode, diags)
}

func resourceZNodeRead(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}

	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Id()

	znode, err := zkClient.Read(znodePath)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to read ZNode '%s': %v", znodePath, err),
		})
	}

	return setResourceDataFromZNode(rscData, &znode, diags)
}

func resourceZNodeUpdate(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}

	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Id()

	if rscData.HasChange(fieldData) {
		znode, err := zkClient.Update(znodePath, getFieldDataFromResourceData(rscData))
		if err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to update ZNode '%s': %v", znodePath, err),
			})
		}

		return setResourceDataFromZNode(rscData, &znode, diags)
	}

	return diags
}

func resourceZNodeDelete(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}

	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Id()

	err := zkClient.Delete(znodePath)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to delete ZNode '%s': %v", znodePath, err),
		})
	}

	return diags
}
