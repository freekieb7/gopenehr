package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ELEMENT_MODEL_NAME string = "ELEMENT"

type ELEMENT struct {
	Type_            util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	NullFlavour      util.Optional[DV_CODED_TEXT]  `json:"null_flavour,omitzero"`
	Value            util.Optional[X_DATA_VALUE]   `json:"value,omitzero"`
	NullReason       util.Optional[X_DV_TEXT]      `json:"null_reason,omitzero"`
}

func (e ELEMENT) isItemModel() {}

func (e ELEMENT) HasModelName() bool {
	return e.Type_.E
}

func (e *ELEMENT) SetModelName() {
	e.Type_ = util.Some(ELEMENT_MODEL_NAME)
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

func (e ELEMENT) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if e.Type_.E && e.Type_.V != ELEMENT_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          ELEMENT_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + ELEMENT_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	errors = append(errors, e.Name.Validate(attrPath)...)

	// Validate uid
	if e.UID.E {
		attrPath = path + ".uid"
		errors = append(errors, e.UID.V.Validate(attrPath)...)
	}

	// Validate links
	if e.Links.E {
		for i := range e.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			errors = append(errors, e.Links.V[i].Validate(attrPath)...)
		}
	}

	// Validate archetype_details
	if e.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errors = append(errors, e.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate feeder_audit
	if e.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errors = append(errors, e.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate null_flavour
	if e.NullFlavour.E {
		attrPath = path + ".null_flavour"
		errors = append(errors, e.NullFlavour.V.Validate(attrPath)...)
	}

	// Validate value
	if e.Value.E {
		attrPath = path + ".value"
		errors = append(errors, e.Value.V.Validate(attrPath)...)
	}

	// Validate null_reason
	if e.NullReason.E {
		attrPath = path + ".null_reason"
		errors = append(errors, e.NullReason.V.Validate(attrPath)...)
	}

	return errors
}
