package models

func DumpKPIBaseSearches(user, password, host string, port int) error {
	auditList := []string{}
	base := NewBase("", "", "itoa_interface", "kpi_base_search")
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
