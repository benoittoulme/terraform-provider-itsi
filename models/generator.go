package models

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
	"text/template"
)

func NilTemplateMarshal(b *Base) (string, error) {
	return "", nil
}

func tfname(name string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	name = reg.ReplaceAllString(name, "-")
	if strings.HasSuffix(name, "-") {
		name = name[:len(name)-1]
	}
	if strings.HasPrefix(name, "-") {
		name = name[1:]
	}
	return name, nil
}

func KPIBaseSearchMarshal(b *Base) (string, error) {
	format := `
	resource "itsi_kpi_base_search" "{{.tfname}}" {
		title                         = "{{.title}}"
		actions                       = "{{.actions}}"
		alert_lag                     = "{{.alert_lag}}"
		alert_period                  = "{{.alert_period}}"
		base_search                   = <<-EOF
		{{.base_search}}
		EOF
		description                   = "{{.description}}"
		{{if .entity_alias_filtering_fields}}entity_alias_filtering_fields =  {{.entity_alias_filtering_fields}}{{end}}
		entity_breakdown_id_fields    = "{{.entity_breakdown_id_fields}}"
		entity_id_fields              = "{{.entity_id_fields}}"
		identifying_name              = "{{.identifying_name}}"
		is_entity_breakdown           =  {{.is_entity_breakdown}}
		is_service_entity_filter      =  {{.is_service_entity_filter}}
		{{if .is_first_time_save_done}}is_first_time_save_done       =  {{.is_first_time_save_done}}{{end}}
		metric_qualifier              = "{{.metric_qualifier}}"
		{{range $index, $element := .metrics}}{{ with $element }}
		metrics {
			aggregate_statop         = "{{.aggregate_statop}}"
			entity_statop            = "{{.entity_statop}}"
			fill_gaps                = "{{.fill_gaps}}"
			gap_custom_alert_value   = "{{.gap_custom_alert_value}}"
			gap_severity             = "{{.gap_severity}}"
			gap_severity_color       = "{{.gap_severity_color}}"
			gap_severity_color_light = "{{.gap_severity_color_light}}"
			gap_severity_value       = "{{.gap_severity_value}}"
			threshold_field          = "{{.threshold_field}}"
			title                    = "{{.title}}"
			unit                     = "{{.unit}}"
		}
		{{end}}{{end}}
	  search_alert_earliest = "{{.search_alert_earliest}}"
	  sec_grp               = "{{.sec_grp}}"
	  source_itsi_da        = "{{.source_itsi_da}}"
	}
	
`
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
	tfName, err := tfname(m["title"].(string))
	if err != nil {
		return "", err
	}
	m["tfname"] = tfName

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, m)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
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
				base_severity_color       = "{{.baseSeverityColor}}"
				base_severity_color_light = "{{.baseSeverityColorLight}}"
				base_severity_label       = "{{.baseSeverityLabel}}"
				base_severity_value       = {{.baseSeverityValue}}
				gauge_max                 = {{.gaugeMax}}
				gauge_min                 = {{.gaugeMin}}
				is_max_static             = {{.isMaxStatic}}
				is_min_static             = {{.isMinStatic}}
				metric_field              = "{{.isMinStatic}}"
				render_boundary_max       = {{.renderBoundaryMax}}
				render_boundary_min       = {{.renderBoundaryMax}}
				search                    = "{{.search}}"
				{{range .thresholdLevels}}
				threshold_levels {
					dynamic_param        = {{.dynamicParam}}
					severity_color       = "{{.severityColor}}"
					severity_color_light = "{{.severityColor}}"
					severity_label       = "{{.severityColor}}"
					severity_value       = {{.severityValue}}
					threshold_value      = {{.thresholdValue}}
				} {{end}}
			} {{end}}

			entity_thresholds { {{with .entity_thresholds}}
				base_severity_color       = "{{.baseSeverityColor}}"
				base_severity_color_light = "{{.baseSeverityColorLight}}"
				base_severity_label       = "{{.baseSeverityLabel}}"
				base_severity_value       = {{.baseSeverityValue}}
				gauge_max                 = {{.gaugeMax}}
				gauge_min                 = {{.gaugeMin}}
				is_max_static             = {{.isMaxStatic}}
				is_min_static             = {{.isMinStatic}}
				metric_field              = "{{.isMinStatic}}"
				render_boundary_max       = {{.renderBoundaryMax}}
				render_boundary_min       = {{.renderBoundaryMax}}
				search                    = "{{.search}}"
				{{range .thresholdLevels}}
				threshold_levels {
					dynamic_param        = {{.dynamicParam}}
					severity_color       = "{{.severityColor}}"
					severity_color_light = "{{.severityColor}}"
					severity_label       = "{{.severityColor}}"
					severity_value       = {{.severityValue}}
					threshold_value      = {{.thresholdValue}}
				} {{end}}
			} {{end}}
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
	tfName, err := tfname(m["title"].(string))
	if err != nil {
		return "", err
	}
	m["tfname"] = tfName

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, m)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
