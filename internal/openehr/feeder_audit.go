package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const FEEDER_AUDIT_MODEL_NAME string = "FEEDER_AUDIT"

type FEEDER_AUDIT struct {
	Type_                    util.Optional[string]               `json:"_type,omitzero"`
	OriginatingSystemItemIDs util.Optional[[]DV_IDENTIFIER]      `json:"originating_system_item_ids,omitzero"`
	FeederSystemItemIDs      util.Optional[[]DV_IDENTIFIER]      `json:"feeder_system_item_ids,omitzero"`
	OriginalContent          util.Optional[X_DV_ENCAPSULATED]    `json:"original_content,omitzero"`
	OriginatingSystemAudit   FEEDER_AUDIT_DETAILS                `json:"originating_system_audit"`
	FeederSystemAudit        util.Optional[FEEDER_AUDIT_DETAILS] `json:"feeder_system_audit,omitzero"`
}

func (f *FEEDER_AUDIT) SetModelName() {
	f.Type_ = util.Some(FEEDER_AUDIT_MODEL_NAME)
	if f.OriginatingSystemItemIDs.E {
		for i := range f.OriginatingSystemItemIDs.V {
			f.OriginatingSystemItemIDs.V[i].SetModelName()
		}
	}
	if f.FeederSystemItemIDs.E {
		for i := range f.FeederSystemItemIDs.V {
			f.FeederSystemItemIDs.V[i].SetModelName()
		}
	}
	if f.OriginalContent.E {
		f.OriginalContent.V.SetModelName()
	}
	f.OriginatingSystemAudit.SetModelName()
	if f.FeederSystemAudit.E {
		f.FeederSystemAudit.V.SetModelName()
	}
}

func (f *FEEDER_AUDIT) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if f.Type_.E && f.Type_.V != FEEDER_AUDIT_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          FEEDER_AUDIT_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to FEEDER_AUDIT",
		})
	}

	// Validate originating_system_audit
	attrPath = path + ".originating_system_audit"
	validateErr.Errs = append(validateErr.Errs, f.OriginatingSystemAudit.Validate(attrPath).Errs...)

	// Validate feeder_system_audit
	if f.FeederSystemAudit.E {
		attrPath = path + ".feeder_system_audit"
		validateErr.Errs = append(validateErr.Errs, f.FeederSystemAudit.V.Validate(attrPath).Errs...)
	}

	return validateErr
}
