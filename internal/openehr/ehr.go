package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const EHR_MODEL_NAME string = "EHR"

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

func (e *EHR) SetModelName() {
	e.Type_ = util.Some(EHR_MODEL_NAME)
	if e.SystemID.E {
		e.SystemID.V.SetModelName()
	}
	e.EHRID.SetModelName()
	if e.Contributions.E {
		for i := range e.Contributions.V {
			e.Contributions.V[i].SetModelName()
		}
	}
	e.EHRStatus.SetModelName()
	e.EHRAccess.SetModelName()
	for i := range e.Compositions.V {
		e.Compositions.V[i].SetModelName()
	}
	if e.Directory.E {
		e.Directory.V.SetModelName()
	}
	e.TimeCreated.SetModelName()
	for i := range e.Folders.V {
		e.Folders.V[i].SetModelName()
	}
	for i := range e.Tags.V {
		e.Tags.V[i].SetModelName()
	}
}

func (e *EHR) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if e.Type_.E && e.Type_.V != EHR_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          EHR_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid EHR _type field: %s", e.Type_.V),
			Recommendation: "Ensure the _type field is set to 'EHR'",
		})
	}

	// Validate system_id
	if e.SystemID.E {
		attrPath = path + ".system_id"
		errs = append(errs, e.SystemID.V.Validate(attrPath)...)
	}

	// Validate ehr_id
	attrPath = path + ".ehr_id"
	errs = append(errs, e.EHRID.Validate(attrPath)...)

	// Validate contributions
	if e.Contributions.E {
		for i := range e.Contributions.V {
			attrPath = path + fmt.Sprintf(".contributions[%d]", i)
			if e.Contributions.V[i].Type != CONTRIBUTION_MODEL_NAME {
				errs = append(errs, util.ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid contribution type: %s", e.Contributions.V[i].Type),
					Recommendation: fmt.Sprintf("Ensure contributions[%d] _type field is set to '%s'", i, CONTRIBUTION_MODEL_NAME),
				})
			}
			errs = append(errs, e.Contributions.V[i].Validate(attrPath)...)
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
	if e.Compositions.E {
		for i := range e.Compositions.V {
			attrPath = path + fmt.Sprintf(".compositions[%d]", i)
			if e.Compositions.V[i].Type != VERSIONED_COMPOSITION_MODEL_NAME {
				errs = append(errs, util.ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid composition type: %s", e.Compositions.V[i].Type),
					Recommendation: fmt.Sprintf("Ensure compositions[%d] _type field is set to '%s'", i, VERSIONED_COMPOSITION_MODEL_NAME),
				})
			}
			errs = append(errs, e.Compositions.V[i].Validate(attrPath)...)
		}
	}

	// Validate directory
	if e.Directory.E {
		directory := e.Directory.V
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
	if e.Folders.E {
		for i := range e.Folders.V {
			attrPath = path + fmt.Sprintf(".folders[%d]", i)
			if e.Folders.V[i].Type != VERSIONED_FOLDER_MODEL_NAME {
				errs = append(errs, util.ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid folder type: %s", e.Folders.V[i].Type),
					Recommendation: fmt.Sprintf("Ensure folders[%d] _type field is set to '%s'", i, VERSIONED_FOLDER_MODEL_NAME),
				})
			}
			errs = append(errs, e.Folders.V[i].Validate(attrPath)...)
		}
	}

	// Validate tags
	if e.Tags.E {
		for i := range e.Tags.V {
			attrPath = path + fmt.Sprintf(".tags[%d]", i)
			errs = append(errs, e.Tags.V[i].Validate(attrPath)...)
		}
	}

	return errs
}
