package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const PARTY_IDENTITY_MODEL_NAME string = "PARTY_IDENTITY"

type PARTY_IDENTITY struct {
	Type_            utils.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                      `json:"name"`
	ArchetypeNodeID  string                         `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Details          X_ITEM_STRUCTURE               `json:"details"`
}

func (p *PARTY_IDENTITY) SetModelName() {
	p.Type_ = utils.Some(PARTY_IDENTITY_MODEL_NAME)
	p.Name.SetModelName()
	if p.UID.E {
		p.UID.V.SetModelName()
	}
	if p.Links.E {
		for i := range p.Links.V {
			p.Links.V[i].SetModelName()
		}
	}
	if p.ArchetypeDetails.E {
		p.ArchetypeDetails.V.SetModelName()
	}
	if p.FeederAudit.E {
		p.FeederAudit.V.SetModelName()
	}
	p.Details.SetModelName()
}

func (p *PARTY_IDENTITY) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if p.Type_.E && p.Type_.V != PARTY_IDENTITY_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          PARTY_IDENTITY_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", PARTY_IDENTITY_MODEL_NAME, p.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", PARTY_IDENTITY_MODEL_NAME),
		})
	}

	// Validate Name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, p.Name.Validate(attrPath).Errs...)

	// Validate ArchetypeNodeID
	attrPath = path + ".archetype_node_id"
	if p.ArchetypeNodeID == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          PARTY_IDENTITY_MODEL_NAME,
			Path:           attrPath,
			Message:        "archetype_node_id is required",
			Recommendation: "Ensure archetype_node_id is set",
		})
	}

	// Validate UID
	if p.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, p.UID.V.Validate(attrPath).Errs...)
	}

	// Validate Links
	if p.Links.E {
		for i, link := range p.Links.V {
			attrPath = path + fmt.Sprintf(".links[%d]", i)
			validateErr.Errs = append(validateErr.Errs, link.Validate(attrPath).Errs...)
		}
	}

	// Validate ArchetypeDetails
	if p.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, p.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate FeederAudit
	if p.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, p.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate Details
	attrPath = path + ".details"
	validateErr.Errs = append(validateErr.Errs, p.Details.Validate(attrPath).Errs...)

	return validateErr
}
