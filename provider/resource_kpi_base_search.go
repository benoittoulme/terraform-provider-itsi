package main

import (
	"encoding/json"

	"github.com/benoittoulme/terraform-provider-itsi/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func kpiBaseSearchBase(key string, title string) *models.Base {
	base := models.NewBase(key, title, "kpi_base_search")
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
			// "_key": {
			// 	Type:         schema.TypeString,
			// 	Optional:     true,
			// 	InputDefault: "",
			// },
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
			"is_first_time_save_done": {
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

func metric(source map[string]interface{}) interface{} {
	m := map[string]interface{}{}
	m["_key"] = source["_key"]
	m["aggregate_statop"] = source["aggregate_statop"]
	m["entity_statop"] = source["entity_statop"]
	m["fill_gaps"] = source["fill_gaps"]
	m["gap_custom_alert_value"] = source["gap_custom_alert_value"]
	m["gap_severity"] = source["gap_severity"]
	m["gap_severity_color"] = source["gap_severity_color"]
	m["gap_severity_color_light"] = source["gap_severity_color_light"]
	m["gap_severity_value"] = source["gap_severity_value"]
	m["threshold_field"] = source["threshold_field"]
	m["title"] = source["title"]
	m["unit"] = source["unit"]
	return m
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
	body["isFirstTimeSaveDone"] = d.Get("is_first_time_save_done").(bool)
	body["metric_qualifier"] = d.Get("metric_qualifier").(string)

	metrics := []interface{}{}
	for _, g := range d.Get("metrics").(*schema.Set).List() {
		metrics = append(metrics, metric(g.(map[string]interface{})))
		if err != nil {
			return nil, err
		}
	}
	body["metrics"] = metrics
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
	b, err := template.Create(client.User, client.Password, client.Host, client.Port)
	if err != nil {
		return err
	}
	b.Read(client.User, client.Password, client.Host, client.Port)
	return populate(b, d)
}

func kpiBaseSearchRead(d *schema.ResourceData, m interface{}) error {
	client := m.(client)

	base := kpiBaseSearchBase(d.Id(), d.Get("title").(string))
	b, err := base.Find(client.User, client.Password, client.Host, client.Port)
	if err != nil {
		return err
	}
	if b == nil {
		d.SetId("")
		return nil
	}
	return populate(b, d)
}

func populateBaseSearch(b *models.Base, d *schema.ResourceData) error {
	by, err := b.RawJson.MarshalJSON()
	if err != nil {
		return err
	}
	var interfaceMap map[string]interface{}
	err = json.Unmarshal(by, &interfaceMap)
	if err != nil {
		return err
	}

	err = d.Set("title", interfaceMap["title"])
	if err != nil {
		return err
	}
	err = d.Set("description", interfaceMap["description"])
	if err != nil {
		return err
	}
	err = d.Set("actions", interfaceMap["actions"])
	if err != nil {
		return err
	}
	err = d.Set("alert_lag", interfaceMap["alert_lag"])
	if err != nil {
		return err
	}
	err = d.Set("alert_period", interfaceMap["alert_period"])
	if err != nil {
		return err
	}
	err = d.Set("base_search", interfaceMap["base_search"])
	if err != nil {
		return err
	}
	err = d.Set("entity_alias_filtering_fields", interfaceMap["entity_alias_filtering_fields"])
	if err != nil {
		return err
	}
	err = d.Set("entity_breakdown_id_fields", interfaceMap["entity_breakdown_id_fields"])
	if err != nil {
		return err
	}
	err = d.Set("entity_id_fields", interfaceMap["entity_id_fields"])
	if err != nil {
		return err
	}
	err = d.Set("identifying_name", interfaceMap["identifying_name"])
	if err != nil {
		return err
	}
	err = d.Set("is_entity_breakdown", interfaceMap["is_entity_breakdown"])
	if err != nil {
		return err
	}
	err = d.Set("is_service_entity_filter", interfaceMap["is_service_entity_filter"])
	if err != nil {
		return err
	}
	err = d.Set("is_first_time_save_done", interfaceMap["isFirstTimeSaveDone"])
	if err != nil {
		return err
	}
	err = d.Set("metric_qualifier", interfaceMap["metric_qualifier"])
	if err != nil {
		return err
	}
	err = d.Set("metrics", interfaceMap["metrics"])
	if err != nil {
		return err
	}
	err = d.Set("search_alert_earliest", interfaceMap["search_alert_earliest"])
	if err != nil {
		return err
	}
	err = d.Set("sec_grp", interfaceMap["sec_grp"])
	if err != nil {
		return err
	}
	err = d.Set("source_itsi_da", interfaceMap["source_itsi_da"])
	if err != nil {
		return err
	}
	d.SetId(b.RESTKey)
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
