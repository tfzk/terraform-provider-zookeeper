package provider

import (
	"encoding/base64"
	"fmt"

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
	if err := rscData.Set("data_base64", base64.StdEncoding.EncodeToString(znode.Data)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
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
