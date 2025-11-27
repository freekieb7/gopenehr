package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ADDRESS_MODEL_NAME string = "ADDRESS"

type ADDRESS struct {
	Type_            utils.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                      `json:"name"`
	ArchetypeNodeID  string                         `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Details          X_ITEM_STRUCTURE               `json:"details"`
}

func (a *ADDRESS) SetModelName() {
	a.Type_ = utils.Some(ADDRESS_MODEL_NAME)
	a.Name.SetModelName()
	if a.UID.E {
		a.UID.V.SetModelName()
	}
	if a.Links.E {
		for i := range a.Links.V {
			a.Links.V[i].SetModelName()
		}
	}
	if a.ArchetypeDetails.E {
		a.ArchetypeDetails.V.SetModelName()
	}
	if a.FeederAudit.E {
		a.FeederAudit.V.SetModelName()
	}
	a.Details.SetModelName()
}

func (a *ADDRESS) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ADDRESS_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ADDRESS_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid " + ADDRESS_MODEL_NAME + " _type field: " + a.Type_.V,
			Recommendation: "Ensure _type field is set to '" + ADDRESS_MODEL_NAME + "'",
		})
	}

	// Validate Name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, a.Name.Validate(attrPath).Errs...)

	// Validate UID
	if a.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, a.UID.V.Validate(attrPath).Errs...)
	}

	// Validate Links
	if a.Links.E {
		for i, link := range a.Links.V {
			attrPath = path + fmt.Sprintf(".links[%d]", i)
			validateErr.Errs = append(validateErr.Errs, link.Validate(attrPath).Errs...)
		}
	}

	// Validate ArchetypeDetails
	if a.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, a.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate FeederAudit
	if a.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, a.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate Details
	attrPath = path + ".details"
	validateErr.Errs = append(validateErr.Errs, a.Details.Validate(attrPath).Errs...)

	return validateErr
}
