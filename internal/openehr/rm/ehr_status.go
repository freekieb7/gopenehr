package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const EHR_STATUS_MODEL_NAME string = "EHR_STATUS"

type EHR_STATUS struct {
	Type_            utils.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                        `json:"name"`
	ArchetypeNodeID  string                           `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Subject          PARTY_SELF                       `json:"subject"`
	IsQueryable      bool                             `json:"is_queryable"`
	IsModifiable     bool                             `json:"is_modifiable"`
	OtherDetails     utils.Optional[X_ITEM_STRUCTURE] `json:"other_details,omitzero"`
}

func (e *EHR_STATUS) isVersionModel() {}

func (e *EHR_STATUS) SetModelName() {
	e.Type_ = utils.Some(EHR_STATUS_MODEL_NAME)
	e.Name.SetModelName()
	if e.UID.E {
		e.UID.V.SetModelName()
	}
	if e.Links.E {
		for i := range e.Links.V {
			e.Links.V[i].SetModelName()
		}
	}
	if e.ArchetypeDetails.E {
		e.ArchetypeDetails.V.SetModelName()
	}
	e.Subject.SetModelName()
	if e.FeederAudit.E {
		e.FeederAudit.V.SetModelName()
	}
	if e.OtherDetails.E {
		e.OtherDetails.V.SetModelName()
	}
}

func (e *EHR_STATUS) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if e.Type_.E && e.Type_.V != "EHR_STATUS" {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          "EHR_STATUS",
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to EHR_STATUS",
		})
	}

	// Validate name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, e.Name.Validate(attrPath).Errs...)

	// Validate uid
	if e.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, e.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if e.Links.E {
		for i := range e.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, e.Links.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if e.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, e.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate subject
	attrPath = path + ".subject"
	validateErr.Errs = append(validateErr.Errs, e.Subject.Validate(attrPath).Errs...)

	// Validate feeder_audit
	if e.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, e.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate other_details
	if e.OtherDetails.E {
		attrPath = path + ".other_details"
		validateErr.Errs = append(validateErr.Errs, e.OtherDetails.V.Validate(attrPath).Errs...)
	}

	return validateErr
}

func (e *EHR_STATUS) ObjectVersionID() OBJECT_VERSION_ID {
	if e.UID.E {
		if objVerID, ok := e.UID.V.Value.(*OBJECT_VERSION_ID); ok {
			return *objVerID
		}
	}
	return OBJECT_VERSION_ID{}
}
