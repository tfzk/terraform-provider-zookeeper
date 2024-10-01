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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Absolute path to the ZNode to create.",
			},
			"data": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"data_base64"},
				Description: "Content to store in the ZNode, as a UTF-8 string. " +
					"Mutually exclusive with `data_base64`.",
			},
			"data_base64": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"data"},
				Description: "Content to store in the ZNode, as Base64 encoded bytes. " +
					"Mutually exclusive with `data`.",
			},
			"stat": statSchema(),
			"acl": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of ACL entries for the ZNode.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheme": {
							Type:     schema.TypeString,
							Required: true,
							Description: "The ACL scheme, such as 'world', 'digest', " +
								"'ip', 'x509'.",
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
							Description: "The ID for the ACL entry. For example, " +
								"user:hash in 'digest' scheme.",
						},
						"permissions": {
							Type:     schema.TypeInt,
							Required: true,
							Description: "The permissions for the ACL entry, " +
								"represented as an integer bitmask.",
						},
					},
				},
			},
		},
		Description: "Manages the lifecycle of a " +
			zNodeLinkForDesc + ". " +
			"This resource manages **Persistent ZNodes**. " +
			"The data can be provided either as UTF-8 string, or as Base64 encoded bytes. " +
			"The ability to create ZNodes is determined by ZooKeeper ACL.",
	}
}

func resourceZNodeCreate(_ context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Get("path").(string)

	dataBytes, err := getDataBytesFromResourceData(rscData)
	if err != nil {
		return diag.FromErr(err)
	}

	acls, err := parseACLsFromResourceData(rscData)
	if err != nil {
		return diag.FromErr(err)
	}

	znode, err := zkClient.Create(znodePath, dataBytes, acls)
	if err != nil {
		return diag.Errorf("Failed to create ZNode '%s': %v", znodePath, err)
	}

	// Terraform will use the ZNode.Path as unique identifier for this Resource
	rscData.SetId(znode.Path)
	rscData.MarkNewResource()

	return setAttributesFromZNode(rscData, znode, diag.Diagnostics{})
}

func resourceZNodeRead(_ context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
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

	return setAttributesFromZNode(rscData, znode, diag.Diagnostics{})
}

func resourceZNodeUpdate(_ context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Id()

	if rscData.HasChanges("data", "data_base64", "acl") {
		dataBytes, err := getDataBytesFromResourceData(rscData)
		if err != nil {
			return diag.FromErr(err)
		}

		acls, err := parseACLsFromResourceData(rscData)
		if err != nil {
			return diag.FromErr(err)
		}

		znode, err := zkClient.Update(znodePath, dataBytes, acls)
		if err != nil {
			return diag.Errorf("Failed to update ZNode '%s': %v", znodePath, err)
		}

		return setAttributesFromZNode(rscData, znode, diag.Diagnostics{})
	}

	return diag.Diagnostics{}
}

func resourceZNodeDelete(_ context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
	zkClient := prvClient.(*client.Client)

	znodePath := rscData.Id()

	err := zkClient.Delete(znodePath)
	if err != nil {
		return diag.Errorf("Failed to delete ZNode '%s': %v", znodePath, err)
	}

	return diag.Diagnostics{}
}
