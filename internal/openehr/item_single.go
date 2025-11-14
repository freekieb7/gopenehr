package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const ITEM_SINGLE_MODEL_NAME string = "ITEM_SINGLE"

type ITEM_SINGLE struct {
	Type_            util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[LINK]           `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Item             ELEMENT                       `json:"item"`
}

func (i ITEM_SINGLE) isItemStructureModel() {}

func (i ITEM_SINGLE) HasModelName() bool {
	return i.Type_.E
}

func (i *ITEM_SINGLE) SetModelName() {
	i.Type_ = util.Some(ITEM_SINGLE_MODEL_NAME)
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

func (i ITEM_SINGLE) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != ITEM_SINGLE_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          ITEM_SINGLE_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + ITEM_SINGLE_MODEL_NAME,
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
		attrPath = path + ".links"
		errors = append(errors, i.Links.V.Validate(attrPath)...)
	}

	// Validate archetype_details
	if i.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errors = append(errors, i.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate feeder_audit
	if i.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errors = append(errors, i.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate item
	attrPath = path + ".item"
	errors = append(errors, i.Item.Validate(attrPath)...)

	return errors
}
