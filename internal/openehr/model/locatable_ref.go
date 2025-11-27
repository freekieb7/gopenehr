package model

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const LOCATABLE_REF_MODEL_NAME string = "LOCATABLE_REF"

type LOCATABLE_REF struct {
	Type_     utils.Optional[string] `json:"_type,omitzero"`
	Namespace string                 `json:"namespace"`
	Type      string                 `json:"type"`
	Path      utils.Optional[string] `json:"path,omitzero"`
	ID        X_UID_BASED_ID         `json:"id"`
}

func (l *LOCATABLE_REF) HasModelName() bool {
	return l.Type_.E
}

func (l *LOCATABLE_REF) SetModelName() {
	l.Type_ = utils.Some(LOCATABLE_REF_MODEL_NAME)
	l.ID.SetModelName()
}

func (l *LOCATABLE_REF) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if l.Type_.E && l.Type_.V != LOCATABLE_REF_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          LOCATABLE_REF_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + LOCATABLE_REF_MODEL_NAME,
			Recommendation: "Set _type to " + LOCATABLE_REF_MODEL_NAME,
		})
	}

	// Validate namespace
	// No validation for string type

	// Validate type
	// No validation for string type

	// Validate path
	// if l.Path.E {
	// 	// No validation for string type
	// }

	// Validate id
	attrPath = path + ".id"
	validateErr.Errs = append(validateErr.Errs, l.ID.Validate(attrPath).Errs...)

	return validateErr
}
