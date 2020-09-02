package models

type CorrelationSearch struct {
	*Base
}

func DumpCorrelationSearches(user, password, host string, port int) error {
	auditList := []string{
		"cron_schedule",
		"description",
		"disabled",
		"is_scheduled",
		"is_visible",
		"max_concurrent",
		"name",
		"search",
	}
	c := CorrelationSearch{
		Base: NewBase("", "", "event_management_interface", "correlation_search"),
	}
	c.RESTKeyField = func() string {
		return "name"
	}
	items, err := c.Dump(user, password, host, port)
	if err != nil {
		return err
	}
	err = c.auditLog(items, auditList)
	if err != nil {
		return err
	}
	return c.auditFields(items)
}
