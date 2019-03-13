package statuspageio

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	resty "gopkg.in/resty.v1"
)

func resourceStatuspageIOComponent() *schema.Resource {
	return &schema.Resource{
		Create: resourceStatuspageIOComponentCreate,
		Read:   resourceStatuspageIOComponentRead,
		Update: resourceStatuspageIOComponentUpdate,
		Delete: resourceStatuspageIOComponentDelete,
		Exists: resourceStatuspageIOComponentExists,
		Importer: &schema.ResourceImporter{
			State: resourceStatuspageIOComponentImport,
		},

		Schema: map[string]*schema.Schema{

			"page": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "operational",
			},

			"showcase": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"only_show_if_degraded": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type Component struct {
	ID                 string `json:"id,omitempty"`
	Page               string `json:"-"`
	Name               string `json:"name"`
	Description        string `json:"description,omitempty"`
	Status             string `json:"status"`
	Showcase           bool   `json:"showcase"`
	OnlyShowIfDegraded bool   `json:"only_show_if_degraded"`
	GroupID            string `json:"group_id,omitempty"`
}

type componentCreateReq struct {
	Component Component `json:"component"`
}

func buildComponentStruct(d *schema.ResourceData) Component {
	var component Component

	if attr, ok := d.GetOk("page"); ok {
		component.Page = attr.(string)
	}

	if attr, ok := d.GetOk("name"); ok {
		component.Name = attr.(string)
	}

	if attr, ok := d.GetOk("description"); ok {
		component.Description = attr.(string)
	}

	if attr, ok := d.GetOk("status"); ok {
		component.Status = attr.(string)
	}

	if attr, ok := d.GetOk("showcase"); ok {
		component.Showcase = attr.(bool)
	}

	if attr, ok := d.GetOk("only_show_if_degraded"); ok {
		component.OnlyShowIfDegraded = attr.(bool)
	}

	if attr, ok := d.GetOk("group_id"); ok {
		component.GroupID = attr.(string)
	}

	return component
}

func refreshComponentResource(d *schema.ResourceData, component Component) {
	d.SetId(component.ID)
	d.Set("name", component.Name)
	d.Set("description", component.Description)
	d.Set("status", component.Status)
	d.Set("group_id", component.GroupID)
	d.Set("showcase", component.Showcase)
	d.Set("only_show_if_degraded", component.OnlyShowIfDegraded)
}

func resourceStatuspageIOComponentExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := meta.(*resty.Client).R()
	resp, err := client.
		SetPathParams(map[string]string{
			"component_id": d.Id(),
			"page":         d.Get("page").(string),
		}).
		SetError(APIError{}).
		Get("pages/{page}/components/{component_id}")

	if err != nil {
		return false, err
	}

	if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			log.Printf("[DEBUG] component: %s not found - %s", d.Id(), resp.Error())
			return false, nil
		}
		return false, fmt.Errorf("error checking existence of component: %s - %s %s", d.Id(), resp.StatusCode(), resp.Error())

	}

	log.Printf("[DEBUG] component: %s found", d.Id())
	return true, nil
}

func resourceStatuspageIOComponentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*resty.Client).R()

	component := buildComponentStruct(d)

	resp, err := client.
		SetBody(componentCreateReq{Component: component}).
		SetResult(&component).
		SetError(APIError{}).
		SetPathParams(map[string]string{
			"page": d.Get("page").(string),
		}).
		Post("pages/{page}/components")

	if err != nil {
		return fmt.Errorf("error creating component: %s", err.Error())
	}

	if resp.IsError() {
		return fmt.Errorf("error creating component: %s", resp.Error())
	}

	refreshComponentResource(d, component)

	return nil
}

func resourceStatuspageIOComponentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*resty.Client).R()
	component := Component{}
	resp, err := client.
		SetPathParams(map[string]string{
			"component_id": d.Id(),
			"page":         d.Get("page").(string),
		}).
		SetResult(&component).
		SetError(APIError{}).
		Get("pages/{page}/components/{component_id}")

	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("error reading component: %s - %s", d.Id(), resp.Error())
	}

	log.Printf("[DEBUG] found and read component: %v", component)
	refreshComponentResource(d, component)

	return nil
}

func resourceStatuspageIOComponentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*resty.Client).R()
	component := buildComponentStruct(d)

	resp, err := client.
		SetPathParams(map[string]string{
			"component_id": d.Id(),
			"page":         d.Get("page").(string),
		}).
		SetBody(componentCreateReq{Component: component}).
		SetResult(&component).
		SetError(APIError{}).
		Patch("pages/{page}/components/{component_id}")

	if err != nil {
		return fmt.Errorf("error updating component: %s", err.Error())
	}

	if resp.IsError() {
		return fmt.Errorf("error updating component: %s - %s", d.Id(), resp.Error())
	}

	refreshComponentResource(d, component)
	return nil
}

func resourceStatuspageIOComponentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*resty.Client).R()

	resp, err := client.
		SetPathParams(map[string]string{
			"component_id": d.Id(),
			"page":         d.Get("page").(string),
		}).
		SetError(APIError{}).
		Delete("pages/{page}/components/{component_id}")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("error deleting component: %s - %s", d.Id(), resp.Error())
	}

	return nil
}

func resourceStatuspageIOComponentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceStatuspageIOComponentRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
