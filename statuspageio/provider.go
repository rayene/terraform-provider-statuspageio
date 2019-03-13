package statuspageio

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	resty "gopkg.in/resty.v1"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:     schema.TypeString,
				Required: true,
				//DefaultFunc: schema.EnvDefaultFunc("STATUSPAGEIO_API_KEY", nil),
			},
			"api_url": {
				Type:     schema.TypeString,
				Required: true,
				//DefaultFunc: schema.EnvDefaultFunc("STATUSPAGEIO_API_URL", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"statuspageio_component":       resourceStatuspageIOComponent(),
			"statuspageio_component_group": resourceStatuspageIOComponentGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := resty.
		SetHostURL(d.Get("api_url").(string)).
		SetHeader("Authorization", "OAuth "+d.Get("api_key").(string)).
		SetDebug(true).
		SetRetryCount(10).
		SetRetryWaitTime(30 * time.Second).
		SetRetryMaxWaitTime(120 * time.Second).
		AddRetryCondition(
			func(r *resty.Response) (bool, error) {
				return r.StatusCode() == 420, nil
			},
		)

	log.Println("[INFO] Statuspage.io client successfully initialized, now validating...")
	resp, err := client.R().
		Get("pages")

	if err != nil {
		log.Printf("[ERROR] Statuspage.io  Client validation error: %v", err)
		return client, err
	}
	if resp.IsError() {
		return client, fmt.Errorf("error reading component: %s - %s", d.Id(), resp.Error())
	}

	log.Printf("[INFO] Statuspage.io Client successfully validated.")

	return client, nil
}

type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
