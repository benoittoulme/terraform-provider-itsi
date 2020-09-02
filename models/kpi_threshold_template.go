package models

func DumpKPIThresholdTemplates(user, password, host string, port int) error {
	auditList := []string{
		// "adaptive_thresholding_training_window",
		// "adaptive_thresholds_is_enabled",
		// "description",
		// "identifying_name",
		// "time_variate_thresholds",
		// "time_variate_thresholds_specification",
		// "title",
	}
	base := NewBase("", "", "itoa_interface", "kpi_threshold_template")
	base.TFIDField = func() string {
		return "title"
	}
	items, err := base.Dump(user, password, host, port)
	if err != nil {
		return err
	}
	err = base.auditLog(items, auditList)
	if err != nil {
		return err
	}
	return base.auditFields(items)
}
