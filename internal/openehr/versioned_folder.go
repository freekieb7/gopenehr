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

func (vf *VERSIONED_FOLDER) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if vf.Type_.E && vf.Type_.V != VERSIONED_FOLDER_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          VERSIONED_FOLDER_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + VERSIONED_FOLDER_MODEL_NAME,
			Recommendation: "Set _type to " + VERSIONED_FOLDER_MODEL_NAME,
		})
	}

	// Validate uid
	attrPath = path + ".uid"
	validateErr.Errs = append(validateErr.Errs, vf.UID.Validate(attrPath).Errs...)

	// Validate owner_id
	attrPath = path + ".owner_id"
	validateErr.Errs = append(validateErr.Errs, vf.OwnerID.Validate(attrPath).Errs...)
	// Validate time_created
	attrPath = path + ".time_created"
	validateErr.Errs = append(validateErr.Errs, vf.TimeCreated.Validate(attrPath).Errs...)

	return validateErr
}
