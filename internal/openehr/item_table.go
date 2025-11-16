package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ITEM_TABLE_MODEL_NAME string = "ITEM_TABLE"

type ITEM_TABLE struct {
	Type_            util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Rows             util.Optional[[]CLUSTER]      `json:"rows,omitzero"`
}

func (i *ITEM_TABLE) isItemStructureModel() {}

func (i *ITEM_TABLE) HasModelName() bool {
	return i.Type_.E
}

func (i *ITEM_TABLE) SetModelName() {
	i.Type_ = util.Some(ITEM_TABLE_MODEL_NAME)
	i.Name.SetModelName()
	if i.UID.E {
		i.UID.V.SetModelName()
	}
	if i.Links.E {
		for j := range i.Links.V {
			i.Links.V[j].SetModelName()
		}
	}
	if i.ArchetypeDetails.E {
		i.ArchetypeDetails.V.SetModelName()
	}
	if i.FeederAudit.E {
		i.FeederAudit.V.SetModelName()
	}
	if i.Rows.E {
		for k := range i.Rows.V {
			i.Rows.V[k].SetModelName()
		}
	}
}

func (i *ITEM_TABLE) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != ITEM_TABLE_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          ITEM_TABLE_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + ITEM_TABLE_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	errors = append(errors, i.Name.Validate(attrPath)...)

	// Validate uid
	if i.UID.E {
		attrPath = path + ".uid"
		errors = append(errors, i.UID.V.Validate(attrPath)...)
	}

	// Validate links
	if i.Links.E {
		for j := range i.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, j)
			errors = append(errors, i.Links.V[j].Validate(attrPath)...)
		}
	}

	// Validate archetype_details
	if i.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errors = append(errors, i.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate rows
	if i.Rows.E {
		for k := range i.Rows.V {
			attrPath = fmt.Sprintf("%s.rows[%d]", path, k)
			errors = append(errors, i.Rows.V[k].Validate(attrPath)...)
		}
	}

	return errors
}
