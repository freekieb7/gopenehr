package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const OBJECT_REF_TYPE = "OBJECT_REF"

const (
	Namespace_local   = "local"
	Namespace_unknown = "unknown"
)

type OBJECT_REF struct {
	Type_     utils.Optional[string] `json:"_type,omitzero"`
	Namespace string                 `json:"namespace"`
	Type      string                 `json:"type"`
	ID        ObjectIDUnion          `json:"id"`
}

func (o *OBJECT_REF) SetModelName() {
	o.Type_ = utils.Some(OBJECT_REF_TYPE)
	o.ID.SetModelName()
}

func (o *OBJECT_REF) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if o.Type_.E && o.Type_.V != OBJECT_REF_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          OBJECT_REF_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", OBJECT_REF_TYPE, o.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", OBJECT_REF_TYPE),
		})
	}

	// Validate namespace
	attrPath = path + ".namespace"
	if o.Namespace == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          "String",
			Path:           attrPath,
			Message:        "invalid namespace: cannot be empty",
			Recommendation: "Fill in a value matching the regex standard documented in the specifications",
		})
	} else {
		if !util.NamespaceRegex.MatchString(o.Namespace) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          "String",
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid namespace: %s", o.Namespace),
				Recommendation: "Fill in a value matching the regex standard documented in the specifications",
			})
		}
	}

	// Validate type
	if o.Type == "" {
		attrPath = path + ".type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          "String",
			Path:           attrPath,
			Message:        "invalid type: cannot be empty",
			Recommendation: "Fill in a value matching the regex standard documented in the specifications",
		})
	}

	// Validate id
	attrPath = path + ".id"
	validateErr.Errs = append(validateErr.Errs, o.ID.Validate(attrPath).Errs...)

	// Validated overal object values
	switch o.Type {
	case EHR_TYPE,
		VERSIONED_EHR_STATUS_TYPE,
		VERSIONED_EHR_ACCESS_TYPE,
		VERSIONED_COMPOSITION_TYPE,
		VERSIONED_FOLDER_TYPE,
		VERSIONED_PARTY_TYPE:
		// Valid type
		_, ok := o.ID.Value.(*HIER_OBJECT_ID)
		if !ok {
			attrPath = path + ".id"
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          OBJECT_REF_TYPE,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid id type for object ref type %s: expected HIER_OBJECT_ID", o.Type),
				Recommendation: fmt.Sprintf("Ensure id is of type HIER_OBJECT_ID for object ref type %s", o.Type),
			})
		}
	case EHR_STATUS_TYPE,
		EHR_ACCESS_TYPE,
		COMPOSITION_TYPE,
		FOLDER_TYPE,
		PERSON_TYPE,
		AGENT_TYPE,
		GROUP_TYPE,
		ORGANISATION_TYPE:
		// Valid type
		_, ok := o.ID.Value.(*OBJECT_VERSION_ID)
		if !ok {
			attrPath = path + ".id"
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          OBJECT_REF_TYPE,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid id type for object ref type %s: expected OBJECT_VERSION_ID", o.Type),
				Recommendation: fmt.Sprintf("Ensure id is of type OBJECT_VERSION_ID for object ref type %s", o.Type),
			})
		}
	}

	return validateErr
}
