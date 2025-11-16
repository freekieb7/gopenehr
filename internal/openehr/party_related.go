package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const PARTY_RELATED_MODEL_NAME string = "PARTY_RELATED"

type PARTY_RELATED struct {
	Type_        util.Optional[string]          `json:"_type,omitzero"`
	ExternalRef  util.Optional[PARTY_REF]       `json:"external_ref,omitzero"`
	Name         util.Optional[string]          `json:"name,omitzero"`
	Identifiers  util.Optional[[]DV_IDENTIFIER] `json:"identifiers,omitzero"`
	Relationship DV_CODED_TEXT                  `json:"relationship"`
}

func (p *PARTY_RELATED) isPartyProxyModel() {}

func (p *PARTY_RELATED) HasModelName() bool {
	return p.Type_.E
}

func (p *PARTY_RELATED) SetModelName() {
	p.Type_ = util.Some(PARTY_RELATED_MODEL_NAME)
	if p.ExternalRef.E {
		p.ExternalRef.V.SetModelName()
	}
	if p.Identifiers.E {
		for i := range p.Identifiers.V {
			p.Identifiers.V[i].SetModelName()
		}
	}
}

func (p *PARTY_RELATED) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if p.Type_.E && p.Type_.V != PARTY_RELATED_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          PARTY_RELATED_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to PARTY_RELATED",
		})
	}

	// Validate external_ref
	if p.ExternalRef.E {
		attrPath = path + ".external_ref"
		errors = append(errors, p.ExternalRef.V.Validate(attrPath)...)
	}

	// Validate identifiers
	if p.Identifiers.E {
		for i := range p.Identifiers.V {
			itemPath := fmt.Sprintf("%s.identifiers[%d]", attrPath, i)
			errors = append(errors, p.Identifiers.V[i].Validate(itemPath)...)
		}
	}

	// Validate relationship
	attrPath = path + ".relationship"
	errors = append(errors, p.Relationship.Validate(attrPath)...)

	return errors
}
