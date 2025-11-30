package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const POINT_EVENT_TYPE string = "POINT_EVENT"

type POINT_EVENT struct {
	Type_            utils.Optional[string]             `json:"_type,omitzero"`
	Name             DvTextUnion                        `json:"name"`
	ArchetypeNodeID  string                             `json:"archetype_node_id"`
	UID              utils.Optional[UIDBasedIDUnion]    `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]             `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]         `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]       `json:"feeder_audit,omitzero"`
	Time             DV_DATE_TIME                       `json:"time"`
	State            utils.Optional[ItemStructureUnion] `json:"state,omitzero"`
	Data             ItemStructureUnion                 `json:"data"`
}

func (p *POINT_EVENT) SetModelName() {
	p.Type_ = utils.Some(POINT_EVENT_TYPE)
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

func (p *POINT_EVENT) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if p.Type_.E && p.Type_.V != POINT_EVENT_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          POINT_EVENT_TYPE,
			Path:           attrPath,
			Message:        "_type must be " + POINT_EVENT_TYPE,
			Recommendation: "Set _type to " + POINT_EVENT_TYPE,
		})
	}

	// Validate name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, p.Name.Validate(attrPath).Errs...)

	// Validate uid
	if p.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, p.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if p.Links.E {
		for i := range p.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, p.Links.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if p.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, p.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate feeder_audit
	if p.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, p.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate time
	attrPath = path + ".time"
	validateErr.Errs = append(validateErr.Errs, p.Time.Validate(attrPath).Errs...)

	// Validate state
	if p.State.E {
		attrPath = path + ".state"
		validateErr.Errs = append(validateErr.Errs, p.State.V.Validate(attrPath).Errs...)
	}

	// Validate data
	attrPath = path + ".data"
	validateErr.Errs = append(validateErr.Errs, p.Data.Validate(attrPath).Errs...)

	return validateErr
}
