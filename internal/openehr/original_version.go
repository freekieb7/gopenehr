package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

type VersionModel interface {
	isVersionModel()
	SetModelName()
	Validate(path string) []util.ValidationError
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

func (ov *ORIGINAL_VERSION) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if ov.Type_.E && ov.Type_.V != "ORIGINAL_VERSION" {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          "ORIGINAL_VERSION",
			Path:           attrPath,
			Message:        "_type must be ORIGINAL_VERSION",
			Recommendation: "Set _type to ORIGINAL_VERSION",
		})
	}

	// Validate uid
	attrPath = path + ".uid"
	errs = append(errs, ov.UID.Validate(attrPath)...)

	// Validate preceding_version_uid
	if ov.PrecedingVersionUID.E {
		attrPath = path + ".preceding_version_uid"
		errs = append(errs, ov.PrecedingVersionUID.V.Validate(attrPath)...)
	}

	// Validate other_input_version_uids
	if ov.OtherInputVersionUIDs.E {
		attrPath = path + ".other_input_version_uids"
		for i := range ov.OtherInputVersionUIDs.V {
			uidPath := fmt.Sprintf("%s[%d]", attrPath, i)
			errs = append(errs, ov.OtherInputVersionUIDs.V[i].Validate(uidPath)...)
		}
	}

	// Validate lifecycle_state
	attrPath = path + ".lifecycle_state"
	errs = append(errs, ov.LifecycleState.Validate(attrPath)...)

	// Validate attestations
	if ov.Attestations.E {
		attrPath = path + ".attestations"
		for i := range ov.Attestations.V {
			attPath := fmt.Sprintf("%s[%d]", attrPath, i)
			errs = append(errs, ov.Attestations.V[i].Validate(attPath)...)
		}
	}

	// Validate data
	attrPath = path + ".data"
	errs = append(errs, ov.Data.Validate(attrPath)...)

	return errs
}
