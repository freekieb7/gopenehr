package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const INSTRUCTION_MODEL_NAME string = "INSTRUCTION"

type INSTRUCTION struct {
	Type_               utils.Optional[string]           `json:"_type,omitzero"`
	Name                X_DV_TEXT                        `json:"name"`
	ArchetypeNodeID     string                           `json:"archetype_node_id"`
	UID                 utils.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links               utils.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails    utils.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit         utils.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Language            CODE_PHRASE                      `json:"language"`
	Encoding            CODE_PHRASE                      `json:"encoding"`
	OtherParticipations utils.Optional[[]PARTICIPATION]  `json:"other_participations,omitzero"`
	WorkflowID          utils.Optional[OBJECT_REF]       `json:"workflow_id,omitzero"`
	Subject             X_PARTY_PROXY                    `json:"subject"`
	Provider            utils.Optional[X_PARTY_PROXY]    `json:"provider,omitzero"`
	Protocol            utils.Optional[X_ITEM_STRUCTURE] `json:"protocol,omitzero"`
	GuidelineID         utils.Optional[OBJECT_REF]       `json:"guideline_id,omitzero"`
	Narrative           X_DV_TEXT                        `json:"narrative"`
	ExpiryTime          utils.Optional[DV_DATE_TIME]     `json:"expiry_time,omitzero"`
	WFDefinition        utils.Optional[DV_PARSABLE]      `json:"wf_definition,omitzero"`
	Activities          utils.Optional[[]ACTIVITY]       `json:"activities,omitzero"`
}

func (i *INSTRUCTION) isContentItemModel() {}

func (i *INSTRUCTION) HasModelName() bool {
	return i.Type_.E
}

func (i *INSTRUCTION) SetModelName() {
	i.Type_ = utils.Some(INSTRUCTION_MODEL_NAME)
	i.Name.SetModelName()
	if i.UID.E {
		i.UID.V.SetModelName()
	}
	if i.Links.E {
		for j := range i.Links.V {
			i.Links.V[j].SetModelName()
		}
	}
	if i.ArchetypeDetails.E {
		i.ArchetypeDetails.V.SetModelName()
	}
	if i.FeederAudit.E {
		i.FeederAudit.V.SetModelName()
	}
	i.Language.SetModelName()
	i.Encoding.SetModelName()
	if i.OtherParticipations.E {
		for j := range i.OtherParticipations.V {
			i.OtherParticipations.V[j].SetModelName()
		}
	}
	if i.WorkflowID.E {
		i.WorkflowID.V.SetModelName()
	}
	i.Subject.SetModelName()
	if i.Provider.E {
		i.Provider.V.SetModelName()
	}
	if i.Protocol.E {
		i.Protocol.V.SetModelName()
	}
	if i.GuidelineID.E {
		i.GuidelineID.V.SetModelName()
	}
	i.Narrative.SetModelName()
	if i.ExpiryTime.E {
		i.ExpiryTime.V.SetModelName()
	}
	if i.WFDefinition.E {
		i.WFDefinition.V.SetModelName()
	}
	if i.Activities.E {
		for j := range i.Activities.V {
			i.Activities.V[j].SetModelName()
		}
	}
}

func (i *INSTRUCTION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != INSTRUCTION_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          INSTRUCTION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + INSTRUCTION_MODEL_NAME,
			Recommendation: "Set _type to " + INSTRUCTION_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, i.Name.Validate(attrPath).Errs...)

	// Validate uid
	if i.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, i.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if i.Links.E {
		for j := range i.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, j)
			validateErr.Errs = append(validateErr.Errs, i.Links.V[j].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if i.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, i.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate feeder_audit
	if i.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, i.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate language
	attrPath = path + ".language"
	validateErr.Errs = append(validateErr.Errs, i.Language.Validate(attrPath).Errs...)

	// Validate encoding
	attrPath = path + ".encoding"
	validateErr.Errs = append(validateErr.Errs, i.Encoding.Validate(attrPath).Errs...)
	// Validate other_participations
	if i.OtherParticipations.E {
		for j := range i.OtherParticipations.V {
			attrPath = fmt.Sprintf("%s.other_participations[%d]", path, j)
			validateErr.Errs = append(validateErr.Errs, i.OtherParticipations.V[j].Validate(attrPath).Errs...)
		}
	}

	// Validate workflow_id
	if i.WorkflowID.E {
		attrPath = path + ".workflow_id"
		validateErr.Errs = append(validateErr.Errs, i.WorkflowID.V.Validate(attrPath).Errs...)
	}

	// Validate subject
	attrPath = path + ".subject"
	validateErr.Errs = append(validateErr.Errs, i.Subject.Validate(attrPath).Errs...)

	// Validate provider
	if i.Provider.E {
		attrPath = path + ".provider"
		validateErr.Errs = append(validateErr.Errs, i.Provider.V.Validate(attrPath).Errs...)
	}

	// Validate protocol
	if i.Protocol.E {
		attrPath = path + ".protocol"
		validateErr.Errs = append(validateErr.Errs, i.Protocol.V.Validate(attrPath).Errs...)
	}

	// Validate guideline_id
	if i.GuidelineID.E {
		attrPath = path + ".guideline_id"
		validateErr.Errs = append(validateErr.Errs, i.GuidelineID.V.Validate(attrPath).Errs...)
	}

	// Validate narrative
	attrPath = path + ".narrative"
	validateErr.Errs = append(validateErr.Errs, i.Narrative.Validate(attrPath).Errs...)

	// Validate expiry_time
	if i.ExpiryTime.E {
		attrPath = path + ".expiry_time"
		validateErr.Errs = append(validateErr.Errs, i.ExpiryTime.V.Validate(attrPath).Errs...)
	}

	// Validate wf_definition
	if i.WFDefinition.E {
		attrPath = path + ".wf_definition"
		validateErr.Errs = append(validateErr.Errs, i.WFDefinition.V.Validate(attrPath).Errs...)
	}

	// Validate activities
	if i.Activities.E {
		for j := range i.Activities.V {
			attrPath = fmt.Sprintf("%s.activities[%d]", path, j)
			validateErr.Errs = append(validateErr.Errs, i.Activities.V[j].Validate(attrPath).Errs...)
		}
	}

	return validateErr
}
