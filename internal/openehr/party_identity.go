package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const PARTY_IDENTITY_MODEL_NAME string = "PARTY_IDENTITY"

type PARTY_IDENTITY struct {
	Type_            util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Details          X_ITEM_STRUCTURE              `json:"details"`
}

func (p *PARTY_IDENTITY) SetModelName() {
	p.Type_ = util.Some(PARTY_IDENTITY_MODEL_NAME)
	p.Name.SetModelName()
	if p.UID.E {
		p.UID.V.SetModelName()
	}
	if p.Links.E {
		for i := range p.Links.V {
			p.Links.V[i].SetModelName()
		}
	}
	if p.ArchetypeDetails.E {
		p.ArchetypeDetails.V.SetModelName()
	}
	if p.FeederAudit.E {
		p.FeederAudit.V.SetModelName()
	}
	p.Details.SetModelName()
}

func (p *PARTY_IDENTITY) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if p.Type_.E && p.Type_.V != PARTY_IDENTITY_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          PARTY_IDENTITY_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", PARTY_IDENTITY_MODEL_NAME, p.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", PARTY_IDENTITY_MODEL_NAME),
		})
	}

	// Validate Name
	attrPath = path + ".name"
	errors = append(errors, p.Name.Validate(attrPath)...)

	// Validate ArchetypeNodeID
	attrPath = path + ".archetype_node_id"
	if p.ArchetypeNodeID == "" {
		errors = append(errors, util.ValidationError{
			Model:          PARTY_IDENTITY_MODEL_NAME,
			Path:           attrPath,
			Message:        "archetype_node_id is required",
			Recommendation: "Ensure archetype_node_id is set",
		})
	}

	// Validate UID
	if p.UID.E {
		attrPath = path + ".uid"
		errors = append(errors, p.UID.V.Validate(attrPath)...)
	}

	// Validate Links
	if p.Links.E {
		for i, link := range p.Links.V {
			attrPath = path + fmt.Sprintf(".links[%d]", i)
			errors = append(errors, link.Validate(attrPath)...)
		}
	}

	// Validate ArchetypeDetails
	if p.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errors = append(errors, p.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate FeederAudit
	if p.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errors = append(errors, p.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate Details
	attrPath = path + ".details"
	errors = append(errors, p.Details.Validate(attrPath)...)

	return errors
}
