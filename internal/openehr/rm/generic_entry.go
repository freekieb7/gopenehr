package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const GENERIC_ENTRY_MODEL_NAME string = "GENERIC_ENTRY"

type GENERIC_ENTRY struct {
	Type_            utils.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                      `json:"name"`
	ArchetypeNodeID  string                         `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Data             X_ITEM                         `json:"data"`
}

func (g *GENERIC_ENTRY) isContentItemModel() {}

func (g *GENERIC_ENTRY) HasModelName() bool {
	return g.Type_.E
}

func (g *GENERIC_ENTRY) SetModelName() {
	g.Type_ = utils.Some(GENERIC_ENTRY_MODEL_NAME)
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

func (g *GENERIC_ENTRY) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if g.Type_.E && g.Type_.V != GENERIC_ENTRY_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          GENERIC_ENTRY_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", GENERIC_ENTRY_MODEL_NAME, g.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", GENERIC_ENTRY_MODEL_NAME),
		})
	}

	// Validate name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, g.Name.Validate(attrPath).Errs...)

	// Validate uid
	if g.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, g.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if g.Links.E {
		for i := range g.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, g.Links.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if g.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, g.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate feeder_audit
	if g.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, g.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate data
	attrPath = path + ".data"
	validateErr.Errs = append(validateErr.Errs, g.Data.Validate(attrPath).Errs...)

	return validateErr
}
