package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const OBJECT_REF_MODEL_NAME = "OBJECT_REF"

var _ util.ReferenceModel = (*OBJECT_REF)(nil)

type OBJECT_REF struct {
	Type_     util.Optional[string] `json:"_type,omitzero"`
	Namespace string                `json:"namespace"`
	Type      string                `json:"type"`
	ID        X_OBJECT_ID           `json:"id"`
}

func (o OBJECT_REF) HasModelName() bool {
	return o.Type_.IsSet()
}

func (o OBJECT_REF) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if o.Type_.IsSet() && o.Type_.Unwrap() != OBJECT_REF_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          OBJECT_REF_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", OBJECT_REF_MODEL_NAME, o.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", OBJECT_REF_MODEL_NAME),
		})
	}

	// Validate namespace
	attrPath = path + ".namespace"
	if o.Namespace == "" {
		errs = append(errs, util.ValidationError{
			Model:          "String",
			Path:           attrPath,
			Message:        "invalid namespace: cannot be empty",
			Recommendation: "Fill in a value matching the regex standard documented in the specifications",
		})
	} else {
		if !util.NamespaceRegex.MatchString(o.Namespace) {
			errs = append(errs, util.ValidationError{
				Model:          "String",
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid namespace: %s", o.Namespace),
				Recommendation: "Fill in a value matching the regex standard documented in the specifications",
			})
		}
	}

	// Validate type
	attrPath = path + ".type"
	if o.Type == "" {
		errs = append(errs, util.ValidationError{
			Model:          "String",
			Path:           attrPath,
			Message:        "invalid type: cannot be empty",
			Recommendation: "Fill in a value matching the regex standard documented in the specifications",
		})
	}

	// Validate id
	attrPath = path + ".id"
	errs = append(errs, o.ID.Validate(attrPath)...)

	return errs
}
