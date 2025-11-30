package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ADMIN_ENTRY_TYPE string = "ADMIN_ENTRY"

type ADMIN_ENTRY struct {
	Type_               utils.Optional[string]          `json:"_type,omitzero"`
	Name                DvTextUnion                     `json:"name"`
	ArchetypeNodeID     string                          `json:"archetype_node_id"`
	UID                 utils.Optional[UIDBasedIDUnion] `json:"uid,omitzero"`
	Links               utils.Optional[[]LINK]          `json:"links,omitzero"`
	ArchetypeDetails    utils.Optional[ARCHETYPED]      `json:"archetype_details,omitzero"`
	FeederAudit         utils.Optional[FEEDER_AUDIT]    `json:"feeder_audit,omitzero"`
	Language            CODE_PHRASE                     `json:"language"`
	Encoding            CODE_PHRASE                     `json:"encoding"`
	OtherParticipations utils.Optional[[]PARTICIPATION] `json:"other_participations,omitzero"`
	WorkflowID          utils.Optional[OBJECT_REF]      `json:"workflow_id,omitzero"`
	Subject             PartyProxyUnion                 `json:"subject"`
	Provider            utils.Optional[PartyProxyUnion] `json:"provider,omitzero"`
	Data                ItemStructureUnion              `json:"data"`
}

func (a *ADMIN_ENTRY) SetModelName() {
	a.Type_ = utils.Some(ADMIN_ENTRY_TYPE)
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
	a.Language.SetModelName()
	a.Encoding.SetModelName()
	if a.OtherParticipations.E {
		for i := range a.OtherParticipations.V {
			a.OtherParticipations.V[i].SetModelName()
		}
	}
	if a.WorkflowID.E {
		a.WorkflowID.V.SetModelName()
	}
	a.Subject.SetModelName()
	if a.Provider.E {
		a.Provider.V.SetModelName()
	}
	a.Data.SetModelName()
}

func (a *ADMIN_ENTRY) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ADMIN_ENTRY_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ADMIN_ENTRY_TYPE,
			Path:           attrPath,
			Message:        "_type must be " + ADMIN_ENTRY_TYPE,
			Recommendation: "Set _type to " + ADMIN_ENTRY_TYPE,
		})
	}

	// Validate name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, a.Name.Validate(attrPath).Errs...)

	// Validate uid
	if a.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, a.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if a.Links.E {
		for i := range a.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, a.Links.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if a.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, a.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate feeder_audit
	if a.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, a.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate language
	attrPath = path + ".language"
	validateErr.Errs = append(validateErr.Errs, a.Language.Validate(attrPath).Errs...)

	// Validate encoding
	attrPath = path + ".encoding"
	validateErr.Errs = append(validateErr.Errs, a.Encoding.Validate(attrPath).Errs...)

	// Validate other_participations
	if a.OtherParticipations.E {
		for i := range a.OtherParticipations.V {
			attrPath = fmt.Sprintf("%s.other_participations[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, a.OtherParticipations.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate workflow_id
	if a.WorkflowID.E {
		attrPath = path + ".workflow_id"
		validateErr.Errs = append(validateErr.Errs, a.WorkflowID.V.Validate(attrPath).Errs...)
	}

	// Validate subject
	attrPath = path + ".subject"
	validateErr.Errs = append(validateErr.Errs, a.Subject.Validate(attrPath).Errs...)

	// Validate provider
	if a.Provider.E {
		attrPath = path + ".provider"
		validateErr.Errs = append(validateErr.Errs, a.Provider.V.Validate(attrPath).Errs...)
	}

	// Validate data
	attrPath = path + ".data"
	validateErr.Errs = append(validateErr.Errs, a.Data.Validate(attrPath).Errs...)

	return validateErr
}
