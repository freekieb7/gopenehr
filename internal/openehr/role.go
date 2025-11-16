package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ROLE_MODEL_NAME string = "ROLE"

type ROLE struct {
	Type_            util.Optional[string]           `json:"_type,omitempty"`
	Name             X_DV_TEXT                       `json:"name"`
	ArchetypeNodeID  string                          `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Details          util.Optional[X_ITEM_STRUCTURE] `json:"details,omitzero"`
	Target           PARTY_REF                       `json:"target"`
	TimeValidity     util.Optional[DV_INTERVAL]      `json:"time_validity,omitzero"`
	Source           PARTY_REF                       `json:"source"`
}

func (a *ROLE) isVersionModel() {}

func (a *ROLE) SetModelName() {
	a.Type_ = util.Some(ROLE_MODEL_NAME)
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
	if a.Details.E {
		a.Details.V.SetModelName()
	}
	a.Target.SetModelName()
	if a.TimeValidity.E {
		a.TimeValidity.V.SetModelName()
	}
	a.Source.SetModelName()
}

func (a *ROLE) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	attrPath = path + "._type"
	if a.Type_.E && a.Type_.V != ROLE_MODEL_NAME {
		errors = append(errors, util.ValidationError{
			Model:          ROLE_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid " + ROLE_MODEL_NAME + " _type field: " + a.Type_.V,
			Recommendation: "Ensure _type field is set to '" + ROLE_MODEL_NAME + "'",
		})
	}

	// Validate Name
	attrPath = path + ".name"
	errors = append(errors, a.Name.Validate(attrPath)...)

	// Validate UID
	attrPath = path + ".uid"
	if a.UID.E {
		errors = append(errors, a.UID.V.Validate(attrPath)...)
	}

	// Validate Links
	attrPath = path + ".links"
	if a.Links.E {
		for i := range a.Links.V {
			linkPath := fmt.Sprintf("%s[%d]", attrPath, i)
			errors = append(errors, a.Links.V[i].Validate(linkPath)...)
		}
	}

	// Validate ArchetypeDetails
	attrPath = path + ".archetype_details"
	if a.ArchetypeDetails.E {
		errors = append(errors, a.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate FeederAudit
	attrPath = path + ".feeder_audit"
	if a.FeederAudit.E {
		errors = append(errors, a.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate Details
	attrPath = path + ".details"
	if a.Details.E {
		errors = append(errors, a.Details.V.Validate(attrPath)...)
	}

	// Validate Target
	attrPath = path + ".target"
	errors = append(errors, a.Target.Validate(attrPath)...)

	// Validate TimeValidity
	attrPath = path + ".time_validity"
	if a.TimeValidity.E {
		errors = append(errors, a.TimeValidity.V.Validate(attrPath)...)
	}

	// Validate Source
	attrPath = path + ".source"
	errors = append(errors, a.Source.Validate(attrPath)...)

	return errors
}
