package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ITEM_SINGLE_MODEL_NAME string = "ITEM_SINGLE"

type ITEM_SINGLE struct {
	Type_            utils.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                      `json:"name"`
	ArchetypeNodeID  string                         `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            utils.Optional[LINK]           `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Item             ELEMENT                        `json:"item"`
}

func (i *ITEM_SINGLE) isItemStructureModel() {}

func (i *ITEM_SINGLE) HasModelName() bool {
	return i.Type_.E
}

func (i *ITEM_SINGLE) SetModelName() {
	i.Type_ = utils.Some(ITEM_SINGLE_MODEL_NAME)
	i.Name.SetModelName()
	if i.UID.E {
		i.UID.V.SetModelName()
	}
	if i.Links.E {
		i.Links.V.SetModelName()
	}
	if i.ArchetypeDetails.E {
		i.ArchetypeDetails.V.SetModelName()
	}
	if i.FeederAudit.E {
		i.FeederAudit.V.SetModelName()
	}
	i.Item.SetModelName()
}

func (i ITEM_SINGLE) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != ITEM_SINGLE_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ITEM_SINGLE_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + ITEM_SINGLE_MODEL_NAME,
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
		attrPath = path + ".links"
		validateErr.Errs = append(validateErr.Errs, i.Links.V.Validate(attrPath).Errs...)
	}

	// Validate archetype_details
	if i.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, i.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate feeder_audit
	if i.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, i.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate item
	attrPath = path + ".item"
	validateErr.Errs = append(validateErr.Errs, i.Item.Validate(attrPath).Errs...)

	return validateErr
}
