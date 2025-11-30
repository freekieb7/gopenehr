package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ACTION_TYPE string = "ACTION"

type ACTION struct {
	Type_               utils.Optional[string]              `json:"_type,omitzero"`
	Name                DvTextUnion                         `json:"name"`
	ArchetypeNodeID     string                              `json:"archetype_node_id"`
	UID                 utils.Optional[UIDBasedIDUnion]     `json:"uid,omitzero"`
	Links               utils.Optional[[]LINK]              `json:"links,omitzero"`
	ArchetypeDetails    utils.Optional[ARCHETYPED]          `json:"archetype_details,omitzero"`
	FeederAudit         utils.Optional[FEEDER_AUDIT]        `json:"feeder_audit,omitzero"`
	Language            CODE_PHRASE                         `json:"language"`
	Encoding            CODE_PHRASE                         `json:"encoding"`
	OtherParticipations utils.Optional[[]PARTICIPATION]     `json:"other_participations,omitzero"`
	WorkflowID          utils.Optional[OBJECT_REF]          `json:"workflow_id,omitzero"`
	Subject             PartyProxyUnion                     `json:"subject"`
	Provider            utils.Optional[PartyProxyUnion]     `json:"provider,omitzero"`
	Protocol            utils.Optional[ItemStructureUnion]  `json:"protocol,omitzero"`
	GuidelineID         utils.Optional[OBJECT_REF]          `json:"guideline_id,omitzero"`
	Time                DV_DATE_TIME                        `json:"time"`
	IsmTransition       ISM_TRANSITION                      `json:"ism_transition"`
	InstructionDetails  utils.Optional[INSTRUCTION_DETAILS] `json:"instruction_details,omitzero"`
	Description         ItemStructureUnion                  `json:"description"`
}

func (a *ACTION) SetModelName() {
	a.Type_ = utils.Some(ACTION_TYPE)
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
	if a.Protocol.E {
		a.Protocol.V.SetModelName()
	}
	if a.GuidelineID.E {
		a.GuidelineID.V.SetModelName()
	}
	a.IsmTransition.SetModelName()
	if a.InstructionDetails.E {
		a.InstructionDetails.V.SetModelName()
	}
	a.Description.SetModelName()
}

func (a *ACTION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ACTION_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          ACTION_TYPE,
			Path:           attrPath,
			Message:        "_type must be " + ACTION_TYPE,
			Recommendation: "Set _type to " + ACTION_TYPE,
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

	// Validate protocol
	if a.Protocol.E {
		attrPath = path + ".protocol"
		validateErr.Errs = append(validateErr.Errs, a.Protocol.V.Validate(attrPath).Errs...)
	}

	// Validate guideline_id
	if a.GuidelineID.E {
		attrPath = path + ".guideline_id"
		validateErr.Errs = append(validateErr.Errs, a.GuidelineID.V.Validate(attrPath).Errs...)
	}

	// Validate ism_transition
	attrPath = path + ".ism_transition"
	validateErr.Errs = append(validateErr.Errs, a.IsmTransition.Validate(attrPath).Errs...)

	// Validate instruction_details
	if a.InstructionDetails.E {
		attrPath = path + ".instruction_details"
		validateErr.Errs = append(validateErr.Errs, a.InstructionDetails.V.Validate(attrPath).Errs...)
	}

	// Validate description
	attrPath = path + ".description"
	validateErr.Errs = append(validateErr.Errs, a.Description.Validate(attrPath).Errs...)

	return validateErr
}
