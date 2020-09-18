package main

import (
	"encoding/json"

	"github.com/benoittoulme/terraform-provider-itsi/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func kpiThresholdTemplateBase(key string, title string) *models.Base {
	base := models.NewBase(key, title, "itoa_interface", "kpi_threshold_template")
	base.TFIDField = func() string {
		return "title"
	}
	return base
}

func resourceThreshold() *schema.Resource {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"base_severity_color": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_severity_color_light": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_severity_label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_severity_value": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"gauge_max": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"gauge_min": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"is_max_static": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"is_min_static": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"metric_field": {
				Type:     schema.TypeString,
				Required: true,
			},
			"render_boundary_max": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"render_boundary_min": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"search": {
				Type:     schema.TypeString,
				Required: true,
			},
			"threshold_levels": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dynamic_param": {
							Type:     schema.TypeString,
							Required: true,
						},
						"severity_color": {
							Type:     schema.TypeString,
							Required: true,
						},
						"severity_color_light": {
							Type:     schema.TypeString,
							Required: true,
						},
						"severity_label": {
							Type:     schema.TypeString,
							Required: true,
						},
						"severity_value": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"threshold_value": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
	return r
}

func resourceKPIThresholdTemplate() *schema.Resource {
	return &schema.Resource{
		Create: kpiThresholdTemplateCreate,
		Read:   kpiThresholdTemplateRead,
		Update: kpiThresholdTemplateUpdate,
		Delete: kpiThresholdTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: kpiThresholdTemplateImport,
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"adaptive_thresholds_is_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"adaptive_thresholding_training_window": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time_variate_thresholds": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"sec_grp": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time_variate_thresholds_specification": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policies": {
							Required: true,
							Type:     schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"title": {
										Type:     schema.TypeString,
										Required: true,
									},
									"policy_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"time_blocks": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cron": {
													Type:     schema.TypeString,
													Required: true,
												},
												"interval": {
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
									},
									"aggregate_thresholds": {
										Required: true,
										Type:     schema.TypeSet,
										Elem:     resourceThreshold(),
									},
									"entity_thresholds": {
										Required: true,
										Type:     schema.TypeSet,
										Elem:     resourceThreshold(),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func threshold(source map[string]interface{}) (interface{}, error) {
	threshold := map[string]interface{}{}
	threshold["baseSeverityColor"] = source["base_severity_color"].(string)
	threshold["baseSeverityColorLight"] = source["base_severity_color_light"].(string)
	threshold["baseSeverityLabel"] = source["base_severity_label"].(string)
	threshold["baseSeverityValue"] = source["base_severity_value"].(int)
	threshold["gaugeMax"] = source["gauge_max"].(int)
	threshold["gaugeMin"] = source["gauge_min"].(int)
	threshold["isMaxStatic"] = source["is_max_static"].(bool)
	threshold["isMinStatic"] = source["is_min_static"].(bool)
	threshold["metricField"] = source["metric_field"].(string)
	threshold["renderBoundaryMax"] = source["render_boundary_max"].(int)
	threshold["renderBoundaryMin"] = source["render_boundary_min"].(int)
	threshold["search"] = source["search"].(string)
	thresholdLevels := []interface{}{}
	for _, h_ := range source["threshold_levels"].(*schema.Set).List() {
		h := h_.(map[string]interface{})
		thresholdLevel := map[string]interface{}{}
		thresholdLevel["dynamicParam"] = h["dynamic_param"].(string)
		thresholdLevel["severityColor"] = h["severity_color"].(string)
		thresholdLevel["severityColorLight"] = h["severity_color_light"].(string)
		thresholdLevel["severityLabel"] = h["severity_label"].(string)
		thresholdLevel["severityValue"] = h["severity_value"].(int)
		thresholdLevel["thresholdValue"] = h["threshold_value"].(int)
		thresholdLevels = append(thresholdLevels, thresholdLevel)
	}
	threshold["thresholdLevels"] = thresholdLevels
	return threshold, nil
}

func kpiThresholdTemplate(d *schema.ResourceData) (config *models.Base, err error) {
	body := map[string]interface{}{}
	body["objectType"] = "kpi_threshold_template"
	body["title"] = d.Get("title").(string)
	body["description"] = d.Get("description").(string)
	body["adaptive_thresholds_is_enabled"] = d.Get("adaptive_thresholds_is_enabled").(bool)
	body["adaptive_thresholding_training_window"] = d.Get("adaptive_thresholding_training_window").(string)
	body["time_variate_thresholds"] = d.Get("time_variate_thresholds").(bool)
	body["sec_grp"] = d.Get("sec_grp").(string)

	time_variate_thresholds_specification := map[string]interface{}{}
	for _, e_ := range d.Get("time_variate_thresholds_specification").(*schema.Set).List() {
		e := e_.(map[string]interface{})
		policies := map[string]interface{}{}
		for _, f_ := range e["policies"].(*schema.Set).List() {
			f := f_.(map[string]interface{})
			policy := map[string]interface{}{}
			policy_name := f["policy_name"].(string)
			policy["title"] = f["title"].(string)
			policy["policy_type"] = f["policy_type"].(string)
			for _, g := range f["aggregate_thresholds"].(*schema.Set).List() {
				aggregate_thresholds, err := threshold(g.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				policy["aggregate_thresholds"] = aggregate_thresholds
			}
			for _, g := range f["entity_thresholds"].(*schema.Set).List() {
				entity_thresholds, err := threshold(g.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				policy["entity_thresholds"] = entity_thresholds
			}
			policies[policy_name] = policy
		}
		time_variate_thresholds_specification["policies"] = policies
	}
	body["time_variate_thresholds_specification"] = time_variate_thresholds_specification
	by, err := json.Marshal(body)
	if err != nil {
		return
	}
	base := kpiThresholdTemplateBase(d.Id(), d.Get("title").(string))
	err = json.Unmarshal(by, &base.RawJson)
	if err != nil {
		return nil, err
	}
	return base, nil
}

func kpiThresholdTemplateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(client)
	template, err := kpiThresholdTemplate(d)
	if err != nil {
		return err
	}
	return template.Create(client.User, client.Password, client.Host, client.Port)
}

func thresholdRead(threshold_data map[string]interface{}) interface{} {
	threshold := map[string]interface{}{}
	threshold["base_severity_color"] = threshold_data["baseSeverityColor"]
	threshold["base_severity_color_light"] = threshold_data["baseSeverityColorLight"]
	threshold["base_severity_label"] = threshold_data["baseSeverityLabel"]
	threshold["base_severity_value"] = threshold_data["baseSeverityValue"]
	threshold["gauge_max"] = threshold_data["gaugeMax"]
	threshold["gauge_min"] = threshold_data["gaugeMin"]
	threshold["is_max_static"] = threshold_data["isMaxStatic"]
	threshold["is_min_static"] = threshold_data["isMinStatic"]
	threshold["metric_field"] = threshold_data["metricField"]
	threshold["render_boundary_max"] = threshold_data["renderBoundaryMax"]
	threshold["render_boundary_min"] = threshold_data["renderBoundaryMin"]
	threshold["search"] = threshold_data["search"]
	thresholdLevels := []interface{}{}
	for _, tdata_ := range threshold_data["thresholdLevels"].([]interface{}) {
		tdata := tdata_.(map[string]interface{})
		thresholdLevel := map[string]interface{}{}
		thresholdLevel["dynamic_param"] = tdata["dynamicParam"]
		thresholdLevel["severity_color"] = tdata["severityColor"]
		thresholdLevel["severity_color_light"] = tdata["severityColorLight"]
		thresholdLevel["severity_label"] = tdata["severityLabel"]
		thresholdLevel["severity_value"] = tdata["severityValue"]
		thresholdLevel["threshold_value"] = tdata["thresholdValue"]
		thresholdLevels = append(thresholdLevels, thresholdLevel)
	}
	threshold["threshold_levels"] = thresholdLevels
	return []interface{}{threshold}
}

func kpiThresholdTemplateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(client)

	base := kpiThresholdTemplateBase(d.Id(), d.Get("title").(string))
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

func populate(b *models.Base, d *schema.ResourceData) error {
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
	err = d.Set("adaptive_thresholds_is_enabled", interfaceMap["adaptive_thresholds_is_enabled"])
	if err != nil {
		return err
	}
	err = d.Set("adaptive_thresholding_training_window", interfaceMap["adaptive_thresholding_training_window"])
	if err != nil {
		return err
	}
	err = d.Set("time_variate_thresholds", interfaceMap["time_variate_thresholds"])
	if err != nil {
		return err
	}
	err = d.Set("sec_grp", interfaceMap["sec_grp"])
	if err != nil {
		return err
	}
	time_variate_thresholds_specification := []map[string]interface{}{}
	policies := []interface{}{}
	time_variate_thresholds_specification_data := interfaceMap["time_variate_thresholds_specification"].(map[string]interface{})
	for policy_name, pdata := range time_variate_thresholds_specification_data["policies"].(map[string]interface{}) {
		policy_data := pdata.(map[string]interface{})
		policy := map[string]interface{}{}
		policy["policy_name"] = policy_name
		policy["title"] = policy_data["title"]
		policy["policy_type"] = policy_data["policy_type"]
		policy["time_blocks"] = policy_data["time_blocks"]
		policy["aggregate_thresholds"] = thresholdRead(policy_data["aggregate_thresholds"].(map[string]interface{}))
		policy["entity_thresholds"] = thresholdRead(policy_data["entity_thresholds"].(map[string]interface{}))
		policies = append(policies, policy)
	}
	policiesMap := map[string]interface{}{
		"policies": policies,
	}
	time_variate_thresholds_specification = append(time_variate_thresholds_specification, policiesMap)
	d.Set("time_variate_thresholds_specification", time_variate_thresholds_specification)
	d.SetId(b.RESTKey)
	return nil
}

func kpiThresholdTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(client)
	base := kpiThresholdTemplateBase(d.Id(), d.Get("title").(string))
	existing, err := base.Find(client.User, client.Password, client.Host, client.Port)
	if err != nil {
		return err
	}
	if existing == nil {
		return kpiThresholdTemplateCreate(d, m)
	}

	template, err := kpiThresholdTemplate(d)
	if err != nil {
		return err
	}
	return template.Update(client.User, client.Password, client.Host, client.Port)
}

func kpiThresholdTemplateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(client)
	base := kpiThresholdTemplateBase(d.Id(), d.Get("title").(string))
	existing, err := base.Find(client.User, client.Password, client.Host, client.Port)
	if err != nil {
		return err
	}
	if existing == nil {
		return nil
	}
	return existing.Delete(client.User, client.Password, client.Host, client.Port)
}

func kpiThresholdTemplateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	client := m.(client)
	b := kpiThresholdTemplateBase(d.Id(), d.Get("title").(string))
	b, err := b.Read(client.User, client.Password, client.Host, client.Port)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, err
	}
	err = populate(b, d)
	if err != nil {
		return nil, err
	}
	if d.Id() == "" {
		return nil, nil
	}
	return []*schema.ResourceData{d}, nil
}
