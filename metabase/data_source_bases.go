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
						"details_host": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"details_port": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"details_db": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"details_user": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"details_password": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"details_ssl": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
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

	var databases []Database

	//body, err := ioutil.ReadAll(r.Body)
	//err = json.Unmarshal(body, &databases)
	err = json.NewDecoder(r.Body).Decode(&databases)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to decode",
			Detail:   fmt.Sprintf("Unable to decode JSON: %s", err.Error()),
		})
		return diags
	}

	flattenned_databases := flattenDatabases(databases)

	if err := d.Set("databases", flattenned_databases); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenDatabases(databases []Database) []interface{} {
	if databases != nil {
		ois := make([]interface{}, len(databases), len(databases))

		for i, database := range databases {
			oi := make(map[string]interface{})

			oi["description"] = database.Description
			oi["features"] = database.Features
			oi["cache_field_values_schedule"] = database.CacheFieldValuesSchedule
			oi["timezone"] = database.Timezone
			oi["auto_run_queries"] = database.AutoRunQueries
			oi["metadata_sync_schedule"] = database.MetadataSyncSchedule
			oi["name"] = database.Name
			oi["caveats"] = database.Caveats
			oi["is_full_sync"] = database.IsFullSync
			oi["updated_at"] = database.UpdatedAt
			oi["native_permissions"] = database.NativePermissions
			oi["details_host"] = database.Details.Host
			oi["details_port"] = database.Details.Port
			oi["details_db"] = database.Details.Db
			oi["details_user"] = database.Details.User
			oi["details_password"] = database.Details.Password
			oi["details_ssl"] = database.Details.Ssl
			oi["is_sample"] = database.IsSample
			oi["id"] = database.Id
			oi["is_on_demand"] = database.IsOnDemand
			oi["options"] = database.Options
			oi["engine"] = database.Engine
			oi["refingerprint"] = database.Refingerprint
			oi["created_at"] = database.CreatedAt
			oi["points_of_interest"] = database.PointsOfInterest

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
