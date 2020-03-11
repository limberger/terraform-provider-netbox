package netbox

import (
	// "fmt"
	// "reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func customFieldFilterSchema(conflicts []string) *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeMap,
		Optional:      true,
		ConflictsWith: conflicts,
		ValidateFunc: func(m interface{}, k string) (ws []string, errors []error) {
			for _, v := range m.(map[string]interface{}) {
				_, err := regexp.Compile(v.(string))
				if err != nil {
					errors = append(errors, err)
				}
			}
			return
		},
	}
}
