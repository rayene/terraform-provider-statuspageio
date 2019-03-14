package statuspageio

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	resty "gopkg.in/resty.v1"
)

func resourceStatuspageIOComponentGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceStatuspageIOComponentGroupCreate,
		Read:   resourceStatuspageIOComponentGroupRead,
		Update: resourceStatuspageIOComponentGroupUpdate,
		Delete: resourceStatuspageIOComponentGroupDelete,
		Exists: resourceStatuspageIOComponentGroupExists,
		Importer: &schema.ResourceImporter{
			State: resourceStatuspageIOComponentGroupImport,
		},

		Schema: map[string]*schema.Schema{
			"page_id": {
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
			"components": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

type ComponentGroup struct {
	ID          string   `json:"id,omitempty"`
	PageID      string   `json:"-"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Components  []string `json:"components"`
}

type componentGroupCreateReq struct {
	ComponentGroup ComponentGroup `json:"component_group"`
}

func buildComponentGroupStruct(d *schema.ResourceData) ComponentGroup {
	var cg ComponentGroup

	if attr, ok := d.GetOk("page_id"); ok {
		cg.PageID = attr.(string)
	}

	if attr, ok := d.GetOk("name"); ok {
		cg.Name = attr.(string)
	}

	if attr, ok := d.GetOk("description"); ok {
		cg.Description = attr.(string)
	}

	if attr, ok := d.GetOk("components"); ok {
		for _, c := range attr.(*schema.Set).List() {
			cg.Components = append(cg.Components, c.(string))
		}
	}
	return cg
}

func refreshComponentGroupResource(d *schema.ResourceData, cg ComponentGroup) {
	d.SetId(cg.ID)
	d.Set("name", cg.Name)
	d.Set("description", cg.Description)
	d.Set("components", cg.Components)
}

func resourceStatuspageIOComponentGroupExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := meta.(*resty.Client).R()
	resp, err := client.
		SetPathParams(map[string]string{
			"id":      d.Id(),
			"page_id": d.Get("page_id").(string),
		}).
		SetError(APIError{}).
		Get("pages/{page_id}/component-groups/{id}")

	if err != nil {
		return false, err
	}

	if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			log.Printf("[DEBUG] component group: %s not found - %s", d.Id(), resp.Error())
			return false, nil
		}
		return false, fmt.Errorf("error checking existence of component group: %s - %d %s", d.Id(), resp.StatusCode(), resp.Error())

	}

	log.Printf("[DEBUG] component group: %s found", d.Id())
	return true, nil
}

func resourceStatuspageIOComponentGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*resty.Client).R()

	cg := buildComponentGroupStruct(d)

	resp, err := client.
		SetBody(componentGroupCreateReq{ComponentGroup: cg}).
		SetResult(&cg).
		SetError(APIError{}).
		SetPathParams(map[string]string{
			"page_id": d.Get("page_id").(string),
		}).
		Post("pages/{page_id}/component-groups")

	if err != nil {
		return fmt.Errorf("error creating component group: %s", err.Error())
	}

	if resp.IsError() {
		return fmt.Errorf("error creating component group: %s", resp.Error())
	}

	refreshComponentGroupResource(d, cg)

	return nil
}

func resourceStatuspageIOComponentGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*resty.Client).R()
	cg := ComponentGroup{}
	resp, err := client.
		SetPathParams(map[string]string{
			"id":      d.Id(),
			"page_id": d.Get("page_id").(string),
		}).
		SetResult(&cg).
		SetError(APIError{}).
		Get("pages/{page_id}/component-groups/{id}")

	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("error reading component group: %s - %s", d.Id(), resp.Error())
	}

	log.Printf("[DEBUG] found and read component group: %v", cg)
	refreshComponentGroupResource(d, cg)

	return nil
}

func resourceStatuspageIOComponentGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*resty.Client).R()
	cg := buildComponentGroupStruct(d)

	resp, err := client.
		SetPathParams(map[string]string{
			"id":      d.Id(),
			"page_id": d.Get("page_id").(string),
		}).
		SetBody(componentGroupCreateReq{ComponentGroup: cg}).
		SetResult(&cg).
		SetError(APIError{}).
		Patch("pages/{page_id}/component-groups/{id}")

	if err != nil {
		return fmt.Errorf("error updating component group: %s", err.Error())
	}

	if resp.IsError() {
		return fmt.Errorf("error updating component group: %s - %s", d.Id(), resp.Error())
	}

	refreshComponentGroupResource(d, cg)
	return nil
}

func resourceStatuspageIOComponentGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*resty.Client).R()

	resp, err := client.
		SetPathParams(map[string]string{
			"id":      d.Id(),
			"page_id": d.Get("page_id").(string),
		}).
		SetError(APIError{}).
		Delete("pages/{page_id}/component-groups/{id}")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("error deleting component group: %s - %s", d.Id(), resp.Error())
	}

	return nil
}

func resourceStatuspageIOComponentGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceStatuspageIOComponentGroupRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
