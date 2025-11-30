package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ROLE_TYPE string = "ROLE"

type ROLE struct {
	Type_            utils.Optional[string]             `json:"_type,omitempty"`
	Name             DvTextUnion                        `json:"name"`
	ArchetypeNodeID  string                             `json:"archetype_node_id"`
	UID              utils.Optional[UIDBasedIDUnion]    `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]             `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]         `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]       `json:"feeder_audit,omitzero"`
	Details          utils.Optional[ItemStructureUnion] `json:"details,omitzero"`
	Target           PARTY_REF                          `json:"target"`
	TimeValidity     utils.Optional[DV_INTERVAL]        `json:"time_validity,omitzero"`
	Source           PARTY_REF                          `json:"source"`
}

func (a *ROLE) SetModelName() {
	a.Type_ = utils.Some(ROLE_TYPE)
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
	if a.Details.E {
		a.Details.V.SetModelName()
	}
	a.Target.SetModelName()
	if a.TimeValidity.E {
		a.TimeValidity.V.SetModelName()
	}
	a.Source.SetModelName()
}

func (a *ROLE) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	attrPath = path + "._type"
	if a.Type_.E && a.Type_.V != ROLE_TYPE {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ROLE_TYPE,
			Path:           attrPath,
			Message:        "invalid " + ROLE_TYPE + " _type field: " + a.Type_.V,
			Recommendation: "Ensure _type field is set to '" + ROLE_TYPE + "'",
		})
	}

	// Validate Name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, a.Name.Validate(attrPath).Errs...)

	// Validate UID
	attrPath = path + ".uid"
	if a.UID.E {
		validateErr.Errs = append(validateErr.Errs, a.UID.V.Validate(attrPath).Errs...)
	}

	// Validate Links
	attrPath = path + ".links"
	if a.Links.E {
		for i := range a.Links.V {
			linkPath := fmt.Sprintf("%s[%d]", attrPath, i)
			validateErr.Errs = append(validateErr.Errs, a.Links.V[i].Validate(linkPath).Errs...)
		}
	}

	// Validate ArchetypeDetails
	attrPath = path + ".archetype_details"
	if a.ArchetypeDetails.E {
		validateErr.Errs = append(validateErr.Errs, a.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate FeederAudit
	attrPath = path + ".feeder_audit"
	if a.FeederAudit.E {
		validateErr.Errs = append(validateErr.Errs, a.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate Details
	attrPath = path + ".details"
	if a.Details.E {
		validateErr.Errs = append(validateErr.Errs, a.Details.V.Validate(attrPath).Errs...)
	}

	// Validate Target
	attrPath = path + ".target"
	validateErr.Errs = append(validateErr.Errs, a.Target.Validate(attrPath).Errs...)

	// Validate TimeValidity
	attrPath = path + ".time_validity"
	if a.TimeValidity.E {
		validateErr.Errs = append(validateErr.Errs, a.TimeValidity.V.Validate(attrPath).Errs...)
	}

	// Validate Source
	attrPath = path + ".source"
	validateErr.Errs = append(validateErr.Errs, a.Source.Validate(attrPath).Errs...)

	return validateErr
}
