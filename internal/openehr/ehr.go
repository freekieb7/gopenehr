package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const EHR_MODEL_NAME string = "EHR"

var _ util.ReferenceModel = (*EHR)(nil)

type EHR struct {
	Type_         util.Optional[string]         `json:"_type,omitzero"`
	SystemID      util.Optional[HIER_OBJECT_ID] `json:"system_id,omitzero"`
	EHRID         HIER_OBJECT_ID                `json:"ehr_id"`
	Contributions util.Optional[[]OBJECT_REF]   `json:"contributions,omitzero"`
	EHRStatus     OBJECT_REF                    `json:"ehr_status"`
	EHRAccess     OBJECT_REF                    `json:"ehr_access"`
	Compositions  util.Optional[[]OBJECT_REF]   `json:"compositions,omitzero"`
	Directory     util.Optional[OBJECT_REF]     `json:"directory,omitzero"`
	TimeCreated   DV_DATE_TIME                  `json:"time_created"`
	Folders       util.Optional[[]OBJECT_REF]   `json:"folders,omitzero"`
	Tags          util.Optional[[]OBJECT_REF]   `json:"tags,omitzero"`
}

func (e EHR) HasModelName() bool {
	return e.Type_.IsSet()
}

func (e EHR) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if e.Type_.IsSet() && e.Type_.Unwrap() != EHR_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          EHR_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid EHR _type field: %s", e.Type_.Unwrap()),
			Recommendation: "Ensure the _type field is set to 'EHR'",
		})
	}

	// Validate system_id
	if e.SystemID.IsSet() {
		attrPath = path + ".system_id"
		errs = append(errs, e.SystemID.Unwrap().Validate(attrPath)...)
	}

	// Validate ehr_id
	attrPath = path + ".ehr_id"
	errs = append(errs, e.EHRID.Validate(attrPath)...)

	// Validate contributions
	if e.Contributions.IsSet() {
		for i, contribRef := range e.Contributions.Unwrap() {
			attrPath = path + fmt.Sprintf(".contributions[%d]", i)
			if contribRef.Type != CONTRIBUTION_MODEL_NAME {
				errs = append(errs, util.ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid contribution type: %s", contribRef.Type),
					Recommendation: fmt.Sprintf("Ensure contributions[%d] _type field is set to '%s'", i, CONTRIBUTION_MODEL_NAME),
				})
			}
			errs = append(errs, contribRef.Validate(attrPath)...)
		}
	}

	// Validate ehr_status
	attrPath = path + ".ehr_status"
	if e.EHRStatus.Type != VERSIONED_EHR_STATUS_MODEL_NAME {
		errs = append(errs, util.ValidationError{
			Model:          EHR_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid EHR status type: %s", e.EHRStatus.Type),
			Recommendation: fmt.Sprintf("Ensure ehr_status _type field is set to '%s'", VERSIONED_EHR_STATUS_MODEL_NAME),
		})

	}
	errs = append(errs, e.EHRStatus.Validate(attrPath)...)

	// Validate ehr_access
	attrPath = path + ".ehr_access"
	if e.EHRAccess.Type != VERSIONED_EHR_ACCESS_MODEL_NAME {
		errs = append(errs, util.ValidationError{
			Model:          EHR_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid EHR access type: %s", e.EHRAccess.Type),
			Recommendation: fmt.Sprintf("Ensure ehr_access _type field is set to '%s'", VERSIONED_EHR_ACCESS_MODEL_NAME),
		})
	}
	errs = append(errs, e.EHRAccess.Validate(attrPath)...)

	// Validate compositions
	if e.Compositions.IsSet() {
		for i, compRef := range e.Compositions.Unwrap() {
			attrPath = path + fmt.Sprintf(".compositions[%d]", i)
			if compRef.Type != VERSIONED_COMPOSITION_MODEL_NAME {
				errs = append(errs, util.ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid composition type: %s", compRef.Type),
					Recommendation: fmt.Sprintf("Ensure compositions[%d] _type field is set to '%s'", i, VERSIONED_COMPOSITION_MODEL_NAME),
				})
			}
			errs = append(errs, compRef.Validate(attrPath)...)
		}
	}

	// Validate directory
	if e.Directory.IsSet() {
		directory := e.Directory.Unwrap()
		attrPath = path + ".directory"
		if directory.Type != VERSIONED_FOLDER_MODEL_NAME {
			errs = append(errs, util.ValidationError{
				Model:          EHR_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid folder type: %s", directory.Type),
				Recommendation: fmt.Sprintf("Ensure directory _type field is set to '%s'", VERSIONED_FOLDER_MODEL_NAME),
			})
		}

		errs = append(errs, directory.Validate(attrPath)...)
	}

	// Validate time_created
	attrPath = path + ".time_created"
	errs = append(errs, e.TimeCreated.Validate(attrPath)...)

	// Validate folders
	if e.Folders.IsSet() {
		for i, folderRef := range e.Folders.Unwrap() {
			attrPath = path + fmt.Sprintf(".folders[%d]", i)
			if folderRef.Type != VERSIONED_FOLDER_MODEL_NAME {
				errs = append(errs, util.ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid folder type: %s", folderRef.Type),
					Recommendation: fmt.Sprintf("Ensure folders[%d] _type field is set to '%s'", i, VERSIONED_FOLDER_MODEL_NAME),
				})
			}
			errs = append(errs, folderRef.Validate(attrPath)...)
		}
	}

	// Validate tags
	if e.Tags.IsSet() {
		for i, tagRef := range e.Tags.Unwrap() {
			attrPath = path + fmt.Sprintf(".tags[%d]", i)
			errs = append(errs, tagRef.Validate(attrPath)...)
		}
	}

	return errs
}
