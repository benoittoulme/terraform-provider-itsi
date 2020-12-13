package main

import (
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/benoittoulme/terraform-provider-itsi/models"
)

func main() {
	parser := argparse.NewParser("ITSI scraper", "Dump ITSI resources via REST interface and format them in a file.")

	user := parser.String("u", "user", &argparse.Options{Required: false, Help: "user", Default: "admin"})
	password := parser.String("p", "password", &argparse.Options{Required: false, Help: "password", Default: "changeme"})
	host := parser.String("s", "host", &argparse.Options{Required: false, Help: "host", Default: "localhost"})
	port := parser.Int("o", "port", &argparse.Options{Required: false, Help: "port", Default: 8089})
	verbose := parser.Selector("v", "verbose", []string{"true", "false"}, &argparse.Options{Required: false, Help: "verbose mode", Default: "false"})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(parser.Usage(err))
	}

	models.Verbose = (*verbose == "true")
	for k, _ := range models.RestConfigs {
		err := dump(*user, *password, *host, *port, k)
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
