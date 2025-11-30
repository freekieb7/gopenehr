package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const LINK_MODEL_NAME string = "LINK"

type LINK struct {
	Type_   utils.Optional[string] `json:"_type,omitzero"`
	Meaning X_DV_TEXT              `json:"meaning"`
	Type    X_DV_TEXT              `json:"type"`
	Target  DV_EHR_URI             `json:"target"`
}

func (l *LINK) SetModelName() {
	l.Type_ = utils.Some(LINK_MODEL_NAME)
	l.Meaning.SetModelName()
	l.Type.SetModelName()
	l.Target.SetModelName()
}

func (l *LINK) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if l.Type_.E && l.Type_.V != LINK_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          LINK_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to LINK",
		})
	}

	return validateErr
}
