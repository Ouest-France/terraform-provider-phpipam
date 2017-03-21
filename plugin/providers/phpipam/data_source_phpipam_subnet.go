package phpipam

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/paybyphone/phpipam-sdk-go/controllers/subnets"
)

func dataSourcePHPIPAMSubnet() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePHPIPAMSubnetRead,
		Schema: dataSourceSubnetSchema(),
	}
}

func dataSourcePHPIPAMSubnetRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	var out subnets.Subnet
	// We need to determine how to get the subnet. An ID search takes priority,
	// and after that subnets.
	switch {
	case d.Get("subnet_id").(int) != 0:
		var err error
		out, err = c.GetSubnetByID(d.Get("subnet_id").(int))
		if err != nil {
			return err
		}
	case d.Get("subnet_address").(string) != "" && d.Get("subnet_mask").(int) != 0:
		v, err := c.GetSubnetsByCIDR(fmt.Sprintf("%s/%d", d.Get("subnet_address"), d.Get("subnet_mask")))
		if err != nil {
			return err
		}
		// This should not happen, as a CIDR search from what I have seen so far
		// generally only returns 1 subnet ever. Nontheless, the API spec and the
		// function return a slice of subnets, so we need to assert that the search
		// has only retuned a single match.
		if len(v) != 1 {
			return errors.New("CIDR search returned either zero or multiple results. Please correct your search and try again")
		}
		out = v[0]
	default:
		return errors.New("subnet_id or subnet_address not defined, cannot proceed with reading resource data")
	}
	flattenSubnet(out, d)
	return nil
}