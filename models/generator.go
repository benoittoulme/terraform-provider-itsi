package models

import (
	"bytes"
	"encoding/json"
	"html/template"
	"strings"
)

func NilTemplateMarshal(b *Base) (string, error) {
	return "", nil
}

func ThresholdTemplateMarshal(b *Base) (string, error) {
	format := `
	resource "itsi_kpi_threshold_template" "{{.tfname}}" {
		title                                 = "{{.title}}"
		description                           = "{{.description}}"
		adaptive_thresholds_is_enabled        = {{.adaptive_thresholds_is_enabled}}
		adaptive_thresholding_training_window = "{{.adaptive_thresholding_training_window}}"
		time_variate_thresholds               = {{.time_variate_thresholds}}
		sec_grp                               = "{{.sec_grp}}"
	  
		time_variate_thresholds_specification {
		{{range $index, $element := .time_variate_thresholds_specification.policies}}{{ with $element }}
		  policies {
			policy_name = "{{$index}}"
			title = "{{.title}}"
			policy_type = "{{.policy_type}}"
			aggregate_thresholds { {{with .aggregate_thresholds}}
					  base_severity_color = "{{.baseSeverityColor}}"
					  base_severity_color_light = "{{.baseSeverityColorLight}}"
					  base_severity_label = "{{.baseSeverityLabel}}"
					  base_severity_value = {{.baseSeverityValue}}
					  gauge_max = {{.gaugeMax}}
					  gauge_min = {{.gaugeMin}}
					  is_max_static = {{.isMaxStatic}}
					  is_min_static = {{.isMinStatic}}
					  metric_field = "{{.isMinStatic}}"
					  render_boundary_max = {{.renderBoundaryMax}}
					  render_boundary_min = {{.renderBoundaryMax}}
					  search = "{{.search}}"
					  {{range .thresholdLevels}}
					  threshold_levels {
						  dynamic_param = {{.dynamicParam}}
						  severity_color = "{{.severityColor}}"
						  severity_color_light = "{{.severityColor}}"
						  severity_label = "{{.severityColor}}"
						  severity_value = {{.severityValue}}
						  threshold_value = {{.thresholdValue}}
					  }
					  {{end}}
			} {{end}}
	  
			entity_thresholds {
				base_severity_color = "#99D18B"
				base_severity_color_light = "#DCEFD7"
				base_severity_label = "normal"
				base_severity_value = 2
				gauge_max = 100
				gauge_min = 0
				is_max_static = false
				is_min_static = true
				metric_field = ""
				render_boundary_max = 100
				render_boundary_min = 0
				search = ""
			}   
		  }
		  {{end}}{{end}}
		}
	  }`

	tmpl, err := template.New("test").Parse(format)
	if err != nil {
		return "", err
	}

	m := map[string]interface{}{}
	by, err := b.RawJson.MarshalJSON()
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(by, &m); err != nil {
		return "", err
	}
	m["tfname"] = strings.ReplaceAll(m["title"].(string), " ", "-")

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, m)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
