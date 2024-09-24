package openehr

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const (
	REGEX_UUID    = `\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`
	REGEX_ISO_OID = `([1-9]+)([0-9]+){0,8}`
	//REGEX_OBJECT_VERSION_ID = `([a-zA-Z0-9]+)::([a-zA-Z0-9]+)::([1-9]([0-9]+)?)(\.[1-9]([0-9]+)?){2}`
	REGEX_INTERNET_ID = `([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})(\.(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))){1,2}`
)

var (
	validationSchemaCollection = map[string]ValidationSchema{
		"EHR": {
			Attributes: map[string]ValidationRules{
				"_type": {
					Type:    "string",
					EqualTo: "EHR",
				},
				"system_id": {
					Required: true,
					Type:     "HIER_OBJECT_ID",
				},
				"ehr_id": {
					Required: true,
					Type:     "HIER_OBJECT_ID",
				},
				"contributions": {
					Type: "LIST<OBJECT_REF>",
				},
				"ehr_status": {
					Required: true,
					Type:     "OBJECT_REF",
				},
				"ehr_access": {
					Required: true,
					Type:     "OBJECT_REF",
				},
				"compositions": {
					Type: "LIST<OBJECT_REF>",
				},
				"directory": {
					Type: "OBJECT_REF",
				},
				"time_created": {
					Required: true,
					Type:     "DV_DATE_TIME",
				},
				"folders": {
					Type: "LIST<OBJECT_REF>",
				},
				"tags": {
					Type: "LIST<OBJECT_REF>",
				},
			},
		},
		"HIER_OBJECT_ID": {
			Inherits: []string{"UID_BASED_ID"},
			Attributes: map[string]ValidationRules{
				"_type": {
					Type:    "string",
					EqualTo: "HIER_OBJECT_ID",
				},
				"value": {
					Type:   "string",
					Regexp: regexp.MustCompile(fmt.Sprintf(`^((%s)|(%s)|(%s))(::(\w){1,36})?$`, REGEX_ISO_OID, REGEX_UUID, REGEX_INTERNET_ID)),
				},
			},
		},
		"UID_BASED_ID": {
			IsAbstract: true,
			Inherits:   []string{"OBJECT_ID"},
			Attributes: map[string]ValidationRules{
				"_type": {
					Type:    "string",
					EqualTo: "UID_BASED_ID",
				},
			},
		},
		"OBJECT_REF": {
			Attributes: map[string]ValidationRules{
				"_type": {
					Type:    "string",
					EqualTo: "OBJECT_REF",
				},
				"namespace": {
					Required: true,
					Type:     "string",
					Regexp:   regexp.MustCompile(`^(local)|(unknown)|([a-zA-Z][a-zA-Z0-9_.:\\/&?=+-]*)$`), // Can be better
				},
				"type": {
					Required: true,
					Type:     "string",
				},
				"id": {
					Required: true,
					Type:     "OBJECT_ID",
				},
			},
		},
		"OBJECT_ID": {
			IsAbstract: true,
			Attributes: map[string]ValidationRules{
				"_type": {
					Type:    "string",
					EqualTo: "OBJECT_ID",
				},
				"value": {
					Type: "string",
				},
			},
		},
		"DV_DATE_TIME": {
			Inherits: []string{"DV_TEMPORAL", "ISO_8601_DATE_TIME"},
			Attributes: map[string]ValidationRules{
				"_type": {
					Type:    "string",
					EqualTo: "DV_DATE_TIME",
				},
				"value": {
					Type:   "string",
					Regexp: regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))$`),
				},
			},
		},
		"OBJECT_VERSION_ID": {
			Inherits: []string{"UID_BASED_ID"},
			Attributes: map[string]ValidationRules{
				"_type": {
					Type:    "string",
					EqualTo: "OBJECT_VERSION_ID",
				},
			},
		},
	}
)

type ValidationSchema struct {
	IsAbstract bool
	Inherits   []string
	Attributes map[string]ValidationRules
}

type ValidationRules struct {
	Required bool
	Type     string
	Regexp   *regexp.Regexp
	EqualTo  string
}

func validate(data map[string]any, currentType string, schema ValidationSchema) error {
	// Abstraction check
	if schema.IsAbstract {
		elementType, found := data["_type"]

		if !found {
			return errors.New("abstract fields need an _type field specified")
		}

		elementTypeAsStr, isStr := elementType.(string)

		if !isStr {
			return errors.New("_type field needs to be a string")
		}

		schema, found = validationSchemaCollection[elementTypeAsStr]

		if !found {
			return fmt.Errorf("object with type %s not found", elementTypeAsStr)
		}

		if schema.IsAbstract {
			return fmt.Errorf("cannot have abstract type %s in data", elementTypeAsStr)
		}

		// Upper level reached
		i := 0
		found = false

		for {
			if i > 100 { // Infinite inherits loop protection
				return errors.New("too many fields found in inherits list")
			}

			if i == len(schema.Inherits) {
				break
			}

			inherit := schema.Inherits[i]

			if inherit == currentType {
				found = true
			}

			inheritSchema, found := validationSchemaCollection[inherit]

			if !found {
				return fmt.Errorf("object with type %s not found", elementTypeAsStr)
			}

			schema.Inherits = append(schema.Inherits, inheritSchema.Inherits...)
			for a, b := range inheritSchema.Attributes {
				if _, found := schema.Attributes[a]; found {
					continue
				}

				schema.Attributes[a] = b
			}

			i++
		}

		if !found {
			return fmt.Errorf("%s does not inherit %s", elementTypeAsStr, currentType)
		}
	}

	// Search for unwanted fields in data
	for attributeName, _ := range data {
		if _, found := schema.Attributes[attributeName]; !found {
			return fmt.Errorf("%s is not included in the expected attribute list", attributeName)
		}
	}

	for attributeName, attributeSchema := range schema.Attributes {
		attributeData, found := data[attributeName]

		// Required check
		if !found {
			if attributeSchema.Required == true {
				return fmt.Errorf("%s is a required field", attributeName)
			}

			continue
		}

		//// Inherits check
		//for _, inheritName := range schema.Inherits {
		//	inheritSchema, inheritFound := validationSchemaCollection[inheritName]
		//
		//	if !inheritFound {
		//		return fmt.Errorf("validation not implemented for type %s", inheritName)
		//	}
		//
		//	for inheritAttrName, inheritAttrVal := range inheritSchema.Attributes {
		//		// Prevent overriding of inherit attr over current schema attribute
		//		if _, exists := schema.Attributes[inheritAttrName]; !exists {
		//			continue
		//		}
		//
		//		// Add to attributes list
		//		schema.Attributes[inheritAttrName] = inheritAttrVal
		//	}
		//}

		// Type check
		if err := validateAttribute(attributeData, attributeSchema); err != nil {
			return fmt.Errorf("%s.%w", attributeName, err)
		}
	}

	return nil
}

func validateAttribute(data any, rules ValidationRules) error {
	switch rules.Type {
	case "string":
		{
			val, ok := data.(string)

			if !ok {
				return errors.New("data is not a string")
			}

			if rules.EqualTo != "" && rules.EqualTo != val {
				return fmt.Errorf("data must be equal to %s", rules.EqualTo)
			}

			if rules.Regexp != nil && !rules.Regexp.MatchString(val) {
				return fmt.Errorf("data does not match expected lexical format")
			}
		}
	default:
		{
			if strings.HasPrefix(rules.Type, "LIST<") && strings.HasSuffix(rules.Type, ">") {
				// Validate LIST<OBJECT>
				valList, ok := data.([]any)

				if !ok {
					return errors.New(" data is not a list")
				}

				listElementType := rules.Type[5 : len(rules.Type)-1]
				listElementSchema, found := validationSchemaCollection[listElementType]

				if !found {
					return fmt.Errorf("validation not implemented for type %s", listElementType)
				}

				for i, val := range valList {
					fmt.Println(reflect.TypeOf(val))
					if err := validate(val.(map[string]interface{}), listElementType, listElementSchema); err != nil {
						return fmt.Errorf("[%d] %w", i, err)
					}
				}

				return nil
			}

			// Validate OBJECT
			val, ok := data.(map[string]any)

			if !ok {
				return errors.New("data is not an object")
			}

			elementSchema, found := validationSchemaCollection[rules.Type]

			if !found {
				return fmt.Errorf("validation not implemented for type %s", rules.Type)
			}

			if err := validate(val, rules.Type, elementSchema); err != nil {
				return err
			}
		}
	}

	return nil
}
