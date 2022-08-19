package endpointmonitor

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"search": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}
