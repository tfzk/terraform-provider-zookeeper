package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tfzk/terraform-provider-zookeeper/internal/client"
)

func resourceSeqZNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSeqZNodeCreate,
		ReadContext:   resourceSeqZNodeRead,
		UpdateContext: resourceSeqZNodeUpdate,
		DeleteContext: resourceSeqZNodeDelete,
		Schema: map[string]*schema.Schema{
			"path_prefix": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stat": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceSeqZNodeCreate(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	zkClient := prvClient.(*client.Client)

	znodePathPrefix := rscData.Get("path_prefix").(string)

	znode, err := zkClient.CreateSequential(znodePathPrefix, getFieldDataFromResourceData(rscData))
	if err != nil {
		return diag.Errorf("Failed to create Sequential ZNode '%s': %v", znodePathPrefix, err)
	}

	// Terraform will use the ZNode.Path as unique identifier for this Resource
	rscData.SetId(znode.Path)
	rscData.MarkNewResource()

	return setResourceDataFromZNode(rscData, znode, diag.Diagnostics{})
}

func resourceSeqZNodeRead(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	return resourceZNodeRead(ctx, rscData, prvClient)
}

func resourceSeqZNodeUpdate(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	return resourceZNodeUpdate(ctx, rscData, prvClient)
}

func resourceSeqZNodeDelete(ctx context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	return resourceZNodeDelete(ctx, rscData, prvClient)
}
