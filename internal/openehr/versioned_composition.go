package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const VERSIONED_COMPOSITION_MODEL_NAME string = "VERSIONED_COMPOSITION"

type VERSIONED_COMPOSITION struct {
	Type_       util.Optional[string] `json:"_type,omitzero"`
	UID         HIER_OBJECT_ID        `json:"uid"`
	OwnerID     OBJECT_REF            `json:"owner_id"`
	TimeCreated DV_DATE_TIME          `json:"time_created"`
}

func (vc *VERSIONED_COMPOSITION) isVersionModel() {}

func (vc *VERSIONED_COMPOSITION) SetModelName() {
	vc.Type_ = util.Some(VERSIONED_COMPOSITION_MODEL_NAME)
	vc.UID.SetModelName()
	vc.OwnerID.SetModelName()
	vc.TimeCreated.SetModelName()
}

func (vc *VERSIONED_COMPOSITION) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if vc.Type_.E && vc.Type_.V != VERSIONED_COMPOSITION_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          VERSIONED_COMPOSITION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + VERSIONED_COMPOSITION_MODEL_NAME,
			Recommendation: "Set _type to " + VERSIONED_COMPOSITION_MODEL_NAME,
		})
	}

	// Validate uid
	if err := vc.UID.Validate(path + ".uid"); len(err) > 0 {
		errors = append(errors, err...)
	}

	// Validate owner_id
	if err := vc.OwnerID.Validate(path + ".owner_id"); len(err) > 0 {
		errors = append(errors, err...)
	}

	// Validate time_created
	if err := vc.TimeCreated.Validate(path + ".time_created"); len(err) > 0 {
		errors = append(errors, err...)
	}

	return errors
}
