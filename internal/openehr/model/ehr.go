package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const EHR_MODEL_NAME string = "EHR"

type EHR struct {
	Type_         utils.Optional[string]       `json:"_type,omitzero"`
	SystemID      HIER_OBJECT_ID               `json:"system_id"`
	EHRID         HIER_OBJECT_ID               `json:"ehr_id"`
	Contributions utils.Optional[[]OBJECT_REF] `json:"contributions,omitzero"`
	EHRStatus     OBJECT_REF                   `json:"ehr_status"`
	EHRAccess     OBJECT_REF                   `json:"ehr_access"`
	Compositions  utils.Optional[[]OBJECT_REF] `json:"compositions,omitzero"`
	Directory     utils.Optional[OBJECT_REF]   `json:"directory,omitzero"`
	TimeCreated   DV_DATE_TIME                 `json:"time_created"`
	Folders       utils.Optional[[]OBJECT_REF] `json:"folders,omitzero"`
	Tags          utils.Optional[[]OBJECT_REF] `json:"tags,omitzero"`
}

func (e *EHR) SetModelName() {
	e.Type_ = utils.Some(EHR_MODEL_NAME)
	e.SystemID.SetModelName()
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

func (e *EHR) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if e.Type_.E && e.Type_.V != EHR_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          EHR_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid EHR _type field: %s", e.Type_.V),
			Recommendation: "Ensure the _type field is set to 'EHR'",
		})
	}

	// Validate system_id
	attrPath = path + ".system_id"
	validateErr.Errs = append(validateErr.Errs, e.SystemID.Validate(attrPath).Errs...)

	// Validate ehr_id
	attrPath = path + ".ehr_id"
	validateErr.Errs = append(validateErr.Errs, e.EHRID.Validate(attrPath).Errs...)
	// Validate contributions
	if e.Contributions.E {
		for i := range e.Contributions.V {
			attrPath = path + fmt.Sprintf(".contributions[%d]", i)
			if e.Contributions.V[i].Type != CONTRIBUTION_MODEL_NAME {
				validateErr.Errs = append(validateErr.Errs, util.ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid contribution type: %s", e.Contributions.V[i].Type),
					Recommendation: fmt.Sprintf("Ensure contributions[%d] _type field is set to '%s'", i, CONTRIBUTION_MODEL_NAME),
				})
			}
			validateErr.Errs = append(validateErr.Errs, e.Contributions.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate ehr_status
	attrPath = path + ".ehr_status"
	if e.EHRStatus.Type != EHR_STATUS_MODEL_NAME && e.EHRStatus.Type != VERSIONED_EHR_STATUS_MODEL_NAME {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          EHR_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid EHR status type: %s", e.EHRStatus.Type),
			Recommendation: fmt.Sprintf("Ensure ehr_status _type field is set to '%s' or '%s'", EHR_STATUS_MODEL_NAME, VERSIONED_EHR_STATUS_MODEL_NAME),
		})

	}
	validateErr.Errs = append(validateErr.Errs, e.EHRStatus.Validate(attrPath).Errs...)

	// Validate ehr_access
	attrPath = path + ".ehr_access"
	if e.EHRAccess.Type != EHR_ACCESS_MODEL_NAME && e.EHRAccess.Type != VERSIONED_EHR_ACCESS_MODEL_NAME {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          EHR_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid EHR access type: %s", e.EHRAccess.Type),
			Recommendation: fmt.Sprintf("Ensure ehr_access _type field is set to '%s' or '%s'", EHR_ACCESS_MODEL_NAME, VERSIONED_EHR_ACCESS_MODEL_NAME),
		})
	}
	validateErr.Errs = append(validateErr.Errs, e.EHRAccess.Validate(attrPath).Errs...)

	// Validate compositions
	if e.Compositions.E {
		for i := range e.Compositions.V {
			attrPath = path + fmt.Sprintf(".compositions[%d]", i)
			if e.Compositions.V[i].Type != COMPOSITION_MODEL_NAME {
				validateErr.Errs = append(validateErr.Errs, util.ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid composition type: %s", e.Compositions.V[i].Type),
					Recommendation: fmt.Sprintf("Ensure compositions[%d] _type field is set to '%s'", i, COMPOSITION_MODEL_NAME),
				})
			}
			validateErr.Errs = append(validateErr.Errs, e.Compositions.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate directory
	if e.Directory.E {
		directory := e.Directory.V
		attrPath = path + ".directory"
		if directory.Type != FOLDER_MODEL_NAME {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          EHR_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid folder type: %s", directory.Type),
				Recommendation: fmt.Sprintf("Ensure directory _type field is set to '%s'", FOLDER_MODEL_NAME),
			})
		}

		validateErr.Errs = append(validateErr.Errs, directory.Validate(attrPath).Errs...)
	}

	// Validate time_created
	attrPath = path + ".time_created"
	validateErr.Errs = append(validateErr.Errs, e.TimeCreated.Validate(attrPath).Errs...)

	// Validate folders
	if e.Folders.E {
		for i := range e.Folders.V {
			attrPath = path + fmt.Sprintf(".folders[%d]", i)
			if e.Folders.V[i].Type != FOLDER_MODEL_NAME {
				validateErr.Errs = append(validateErr.Errs, util.ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid folder type: %s", e.Folders.V[i].Type),
					Recommendation: fmt.Sprintf("Ensure folders[%d] _type field is set to '%s'", i, FOLDER_MODEL_NAME),
				})
			}
			validateErr.Errs = append(validateErr.Errs, e.Folders.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate tags
	if e.Tags.E {
		for i := range e.Tags.V {
			attrPath = path + fmt.Sprintf(".tags[%d]", i)
			validateErr.Errs = append(validateErr.Errs, e.Tags.V[i].Validate(attrPath).Errs...)
		}
	}

	return validateErr
}
