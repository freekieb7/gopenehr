package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const PARTY_SELF_MODEL_NAME string = "PARTY_SELF"

type PARTY_SELF struct {
	Type_       utils.Optional[string]    `json:"_type,omitzero"`
	ExternalRef utils.Optional[PARTY_REF] `json:"external_ref,omitzero"`
}

func (p *PARTY_SELF) isPartyProxyModel() {}

func (p *PARTY_SELF) HasModelName() bool {
	return p.Type_.E
}

func (p *PARTY_SELF) SetModelName() {
	p.Type_ = utils.Some(PARTY_SELF_MODEL_NAME)
	if p.ExternalRef.E {
		p.ExternalRef.V.SetModelName()
	}
}

func (p *PARTY_SELF) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if p.Type_.E && p.Type_.V != PARTY_SELF_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          PARTY_SELF_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to PARTY_SELF",
		})
	}

	// Validate external_ref
	if p.ExternalRef.E {
		attrPath = path + ".external_ref"
		validateErr.Errs = append(validateErr.Errs, p.ExternalRef.V.Validate(attrPath).Errs...)
	}

	return validateErr
}
