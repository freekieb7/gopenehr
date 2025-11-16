package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const AGENT_MODEL_NAME string = "AGENT"

type AGENT struct {
	Type_                util.Optional[string]               `json:"_type,omitzero"`
	Name                 X_DV_TEXT                           `json:"name"`
	ArchetypeNodeID      string                              `json:"archetype_node_id"`
	UID                  util.Optional[X_UID_BASED_ID]       `json:"uid,omitzero"`
	Links                util.Optional[[]LINK]               `json:"links,omitzero"`
	ArchetypeDetails     util.Optional[ARCHETYPED]           `json:"archetype_details,omitzero"`
	FeederAudit          util.Optional[FEEDER_AUDIT]         `json:"feeder_audit,omitzero"`
	Identities           []PARTY_IDENTITY                    `json:"identities"`
	Contacts             util.Optional[[]CONTACT]            `json:"contacts,omitzero"`
	Details              util.Optional[X_ITEM_STRUCTURE]     `json:"details,omitzero"`
	ReverseRelationships util.Optional[[]PARTY_RELATIONSHIP] `json:"reverse_relationships,omitzero"`
	Relationships        util.Optional[[]PARTY_RELATIONSHIP] `json:"relationships,omitzero"`
	Languages            util.Optional[[]X_DV_TEXT]          `json:"languages,omitzero"`
	Roles                util.Optional[[]PARTY_REF]          `json:"roles,omitzero"`
}

func (a *AGENT) isVersionModel() {}

func (a *AGENT) SetModelName() {
	a.Type_ = util.Some(AGENT_MODEL_NAME)
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
	for i := range a.Identities {
		a.Identities[i].SetModelName()
	}
	if a.Contacts.E {
		for i := range a.Contacts.V {
			a.Contacts.V[i].SetModelName()
		}
	}
	if a.Details.E {
		a.Details.V.SetModelName()
	}
	if a.ReverseRelationships.E {
		for i := range a.ReverseRelationships.V {
			a.ReverseRelationships.V[i].SetModelName()
		}
	}
	if a.Relationships.E {
		for i := range a.Relationships.V {
			a.Relationships.V[i].SetModelName()
		}
	}
	if a.Languages.E {
		for i := range a.Languages.V {
			a.Languages.V[i].SetModelName()
		}
	}
	if a.Roles.E {
		for i := range a.Roles.V {
			a.Roles.V[i].SetModelName()
		}
	}
}

func (a *AGENT) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	attrPath = path + "._type"
	if a.Type_.E && a.Type_.V != AGENT_MODEL_NAME {
		errors = append(errors, util.ValidationError{
			Model:          AGENT_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid " + AGENT_MODEL_NAME + " _type field: " + a.Type_.V,
			Recommendation: "Ensure _type field is set to '" + AGENT_MODEL_NAME + "'",
		})
	}

	// Validate Name
	attrPath = path + ".name"
	errors = append(errors, a.Name.Validate(attrPath)...)

	// Validate ArchetypeNodeID
	attrPath = path + ".archetype_node_id"
	if a.ArchetypeNodeID == "" {
		errors = append(errors, util.ValidationError{
			Model:          AGENT_MODEL_NAME,
			Path:           attrPath,
			Message:        "archetype_node_id is required",
			Recommendation: "Ensure archetype_node_id is set",
		})
	}

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

	// Validate Identities
	if len(a.Identities) == 0 {
		attrPath = path + ".identities"
		errors = append(errors, util.ValidationError{
			Model:          AGENT_MODEL_NAME,
			Path:           attrPath,
			Message:        "identities is required",
			Recommendation: "Ensure identities is set",
		})
	} else {
		for i, identity := range a.Identities {
			attrPath = path + fmt.Sprintf(".identities[%d]", i)
			errors = append(errors, identity.Validate(attrPath)...)
		}
	}

	// Validate Contacts
	if a.Contacts.E {
		for i, contact := range a.Contacts.V {
			attrPath = path + fmt.Sprintf(".contacts[%d]", i)
			errors = append(errors, contact.Validate(attrPath)...)
		}
	}

	// Validate Details
	if a.Details.E {
		attrPath = path + ".details"
		errors = append(errors, a.Details.V.Validate(attrPath)...)
	}

	// Validate ReverseRelationships
	if a.ReverseRelationships.E {
		for i, rel := range a.ReverseRelationships.V {
			attrPath = path + fmt.Sprintf(".reverse_relationships[%d]", i)
			errors = append(errors, rel.Validate(attrPath)...)
		}
	}

	// Validate Relationships
	if a.Relationships.E {
		for i, rel := range a.Relationships.V {
			attrPath = path + fmt.Sprintf(".relationships[%d]", i)
			errors = append(errors, rel.Validate(attrPath)...)
		}
	}

	// Validate Languages
	if a.Languages.E {
		for i, lang := range a.Languages.V {
			attrPath = path + fmt.Sprintf(".languages[%d]", i)
			errors = append(errors, lang.Validate(attrPath)...)
		}
	}

	// Validate Roles
	if a.Roles.E {
		for i := range a.Roles.V {
			attrPath = path + fmt.Sprintf(".roles[%d]", i)
			errors = append(errors, a.Roles.V[i].Validate(attrPath)...)
		}
	}

	return errors
}
