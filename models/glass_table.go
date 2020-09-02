package models

func DumpGlassTables(user, password, host string, port int) error {
	auditList := []string{}
	base := NewBase("", "", "itoa_interface", "glass_table")
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
