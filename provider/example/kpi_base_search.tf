


resource "itsi_kpi_base_search" "my_kpi_base_search" {
    title = "this_is_a_test2"
    actions = ""
    alert_lag = "31"
    alert_period = "5"
    base_search = "index=_internal source=*metrics.log group=tcpin_connections \n | eval sourceHost=if(isnull(hostname), sourceHost,hostname)"
    description = "This is a description for a KPI base search"
    entity_alias_filtering_fields = null
    entity_breakdown_id_fields = "host"
    entity_id_fields = "host"
    identifying_name = "this_is_a_test2"
    is_entity_breakdown = false
    is_service_entity_filter = false
    is_first_time_save_done = true
    metric_qualifier = ""
    metrics {
        aggregate_statop = "dc"
        entity_statop = "avg"
        fill_gaps = "null_value"
        gap_custom_alert_value = "0"
        gap_severity = "unknown"
        gap_severity_color = "#CCCCCC"
        gap_severity_color_light = "#EEEEEE"
        gap_severity_value = "-1"
        threshold_field = "sourceIp"
        title = "Forwarder Count"
        unit = ""
    }
  search_alert_earliest = "5"
  sec_grp = "default_itsi_security_group"
  source_itsi_da = "itsi"
}

#   description = ""
#   entity_alias_filtering_fields = null
#   entity_breakdown_id_fields = host
#   entity_id_fields = host
#   identifying_name = splkaas-splk-forwarder.inventory
#   is_entity_breakdown = false
#   is_service_entity_filter = false
#   isFirstTimeSaveDone = true
#   metric_qualifier = ""
#   metrics =
#     - _key = 196244dad6e0635a5661fe90
#       aggregate_statop = dc
#       entity_statop = avg
#       fill_gaps = null_value
#       gap_custom_alert_value = "0"
#       gap_severity = unknown
#       gap_severity_color = '#CCCCCC'
#       gap_severity_color_light = '#EEEEEE'
#       gap_severity_value = "-1"
#       threshold_field = sourceIp
#       title = Forwarder Count
#       unit = ""
#   mod_source = REST
#   mod_time = ""
#   mod_timestamp = "2019-10-29T18 =26 =11.966043+00 =00"
#   object_type = kpi_base_search
#   permissions =
#     delete = false
#     group =
#         delete = false
#         read = true
#         write = false
#     read = true
#     user = rest-user
#     write = false
#   search_alert_earliest = "5"
#   sec_grp = default_itsi_security_group
#   source_itsi_da = itsi
#   title = SPLKaaS-SPLK-Forwarder.Inventory