package models

func DumpServices(user, password, host string, port int) error {
	auditList := []string{
		"services_depends_on",
		"base_service_template_id",
		"tags",
		"enabled",
		"identifying_name",
		"services_depending_on_me",
		"title",
		"backfill_enabled",
		"description",
		"kpis",
	}

	base := NewBase("", "itoa_interface", "service", auditList)
	items, err := base.Dump(user, password, host, port)
	if err != nil {
		return err
	}
	err = base.auditLog(items)
	if err != nil {
		return err
	}
	return base.auditFields(items)
}
