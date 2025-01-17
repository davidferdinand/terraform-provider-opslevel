package opslevel

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go"
)

// Tag key names are stored in OpsLevel as lowercase so we need to ensure the configuration is written as lowercase
var TagKeyRegex = regexp.MustCompile(`\A[a-z][0-9a-z_\.\/\\-]*\z`)
var TagKeyErrorMsg = "tag key name must start with a letter and be only lowercase alphanumerics, underscores, hyphens, periods, and slashes."

func resourceServiceTag() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a service tag",
		Create:      wrap(resourceServiceTagCreate),
		Read:        wrap(resourceServiceTagRead),
		Update:      wrap(resourceServiceTagUpdate),
		Delete:      wrap(resourceServiceTagDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service": {
				Type:        schema.TypeString,
				Description: "The id of the service that this will be added to.",
				ForceNew:    true,
				Optional:    true,
			},
			"service_alias": {
				Type:        schema.TypeString,
				Description: "The alias of the service that this will be added to.",
				ForceNew:    true,
				Optional:    true,
			},
			"key": {
				Type:         schema.TypeString,
				Description:  "The tag's key.",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringMatch(TagKeyRegex, TagKeyErrorMsg),
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The tag's value.",
				ForceNew:    false,
				Required:    true,
			},
		},
	}
}

func resourceServiceTagCreate(d *schema.ResourceData, client *opslevel.Client) error {
	service, err := findService("service_alias", "service", d, client)
	if err != nil {
		return err
	}

	input := opslevel.TagCreateInput{
		Id: service.Id,

		Key:   d.Get("key").(string),
		Value: d.Get("value").(string),
	}
	resource, err := client.CreateTag(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	if err := d.Set("key", resource.Key); err != nil {
		return err
	}
	if err := d.Set("value", resource.Value); err != nil {
		return err
	}

	return nil
}

func resourceServiceTagRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	// Handle Import by spliting the ID into the 2 parts
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		d.Set("service", parts[0])
		id = parts[1]
		d.SetId(id)
	}

	service, err := findService("service_alias", "service", d, client)
	if err != nil {
		return err
	}

	var resource *opslevel.Tag
	for _, t := range service.Tags.Nodes {
		if t.Id == id {
			resource = &t
			break
		}
	}
	if resource == nil {
		return fmt.Errorf("unable to find tag with id '%s' on service '%s'", id, service.Aliases[0])
	}

	if err := d.Set("key", resource.Key); err != nil {
		return err
	}
	if err := d.Set("value", resource.Value); err != nil {
		return err
	}

	return nil
}

func resourceServiceTagUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.TagUpdateInput{
		Id: d.Id(),
	}

	if d.HasChange("key") {
		input.Key = d.Get("key").(string)
	}
	if d.HasChange("value") {
		input.Value = d.Get("value").(string)
	}

	resource, err := client.UpdateTag(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())

	if err := d.Set("key", resource.Key); err != nil {
		return err
	}
	if err := d.Set("value", resource.Value); err != nil {
		return err
	}
	return nil
}

func resourceServiceTagDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteTag(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
