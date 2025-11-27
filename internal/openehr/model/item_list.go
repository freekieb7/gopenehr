package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ITEM_LIST_MODEL_NAME string = "ITEM_LIST"

type ITEM_LIST struct {
	Type_            utils.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                      `json:"name"`
	ArchetypeNodeID  string                         `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Items            utils.Optional[[]ELEMENT]      `json:"items,omitzero"`
}

func (i *ITEM_LIST) isItemStructureModel() {}

func (i *ITEM_LIST) HasModelName() bool {
	return i.Type_.E
}

func (i *ITEM_LIST) SetModelName() {
	i.Type_ = utils.Some(ITEM_LIST_MODEL_NAME)
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

func (i ITEM_LIST) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != ITEM_LIST_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ITEM_LIST_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + ITEM_LIST_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, i.Name.Validate(attrPath).Errs...)

	// Validate uid
	if i.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, i.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if i.Links.E {
		for j := range i.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, j)
			validateErr.Errs = append(validateErr.Errs, i.Links.V[j].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if i.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, i.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate items
	if i.Items.E {
		for k := range i.Items.V {
			attrPath = fmt.Sprintf("%s.items[%d]", path, k)
			validateErr.Errs = append(validateErr.Errs, i.Items.V[k].Validate(attrPath).Errs...)
		}
	}

	return validateErr
}
