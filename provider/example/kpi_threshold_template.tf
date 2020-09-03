resource "itsi_kpi_threshold_template" "my_kpi_threshold_template" {
  title                                 = "BLAH 3-hour blocks work week (adaptive/stdev)"
  description                           = "Work week 3hr. Monday to Friday 3 hour chunks. Sat Sun 1 chunk"
  adaptive_thresholds_is_enabled        = true
  adaptive_thresholding_training_window = "-7d"
  time_variate_thresholds               = true
  sec_grp                               = "default_itsi_security_group"

  time_variate_thresholds_specification {
    policies {
      policy_name = "default_policy"
      title = "Default"
      policy_type = "static"
      aggregate_thresholds {
                base_severity_color = "#B50101"
                base_severity_color_light = "#E5A6A6"
                base_severity_label = "critical"
                base_severity_value = 6
                gauge_max = 100
                gauge_min = 0
                is_max_static = false
                is_min_static = true
                metric_field = ""
                render_boundary_max = 100
                render_boundary_min = 0
                search = ""
                threshold_levels {
                    dynamic_param = -2
                    severity_color = "#FCB64E"
                    severity_color_light = "#FEE6C1"
                    severity_label = "medium"
                    severity_value = 4
                    threshold_value = 20
                }
                threshold_levels {
                    dynamic_param = -1
                    severity_color = "#99D18B"
                    severity_color_light = "#DCEFD7"
                    severity_label = "normal"
                    severity_value = 2
                    threshold_value = 40
                }
      }

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
  }
}