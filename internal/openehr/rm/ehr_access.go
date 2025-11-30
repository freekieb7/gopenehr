package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const EHR_ACCESS_TYPE string = "EHR_ACCESS"

type EHR_ACCESS struct {
	Type_            utils.Optional[string]          `json:"_type,omitzero"`
	Name             DvTextUnion                     `json:"name"`
	ArchetypeNodeID  string                          `json:"archetype_node_id"`
	UID              utils.Optional[UIDBasedIDUnion] `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]          `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]      `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]    `json:"feeder_audit,omitzero"`
	// Settings         utils.Optional[ACCESS_CONTROL_SETTINGS] `json:"settings,omitzero"`
}

func (e *EHR_ACCESS) SetModelName() {
	e.Type_ = utils.Some(EHR_ACCESS_TYPE)
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

func (e *EHR_ACCESS) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if e.Type_.E && e.Type_.V != EHR_ACCESS_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          EHR_ACCESS_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to EHR_ACCESS",
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

	return validateErr
}
