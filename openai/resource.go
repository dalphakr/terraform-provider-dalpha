package openai

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"dalpha_openai_vault_apikey": {
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
				},
				"project_id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"project_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"entity": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
			Create: func(data *schema.ResourceData, i interface{}) error {
				// Create http request to openAI server
				comm := i.(Common)
				id, apiKey, err := CreateServiceAccount(data.Get("project_id").(string), data.Get("entity").(string))
				if err != nil {
					return fmt.Errorf("failed to create service account: %s", err)
				}

				if err = InsertVault(comm, data.Get("project_name").(string), data.Get("entity").(string), apiKey); err != nil {
					return fmt.Errorf("failed to insert into vault: %s", err)
				}

				data.SetId(id)
				return nil
			},
			Read: func(data *schema.ResourceData, i interface{}) error {
				comm := i.(Common)
				if ok, err := ExistsVault(comm, data.Get("project_name").(string), data.Get("entity").(string)); err != nil {
					return fmt.Errorf("failed to check secret exists: %s", err)
				} else if !ok {
					return fmt.Errorf("secret not found")
				}

				sa, err := FindServiceAccount(data.Get("project_id").(string), data.Get("entity").(string))
				if err != nil {
					return fmt.Errorf("failed to find service account: %s", err)
				}

				data.Set("type", "service_account")
				data.Set("project_id", data.Get("project_id").(string))
				data.Set("project_name", data.Get("project_name").(string))
				data.Set("entity", sa.Name)
				data.SetId(sa.Id)
				return nil
			},
			Update: func(data *schema.ResourceData, i interface{}) error {
				return nil
			},
			Delete: func(data *schema.ResourceData, i interface{}) error {
				return nil
			},
		},
	}
}
