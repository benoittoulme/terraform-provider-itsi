package models

func DumpEntities(user, password, host string, port int) error {
	auditList := []string{
		"title",
		"services",
		"entity",
		"mod_timestamp",
		"description",
		"informational",
		"identifying_name",
	}
	base := NewBase("", "itoa_interface", "entity", auditList)
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
