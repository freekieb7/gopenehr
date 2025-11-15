package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const EHR_ACCESS_MODEL_NAME string = "EHR_ACCESS"

type EHR_ACCESS struct {
	MetaType         util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	// Settings         util.Optional[ACCESS_CONTROL_SETTINGS] `json:"settings,omitzero"`
}

func (e EHR_ACCESS) isVersionModel() {}

func (e *EHR_ACCESS) SetModelName() {
	e.MetaType = util.Some(EHR_ACCESS_MODEL_NAME)
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
}

func (e EHR_ACCESS) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if e.MetaType.E && e.MetaType.V != "EHR_ACCESS" {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          "EHR_ACCESS",
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to EHR_ACCESS",
		})
	}

	// Validate name
	if err := e.Name.Validate(path + ".name"); len(err) > 0 {
		errors = append(errors, err...)
	}

	// Validate uid
	if e.UID.E {
		if err := e.UID.V.Validate(path + ".uid"); len(err) > 0 {
			errors = append(errors, err...)
		}
	}

	// Validate links
	if e.Links.E {
		for i := range e.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			if err := e.Links.V[i].Validate(attrPath); len(err) > 0 {
				errors = append(errors, err...)
			}
		}
	}

	// Validate archetype_details
	if e.ArchetypeDetails.E {
		if err := e.ArchetypeDetails.V.Validate(path + ".archetype_details"); len(err) > 0 {
			errors = append(errors, err...)
		}
	}

	// Validate feeder_audit
	if e.FeederAudit.E {
		if err := e.FeederAudit.V.Validate(path + ".feeder_audit"); len(err) > 0 {
			errors = append(errors, err...)
		}
	}

	return errors
}
