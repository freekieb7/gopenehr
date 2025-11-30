package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const FEEDER_AUDIT_TYPE string = "FEEDER_AUDIT"

type FEEDER_AUDIT struct {
	Type_                    utils.Optional[string]               `json:"_type,omitzero"`
	OriginatingSystemItemIDs utils.Optional[[]DV_IDENTIFIER]      `json:"originating_system_item_ids,omitzero"`
	FeederSystemItemIDs      utils.Optional[[]DV_IDENTIFIER]      `json:"feeder_system_item_ids,omitzero"`
	OriginalContent          utils.Optional[DvEncapsulatedUnion]  `json:"original_content,omitzero"`
	OriginatingSystemAudit   FEEDER_AUDIT_DETAILS                 `json:"originating_system_audit"`
	FeederSystemAudit        utils.Optional[FEEDER_AUDIT_DETAILS] `json:"feeder_system_audit,omitzero"`
}

func (f *FEEDER_AUDIT) SetModelName() {
	f.Type_ = utils.Some(FEEDER_AUDIT_TYPE)
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
	if f.Type_.E && f.Type_.V != FEEDER_AUDIT_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          FEEDER_AUDIT_TYPE,
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
