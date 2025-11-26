package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const PARTY_IDENTIFIED_MODEL_NAME string = "PARTY_IDENTIFIED"

type PARTY_IDENTIFIED struct {
	Type_       util.Optional[string]          `json:"_type,omitzero"`
	ExternalRef util.Optional[PARTY_REF]       `json:"external_ref,omitzero"`
	Name        util.Optional[string]          `json:"name,omitzero"`
	Identifiers util.Optional[[]DV_IDENTIFIER] `json:"identifiers,omitzero"`
}

func (p *PARTY_IDENTIFIED) isPartyProxyModel() {}

func (p *PARTY_IDENTIFIED) HasModelName() bool {
	return p.Type_.E
}

func (p *PARTY_IDENTIFIED) SetModelName() {
	p.Type_ = util.Some(PARTY_IDENTIFIED_MODEL_NAME)
	if p.ExternalRef.E {
		p.ExternalRef.V.SetModelName()
	}
	if p.Identifiers.E {
		for i := range p.Identifiers.V {
			p.Identifiers.V[i].SetModelName()
		}
	}
}

func (p *PARTY_IDENTIFIED) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if p.Type_.E && p.Type_.V != PARTY_IDENTIFIED_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          PARTY_IDENTIFIED_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to PARTY_IDENTIFIED",
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

	return validateErr
}
