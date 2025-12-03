package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const PARTY_PROXY_TYPE string = "PARTY_PROXY"

type PartyProxyKind int

const (
	PARTY_PROXY_kind_unknown PartyProxyKind = iota
	PARTY_PROXY_kind_PARTY_SELF
	PARTY_PROXY_kind_PARTY_IDENTIFIED
	PARTY_PROXY_kind_PARTY_RELATED
)

type PartyProxyUnion struct {
	Kind  PartyProxyKind
	Value any
}

func (p *PartyProxyUnion) SetModelName() {
	switch p.Kind {
	case PARTY_PROXY_kind_PARTY_SELF:
		p.Value.(*PARTY_SELF).SetModelName()
	case PARTY_PROXY_kind_PARTY_IDENTIFIED:
		p.Value.(*PARTY_IDENTIFIED).SetModelName()
	case PARTY_PROXY_kind_PARTY_RELATED:
		p.Value.(*PARTY_RELATED).SetModelName()
	}
}

func (p *PartyProxyUnion) Validate(path string) util.ValidateError {
	switch p.Kind {
	case PARTY_PROXY_kind_PARTY_SELF:
		return p.Value.(*PARTY_SELF).Validate(path)
	case PARTY_PROXY_kind_PARTY_IDENTIFIED:
		return p.Value.(*PARTY_IDENTIFIED).Validate(path)
	case PARTY_PROXY_kind_PARTY_RELATED:
		return p.Value.(*PARTY_RELATED).Validate(path)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          PARTY_PROXY_TYPE,
					Path:           path,
					Message:        "value is not known PARTY_PROXY subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}
}

func (p PartyProxyUnion) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Value)
}

func (p *PartyProxyUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case PARTY_SELF_TYPE:
		p.Kind = PARTY_PROXY_kind_PARTY_SELF
		p.Value = &PARTY_SELF{}
	case PARTY_IDENTIFIED_TYPE:
		p.Kind = PARTY_PROXY_kind_PARTY_IDENTIFIED
		p.Value = &PARTY_IDENTIFIED{}
	case PARTY_RELATED_TYPE:
		p.Kind = PARTY_PROXY_kind_PARTY_RELATED
		p.Value = &PARTY_RELATED{}
	default:
		p.Kind = PARTY_PROXY_kind_unknown
		return nil
	}

	return json.Unmarshal(data, p.Value)
}

func (p *PartyProxyUnion) PARTY_SELF() *PARTY_SELF {
	if p.Kind == PARTY_PROXY_kind_PARTY_SELF {
		return p.Value.(*PARTY_SELF)
	}
	return nil
}

func (p *PartyProxyUnion) PARTY_IDENTIFIED() *PARTY_IDENTIFIED {
	if p.Kind == PARTY_PROXY_kind_PARTY_IDENTIFIED {
		return p.Value.(*PARTY_IDENTIFIED)
	}
	return nil
}

func (p *PartyProxyUnion) PARTY_RELATED() *PARTY_RELATED {
	if p.Kind == PARTY_PROXY_kind_PARTY_RELATED {
		return p.Value.(*PARTY_RELATED)
	}
	return nil
}

func PARTY_PROXY_from_PARTY_SELF(partySelf PARTY_SELF) PartyProxyUnion {
	partySelf.Type_ = utils.Some(PARTY_SELF_TYPE)
	return PartyProxyUnion{
		Kind:  PARTY_PROXY_kind_PARTY_SELF,
		Value: &partySelf,
	}
}

func PARTY_PROXY_from_PARTY_IDENTIFIED(partyIdentified PARTY_IDENTIFIED) PartyProxyUnion {
	partyIdentified.Type_ = utils.Some(PARTY_IDENTIFIED_TYPE)
	return PartyProxyUnion{
		Kind:  PARTY_PROXY_kind_PARTY_IDENTIFIED,
		Value: &partyIdentified,
	}
}

func PARTY_PROXY_from_PARTY_RELATED(partyRelated PARTY_RELATED) PartyProxyUnion {
	partyRelated.Type_ = utils.Some(PARTY_RELATED_TYPE)
	return PartyProxyUnion{
		Kind:  PARTY_PROXY_kind_PARTY_RELATED,
		Value: &partyRelated,
	}
}
