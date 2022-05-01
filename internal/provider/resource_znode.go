package provider

import (
	"context"
	"errors"

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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stat": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceZNodeCreate(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Get("path").(string)

	znode, err := zkClient.Create(znodePath, getFieldDataFromResourceData(rscData))
	if err != nil {
		return diag.Errorf("Failed to create ZNode '%s': %v", znodePath, err)
	}

	// Terraform will use the ZNode.Path as unique identifier for this Resource
	rscData.SetId(znode.Path)
	rscData.MarkNewResource()

	return setResourceDataFromZNode(rscData, znode, diag.Diagnostics{})
}

func resourceZNodeRead(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Id()

	znode, err := zkClient.Read(znodePath)
	if err != nil {
		// If the ZNode is not found, it means it was changed outside of Terraform.
		// We set the ID to blank, so it's state will be removed.
		if errors.Is(err, client.ErrorZNodeDoesNotExist) {
			rscData.SetId("")
			return diag.Diagnostics{}
		}

		return diag.Errorf("Failed to read ZNode '%s': %v", znodePath, err)
	}

	return setResourceDataFromZNode(rscData, znode, diag.Diagnostics{})
}

func resourceZNodeUpdate(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Id()

	if rscData.HasChange("data") {
		znode, err := zkClient.Update(znodePath, getFieldDataFromResourceData(rscData))
		if err != nil {
			return diag.Errorf("Failed to update ZNode '%s': %v", znodePath, err)
		}

		return setResourceDataFromZNode(rscData, znode, diag.Diagnostics{})
	}

	return diag.Diagnostics{}
}

func resourceZNodeDelete(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Id()

	err := zkClient.Delete(znodePath)
	if err != nil {
		return diag.Errorf("Failed to delete ZNode '%s': %v", znodePath, err)
	}

	return diag.Diagnostics{}
}
