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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Absolute path to the ZNode to read.",
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content of the ZNode. Use this if content is a UTF-8 string.",
			},
			"data_base64": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Content of the ZNode, encoded in Base64. " +
					"Use this if content is binary (i.e. sequence of bytes).",
			},
			"stat": statSchema(),
			"acl": {
				Type:        schema.TypeList,
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
		Description: "Provides access to the content of a " +
			zNodeLinkForDesc + ". " +
			"The data is loaded both as UTF-8 string, as well as Base64 encoded bytes. " +
			"The ability to access ZNodes is determined by ZooKeeper ACL.",
	}
}

func dataSourceZNodeRead(_ context.Context, rscData *schema.ResourceData, prvClient interface{}) diag.Diagnostics {
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
