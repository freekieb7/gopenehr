package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const HISTORY_TYPE string = "HISTORY"

type HISTORY struct {
	Type_            utils.Optional[string]             `json:"_type,omitzero"`
	Name             DvTextUnion                        `json:"name"`
	ArchetypeNodeID  string                             `json:"archetype_node_id"`
	UID              utils.Optional[UIDBasedIDUnion]    `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]             `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]         `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]       `json:"feeder_audit,omitzero"`
	Origin           DV_DATE_TIME                       `json:"origin"`
	Period           utils.Optional[DV_DURATION]        `json:"period,omitzero"`
	Duration         utils.Optional[DV_DURATION]        `json:"duration,omitzero"`
	Summary          utils.Optional[ItemStructureUnion] `json:"summary,omitzero"`
	Events           utils.Optional[[]EventUnion]       `json:"events,omitzero"`
}

func (h *HISTORY) SetModelName() {
	h.Type_ = utils.Some(HISTORY_TYPE)
	h.Name.SetModelName()
	if h.UID.E {
		h.UID.V.SetModelName()
	}
	if h.Links.E {
		for i := range h.Links.V {
			h.Links.V[i].SetModelName()
		}
	}
	if h.ArchetypeDetails.E {
		h.ArchetypeDetails.V.SetModelName()
	}
	if h.FeederAudit.E {
		h.FeederAudit.V.SetModelName()
	}
	h.Origin.SetModelName()
	if h.Period.E {
		h.Period.V.SetModelName()
	}
	if h.Duration.E {
		h.Duration.V.SetModelName()
	}
	if h.Summary.E {
		h.Summary.V.SetModelName()
	}
	if h.Events.E {
		for i := range h.Events.V {
			h.Events.V[i].SetModelName()
		}
	}
}

func (h *HISTORY) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if h.Type_.E && h.Type_.V != HISTORY_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          HISTORY_TYPE,
			Path:           attrPath,
			Message:        "_type must be " + HISTORY_TYPE,
			Recommendation: "Set _type to " + HISTORY_TYPE,
		})
	}

	// Validate name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, h.Name.Validate(attrPath).Errs...)

	// Validate uid
	if h.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, h.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if h.Links.E {
		for i := range h.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, h.Links.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if h.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, h.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate feeder_audit
	if h.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, h.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate origin
	attrPath = path + ".origin"
	validateErr.Errs = append(validateErr.Errs, h.Origin.Validate(attrPath).Errs...)

	// Validate period
	if h.Period.E {
		attrPath = path + ".period"
		validateErr.Errs = append(validateErr.Errs, h.Period.V.Validate(attrPath).Errs...)
	}

	// Validate duration
	if h.Duration.E {
		attrPath = path + ".duration"
		validateErr.Errs = append(validateErr.Errs, h.Duration.V.Validate(attrPath).Errs...)
	}

	// Validate summary
	if h.Summary.E {
		attrPath = path + ".summary"
		validateErr.Errs = append(validateErr.Errs, h.Summary.V.Validate(attrPath).Errs...)
	}

	// Validate events
	if h.Events.E {
		for i := range h.Events.V {
			attrPath = fmt.Sprintf("%s.events[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, h.Events.V[i].Validate(attrPath).Errs...)
		}
	}

	return validateErr
}
