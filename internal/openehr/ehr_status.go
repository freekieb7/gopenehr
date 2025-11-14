package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const EHR_STATUS_MODEL_NAME string = "EHR_STATUS"

type EHR_STATUS struct {
	Type_            util.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                       `json:"name"`
	ArchetypeNodeID  string                          `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Subject          PARTY_SELF                      `json:"subject"`
	IsQueryable      bool                            `json:"is_queryable"`
	IsModifiable     bool                            `json:"is_modifiable"`
	OtherDetails     util.Optional[X_ITEM_STRUCTURE] `json:"other_details,omitzero"`
}

func (e *EHR_STATUS) SetModelName() {
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
	e.Subject.SetModelName()
	if e.FeederAudit.E {
		e.FeederAudit.V.SetModelName()
	}
	if e.OtherDetails.E {
		e.OtherDetails.V.SetModelName()
	}
}

func (e EHR_STATUS) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if e.Type_.E && e.Type_.V != "EHR_STATUS" {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          "EHR_STATUS",
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to EHR_STATUS",
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
			itemPath := fmt.Sprintf("%s.links[%d]", path, i)
			errors = append(errors, e.Links.V[i].Validate(itemPath)...)
		}
	}

	// Validate archetype_details
	if e.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errors = append(errors, e.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate subject
	attrPath = path + ".subject"
	errors = append(errors, e.Subject.Validate(attrPath)...)

	// Validate feeder_audit
	if e.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errors = append(errors, e.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate other_details
	if e.OtherDetails.E {
		attrPath = path + ".other_details"
		errors = append(errors, e.OtherDetails.V.Validate(attrPath)...)
	}

	return errors
}
