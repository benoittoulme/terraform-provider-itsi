package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type client struct {
	User     string
	Password string
	Host     string
	Port     int
}

// Provider returns the scooby provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"itsi_kpi_threshold_template": resourceKPIThresholdTemplate(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := client{}
	client.User = d.Get("user").(string)
	client.Password = d.Get("password").(string)
	client.Host = d.Get("host").(string)
	client.Port = d.Get("port").(int)
	return client, nil
}
