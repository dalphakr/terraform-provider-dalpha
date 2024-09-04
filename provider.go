package main

import (
	"github.com/dalphakr/dalpha-terraform-provider/openai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: func(data *schema.ResourceData) (interface{}, error) {
			return openai.Common{
				VaultAddr: data.Get("vault_addr").(string),
			}, nil
		},
		Schema: map[string]*schema.Schema{
			"vault_addr": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_ADDR", ""),
			},
		},
		ResourcesMap:   openai.Resources(),
		DataSourcesMap: map[string]*schema.Resource{},
	}
}
