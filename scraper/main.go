package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/benoittoulme/terraform-provider-itsi/models"
)

func main() {
	parser := argparse.NewParser("ITSI scraper", "Dump ITSI resources via REST interface and format them in a file.")

	user := parser.String("u", "user", &argparse.Options{Required: false, Help: "user", Default: "admin"})
	password := parser.String("p", "password", &argparse.Options{Required: false, Help: "password", Default: "changeme"})
	host := parser.String("t", "host", &argparse.Options{Required: false, Help: "host", Default: "localhost"})
	port := parser.Int("o", "port", &argparse.Options{Required: false, Help: "port", Default: 8089})
	verbose := parser.Selector("v", "verbose", []string{"true", "false"}, &argparse.Options{Required: false, Help: "verbose mode", Default: "false"})
	skipTLS := parser.Selector("s", "skip-tls", []string{"true", "false"}, &argparse.Options{Required: false, Help: "skip TLS check", Default: "false"})
	format := parser.Selector("f", "format", []string{"json", "yaml"}, &argparse.Options{Required: false, Help: "output format. json|yaml", Default: "yaml"})

	objectTypes := []string{}
	for k, _ := range models.RestConfigs {
		objectTypes = append(objectTypes, k)
	}
	objs := parser.StringList("b", "obj", &argparse.Options{Required: false, Help: "object types", Default: objectTypes})
	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(parser.Usage(err))
	}

	models.Verbose = (*verbose == "true")
	models.SkipTLS = (*skipTLS == "true")
	for _, k := range *objs {
		fmt.Printf("scraping %s...\n", k)
		err := dump(*user, *password, *host, *port, k, *format)
		if err != nil {
			panic(err)
		}
	}
}

func dump(user, password, host string, port int, objectType, format string) error {
	base := models.NewBase("", "", objectType)
	items, err := base.Dump(user, password, host, port)
	if err != nil {
		return err
	}
	err = base.AuditLog(items, nil, format)
	if err != nil {
		return err
	}
	return base.AuditFields(items)
}
