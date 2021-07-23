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

type DetailsDb struct {
	Db string `json:"db"`
}

type CacheFieldValues struct {
	ScheduleMinute int    `json:"schedule_minute"`
	ScheduleDay    string `json:"schedule_day"`
	ScheduleFrame  string `json:"schedule_frame"`
	ScheduleHour   int    `json:"schedule_hour"`
	ScheduleType   string `json:"schedule_type"`
}

type Schedules struct {
	CacheFieldValues CacheFieldValues `json:"cache_field_values"`
	MetadataSync     CacheFieldValues `json:"metadata_sync"`
}

type DatabaseRead struct {
	Description              string    `json:"description"`
	Features                 []string  `json:"features"`
	CacheFieldValuesSchedule string    `json:"cache_field_values_schedule"`
	Timezone                 string    `json:"timezone"`
	AutoRunQueries           bool      `json:"auto_run_queries"`
	MetadataSyncSchedule     string    `json:"metadata_sync_schedule"`
	Name                     string    `json:"name"`
	Caveats                  string    `json:"caveats"`
	IsFullSync               bool      `json:"is_full_sync"`
	UpdatedAt                string    `json:"updated_at"`
	Details                  DetailsDb `json:"details"`
	IsSample                 bool      `json:"is_sample"`
	Id                       int       `json:"id"`
	IsOnDemand               bool      `json:"is_on_demand"`
	Options                  string    `json:"options"`
	Schedules                Schedules `json:"schedules"`
	Engine                   string    `json:"engine"`
	Refingerprint            string    `json:"refingerprint"`
	CreatedAt                string    `json:"created_at"`
	PointsOfInterest         string    `json:"points_of_interest"`
}

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
	client := &http.Client{Timeout: 10 * time.Second}

	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	dbId := strconv.Itoa(d.Get("id").(int))

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/database/%s", "http://localhost:3000", dbId), nil)
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
