package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Description:   "`metabase_database` resource can be used for managing databases (CRUD).\n\n",
		CreateContext: resourceDatabaseCreate,
		ReadContext:   resourceDatabaseRead,
		UpdateContext: resourceDatabaseUpdate,
		DeleteContext: resourceDatabaseDelete,
		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Description: "Description of a source in Metabase",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": &schema.Schema{
				Description: "Name of the source in Metabase",
				Type:        schema.TypeString,
				Required:    true,
			},
			"engine": &schema.Schema{
				Description: "Engine of a database. See [Officially supported databases](https://github.com/metabase/metabase/blob/master/docs/administration-guide/01-managing-databases.md)",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "postgres",
			},
			"host": &schema.Schema{
				Description: "Database host: IP or hostname",
				Type:        schema.TypeString,
				Required:    true,
			},
			"port": &schema.Schema{
				Description: "Database port",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5432,
			},
			"db": &schema.Schema{
				Description: "Database name inside an engine",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user": &schema.Schema{
				Description: "User name to connect to a database",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": &schema.Schema{
				Description: "Password to connect to a database",
				Type:        schema.TypeString,
				Required:    true,
			},
			"last_updated": &schema.Schema{
				Description: "Timestamp when a database has been updated last time",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	client := c.HTTPClient

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

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/database", c.HostURL), bytes.NewBuffer(jsonBody))
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

	database := &Database{}

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
	c := m.(*Client)
	client := c.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	dbId := d.Id()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/database/%s", c.HostURL, dbId), nil)
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

	body, _ := ioutil.ReadAll(r.Body)
	if r.StatusCode == http.StatusNotFound {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Database not found in func resourceDatabaseRead(): %s", body),
			Detail:   fmt.Sprintf("Database not found in func resourceDatabaseRead(): %s", body),
		})
		return diags
	}

	if r.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Return code is not 200 in resourceDatabaseRead(): %s", body),
			Detail:   fmt.Sprintf("Return code is not 200 in resourceDatabaseRead(): %s", body),
		})
		return diags
	}

	database := &DatabaseRead{}
	err = json.Unmarshal(body, &database)
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
	c := m.(*Client)
	client := c.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	dbId := d.Id()

	if d.HasChanges("engine", "name", "host", "port", "db", "user", "password") {
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
				Summary:  "Unable to Marshal JSON in resourceDatabaseUpdate()",
				Detail:   "Unable to Marshal JSON in resourceDatabaseUpdate()",
			})
			return diags
		}

		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/database/%s", c.HostURL, dbId), bytes.NewBuffer(jsonBody))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to make request in resourceDatabaseUpdate()",
				Detail:   "Unable to make request in resourceDatabaseUpdate()",
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
				Summary:  "Unable to POST JSON in resourceDatabaseUpdate()",
				Detail:   "Unable to POST JSON in resourceDatabaseUpdate()",
			})
			return diags
		}
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if r.StatusCode != 200 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Return code is not 200 in resourceDatabaseUpdate(): %s", body),
				Detail:   fmt.Sprintf("Return code is not 200 in resourceDatabaseUpdate(): %s", body),
			})
			return diags
		}

		database := &Database{}

		err = json.Unmarshal(body, &database)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to decode JSON in resourceDatabaseUpdate()",
				Detail:   "Unable to decode JSON in resourceDatabaseUpdate()",
			})
			return diags
		}

		d.Set("last_updated", database.UpdatedAt)
	}

	return resourceDatabaseRead(ctx, d, m)
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	client := c.HTTPClient

	dbId := d.Id()

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/database/%s", c.HostURL, dbId), nil)
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
			Summary:  "Unable to DELETE in resourceDatabaseCreate()",
			Detail:   "Unable to DELETE in resourceDatabaseCreate()",
		})
		return diags
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if r.StatusCode != http.StatusNoContent {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Return code is not 204 in resourceDatabaseCreate(): %s", body),
			Detail:   fmt.Sprintf("Return code is not 204 in resourceDatabaseCreate(): %s", body),
		})
		return diags
	}

	d.SetId("")

	return diags
}
