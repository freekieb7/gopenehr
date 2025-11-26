package openehr

import (
	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const AUDIT_DETAILS_MODEL_NAME = "AUDIT_DETAILS"

type AUDIT_DETAILS struct {
	Type_         util.Optional[string]  `json:"_type,omitzero"`
	SystemID      string                 `json:"system_id"`
	TimeCommitted DV_DATE_TIME           `json:"time_committed"`
	ChangeType    DV_CODED_TEXT          `json:"change_type"`
	Description   util.Optional[DV_TEXT] `json:"description,omitzero"`
	Committer     X_PARTY_PROXY          `json:"committer"`
}

func (a *AUDIT_DETAILS) SetModelName() {
	a.Type_ = util.Some(AUDIT_DETAILS_MODEL_NAME)
	a.TimeCommitted.SetModelName()
	a.ChangeType.SetModelName()
	if a.Description.E {
		a.Description.V.SetModelName()
	}
	a.Committer.SetModelName()
}

func (a *AUDIT_DETAILS) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != AUDIT_DETAILS_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          AUDIT_DETAILS_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + AUDIT_DETAILS_MODEL_NAME,
			Recommendation: "Set _type to " + AUDIT_DETAILS_MODEL_NAME,
		})
	}

	// Validate system_id
	attrPath = path + ".system_id"
	if a.SystemID == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          AUDIT_DETAILS_MODEL_NAME,
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
	if !terminology.IsValidAuditChangeTypeCode(a.ChangeType.DefiningCode.CodeString) {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          AUDIT_DETAILS_MODEL_NAME,
			Path:           attrPath,
			Message:        "Invalid change_type code",
			Recommendation: "Provide a valid change_type code",
		})
	}
	if terminology.GetAuditChangeTypeName(a.ChangeType.DefiningCode.CodeString) != a.ChangeType.Value {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          AUDIT_DETAILS_MODEL_NAME,
			Path:           attrPath,
			Message:        "change_type value does not match code",
			Recommendation: "Provide a matching change_type value for the given code",
		})
	}

	validateErr.Errs = append(validateErr.Errs, a.ChangeType.Validate(attrPath).Errs...)

	// Validate committer
	attrPath = path + ".committer"
	validateErr.Errs = append(validateErr.Errs, a.Committer.Validate(attrPath).Errs...)
	return validateErr
}
