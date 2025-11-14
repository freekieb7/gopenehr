package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ITEM_TREE_MODEL_NAME string = "ITEM_TREE"

type ITEM_TREE struct {
	Type_            util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Items            util.Optional[[]X_ITEM]       `json:"items,omitzero"`
}

func (i ITEM_TREE) isItemStructureModel() {}

func (i ITEM_TREE) HasModelName() bool {
	return i.Type_.E
}

func (i *ITEM_TREE) SetModelName() {
	i.Type_ = util.Some(ITEM_TREE_MODEL_NAME)
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
	if i.Items.E {
		for k := range i.Items.V {
			i.Items.V[k].SetModelName()
		}
	}
}

func (i ITEM_TREE) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != ITEM_TREE_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          ITEM_TREE_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + ITEM_TREE_MODEL_NAME,
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

	// Validate items
	if i.Items.E {
		for k := range i.Items.V {
			attrPath = fmt.Sprintf("%s.items[%d]", path, k)
			errors = append(errors, i.Items.V[k].Validate(attrPath)...)
		}
	}

	return errors
}
