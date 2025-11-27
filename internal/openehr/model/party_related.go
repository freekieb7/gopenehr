package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const PARTY_RELATED_MODEL_NAME string = "PARTY_RELATED"

type PARTY_RELATED struct {
	Type_        utils.Optional[string]          `json:"_type,omitzero"`
	ExternalRef  utils.Optional[PARTY_REF]       `json:"external_ref,omitzero"`
	Name         utils.Optional[string]          `json:"name,omitzero"`
	Identifiers  utils.Optional[[]DV_IDENTIFIER] `json:"identifiers,omitzero"`
	Relationship DV_CODED_TEXT                   `json:"relationship"`
}

func (p *PARTY_RELATED) isPartyProxyModel() {}

func (p *PARTY_RELATED) HasModelName() bool {
	return p.Type_.E
}

func (p *PARTY_RELATED) SetModelName() {
	p.Type_ = utils.Some(PARTY_RELATED_MODEL_NAME)
	if p.ExternalRef.E {
		p.ExternalRef.V.SetModelName()
	}
	if p.Identifiers.E {
		for i := range p.Identifiers.V {
			p.Identifiers.V[i].SetModelName()
		}
	}
}

func (p *PARTY_RELATED) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if p.Type_.E && p.Type_.V != PARTY_RELATED_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          PARTY_RELATED_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to PARTY_RELATED",
		})
	}

	// Validate external_ref
	if p.ExternalRef.E {
		attrPath = path + ".external_ref"
		validateErr.Errs = append(validateErr.Errs, p.ExternalRef.V.Validate(attrPath).Errs...)
	}

	// Validate identifiers
	if p.Identifiers.E {
		for i := range p.Identifiers.V {
			itemPath := fmt.Sprintf("%s.identifiers[%d]", attrPath, i)
			validateErr.Errs = append(validateErr.Errs, p.Identifiers.V[i].Validate(itemPath).Errs...)
		}
	}

	// Validate relationship
	attrPath = path + ".relationship"
	validateErr.Errs = append(validateErr.Errs, p.Relationship.Validate(attrPath).Errs...)

	return validateErr
}
