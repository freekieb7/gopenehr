package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const POINT_EVENT_MODEL_NAME string = "POINT_EVENT"

type POINT_EVENT struct {
	Type_            util.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                       `json:"name"`
	ArchetypeNodeID  string                          `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Time             DV_DATE_TIME                    `json:"time"`
	State            util.Optional[X_ITEM_STRUCTURE] `json:"state,omitzero"`
	Data             X_ITEM_STRUCTURE                `json:"data"`
}

func (p POINT_EVENT) isEventModel() {}

func (p POINT_EVENT) HasModelName() bool {
	return p.Type_.E
}

func (p *POINT_EVENT) SetModelName() {
	p.Type_ = util.Some(POINT_EVENT_MODEL_NAME)
	p.Name.SetModelName()
	if p.UID.E {
		p.UID.V.SetModelName()
	}
	if p.Links.E {
		for i := range p.Links.V {
			p.Links.V[i].SetModelName()
		}
	}
	if p.ArchetypeDetails.E {
		p.ArchetypeDetails.V.SetModelName()
	}
	if p.FeederAudit.E {
		p.FeederAudit.V.SetModelName()
	}
	p.Time.SetModelName()
	if p.State.E {
		p.State.V.SetModelName()
	}
	p.Data.SetModelName()
}

func (p *POINT_EVENT) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if p.Type_.E && p.Type_.V != POINT_EVENT_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          POINT_EVENT_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + POINT_EVENT_MODEL_NAME,
			Recommendation: "Set _type to " + POINT_EVENT_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	errs = append(errs, p.Name.Validate(attrPath)...)

	// Validate uid
	if p.UID.E {
		attrPath = path + ".uid"
		errs = append(errs, p.UID.V.Validate(attrPath)...)
	}

	// Validate links
	if p.Links.E {
		for i := range p.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			errs = append(errs, p.Links.V[i].Validate(attrPath)...)
		}
	}

	// Validate archetype_details
	if p.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errs = append(errs, p.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate feeder_audit
	if p.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errs = append(errs, p.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate time
	attrPath = path + ".time"
	errs = append(errs, p.Time.Validate(attrPath)...)

	// Validate state
	if p.State.E {
		attrPath = path + ".state"
		errs = append(errs, p.State.V.Validate(attrPath)...)
	}

	// Validate data
	attrPath = path + ".data"
	errs = append(errs, p.Data.Validate(attrPath)...)

	return errs
}
