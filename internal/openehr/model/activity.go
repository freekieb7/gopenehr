package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ACTIVITY_MODEL_NAME string = "ACTIVITY"

type ACTIVITY struct {
	Type_             utils.Optional[string]         `json:"_type,omitzero"`
	Name              X_DV_TEXT                      `json:"name"`
	ArchetypeNodeID   string                         `json:"archetype_node_id"`
	UID               utils.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links             utils.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails  utils.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit       utils.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Timing            utils.Optional[DV_PARSABLE]    `json:"timing,omitzero"`
	ActionArchetypeID string                         `json:"action_archetype_id"`
	Description       X_ITEM_STRUCTURE               `json:"description"`
}

func (a *ACTIVITY) isContentItemModel() {}

func (a *ACTIVITY) HasModelName() bool {
	return a.Type_.E
}

func (a *ACTIVITY) SetModelName() {
	a.Type_ = utils.Some(ACTIVITY_MODEL_NAME)
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
	if a.Timing.E {
		a.Timing.V.SetModelName()
	}
	a.Description.SetModelName()
}

func (a *ACTIVITY) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ACTIVITY_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ACTIVITY_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + ACTIVITY_MODEL_NAME,
			Recommendation: "Set _type to " + ACTIVITY_MODEL_NAME,
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
		for i := range a.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, a.Links.V[i].Validate(attrPath).Errs...)
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

	// Validate Timing
	if a.Timing.E {
		attrPath = path + ".timing"
		validateErr.Errs = append(validateErr.Errs, a.Timing.V.Validate(attrPath).Errs...)
	}

	// Validate Description
	attrPath = path + ".description"
	validateErr.Errs = append(validateErr.Errs, a.Description.Validate(attrPath).Errs...)

	return validateErr
}
