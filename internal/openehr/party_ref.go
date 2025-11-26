package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const PARTY_REF_MODEL_NAME string = "PARTY_REF"

type PARTY_REF struct {
	Type_     util.Optional[string] `json:"_type,omitzero"`
	Namespace string                `json:"namespace"`
	Type      string                `json:"type"`
	ID        X_OBJECT_ID           `json:"id"`
}

func (p *PARTY_REF) SetModelName() {
	p.Type_ = util.Some(PARTY_REF_MODEL_NAME)
	p.ID.SetModelName()
}

func (p *PARTY_REF) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if p.Type_.E && p.Type_.V != PARTY_REF_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          PARTY_REF_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to PARTY_REF",
		})
	}

	// Validate namespace
	if p.Namespace == "" {
		attrPath = path + ".namespace"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          PARTY_REF_MODEL_NAME,
			Path:           attrPath,
			Message:        "namespace cannot be empty",
			Recommendation: "Provide a valid namespace",
		})
	}

	// Validate type
	if p.Type == "" {
		attrPath = path + ".type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          PARTY_REF_MODEL_NAME,
			Path:           attrPath,
			Message:        "type cannot be empty",
			Recommendation: "Provide a valid type",
		})
	}

	// Validate id
	attrPath = path + ".id"
	validateErr.Errs = append(validateErr.Errs, p.ID.Validate(attrPath).Errs...)

	return validateErr
}
