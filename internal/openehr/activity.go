package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ACTIVITY_MODEL_NAME string = "ACTIVITY"

type ACTIVITY struct {
	Type_             util.Optional[string]         `json:"_type,omitzero"`
	Name              X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID   string                        `json:"archetype_node_id"`
	UID               util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links             util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails  util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit       util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Timing            util.Optional[DV_PARSABLE]    `json:"timing,omitzero"`
	ActionArchetypeID string                        `json:"action_archetype_id"`
	Description       X_ITEM_STRUCTURE              `json:"description"`
}

func (a *ACTIVITY) isContentItemModel() {}

func (a *ACTIVITY) HasModelName() bool {
	return a.Type_.E
}

func (a *ACTIVITY) SetModelName() {
	a.Type_ = util.Some(ACTIVITY_MODEL_NAME)
	a.Name.SetModelName()
	if a.UID.E {
		a.UID.V.SetModelName()
	}
	if a.Links.E {
		for i := range a.Links.V {
			a.Links.V[i].SetModelName()
		}
	}
	if a.ArchetypeDetails.E {
		a.ArchetypeDetails.V.SetModelName()
	}
	if a.FeederAudit.E {
		a.FeederAudit.V.SetModelName()
	}
	if a.Timing.E {
		a.Timing.V.SetModelName()
	}
	a.Description.SetModelName()
}

func (a *ACTIVITY) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ACTIVITY_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          ACTIVITY_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + ACTIVITY_MODEL_NAME,
			Recommendation: "Set _type to " + ACTIVITY_MODEL_NAME,
		})
	}

	// Validate Name
	attrPath = path + ".name"
	errs = append(errs, a.Name.Validate(attrPath)...)

	// Validate UID
	if a.UID.E {
		attrPath = path + ".uid"
		errs = append(errs, a.UID.V.Validate(attrPath)...)
	}

	// Validate Links
	if a.Links.E {
		for i := range a.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			errs = append(errs, a.Links.V[i].Validate(attrPath)...)
		}
	}

	// Validate ArchetypeDetails
	if a.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errs = append(errs, a.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate FeederAudit
	if a.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errs = append(errs, a.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate Timing
	if a.Timing.E {
		attrPath = path + ".timing"
		errs = append(errs, a.Timing.V.Validate(attrPath)...)
	}

	// Validate Description
	attrPath = path + ".description"
	errs = append(errs, a.Description.Validate(attrPath)...)

	return errs
}
