package config

import (
	"testconfigload/logger"

	"github.com/xeipuuv/gojsonschema"
)

const schema = `{
  "description": "number",
  "type": "object",
  "required": [ "Num"],
  "additionalProperties": false,
  "properties": {
    "Num": {
	  "type": "integer",
	  "minimum": 1,
      "description": "params"
	}
}
}`

func VerifyConfig(conf config) bool {
	configLoader := gojsonschema.NewGoLoader(conf)
	schemaLoader := gojsonschema.NewStringLoader(schema)
	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		logger.Logger.Error("Validate error: %s", err)
		return false
	} else {
		if result.Valid() {
			logger.Logger.Info("The document is valid")
			return true
		} else {
			logger.Logger.Info("The document is not valid. see errors :\n")
			for _, err := range result.Errors() {
				// Err implements the ResultError interface
				logger.Logger.Error("- %s", err)
			}
			return false
		}
	}
}
