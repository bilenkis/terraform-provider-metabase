package metabase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBaseRead,
		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	client := c.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	dbId := strconv.Itoa(d.Get("id").(int))

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/database/%s", c.HostURL, dbId), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Metabase-Session", c.Token)

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	database := &DatabaseRead{}
	err = json.NewDecoder(r.Body).Decode(&database)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to decode JSON in dataSourceBaseRead()",
			Detail:   "Unable to decode JSON in dataSourceBaseRead()",
		})
		return diags
	}

	if err := d.Set("id", database.Id); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set database.Id to d.id in dataSourceBaseRead()",
			Detail:   "Unable to set database.Id to d.id in dataSourceBaseRead()",
		})
		return diags
	}

	if err := d.Set("name", database.Name); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to assign database.Name to d.name in dataSourceBaseRead()",
			Detail:   "Unable to assign database.Name to d.name in dataSourceBaseRead()",
		})
		return diags
	}

	if err := d.Set("description", database.Description); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to assign database.Description to d.description in dataSourceBaseRead()",
			Detail:   "Unable to assign database.Description to d.description in dataSourceBaseRead()",
		})
		return diags
	}

	d.SetId(strconv.Itoa(database.Id))

	return diags
}
