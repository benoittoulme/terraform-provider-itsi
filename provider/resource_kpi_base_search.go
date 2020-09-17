package main

import (
	"encoding/json"

	"github.com/benoittoulme/terraform-provider-itsi/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func kpiBaseSearchBase(key string, title string) *models.Base {
	base := models.NewBase(key, title, "itoa_interface", "kpi_base_search")
	base.TFIDField = func() string {
		return "title"
	}
	return base
}

func resourceKPIBaseSearch() *schema.Resource {
	return &schema.Resource{
		Create: kpiBaseSearchCreate,
		Read:   kpiBaseSearchRead,
		Update: kpiBaseSearchUpdate,
		Delete: kpiBaseSearchDelete,
		Importer: &schema.ResourceImporter{
			State: kpiBaseSearchImport,
		},
		Schema: map[string]*schema.Schema{
			"_key": {
				Type:         schema.TypeString,
				Optional:     true,
				InputDefault: "",
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"actions": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alert_lag": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alert_period": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_search": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entity_alias_filtering_fields": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entity_breakdown_id_fields": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entity_id_fields": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identifying_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_entity_breakdown": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"is_service_entity_filter": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"isFirstTimeSaveDone": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"metric_qualifier": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metrics": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"_key": {
							Type:         schema.TypeString,
							Optional:     true,
							InputDefault: "",
						},
						"aggregate_statop": {
							Type:     schema.TypeString,
							Required: true,
						},
						"entity_statop": {
							Type:     schema.TypeString,
							Required: true,
						},
						"fill_gaps": {
							Type:     schema.TypeString,
							Required: true,
						},
						"gap_custom_alert_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"gap_severity": {
							Type:     schema.TypeString,
							Required: true,
						},
						"gap_severity_color": {
							Type:     schema.TypeString,
							Required: true,
						},
						"gap_severity_color_light": {
							Type:     schema.TypeString,
							Required: true,
						},
						"gap_severity_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"threshold_field": {
							Type:     schema.TypeString,
							Required: true,
						},
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"search_alert_earliest": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sec_grp": {
				Type:     schema.TypeString,
				Required: true,
			},

			"source_itsi_da": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func kpiBaseSearch(d *schema.ResourceData) (config *models.Base, err error) {
	body := map[string]interface{}{}
	body["objectType"] = "kpi_base_search"
	body["title"] = d.Get("title").(string)
	body["description"] = d.Get("description").(string)

	body["actions"] = d.Get("actions").(string)
	body["alert_lag"] = d.Get("alert_lag").(string)
	body["alert_period"] = d.Get("alert_period").(string)
	body["base_search"] = d.Get("base_search").(string)
	body["entity_alias_filtering_fields"] = d.Get("entity_alias_filtering_fields").(string)
	body["entity_breakdown_id_fields"] = d.Get("entity_breakdown_id_fields").(string)
	body["entity_id_fields"] = d.Get("entity_id_fields").(string)
	body["identifying_name"] = d.Get("identifying_name").(string)
	body["is_entity_breakdown"] = d.Get("is_entity_breakdown").(bool)
	body["is_service_entity_filter"] = d.Get("is_service_entity_filter").(bool)
	body["isFirstTimeSaveDone"] = d.Get("isFirstTimeSaveDone").(bool)
	body["metric_qualifier"] = d.Get("metric_qualifier").(string)

	// body["metrics"] = d.Get("metrics").(string)

	body["search_alert_earliest"] = d.Get("search_alert_earliest").(bool)
	body["sec_grp"] = d.Get("sec_grp").(string)
	body["source_itsi_da"] = d.Get("source_itsi_da").(bool)

	by, err := json.Marshal(body)
	if err != nil {
		return
	}
	base := kpiBaseSearchBase(d.Id(), d.Get("title").(string))
	err = json.Unmarshal(by, &base.RawJson)
	if err != nil {
		return nil, err
	}
	return base, nil
}

func kpiBaseSearchCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(client)
	template, err := kpiBaseSearch(d)
	if err != nil {
		return err
	}
	return template.Create(client.User, client.Password, client.Host, client.Port)
}

func kpiBaseSearchRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func kpiBaseSearchUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func kpiBaseSearchDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func kpiBaseSearchImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}
