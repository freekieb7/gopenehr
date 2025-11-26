package openehr

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const PARTY_PROXY_MODEL_NAME string = "PARTY_PROXY"

// Abstract
type PARTY_PROXY struct {
	Type_       util.Optional[string]    `json:"_type,omitzero"`
	ExternalRef util.Optional[PARTY_REF] `json:"external_ref,omitzero"`
}

func (p PARTY_PROXY) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot marshal abstract PARTY_PROXY type")
}

func (p *PARTY_PROXY) UnmarshalJSON(data []byte) error {
	return fmt.Errorf("cannot unmarshal abstract PARTY_PROXY type")
}

// ========== Union of PARTY_PROXY ==========

type PartyProxyModel interface {
	isPartyProxyModel()
	HasModelName() bool
	SetModelName()
	Validate(path string) util.ValidateError
}

type X_PARTY_PROXY struct {
	Value PartyProxyModel
}

func (x *X_PARTY_PROXY) SetModelName() {
	x.Value.SetModelName()
}

func (x *X_PARTY_PROXY) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Abstract model requires _type to be defined
	if !x.Value.HasModelName() {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:   PARTY_PROXY_MODEL_NAME,
			Path:    attrPath,
			Message: "missing _type field for abstract model",
		})
	}

	return validateErr
}

func (x X_PARTY_PROXY) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Value)
}

func (x *X_PARTY_PROXY) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.Type_
	switch t {
	case PARTY_SELF_MODEL_NAME:
		x.Value = new(PARTY_SELF)
	case PARTY_IDENTIFIED_MODEL_NAME:
		x.Value = new(PARTY_IDENTIFIED)
	case PARTY_RELATED_MODEL_NAME:
		x.Value = new(PARTY_RELATED)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          PARTY_PROXY_MODEL_NAME,
					Path:           "$.**._type",
					Message:        fmt.Sprintf("unexpected PARTY_PROXY _type %s", t),
					Recommendation: "Ensure _type field is one of the known PARTY_PROXY subtypes",
				},
			},
		}
	}

	return json.Unmarshal(data, x.Value)
}
