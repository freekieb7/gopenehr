package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const INSTRUCTION_DETAILS_MODEL_NAME string = "INSTRUCTION_DETAILS"

type INSTRUCTION_DETAILS struct {
	Type_         utils.Optional[string]           `json:"_type,omitzero"`
	InstructionID LOCATABLE_REF                    `json:"instruction_id"`
	ActivityID    string                           `json:"activity"`
	WfDetails     utils.Optional[X_ITEM_STRUCTURE] `json:"wf_details,omitzero"`
}

func (i *INSTRUCTION_DETAILS) SetModelName() {
	i.Type_ = utils.Some(INSTRUCTION_DETAILS_MODEL_NAME)
	i.InstructionID.SetModelName()
	if i.WfDetails.E {
		i.WfDetails.V.SetModelName()
	}
}

func (i *INSTRUCTION_DETAILS) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != INSTRUCTION_DETAILS_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          INSTRUCTION_DETAILS_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + INSTRUCTION_DETAILS_MODEL_NAME,
			Recommendation: "Set _type to " + INSTRUCTION_DETAILS_MODEL_NAME,
		})
	}

	// Validate instruction_id
	attrPath = path + ".instruction_id"
	validateErr.Errs = append(validateErr.Errs, i.InstructionID.Validate(attrPath).Errs...)

	// Validate activity
	// No validation for string type

	// Validate wf_details
	if i.WfDetails.E {
		attrPath = path + ".wf_details"
		validateErr.Errs = append(validateErr.Errs, i.WfDetails.V.Validate(attrPath).Errs...)
	}

	return validateErr
}
