package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

type VersionModel interface {
	isVersionModel()
	SetModelName()
	Validate(path string) util.ValidateError
}

const ORIGINAL_VERSION_MODEL_NAME = "ORIGINAL_VERSION"

type ORIGINAL_VERSION struct {
	Type_                 util.Optional[string]              `json:"_type,omitzero"`
	UID                   OBJECT_VERSION_ID                  `json:"uid"`
	PrecedingVersionUID   util.Optional[OBJECT_VERSION_ID]   `json:"preceding_version_uid,omitzero"`
	OtherInputVersionUIDs util.Optional[[]OBJECT_VERSION_ID] `json:"other_input_version_uids,omitzero"`
	LifecycleState        DV_CODED_TEXT                      `json:"lifecycle_state"`
	Attestations          util.Optional[[]ATTESTATION]       `json:"attestations,omitzero"`
	Data                  VersionModel                       `json:"data"`
}

func (ov *ORIGINAL_VERSION) SetModelName() {
	ov.Type_ = util.Some(ORIGINAL_VERSION_MODEL_NAME)
	ov.UID.SetModelName()
	if ov.PrecedingVersionUID.E {
		ov.PrecedingVersionUID.V.SetModelName()
	}
	if ov.OtherInputVersionUIDs.E {
		for i := range ov.OtherInputVersionUIDs.V {
			ov.OtherInputVersionUIDs.V[i].SetModelName()
		}
	}
	ov.LifecycleState.SetModelName()
	if ov.Attestations.E {
		for i := range ov.Attestations.V {
			ov.Attestations.V[i].SetModelName()
		}
	}
	ov.Data.SetModelName()
}

func (ov *ORIGINAL_VERSION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if ov.Type_.E && ov.Type_.V != "ORIGINAL_VERSION" {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          "ORIGINAL_VERSION",
			Path:           attrPath,
			Message:        "_type must be ORIGINAL_VERSION",
			Recommendation: "Set _type to ORIGINAL_VERSION",
		})
	}

	// Validate uid
	attrPath = path + ".uid"
	validateErr.Errs = append(validateErr.Errs, ov.UID.Validate(attrPath).Errs...)

	// Validate preceding_version_uid
	if ov.PrecedingVersionUID.E {
		attrPath = path + ".preceding_version_uid"
		validateErr.Errs = append(validateErr.Errs, ov.PrecedingVersionUID.V.Validate(attrPath).Errs...)
	}

	// Validate other_input_version_uids
	if ov.OtherInputVersionUIDs.E {
		attrPath = path + ".other_input_version_uids"
		for i := range ov.OtherInputVersionUIDs.V {
			itemPath := fmt.Sprintf("%s[%d]", attrPath, i)
			validateErr.Errs = append(validateErr.Errs, ov.OtherInputVersionUIDs.V[i].Validate(itemPath).Errs...)
		}
	}

	// Validate lifecycle_state
	attrPath = path + ".lifecycle_state"
	validateErr.Errs = append(validateErr.Errs, ov.LifecycleState.Validate(attrPath).Errs...)

	// Validate attestations
	if ov.Attestations.E {
		attrPath = path + ".attestations"
		for i := range ov.Attestations.V {
			itemPath := fmt.Sprintf("%s[%d]", attrPath, i)
			validateErr.Errs = append(validateErr.Errs, ov.Attestations.V[i].Validate(itemPath).Errs...)
		}
	}

	// Validate data
	attrPath = path + ".data"
	validateErr.Errs = append(validateErr.Errs, ov.Data.Validate(attrPath).Errs...)

	return validateErr
}
