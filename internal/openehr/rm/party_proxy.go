package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const PARTY_PROXY_TYPE string = "PARTY_PROXY"

type PartyProxyKind int

const (
	PartyProxyKind_Unknown PartyProxyKind = iota
	PartyProxyKind_PARTY_SELF
	PartyProxyKind_PARTY_IDENTIFIED
	PartyProxyKind_PARTY_RELATED
)

type PartyProxyUnion struct {
	Kind  PartyProxyKind
	Value any
}

func (p *PartyProxyUnion) SetModelName() {
	switch p.Kind {
	case PartyProxyKind_PARTY_SELF:
		p.Value.(*PARTY_SELF).SetModelName()
	case PartyProxyKind_PARTY_IDENTIFIED:
		p.Value.(*PARTY_IDENTIFIED).SetModelName()
	case PartyProxyKind_PARTY_RELATED:
		p.Value.(*PARTY_RELATED).SetModelName()
	}
}

func (p *PartyProxyUnion) Validate(path string) util.ValidateError {
	switch p.Kind {
	case PartyProxyKind_PARTY_SELF:
		return p.Value.(*PARTY_SELF).Validate(path)
	case PartyProxyKind_PARTY_IDENTIFIED:
		return p.Value.(*PARTY_IDENTIFIED).Validate(path)
	case PartyProxyKind_PARTY_RELATED:
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
		p.Kind = PartyProxyKind_PARTY_SELF
		p.Value = new(PARTY_SELF)
	case PARTY_IDENTIFIED_TYPE:
		p.Kind = PartyProxyKind_PARTY_IDENTIFIED
		p.Value = new(PARTY_IDENTIFIED)
	case PARTY_RELATED_TYPE:
		p.Kind = PartyProxyKind_PARTY_RELATED
		p.Value = new(PARTY_RELATED)
	default:
		p.Kind = PartyProxyKind_Unknown
		return nil
	}

	return json.Unmarshal(data, p.Value)
}

func (p *PartyProxyUnion) PartySelf() *PARTY_SELF {
	if p.Kind == PartyProxyKind_PARTY_SELF {
		return p.Value.(*PARTY_SELF)
	}
	return nil
}

func (p *PartyProxyUnion) PartyIdentified() *PARTY_IDENTIFIED {
	if p.Kind == PartyProxyKind_PARTY_IDENTIFIED {
		return p.Value.(*PARTY_IDENTIFIED)
	}
	return nil
}

func (p *PartyProxyUnion) PartyRelated() *PARTY_RELATED {
	if p.Kind == PartyProxyKind_PARTY_RELATED {
		return p.Value.(*PARTY_RELATED)
	}
	return nil
}

func PartyProxyFromPartySelf(partySelf PARTY_SELF) PartyProxyUnion {
	partySelf.Type_ = utils.Some(PARTY_SELF_TYPE)
	return PartyProxyUnion{
		Kind:  PartyProxyKind_PARTY_SELF,
		Value: &partySelf,
	}
}

func PartyProxyFromPartyIdentified(partyIdentified PARTY_IDENTIFIED) PartyProxyUnion {
	partyIdentified.Type_ = utils.Some(PARTY_IDENTIFIED_TYPE)
	return PartyProxyUnion{
		Kind:  PartyProxyKind_PARTY_IDENTIFIED,
		Value: &partyIdentified,
	}
}

func PartyProxyFromPartyRelated(partyRelated PARTY_RELATED) PartyProxyUnion {
	partyRelated.Type_ = utils.Some(PARTY_RELATED_TYPE)
	return PartyProxyUnion{
		Kind:  PartyProxyKind_PARTY_RELATED,
		Value: &partyRelated,
	}
}
