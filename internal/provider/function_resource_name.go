package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &resourceNameFunction{}

type resourceNameFunction struct {
}

func (f *resourceNameFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "resource_name"
}

func (f *resourceNameFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Return Azure dx resources naming convention",
		Description: "Given a name, a resource name, an instance number and a resource type, returns the Azure dx resources naming convention.",

		// L'ordine di inserimento è questo
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "prefix",
				Description: "The default prefix for all resources that will be created.",
			},
			function.StringParameter{
				Name:        "name",
				Description: "The resource distintive name.",
			},
			function.StringParameter{
				Name:        "resource_type",
				Description: "Resource type (app or cosmos)",
			},
			function.Int64Parameter{
				Name:        "instance_number",
				Description: "Instance number (1-99)",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *resourceNameFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var prefix, name, resourceType string
	var instance int64
	var result string

	var resourceAbbreviations = map[string]string{
		// Compute
		"virtual_machine":   "vm",
		"container_app_job": "caj",

		// Storage
		"storage_account":          "st",
		"blob_storage":             "blob",
		"queue_storage":            "queue",
		"table_storage":            "table",
		"file_storage":             "file",
		"function_storage_account": "stfn",

		// Networking
		"api_management":              "apim",
		"api_management_autoscale":    "apim-as",
		"virtual_network":             "vnet",
		"network_security_group":      "nsg",
		"apim_network_security_group": "apim-nsg",
		"app_gateway":                 "agw",

		// Private Endpoints
		"cosmos_private_endpoint":          "cosno-pep",
		"postgre_private_endpoint":         "psql-pep",
		"postgre_replica_private_endpoint": "psql-pep-replica",
		"app_private_endpoint":             "app-pep",
		"app_slot_private_endpoint":        "staging-app-pep",
		"function_private_endpoint":        "func-pep",
		"function_slot_private_endpoint":   "staging-func-pep",
		"blob_private_endpoint":            "blob-pep",
		"queue_private_endpoint":           "queue-pep",
		"file_private_endpoint":            "file-pep",
		"table_private_endpoint":           "table-pep",
		"eventhub_private_endpoint":        "evhns-pep",

		// Subnets
		"app_subnet":      "app-snet",
		"apim_subnet":     "apim-snet",
		"function_subnet": "func-snet",

		// Databases
		"cosmos_db":         "cosmos",
		"cosmos_db_nosql":   "cosno",
		"postgresql":        "psql",
		"postgresq_replica": "psql-replica",

		// Integration
		"eventhub_namespace": "evhns",
		"function_app":       "func",
		"app_service":        "app",
		"app_service_plan":   "asp",

		// Security
		"key_vault": "kv",

		// Monitoring
		"application_insights": "appi",

		// Miscellaneous
		"resource_group": "rg",
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &prefix, &name, &resourceType, &instance))

	if resp.Error != nil {
		return
	}

	// Validate provider Prefix configuration
	pattern := `^[a-z]{2}-(d|u|p)-(itn|weu)(-[a-z]+)?$`
	matched, err := regexp.MatchString(pattern, prefix)
	if err != nil || !matched {
		fmt.Println("Regex error:", err)
		resp.Error = function.NewFuncError("prefix must be in the form of 'xx-(d|u|p)-(itn|weu)'")
	}

	if instance < 1 || instance > 99 {
		resp.Error = function.NewFuncError("Instance must be between 1 and 99")
		return
	}

	abbreviation, exists := resourceAbbreviations[resourceType]
	if !exists {
		validKeys := make([]string, 0, len(resourceAbbreviations))
		for key := range resourceAbbreviations {
			validKeys = append(validKeys, key)
		}
		resp.Error = function.NewFuncError(fmt.Sprintf("resource '%s' not found. Accepted values are: %s", resourceType, strings.Join(validKeys, ", ")))
		return
	}

	if name == "" {
		resp.Error = function.NewFuncError("Resource name cannot be empty")
		return
	}

	if strings.Contains(resourceType, "storage_account") {
		result = strings.ToLower(fmt.Sprintf("%s%s%s%02d",
			strings.Replace(prefix, "-", "", -1),
			name,
			abbreviation,
			instance))
	} else {
		result = strings.ToLower(fmt.Sprintf("%s-%s-%s-%02d",
			prefix,
			name,
			abbreviation,
			instance))
	}

	// Check total length (optional, adjust max length as needed)
	// Verify if is better here this check or in the specific module
	// if len(result) > 64 {
	// 	resp.Error = function.NewFuncError("Generated name exceeds maximum length of 64 characters")
	// 	return
	// }

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
