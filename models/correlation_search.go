package models

type CorrelationSearch struct {
	*Base
}

func DumpCorrelationSearches(user, password, host string, port int) error {
	auditList := []string{
		// "cron_schedule",
		// "description",
		// "disabled",
		// "is_scheduled",
		// "is_visible",
		// "max_concurrent",
		// "name",
		// "search",
	}
	b := NewBase("", "", "event_management_interface", "correlation_search")
	b.RESTKeyField = func() string {
		return "name"
	}
	b.TFIDField = func() string {
		return "name"
	}
	items, err := b.Dump(user, password, host, port)
	if err != nil {
		return err
	}
	err = b.auditLog(items, auditList)
	if err != nil {
		return err
	}
	return b.auditFields(items)
}
