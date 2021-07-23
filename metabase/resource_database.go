package metabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Details struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Db       string `json:"db"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type DetailsAdd struct {
	*Details
	Ssl bool `json:"ssl"`
}

type DatabaseCreate struct {
	Engine  string  `json:"engine"`
	Name    string  `json:"name"`
	Details Details `json:"details"`
}

type DatabaseAdd struct {
	Description              string     `json:"description"`
	Features                 []string   `json:"features"`
	CacheFieldValuesSchedule string     `json:"cache_field_values_schedule"`
	Timezone                 string     `json:"timezone"`
	AutoRunQueries           bool       `json:"auto_run_queries"`
	MetadataSyncSchedule     string     `json:"metadata_sync_schedule"`
	Name                     string     `json:"name"`
	Caveats                  string     `json:"caveats"`
	IsFullSync               bool       `json:"is_full_sync"`
	UpdatedAt                string     `json:"updated_at"`
	Details                  DetailsAdd `json:"details"`
	IsSample                 bool       `json:"is_sample"`
	Id                       int        `json:"id"`
	IsOnDemand               bool       `json:"is_on_demand"`
	Options                  string     `json:"options"`
	Engine                   string     `json:"engine"`
	Refingerprint            string     `json:"refingerprint"`
	CreatedAt                string     `json:"created_at"`
	PointsOfInterest         string     `json:"points_of_interest"`
}

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseCreate,
		ReadContext:   resourceDatabaseRead,
		UpdateContext: resourceDatabaseUpdate,
		DeleteContext: resourceDatabaseDelete,
		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"engine": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"db": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	jsonStr := DatabaseCreate{
		Engine: d.Get("engine").(string),
		Name:   d.Get("name").(string),
		Details: Details{
			Host:     d.Get("host").(string),
			Port:     d.Get("port").(int),
			Db:       d.Get("db").(string),
			User:     d.Get("user").(string),
			Password: d.Get("password").(string),
		},
	}

	jsonBody, err := json.Marshal(jsonStr)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to Marshal JSON in resourceDatabaseCreate()",
			Detail:   "Unable to Marshal JSON in resourceDatabaseCreate()",
		})
		return diags
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/database", "http://localhost:3000"), bytes.NewBuffer(jsonBody))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to make request in resourceDatabaseCreate()",
			Detail:   "Unable to make request in resourceDatabaseCreate()",
		})
		return diags
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Metabase-Session", c.Token)
	// disable gzip
	req.Header.Set("Accept-Encoding", "identity")

	r, err := client.Do(req)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to POST JSON in resourceDatabaseCreate()",
			Detail:   "Unable to POST JSON in resourceDatabaseCreate()",
		})
		return diags
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if r.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Return code is not 200 in resourceDatabaseCreate(): %s", body),
			Detail:   fmt.Sprintf("Return code is not 200 in resourceDatabaseCreate(): %s", body),
		})
		return diags
	}

	database := &DatabaseAdd{}

	err = json.Unmarshal(body, &database)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to decode JSON in resourceDatabaseCreate()",
			Detail:   "Unable to decode JSON in resourceDatabaseCreate()",
		})
		return diags
	}

	d.SetId(strconv.Itoa(database.Id))

	resourceDatabaseRead(ctx, d, m)

	return diags
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	dbId := d.Id()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/database/%s", "http://localhost:3000", dbId), nil)
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

	database := &DatabaseRead{}
	err = json.NewDecoder(r.Body).Decode(&database)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to decode JSON in resourceDatabaseRead()",
			Detail:   "Unable to decode JSON in resourceDatabaseRead()",
		})
		return diags
	}

	if err := d.Set("description", database.Description); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to assign database.Description to d.description in resourceDatabaseRead()",
			Detail:   "Unable to assign database.Description to d.description in resourceDatabaseRead()",
		})
		return diags
	}

	return diags
}

func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceDatabaseUpdate(ctx, d, m)
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
