package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ATTESTATION_MODEL_NAME = "ATTESTATION"

type ATTESTATION struct {
	Type_         util.Optional[string]        `json:"_type,omitzero"`
	SystemID      string                       `json:"system_id"`
	TimeCommitted DV_DATE_TIME                 `json:"time_committed"`
	ChangeType    DV_CODED_TEXT                `json:"change_type"`
	Description   util.Optional[X_DV_TEXT]     `json:"description,omitzero"`
	Committer     X_PARTY_PROXY                `json:"committer"`
	AttestedView  util.Optional[DV_MULTIMEDIA] `json:"attested_view,omitzero"`
	Proof         util.Optional[string]        `json:"proof,omitzero"`
	Items         util.Optional[[]DV_EHR_URI]  `json:"items,omitzero"`
	Reason        X_DV_TEXT                    `json:"reason"`
	IsPending     bool                         `json:"is_pending"`
}

func (a *ATTESTATION) SetModelName() {
	a.Type_ = util.Some(ATTESTATION_MODEL_NAME)
	a.TimeCommitted.SetModelName()
	a.ChangeType.SetModelName()
	if a.Description.E {
		a.Description.V.SetModelName()
	}
	a.Committer.SetModelName()
	if a.AttestedView.E {
		a.AttestedView.V.SetModelName()
	}
	for i := range a.Items.V {
		a.Items.V[i].SetModelName()
	}
	a.Reason.SetModelName()
}

func (a *ATTESTATION) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ATTESTATION_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          ATTESTATION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + ATTESTATION_MODEL_NAME,
			Recommendation: "Set _type to " + ATTESTATION_MODEL_NAME,
		})
	}

	// Validate system_id
	attrPath = path + ".system_id"
	if a.SystemID == "" {
		errs = append(errs, util.ValidationError{
			Model:          ATTESTATION_MODEL_NAME,
			Path:           attrPath,
			Message:        "system_id cannot be empty",
			Recommendation: "Provide a valid system_id",
		})
	}

	// Validate time_committed
	attrPath = path + ".time_committed"
	errs = append(errs, a.TimeCommitted.Validate(attrPath)...)

	// Validate change_type
	attrPath = path + ".change_type"
	errs = append(errs, a.ChangeType.Validate(attrPath)...)

	// Validate description
	if a.Description.E {
		attrPath = path + ".description"
		errs = append(errs, a.Description.V.Validate(attrPath)...)
	}

	// Validate committer
	attrPath = path + ".committer"
	errs = append(errs, a.Committer.Validate(attrPath)...)

	// Validate attested_view
	if a.AttestedView.E {
		attrPath = path + ".attested_view"
		errs = append(errs, a.AttestedView.V.Validate(attrPath)...)
	}

	// Validate items
	if a.Items.E {
		attrPath = path + ".items"
		for i := range a.Items.V {
			itemPath := fmt.Sprintf("%s[%d]", attrPath, i)
			errs = append(errs, a.Items.V[i].Validate(itemPath)...)
		}
	}

	// Validate reason
	attrPath = path + ".reason"
	errs = append(errs, a.Reason.Validate(attrPath)...)

	return errs
}
