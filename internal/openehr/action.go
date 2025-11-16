package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const ACTION_MODEL_NAME string = "ACTION"

type ACTION struct {
	Type_               util.Optional[string]              `json:"_type,omitzero"`
	Name                X_DV_TEXT                          `json:"name"`
	ArchetypeNodeID     string                             `json:"archetype_node_id"`
	UID                 util.Optional[X_UID_BASED_ID]      `json:"uid,omitzero"`
	Links               util.Optional[[]LINK]              `json:"links,omitzero"`
	ArchetypeDetails    util.Optional[ARCHETYPED]          `json:"archetype_details,omitzero"`
	FeederAudit         util.Optional[FEEDER_AUDIT]        `json:"feeder_audit,omitzero"`
	Language            CODE_PHRASE                        `json:"language"`
	Encoding            CODE_PHRASE                        `json:"encoding"`
	OtherParticipations util.Optional[[]PARTICIPATION]     `json:"other_participations,omitzero"`
	WorkflowID          util.Optional[OBJECT_REF]          `json:"workflow_id,omitzero"`
	Subject             X_PARTY_PROXY                      `json:"subject"`
	Provider            util.Optional[X_PARTY_PROXY]       `json:"provider,omitzero"`
	Protocol            util.Optional[X_ITEM_STRUCTURE]    `json:"protocol,omitzero"`
	GuidelineID         util.Optional[OBJECT_REF]          `json:"guideline_id,omitzero"`
	Time                DV_DATE_TIME                       `json:"time"`
	IsmTransition       ISM_TRANSITION                     `json:"ism_transition"`
	InstructionDetails  util.Optional[INSTRUCTION_DETAILS] `json:"instruction_details,omitzero"`
	Description         X_ITEM_STRUCTURE                   `json:"description"`
}

func (a *ACTION) isContentItemModel() {}

func (a *ACTION) HasModelName() bool {
	return a.Type_.E
}

func (a *ACTION) SetModelName() {
	a.Type_ = util.Some(ACTION_MODEL_NAME)
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

func (a *ACTION) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if a.Type_.E && a.Type_.V != ACTION_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          ACTION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + ACTION_MODEL_NAME,
			Recommendation: "Set _type to " + ACTION_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	errs = append(errs, a.Name.Validate(attrPath)...)

	// Validate uid
	if a.UID.E {
		attrPath = path + ".uid"
		errs = append(errs, a.UID.V.Validate(attrPath)...)
	}

	// Validate links
	if a.Links.E {
		for i := range a.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			errs = append(errs, a.Links.V[i].Validate(attrPath)...)
		}
	}

	// Validate archetype_details
	if a.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errs = append(errs, a.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate feeder_audit
	if a.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errs = append(errs, a.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate other_participations
	if a.OtherParticipations.E {
		for i := range a.OtherParticipations.V {
			attrPath = fmt.Sprintf("%s.other_participations[%d]", path, i)
			errs = append(errs, a.OtherParticipations.V[i].Validate(attrPath)...)
		}
	}

	// Validate workflow_id
	if a.WorkflowID.E {
		attrPath = path + ".workflow_id"
		errs = append(errs, a.WorkflowID.V.Validate(attrPath)...)
	}

	// Validate subject
	attrPath = path + ".subject"
	errs = append(errs, a.Subject.Validate(attrPath)...)

	// Validate provider
	if a.Provider.E {
		attrPath = path + ".provider"
		errs = append(errs, a.Provider.V.Validate(attrPath)...)
	}

	// Validate protocol
	if a.Protocol.E {
		attrPath = path + ".protocol"
		errs = append(errs, a.Protocol.V.Validate(attrPath)...)
	}

	// Validate guideline_id
	if a.GuidelineID.E {
		attrPath = path + ".guideline_id"
		errs = append(errs, a.GuidelineID.V.Validate(attrPath)...)
	}

	// Validate ism_transition
	attrPath = path + ".ism_transition"
	errs = append(errs, a.IsmTransition.Validate(attrPath)...)

	// Validate instruction_details
	if a.InstructionDetails.E {
		attrPath = path + ".instruction_details"
		errs = append(errs, a.InstructionDetails.V.Validate(attrPath)...)
	}

	// Validate description
	attrPath = path + ".description"
	errs = append(errs, a.Description.Validate(attrPath)...)

	return errs
}
