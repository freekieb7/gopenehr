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
	Validate(path string) []util.ValidationError
}

type X_PARTY_PROXY struct {
	Value PartyProxyModel
}

func (x *X_PARTY_PROXY) SetModelName() {
	x.Value.SetModelName()
}

func (x X_PARTY_PROXY) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Abstract model requires _type to be defined
	if !x.Value.HasModelName() {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:   PARTY_PROXY_MODEL_NAME,
			Path:    attrPath,
			Message: "missing _type field for abstract model",
		})
	}

	return errs
}

func (x X_PARTY_PROXY) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Value)
}

func (x *X_PARTY_PROXY) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.MetaType
	switch t {
	case PARTY_SELF_MODEL_NAME:
		x.Value = new(PARTY_SELF)
	case PARTY_IDENTIFIED_MODEL_NAME:
		x.Value = new(PARTY_IDENTIFIED)
	case PARTY_RELATED_MODEL_NAME:
		x.Value = new(PARTY_RELATED)
	case "":
		return fmt.Errorf("missing PARTY_PROXY _type field")
	default:
		return fmt.Errorf("PARTY_PROXY unexpected _type %s", t)
	}

	return json.Unmarshal(data, x.Value)
}
