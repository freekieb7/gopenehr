package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const INSTRUCTION_DETAILS_MODEL_NAME string = "INSTRUCTION_DETAILS"

type INSTRUCTION_DETAILS struct {
	Type_         util.Optional[string]           `json:"_type,omitzero"`
	InstructionID LOCATABLE_REF                   `json:"instruction_id"`
	ActivityID    string                          `json:"activity"`
	WfDetails     util.Optional[X_ITEM_STRUCTURE] `json:"wf_details,omitzero"`
}

func (i *INSTRUCTION_DETAILS) SetModelName() {
	i.Type_ = util.Some(INSTRUCTION_DETAILS_MODEL_NAME)
	i.InstructionID.SetModelName()
	if i.WfDetails.E {
		i.WfDetails.V.SetModelName()
	}
}

func (i *INSTRUCTION_DETAILS) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != INSTRUCTION_DETAILS_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          INSTRUCTION_DETAILS_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + INSTRUCTION_DETAILS_MODEL_NAME,
			Recommendation: "Set _type to " + INSTRUCTION_DETAILS_MODEL_NAME,
		})
	}

	// Validate instruction_id
	attrPath = path + ".instruction_id"
	errs = append(errs, i.InstructionID.Validate(attrPath)...)

	// Validate activity
	// No validation for string type

	// Validate wf_details
	if i.WfDetails.E {
		attrPath = path + ".wf_details"
		errs = append(errs, i.WfDetails.V.Validate(attrPath)...)
	}

	return errs
}
