package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const VERSIONED_PARTY_TYPE string = "VERSIONED_PARTY"

type VERSIONED_PARTY struct {
	Type_       utils.Optional[string] `json:"_type,omitzero"`
	UID         HIER_OBJECT_ID         `json:"uid"`
	OwnerID     OBJECT_REF             `json:"owner_id"`
	TimeCreated DV_DATE_TIME           `json:"time_created"`
}

func (v *VERSIONED_PARTY) SetModelName() {
	v.Type_ = utils.Some(VERSIONED_PARTY_TYPE)
	v.UID.SetModelName()
	v.OwnerID.SetModelName()
	v.TimeCreated.SetModelName()
}

func (v *VERSIONED_PARTY) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if v.Type_.E && v.Type_.V != VERSIONED_PARTY_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          VERSIONED_PARTY_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + VERSIONED_PARTY_TYPE,
		})
	}

	// Validate uid
	attrPath = path + ".uid"
	validateErr.Errs = append(validateErr.Errs, v.UID.Validate(attrPath).Errs...)

	// Validate owner_id
	attrPath = path + ".owner_id"
	validateErr.Errs = append(validateErr.Errs, v.OwnerID.Validate(attrPath).Errs...)

	// Validate time_created
	attrPath = path + ".time_created"
	validateErr.Errs = append(validateErr.Errs, v.TimeCreated.Validate(attrPath).Errs...)

	return validateErr
}
