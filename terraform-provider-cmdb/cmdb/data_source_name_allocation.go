package cmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// initNameAllocationSchema is where we define the schema of the Terraform data source
func initNameAllocationSchema() *schema.Resource {
	return &schema.Resource{
		Read: initNameDataSourceRead,

		Schema: map[string]*schema.Schema{
			"raw": {
				Type:     schema.TypeString,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// initNameDataSourceRead tells Terraform how to contact our microservice and retrieve the necessary data
func initNameDataSourceRead(d *schema.ResourceData, meta interface{}) (err error) {
	provider := meta.(ProviderClient)
	client := provider.Client

	header := make(http.Header)
	headers, exists := d.GetOk("headers")
	if exists {
		for name, value := range headers.(map[string]interface{}) {
			header.Set(name, value.(string))
		}
	}

	resourceType := d.Get("resource_type").(string)
	if resourceType == "" {
		return fmt.Errorf("Invalid resource type specified")
	}
	region := d.Get("region").(string)
	if region == "" {
		return fmt.Errorf("Invalid region specified")
	}
	b, err := client.doAllocateName(client.BaseUrl.String(), resourceType, region)
	if err != nil {
		return
	}
	outputs, err := flattenNameAllocationResponse(b)
	if err != nil {
		return
	}
	marshalData(d, outputs)

	return
}

func flattenNameAllocationResponse(b []byte) (outputs map[string]interface{}, err error) {
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		err = fmt.Errorf("Cannot unmarshal json of API response: %v", err)
		return
	} else if data["result"] == "" {
		err = fmt.Errorf("missing result key in API response: %v", err)
		return
	}

	outputs = make(map[string]interface{})
	outputs["id"] = time.Now().UTC().String()
	outputs["raw"] = string(b)
	outputs["name"] = data["Name"]

	return
}