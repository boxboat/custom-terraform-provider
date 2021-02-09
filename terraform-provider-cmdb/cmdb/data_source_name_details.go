package cmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// getDetailsForNameSchema is where we define the schema of the Terraform data source
func getDetailsForNameSchema() *schema.Resource {
	return &schema.Resource{
		Read: getDetailsDataSourceRead,

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
				Computed: false,
				Required: true,
				Elem: schema.TypeString,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Elem: schema.TypeString,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
				Elem: schema.TypeString,
			},
		},
	}
}

// dataSourceRead tells Terraform how to contact our microservice and retrieve the necessary data
func getDetailsDataSourceRead(d *schema.ResourceData, meta interface{}) (err error) {
	provider := meta.(ProviderClient)
	client := provider.Client

	header := make(http.Header)
	headers, exists := d.GetOk("headers")
	if exists {
		for name, value := range headers.(map[string]interface{}) {
			header.Set(name, value.(string))
		}
	}

	resourceName := d.Get("name").(string)
	if resourceName == "" {
		return fmt.Errorf("Invalid resource type specified")
	}
	b, err := client.doGetDetailsByName(client.BaseUrl.String(), resourceName)
	if err != nil {
		return
	}
	outputs, err := flattenNameDetailsResponse(b)
	if err != nil {
		return
	}
	marshalData(d, outputs)

	return
}

func flattenNameDetailsResponse(b []byte) (outputs map[string]interface{}, err error) {
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
	outputs["type"] = data["Type"]
	outputs["region"] = data["Region"]

	return
}