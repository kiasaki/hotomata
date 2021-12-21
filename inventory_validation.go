package hotomata

import (
	"github.com/xeipuuv/gojsonschema"
)

const inventorySchema = `
{
  "$schema": "http://json-schema.org/schema#",
  "definitions": {
	"reservedParams": {
	  "properties": {
		"ssh_hostname": {"type": "string"},
		"ssh_username": {"type": "string"},
		"ssh_password": {"type": "string"},
		"ssh_port": {"type": "integer", "minimum": 1, "maximum": 65000},
		"ssh_key": {"type": "string"}
	  }
	},
	"machine": {
	  "type": "object",
	  "$ref": "#/definitions/reservedParams",
	  "properties": {
		"name": {"type": "string"}
	  },
	  "required": ["name"]
	},
	"group": {
	  "type": "object",
	  "$ref": "#/definitions/reservedParams",
	  "properties": {
		"group_name": {"type": "string"},
		"machines": {
		  "type": "array",
		  "items": {
			"anyOf": [
			  {"$ref": "#/definitions/nestedGroup"},
			  {"$ref": "#/definitions/machine"}
			]
		  }
		}
	  },
	  "required": ["group_name", "machines"]
	},
	"nestedGroup": {
	  "type": "object",
	  "$ref": "#/definitions/reservedParams",
	  "properties": {
		"group_name": {"type": "string"},
		"machines": {
		  "type": "array",
		  "items": {"$ref": "#/definitions/machine"}
		}
	  },
	  "required": ["group_name", "machines"]
	}
  },
  "type": "array",
  "items": {
	"anyOf": [
	  {"$ref": "#/definitions/group"},
	  {"$ref": "#/definitions/machine"}
	]
  },
  "minItems": 1
}
`

func ValidateInventory(inventory string) (*gojsonschema.Result, error) {
	schemaLoader := gojsonschema.NewStringLoader(inventorySchema)
	documentLoader := gojsonschema.NewStringLoader(inventory)

	return gojsonschema.Validate(schemaLoader, documentLoader)
}
