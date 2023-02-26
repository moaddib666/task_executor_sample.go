package tasks

import (
	"fmt"
)

func validateJSONAgainstSchema(data map[string]interface{}, schema *JsonSchema) error {
	if err := validateProperties(data, schema.Properties); err != nil {
		return err
	}

	if err := validateRequired(data, schema.Required); err != nil {
		return err
	}

	return nil
}

func validateProperties(data map[string]interface{}, schemaProperties map[string]*JsonSchema) error {
	for propName, propSchema := range schemaProperties {
		propValue, ok := data[propName]
		if !ok {
			continue
		}

		switch propSchema.Type {
		case "object":
			nestedData, ok := propValue.(map[string]interface{})
			if !ok {
				return fmt.Errorf("property %s is not of type object", propName)
			}

			if err := validateProperties(nestedData, propSchema.Properties); err != nil {
				return err
			}
		case "array":
			nestedData, ok := propValue.([]interface{})
			if !ok {
				return fmt.Errorf("property %s is not of type array", propName)
			}

			for _, nestedItem := range nestedData {
				nestedItemData, ok := nestedItem.(map[string]interface{})
				if !ok {
					return fmt.Errorf("item in property %s array is not of type object", propName)
				}

				if err := validateProperties(nestedItemData, propSchema.Properties); err != nil {
					return err
				}
			}
		case "string":
			_, ok := propValue.(string)
			if !ok {
				return fmt.Errorf("property %s is not of type string", propName)
			}
		case "integer":
			_, ok := propValue.(int)
			if !ok {
				return fmt.Errorf("property %s is not of type integer", propName)
			}
		default:
			return fmt.Errorf("unknown property type: %s", propSchema.Type)
		}
	}

	return nil
}

func validateRequired(data map[string]interface{}, requiredProps []string) error {
	for _, propName := range requiredProps {
		_, ok := data[propName]
		if !ok {
			return fmt.Errorf("property %s is required", propName)
		}
	}

	return nil
}
