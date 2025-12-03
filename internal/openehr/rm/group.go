package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const GROUP_TYPE string = "GROUP"

type GROUP struct {
	Type_                utils.Optional[string]               `json:"_type,omitzero"`
	Name                 DvTextUnion                          `json:"name"`
	ArchetypeNodeID      string                               `json:"archetype_node_id"`
	UID                  utils.Optional[UIDBasedIDUnion]      `json:"uid,omitzero"`
	Links                utils.Optional[[]LINK]               `json:"links,omitzero"`
	ArchetypeDetails     utils.Optional[ARCHETYPED]           `json:"archetype_details,omitzero"`
	FeederAudit          utils.Optional[FEEDER_AUDIT]         `json:"feeder_audit,omitzero"`
	Identities           []PARTY_IDENTITY                     `json:"identities"`
	Contacts             utils.Optional[[]CONTACT]            `json:"contacts,omitzero"`
	Details              utils.Optional[ItemStructureUnion]   `json:"details,omitzero"`
	ReverseRelationships utils.Optional[[]PARTY_RELATIONSHIP] `json:"reverse_relationships,omitzero"`
	Relationships        utils.Optional[[]PARTY_RELATIONSHIP] `json:"relationships,omitzero"`
	Languages            utils.Optional[[]DvTextUnion]        `json:"languages,omitzero"`
	Roles                utils.Optional[[]PARTY_REF]          `json:"roles,omitzero"`
}

func (a *GROUP) SetModelName() {
	a.Type_ = utils.Some(GROUP_TYPE)
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

func (a *GROUP) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	attrPath = path + "._type"
	if a.Type_.E && a.Type_.V != GROUP_TYPE {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          GROUP_TYPE,
			Path:           attrPath,
			Message:        "invalid " + GROUP_TYPE + " _type field: " + a.Type_.V,
			Recommendation: "Ensure _type field is set to '" + GROUP_TYPE + "'",
		})
	}

	// Validate Name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, a.Name.Validate(attrPath).Errs...)

	// Validate ArchetypeNodeID
	attrPath = path + ".archetype_node_id"
	if a.ArchetypeNodeID == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          AGENT_TYPE,
			Path:           attrPath,
			Message:        "archetype_node_id is required",
			Recommendation: "Ensure archetype_node_id is set",
		})
	}

	// Validate UID
	if a.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, a.UID.V.Validate(attrPath).Errs...)
	}

	// Validate Links
	if a.Links.E {
		for i, link := range a.Links.V {
			attrPath = path + fmt.Sprintf(".links[%d]", i)
			validateErr.Errs = append(validateErr.Errs, link.Validate(attrPath).Errs...)
		}
	}

	// Validate ArchetypeDetails
	if a.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, a.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate FeederAudit
	if a.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, a.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate Identities
	if len(a.Identities) == 0 {
		attrPath = path + ".identities"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          AGENT_TYPE,
			Path:           attrPath,
			Message:        "identities is required",
			Recommendation: "Ensure identities is set",
		})
	} else {
		for i, identity := range a.Identities {
			attrPath = path + fmt.Sprintf(".identities[%d]", i)
			validateErr.Errs = append(validateErr.Errs, identity.Validate(attrPath).Errs...)
		}
	}

	// Validate Contacts
	if a.Contacts.E {
		for i, contact := range a.Contacts.V {
			attrPath = path + fmt.Sprintf(".contacts[%d]", i)
			validateErr.Errs = append(validateErr.Errs, contact.Validate(attrPath).Errs...)
		}
	}

	// Validate Details
	if a.Details.E {
		attrPath = path + ".details"
		validateErr.Errs = append(validateErr.Errs, a.Details.V.Validate(attrPath).Errs...)
	}

	// Validate ReverseRelationships
	if a.ReverseRelationships.E {
		for i, rel := range a.ReverseRelationships.V {
			attrPath = path + fmt.Sprintf(".reverse_relationships[%d]", i)
			validateErr.Errs = append(validateErr.Errs, rel.Validate(attrPath).Errs...)
		}
	}

	// Validate Relationships
	if a.Relationships.E {
		for i, rel := range a.Relationships.V {
			attrPath = path + fmt.Sprintf(".relationships[%d]", i)
			validateErr.Errs = append(validateErr.Errs, rel.Validate(attrPath).Errs...)
		}
	}

	// Validate Languages
	if a.Languages.E {
		for i, lang := range a.Languages.V {
			attrPath = path + fmt.Sprintf(".languages[%d]", i)
			validateErr.Errs = append(validateErr.Errs, lang.Validate(attrPath).Errs...)
		}
	}

	// Validate Roles
	if a.Roles.E {
		for i := range a.Roles.V {
			attrPath = path + fmt.Sprintf(".roles[%d]", i)
			validateErr.Errs = append(validateErr.Errs, a.Roles.V[i].Validate(attrPath).Errs...)
		}
	}

	return validateErr
}
