package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ATTESTATION_MODEL_NAME = "ATTESTATION"

type ATTESTATION struct {
	Type_         utils.Optional[string]        `json:"_type,omitzero"`
	SystemID      string                        `json:"system_id"`
	TimeCommitted DV_DATE_TIME                  `json:"time_committed"`
	ChangeType    DV_CODED_TEXT                 `json:"change_type"`
	Description   utils.Optional[X_DV_TEXT]     `json:"description,omitzero"`
	Committer     X_PARTY_PROXY                 `json:"committer"`
	AttestedView  utils.Optional[DV_MULTIMEDIA] `json:"attested_view,omitzero"`
	Proof         utils.Optional[string]        `json:"proof,omitzero"`
	Items         utils.Optional[[]DV_EHR_URI]  `json:"items,omitzero"`
	Reason        X_DV_TEXT                     `json:"reason"`
	IsPending     bool                          `json:"is_pending"`
}

func (a *ATTESTATION) SetModelName() {
	a.Type_ = utils.Some(ATTESTATION_MODEL_NAME)
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

func (a *ATTESTATION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ATTESTATION_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ATTESTATION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + ATTESTATION_MODEL_NAME,
			Recommendation: "Set _type to " + ATTESTATION_MODEL_NAME,
		})
	}

	// Validate system_id
	attrPath = path + ".system_id"
	if a.SystemID == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ATTESTATION_MODEL_NAME,
			Path:           attrPath,
			Message:        "system_id cannot be empty",
			Recommendation: "Provide a valid system_id",
		})
	}

	// Validate time_committed
	attrPath = path + ".time_committed"
	validateErr.Errs = append(validateErr.Errs, a.TimeCommitted.Validate(attrPath).Errs...)

	// Validate change_type
	attrPath = path + ".change_type"
	validateErr.Errs = append(validateErr.Errs, a.ChangeType.Validate(attrPath).Errs...)
	// Validate description
	if a.Description.E {
		attrPath = path + ".description"
		validateErr.Errs = append(validateErr.Errs, a.Description.V.Validate(attrPath).Errs...)
	}

	// Validate committer
	attrPath = path + ".committer"
	validateErr.Errs = append(validateErr.Errs, a.Committer.Validate(attrPath).Errs...)

	// Validate attested_view
	if a.AttestedView.E {
		attrPath = path + ".attested_view"
		validateErr.Errs = append(validateErr.Errs, a.AttestedView.V.Validate(attrPath).Errs...)
	}

	// Validate items
	if a.Items.E {
		attrPath = path + ".items"
		for i := range a.Items.V {
			itemPath := fmt.Sprintf("%s[%d]", attrPath, i)
			validateErr.Errs = append(validateErr.Errs, a.Items.V[i].Validate(itemPath).Errs...)
		}
	}

	// Validate reason
	attrPath = path + ".reason"
	validateErr.Errs = append(validateErr.Errs, a.Reason.Validate(attrPath).Errs...)

	return validateErr
}
