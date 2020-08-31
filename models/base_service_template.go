package models

func DumpBaseServiceTemplates(user, password, host string, port int) error {
	auditList := []string{
		"description",
		"identifying_name",
		"kpis",
	}
	base := NewBase("", "itoa_interface", "base_service_template", auditList)
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
