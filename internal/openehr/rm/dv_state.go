package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_STATE_MODEL_NAME string = "DV_STATE"

type DV_STATE struct {
	Type_      utils.Optional[string] `json:"_type,omitzero"`
	Value      DV_CODED_TEXT          `json:"value"`
	IsTerminal bool                   `json:"is_terminal"`
}

func (d *DV_STATE) isDataValueModel() {}

func (d *DV_STATE) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_STATE) SetModelName() {
	d.Type_ = utils.Some(DV_STATE_MODEL_NAME)
	d.Value.SetModelName()
}

func (d *DV_STATE) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_STATE_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_STATE_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_STATE",
		})
	}

	// Validate value
	attrPath = path + ".value"
	validateErr.Errs = append(validateErr.Errs, d.Value.Validate(attrPath).Errs...)

	return validateErr
}
