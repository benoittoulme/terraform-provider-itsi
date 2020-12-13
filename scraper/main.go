package main

import (
	"github.com/benoittoulme/terraform-provider-itsi/models"
)

var USER string = "admin"
var PASSWORD string = "changeme"
var HOST string = "localhost"
var PORT int = 18089

func main() {
	models.Verbose = true
	for k, _ := range models.RestConfigs {
		err := dump(USER, PASSWORD, HOST, PORT, k)
		if err != nil {
			panic(err)
		}
	}
}

func dump(user, password, host string, port int, objectType string) error {
	base := models.NewBase("", "", objectType)
	items, err := base.Dump(user, password, host, port)
	if err != nil {
		return err
	}
	err = base.AuditLog(items, nil)
	if err != nil {
		return err
	}
	return base.AuditFields(items)
}
