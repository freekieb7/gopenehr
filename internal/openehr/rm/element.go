package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ELEMENT_MODEL_NAME string = "ELEMENT"

type ELEMENT struct {
	Type_            utils.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                      `json:"name"`
	ArchetypeNodeID  string                         `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	NullFlavour      utils.Optional[DV_CODED_TEXT]  `json:"null_flavour,omitzero"`
	Value            utils.Optional[X_DATA_VALUE]   `json:"value,omitzero"`
	NullReason       utils.Optional[X_DV_TEXT]      `json:"null_reason,omitzero"`
}

func (e *ELEMENT) isItemModel() {}

func (e *ELEMENT) HasModelName() bool {
	return e.Type_.E
}

func (e *ELEMENT) SetModelName() {
	e.Type_ = utils.Some(ELEMENT_MODEL_NAME)
	e.Name.SetModelName()
	if e.UID.E {
		e.UID.V.SetModelName()
	}
	if e.Links.E {
		for i := range e.Links.V {
			e.Links.V[i].SetModelName()
		}
	}
	if e.ArchetypeDetails.E {
		e.ArchetypeDetails.V.SetModelName()
	}
	if e.FeederAudit.E {
		e.FeederAudit.V.SetModelName()
	}
	if e.NullFlavour.E {
		e.NullFlavour.V.SetModelName()
	}
	if e.Value.E {
		e.Value.V.SetModelName()
	}
	if e.NullReason.E {
		e.NullReason.V.SetModelName()
	}
}

func (e *ELEMENT) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if e.Type_.E && e.Type_.V != ELEMENT_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ELEMENT_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + ELEMENT_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, e.Name.Validate(attrPath).Errs...)

	// Validate uid
	if e.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, e.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if e.Links.E {
		for i := range e.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, e.Links.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if e.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, e.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate feeder_audit
	if e.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, e.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate null_flavour
	if e.NullFlavour.E {
		attrPath = path + ".null_flavour"
		validateErr.Errs = append(validateErr.Errs, e.NullFlavour.V.Validate(attrPath).Errs...)
	}

	// Validate value
	if e.Value.E {
		attrPath = path + ".value"
		validateErr.Errs = append(validateErr.Errs, e.Value.V.Validate(attrPath).Errs...)
	}

	// Validate null_reason
	if e.NullReason.E {
		attrPath = path + ".null_reason"
		validateErr.Errs = append(validateErr.Errs, e.NullReason.V.Validate(attrPath).Errs...)
	}

	return validateErr
}
