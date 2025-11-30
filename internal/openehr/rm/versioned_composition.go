package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const VERSIONED_COMPOSITION_TYPE string = "VERSIONED_COMPOSITION"

type VERSIONED_COMPOSITION struct {
	Type_       utils.Optional[string] `json:"_type,omitzero"`
	UID         HIER_OBJECT_ID         `json:"uid"`
	OwnerID     OBJECT_REF             `json:"owner_id"`
	TimeCreated DV_DATE_TIME           `json:"time_created"`
}

func (vc *VERSIONED_COMPOSITION) SetModelName() {
	vc.Type_ = utils.Some(VERSIONED_COMPOSITION_TYPE)
	vc.UID.SetModelName()
	vc.OwnerID.SetModelName()
	vc.TimeCreated.SetModelName()
}

func (vc *VERSIONED_COMPOSITION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if vc.Type_.E && vc.Type_.V != VERSIONED_COMPOSITION_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          VERSIONED_COMPOSITION_TYPE,
			Path:           attrPath,
			Message:        "_type must be " + VERSIONED_COMPOSITION_TYPE,
			Recommendation: "Set _type to " + VERSIONED_COMPOSITION_TYPE,
		})
	}

	// Validate uid
	attrPath = path + ".uid"
	validateErr.Errs = append(validateErr.Errs, vc.UID.Validate(attrPath).Errs...)

	// Validate owner_id
	attrPath = path + ".owner_id"
	validateErr.Errs = append(validateErr.Errs, vc.OwnerID.Validate(attrPath).Errs...)

	// Validate time_created
	attrPath = path + ".time_created"
	validateErr.Errs = append(validateErr.Errs, vc.TimeCreated.Validate(attrPath).Errs...)

	return validateErr
}
