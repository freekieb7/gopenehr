package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const DV_GENERAL_TIME_SPECIFICATION_MODEL_NAME string = "DV_GENERAL_TIME_SPECIFICATION"

type DV_GENERAL_TIME_SPECIFICATION struct {
	Type_ util.Optional[string] `json:"_type,omitzero"`
	Value DV_PARSABLE           `json:"value"`
}

func (d DV_GENERAL_TIME_SPECIFICATION) isDataValueModel() {}

func (d DV_GENERAL_TIME_SPECIFICATION) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_GENERAL_TIME_SPECIFICATION) SetModelName() {
	d.Type_ = util.Some(DV_GENERAL_TIME_SPECIFICATION_MODEL_NAME)
	d.Value.SetModelName()
}

func (d DV_GENERAL_TIME_SPECIFICATION) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_GENERAL_TIME_SPECIFICATION_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          DV_GENERAL_TIME_SPECIFICATION_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_GENERAL_TIME_SPECIFICATION",
		})
	}

	// Validate value
	attrPath = path + ".value"
	errors = append(errors, d.Value.Validate(attrPath)...)

	return errors
}
