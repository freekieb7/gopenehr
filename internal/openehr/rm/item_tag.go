package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ITEM_TAG_TYPE = "ITEM_TAG"

type ITEM_TAG struct {
	Type_      utils.Optional[string] `json:"_type,omitzero"`
	Key        string                 `json:"key"`
	Value      utils.Optional[string] `json:"value,omitzero"`
	Target     UIDBasedIDUnion        `json:"target"`
	TargetPath utils.Optional[string] `json:"target_path,omitzero"`
	OwnerID    OBJECT_REF             `json:"owner_id"`
}

func (i *ITEM_TAG) SetModelName() {
	i.Type_ = utils.Some(ITEM_TAG_TYPE)
	i.Target.SetModelName()
	i.OwnerID.SetModelName()
}

func (i ITEM_TAG) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != ITEM_TAG_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ITEM_TAG_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + ITEM_TAG_TYPE,
		})
	}

	// Validate key
	attrPath = path + ".key"
	if i.Key == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ITEM_TAG_TYPE,
			Path:           attrPath,
			Message:        "key is required and cannot be empty",
			Recommendation: "Provide a valid key for the ITEM_TAG",
		})
	}

	// Validate value
	if i.Value.E {
		attrPath = path + ".value"
		if i.Value.V == "" {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          ITEM_TAG_TYPE,
				Path:           attrPath,
				Message:        "value cannot be empty if provided",
				Recommendation: "Provide a valid value for the ITEM_TAG or omit the field",
			})
		}
	}

	// Validate target
	attrPath = path + ".target"
	validateErr.Errs = append(validateErr.Errs, i.Target.Validate(attrPath).Errs...)

	// Validate target_path
	if i.TargetPath.E {
		attrPath = path + ".target_path"
		if i.TargetPath.V == "" {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          ITEM_TAG_TYPE,
				Path:           attrPath,
				Message:        "target_path cannot be empty if provided",
				Recommendation: "Provide a valid target_path for the ITEM_TAG or omit the field",
			})
		}
	}

	// Validate owner_id
	attrPath = path + ".owner_id"
	validateErr.Errs = append(validateErr.Errs, i.OwnerID.Validate(attrPath).Errs...)

	return validateErr
}
