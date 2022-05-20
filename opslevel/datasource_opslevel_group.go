package opslevel

import (
	"github.com/opslevel/opslevel-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: wrap(datasourceGroupRead),
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Description: "The id or alias of the group to find.",
				ForceNew:    true,
				Optional:    true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"teams": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func datasourceGroupRead(d *schema.ResourceData, client *opslevel.Client) error {
	identifier := d.Get("identifier").(string)
	var resource *opslevel.Group
	var err error
	if opslevel.IsID(identifier) {
		resource, err = client.GetGroup(identifier)
	} else {
		resource, err = client.GetGroupWithAlias(identifier)
	}
	if err != nil {
		return err
	}

	var parent string
	if resource.Parent.Alias != "" {
		parent = resource.Parent.Alias
	}

	members, err := resource.Members(client)
	if err != nil {
		return err
	}
	members_list := []string{}
	if len(members) > 0 {
		for _, member := range members {
			members_list = append(members_list, member.Email)
		}
	}

	childTeams, err := resource.ChildTeams(client)
	if err != nil {
		return err
	}
	teams := []string{}
	if len(childTeams) > 0 {
		for _, team := range childTeams {
			teams = append(teams, team.Alias)
		}
	}

	d.SetId(resource.Id.(string))
	d.Set("name", resource.Name)
	d.Set("description", resource.Description)
	d.Set("parent", parent)
	d.Set("members", members_list)
	d.Set("teams", teams)

	return nil
}
