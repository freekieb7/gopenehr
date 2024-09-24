package openehr

import (
	"errors"
	"fmt"
)

type Ehr struct {
	Json map[string]any
}

type EhrFactory struct{}

func (factory EhrFactory) Produce(json map[string]any) (Ehr, error) {
	var e Ehr

	currentType := "EHR"

	err := validate(json, currentType, validationSchemaCollection[currentType])

	if err != nil {
		return e, errors.New("EHR Validation Error: " + err.Error())
	}

	// Additional rules
	if v, found := json["contributions"]; found {
		// type check
		v := v.([]map[string]any)
		for i, el := range v {
			if el["type"] != "CONTRIBUTION" {
				return e, fmt.Errorf("contributions[%d].type does not match CONTRIBUTION", i)
			}
		}
	}

	if v, found := json["ehr_status"]; found {
		// type check
		v := v.(map[string]any)
		if v["type"] != "EHR_STATUS" {
			return e, errors.New("ehr_status.type does not match EHR_STATUS")
		}
	}

	if v, found := json["ehr_access"]; found {
		// type check
		v := v.(map[string]any)
		if v["type"] != "EHR_ACCESS" {
			return e, errors.New("ehr_status.type does not match EHR_STATUS")
		}
	}

	if v, found := json["compositions"]; found {
		// type check
		v := v.([]map[string]any)
		for i, el := range v {
			if el["type"] != "COMPOSITION" {
				return e, fmt.Errorf("compositions[%d].type does not match COMPOSITION", i)
			}
		}
	}

	if _, found := json["directory"]; found {
		// type check
		return e, errors.New("directory field is deprecated, use folders where the first is the \"directory\"")
	}

	if v, found := json["folders"]; found {
		// type check
		v := v.([]any)
		for i, el := range v {
			el := el.(map[string]any)
			if el["type"] != "FOLDER" {
				return e, fmt.Errorf("folders[%d].type does not match FOLDER", i)
			}
		}
	}

	e.Json = json
	return e, nil
}
