package metabase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("METABASE_URL", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("METABASE_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("METABASE_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"metabase_database": resourceDatabase(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"metabase_bases": dataSourceBases(),
			"metabase_base":  dataSourceBase(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (username != "") && (password != "") && (url != "") {
		c, err := NewClient(&url, &username, &password)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}

	c, err := NewClient(nil, nil, nil)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
