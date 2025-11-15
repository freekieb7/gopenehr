package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const VERSIONED_FOLDER_MODEL_NAME string = "VERSIONED_FOLDER"

type VERSIONED_FOLDER struct {
	Type_       util.Optional[string] `json:"_type,omitzero"`
	UID         HIER_OBJECT_ID        `json:"uid"`
	OwnerID     OBJECT_REF            `json:"owner_id"`
	TimeCreated DV_DATE_TIME          `json:"time_created"`
}

func (vf *VERSIONED_FOLDER) isVersionModel() {}

func (vf *VERSIONED_FOLDER) SetModelName() {
	vf.Type_ = util.Some(VERSIONED_FOLDER_MODEL_NAME)
	vf.UID.SetModelName()
	vf.OwnerID.SetModelName()
	vf.TimeCreated.SetModelName()
}

func (vf *VERSIONED_FOLDER) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if vf.Type_.E && vf.Type_.V != VERSIONED_FOLDER_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          VERSIONED_FOLDER_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + VERSIONED_FOLDER_MODEL_NAME,
			Recommendation: "Set _type to " + VERSIONED_FOLDER_MODEL_NAME,
		})
	}

	// Validate uid
	if err := vf.UID.Validate(path + ".uid"); len(err) > 0 {
		errors = append(errors, err...)
	}

	// Validate owner_id
	if err := vf.OwnerID.Validate(path + ".owner_id"); len(err) > 0 {
		errors = append(errors, err...)
	}

	// Validate time_created
	if err := vf.TimeCreated.Validate(path + ".time_created"); len(err) > 0 {
		errors = append(errors, err...)
	}

	return errors
}
