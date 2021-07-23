package metabase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBasesRead,
		Schema: map[string]*schema.Schema{
			"databases": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"features": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cache_field_values_schedule": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"timezone": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_run_queries": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"metadata_sync_schedule": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"caveats": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_full_sync": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"updated_at": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"native_permissions": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"details": &schema.Schema{
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"is_sample": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_on_demand": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"options": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"refingerprint": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"points_of_interest": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBasesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	client := c.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/database", c.HostURL), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Metabase-Session", c.Token)
	// disable gzip
	req.Header.Set("Accept-Encoding", "identity")

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	databases := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&databases)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to decode",
			Detail:   "Unable to decode JSON",
		})
		return diags
	}

	if err := d.Set("databases", databases); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
