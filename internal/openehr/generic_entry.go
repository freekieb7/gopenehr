package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const GENERIC_ENTRY_MODEL_NAME string = "GENERIC_ENTRY"

type GENERIC_ENTRY struct {
	Type_            util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Data             X_ITEM                        `json:"data"`
}

func (g GENERIC_ENTRY) isContentItemModel() {}

func (g GENERIC_ENTRY) HasModelName() bool {
	return g.Type_.E
}

func (g *GENERIC_ENTRY) SetModelName() {
	g.Type_ = util.Some(GENERIC_ENTRY_MODEL_NAME)
	g.Name.SetModelName()
	if g.UID.E {
		g.UID.V.SetModelName()
	}
	if g.Links.E {
		for i := range g.Links.V {
			g.Links.V[i].SetModelName()
		}
	}
	if g.ArchetypeDetails.E {
		g.ArchetypeDetails.V.SetModelName()
	}
	if g.FeederAudit.E {
		g.FeederAudit.V.SetModelName()
	}
	g.Data.SetModelName()
}

func (g *GENERIC_ENTRY) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if g.Type_.E && g.Type_.V != GENERIC_ENTRY_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          GENERIC_ENTRY_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", GENERIC_ENTRY_MODEL_NAME, g.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", GENERIC_ENTRY_MODEL_NAME),
		})
	}

	// Validate name
	attrPath = path + ".name"
	errs = append(errs, g.Name.Validate(attrPath)...)

	// Validate uid
	if g.UID.E {
		attrPath = path + ".uid"
		errs = append(errs, g.UID.V.Validate(attrPath)...)
	}

	// Validate links
	if g.Links.E {
		for i := range g.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			errs = append(errs, g.Links.V[i].Validate(attrPath)...)
		}
	}

	// Validate archetype_details
	if g.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errs = append(errs, g.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate feeder_audit
	if g.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errs = append(errs, g.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate data
	attrPath = path + ".data"
	errs = append(errs, g.Data.Validate(attrPath)...)

	return errs
}
