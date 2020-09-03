package models

func DumpKPITemplates(user, password, host string, port int) error {
	auditList := []string{}

	base := NewBase("", "", "itoa_interface", "kpi_template")
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