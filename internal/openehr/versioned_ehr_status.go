package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const VERSIONED_EHR_STATUS_MODEL_NAME string = "VERSIONED_EHR_STATUS"

type VERSIONED_EHR_STATUS struct {
	Type_       util.Optional[string] `json:"_type,omitzero"`
	UID         HIER_OBJECT_ID        `json:"uid"`
	OwnerID     OBJECT_REF            `json:"owner_id"`
	TimeCreated DV_DATE_TIME          `json:"time_created"`
}

func (v *VERSIONED_EHR_STATUS) SetModelName() {
	v.Type_ = util.Some(VERSIONED_EHR_STATUS_MODEL_NAME)
	v.UID.SetModelName()
	v.OwnerID.SetModelName()
	v.TimeCreated.SetModelName()
}

func (v VERSIONED_EHR_STATUS) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if v.Type_.E && v.Type_.V != VERSIONED_EHR_STATUS_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          VERSIONED_EHR_STATUS_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + VERSIONED_EHR_STATUS_MODEL_NAME,
		})
	}

	// Validate uid
	if err := v.UID.Validate(path + ".uid"); len(err) > 0 {
		errors = append(errors, err...)
	}

	// Validate owner_id
	if err := v.OwnerID.Validate(path + ".owner_id"); len(err) > 0 {
		errors = append(errors, err...)
	}

	// Validate time_created
	if err := v.TimeCreated.Validate(path + ".time_created"); len(err) > 0 {
		errors = append(errors, err...)
	}

	return errors
}
