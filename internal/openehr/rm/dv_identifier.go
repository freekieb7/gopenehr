package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_IDENTIFIER_TYPE string = "DV_IDENTIFIER"

type DV_IDENTIFIER struct {
	Type_    utils.Optional[string] `json:"_type,omitzero"`
	Issuer   utils.Optional[string] `json:"issuer,omitzero"`
	Assigner utils.Optional[string] `json:"assigner,omitzero"`
	ID       string                 `json:"id"`
	Type     utils.Optional[string] `json:"type,omitzero"`
}

func (d *DV_IDENTIFIER) SetModelName() {
	d.Type_ = utils.Some(DV_IDENTIFIER_TYPE)
}

func (d *DV_IDENTIFIER) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_IDENTIFIER_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_IDENTIFIER_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_IDENTIFIER",
		})
	}

	// Validate id
	attrPath = path + ".id"
	if d.ID == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_IDENTIFIER_TYPE,
			Path:           attrPath,
			Message:        "id field cannot be empty",
			Recommendation: "Ensure id field is not empty",
		})
	}

	return validateErr
}
