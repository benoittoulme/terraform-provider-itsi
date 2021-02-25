package models

import (
	"bytes"
	"encoding/json"
	"html/template"
	"regexp"
	"strings"
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
