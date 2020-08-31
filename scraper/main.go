package main

import (
	"github.com/benoittoulme/terraform-provider-itsi/models"
)

var USER string = "admin"
var PASSWORD string = "changeme"
var HOST string = "localhost"
var PORT int = 18089

func main() {
	err := models.DumpBaseServiceTemplates(USER, PASSWORD, HOST, PORT)
	if err != nil {
		panic(err)
	}
	err = models.DumpCorrelationSearches(USER, PASSWORD, HOST, PORT)
	if err != nil {
		panic(err)
	}
	err = models.DumpEntities(USER, PASSWORD, HOST, PORT)
	if err != nil {
		panic(err)
	}
	err = models.DumpGlassTables(USER, PASSWORD, HOST, PORT)
	if err != nil {
		panic(err)
	}
	err = models.DumpKPIBaseSearches(USER, PASSWORD, HOST, PORT)
	if err != nil {
		panic(err)
	}
	err = models.DumpKPITemplates(USER, PASSWORD, HOST, PORT)
	if err != nil {
		panic(err)
	}
	err = models.DumpKPIThresholdTemplates(USER, PASSWORD, HOST, PORT)
	if err != nil {
		panic(err)
	}
	err = models.DumpServices(USER, PASSWORD, HOST, PORT)
	if err != nil {
		panic(err)
	}
}
