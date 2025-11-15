package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const HISTORY_MODEL_NAME string = "HISTORY"

type HistoryEventModel interface {
	isHistoryEventModel()
	SetModelName()
	Validate(path string) []util.ValidationError
}

type HISTORY struct {
	Type_            util.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                       `json:"name"`
	ArchetypeNodeID  string                          `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Origin           DV_DATE_TIME                    `json:"origin"`
	Period           util.Optional[DV_DURATION]      `json:"period,omitzero"`
	Duration         util.Optional[DV_DURATION]      `json:"duration,omitzero"`
	Summary          util.Optional[X_ITEM_STRUCTURE] `json:"summary,omitzero"`
	Events           util.Optional[[]X_EVENT]        `json:"events,omitzero"`
}

func (h *HISTORY) SetModelName() {
	h.Type_ = util.Some(HISTORY_MODEL_NAME)
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

func (h *HISTORY) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if h.Type_.E && h.Type_.V != HISTORY_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          HISTORY_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + HISTORY_MODEL_NAME,
			Recommendation: "Set _type to " + HISTORY_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	errs = append(errs, h.Name.Validate(attrPath)...)

	// Validate uid
	if h.UID.E {
		attrPath = path + ".uid"
		errs = append(errs, h.UID.V.Validate(attrPath)...)
	}

	// Validate links
	if h.Links.E {
		for i := range h.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			errs = append(errs, h.Links.V[i].Validate(attrPath)...)
		}
	}

	// Validate archetype_details
	if h.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errs = append(errs, h.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate feeder_audit
	if h.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errs = append(errs, h.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate origin
	attrPath = path + ".origin"
	errs = append(errs, h.Origin.Validate(attrPath)...)

	// Validate period
	if h.Period.E {
		attrPath = path + ".period"
		errs = append(errs, h.Period.V.Validate(attrPath)...)
	}

	// Validate duration
	if h.Duration.E {
		attrPath = path + ".duration"
		errs = append(errs, h.Duration.V.Validate(attrPath)...)
	}

	// Validate summary
	if h.Summary.E {
		attrPath = path + ".summary"
		errs = append(errs, h.Summary.V.Validate(attrPath)...)
	}

	// Validate events
	if h.Events.E {
		for i := range h.Events.V {
			attrPath = fmt.Sprintf("%s.events[%d]", path, i)
			errs = append(errs, h.Events.V[i].Validate(attrPath)...)
		}
	}

	return errs
}
