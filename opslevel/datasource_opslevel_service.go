package opslevel

import (
	"github.com/opslevel/opslevel-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceService() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceServiceRead),
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:        schema.TypeString,
				Description: "An alias of the service to find by.",
				ForceNew:    true,
				Optional:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The id of the service to find.",
				ForceNew:    true,
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the service.",
				Computed:    true,
			},
			"product": {
				Type:        schema.TypeString,
				Description: "A product is an application that your end user interacts with. Multiple services can work together to power a single product.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A brief description of the service.",
				Computed:    true,
			},
			"language": {
				Type:        schema.TypeString,
				Description: "The primary programming language that the service is written in.",
				Computed:    true,
			},
			"framework": {
				Type:        schema.TypeString,
				Description: "The primary software development framework that the service uses.",
				Computed:    true,
			},
			"tier_alias": {
				Type:        schema.TypeString,
				Description: "The software tier that the service belongs to.",
				Computed:    true,
			},
			"owner_alias": {
				Type:        schema.TypeString,
				Description: "The team that owns the service.",
				Computed:    true,
			},
			"lifecycle_alias": {
				Type:        schema.TypeString,
				Description: "The lifecycle stage of the service.",
				Computed:    true,
			},
			"aliases": {
				Type:        schema.TypeList,
				Description: "A list of human-friendly, unique identifiers for the service",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "A list of tags applied to the service.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func datasourceServiceRead(d *schema.ResourceData, client *opslevel.Client) error {
	resource, err := findService("alias", "id", d, client)
	if err != nil {
		return err
	}

	d.SetId(resource.Id.(string))
	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("product", resource.Product); err != nil {
		return err
	}
	if err := d.Set("description", resource.Description); err != nil {
		return err
	}
	if err := d.Set("language", resource.Language); err != nil {
		return err
	}
	if err := d.Set("framework", resource.Framework); err != nil {
		return err
	}
	if err := d.Set("tier_alias", resource.Tier.Alias); err != nil {
		return err
	}
	if err := d.Set("owner_alias", resource.Owner.Alias); err != nil {
		return err
	}
	if err := d.Set("lifecycle_alias", resource.Lifecycle.Alias); err != nil {
		return err
	}

	if err := d.Set("aliases", resource.Aliases); err != nil {
		return err
	}
	if err := d.Set("tags", flattenTagArray(resource.Tags.Nodes)); err != nil {
		return err
	}

	return nil
}
