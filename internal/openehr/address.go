package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ADDRESS_MODEL_NAME string = "ADDRESS"

type ADDRESS struct {
	Type_            util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Details          X_ITEM_STRUCTURE              `json:"details"`
}

func (a *ADDRESS) SetModelName() {
	a.Type_ = util.Some(ADDRESS_MODEL_NAME)
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
	a.Details.SetModelName()
}

func (a *ADDRESS) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ADDRESS_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          ADDRESS_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid " + ADDRESS_MODEL_NAME + " _type field: " + a.Type_.V,
			Recommendation: "Ensure _type field is set to '" + ADDRESS_MODEL_NAME + "'",
		})
	}

	// Validate Name
	attrPath = path + ".name"
	errors = append(errors, a.Name.Validate(attrPath)...)

	// Validate UID
	if a.UID.E {
		attrPath = path + ".uid"
		errors = append(errors, a.UID.V.Validate(attrPath)...)
	}

	// Validate Links
	if a.Links.E {
		for i, link := range a.Links.V {
			attrPath = path + fmt.Sprintf(".links[%d]", i)
			errors = append(errors, link.Validate(attrPath)...)
		}
	}

	// Validate ArchetypeDetails
	if a.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errors = append(errors, a.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate FeederAudit
	if a.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errors = append(errors, a.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate Details
	attrPath = path + ".details"
	errors = append(errors, a.Details.Validate(attrPath)...)

	return errors
}
