package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const INSTRUCTION_MODEL_NAME string = "INSTRUCTION"

type INSTRUCTION struct {
	Type_               util.Optional[string]           `json:"_type,omitzero"`
	Name                X_DV_TEXT                       `json:"name"`
	ArchetypeNodeID     string                          `json:"archetype_node_id"`
	UID                 util.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links               util.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails    util.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit         util.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Language            CODE_PHRASE                     `json:"language"`
	Encoding            CODE_PHRASE                     `json:"encoding"`
	OtherParticipations util.Optional[[]PARTICIPATION]  `json:"other_participations,omitzero"`
	WorkflowID          util.Optional[OBJECT_REF]       `json:"workflow_id,omitzero"`
	Subject             X_PARTY_PROXY                   `json:"subject"`
	Provider            util.Optional[X_PARTY_PROXY]    `json:"provider,omitzero"`
	Protocol            util.Optional[X_ITEM_STRUCTURE] `json:"protocol,omitzero"`
	GuidelineID         util.Optional[OBJECT_REF]       `json:"guideline_id,omitzero"`
	Narrative           X_DV_TEXT                       `json:"narrative"`
	ExpiryTime          util.Optional[DV_DATE_TIME]     `json:"expiry_time,omitzero"`
	WFDefinition        util.Optional[DV_PARSABLE]      `json:"wf_definition,omitzero"`
	Activities          util.Optional[[]ACTIVITY]       `json:"activities,omitzero"`
}

func (i *INSTRUCTION) isContentItemModel() {}

func (i *INSTRUCTION) HasModelName() bool {
	return i.Type_.E
}

func (i *INSTRUCTION) SetModelName() {
	i.Type_ = util.Some(INSTRUCTION_MODEL_NAME)
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

func (i *INSTRUCTION) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != INSTRUCTION_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          INSTRUCTION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + INSTRUCTION_MODEL_NAME,
			Recommendation: "Set _type to " + INSTRUCTION_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	errs = append(errs, i.Name.Validate(attrPath)...)

	// Validate uid
	if i.UID.E {
		attrPath = path + ".uid"
		errs = append(errs, i.UID.V.Validate(attrPath)...)
	}

	// Validate links
	if i.Links.E {
		for j := range i.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, j)
			errs = append(errs, i.Links.V[j].Validate(attrPath)...)
		}
	}

	// Validate archetype_details
	if i.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errs = append(errs, i.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate feeder_audit
	if i.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errs = append(errs, i.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate language
	attrPath = path + ".language"
	errs = append(errs, i.Language.Validate(attrPath)...)

	// Validate encoding
	attrPath = path + ".encoding"
	errs = append(errs, i.Encoding.Validate(attrPath)...)

	// Validate other_participations
	if i.OtherParticipations.E {
		for j := range i.OtherParticipations.V {
			attrPath = fmt.Sprintf("%s.other_participations[%d]", path, j)
			errs = append(errs, i.OtherParticipations.V[j].Validate(attrPath)...)
		}
	}

	// Validate workflow_id
	if i.WorkflowID.E {
		attrPath = path + ".workflow_id"
		errs = append(errs, i.WorkflowID.V.Validate(attrPath)...)
	}

	// Validate subject
	attrPath = path + ".subject"
	errs = append(errs, i.Subject.Validate(attrPath)...)

	// Validate provider
	if i.Provider.E {
		attrPath = path + ".provider"
		errs = append(errs, i.Provider.V.Validate(attrPath)...)
	}

	// Validate protocol
	if i.Protocol.E {
		attrPath = path + ".protocol"
		errs = append(errs, i.Protocol.V.Validate(attrPath)...)
	}

	// Validate guideline_id
	if i.GuidelineID.E {
		attrPath = path + ".guideline_id"
		errs = append(errs, i.GuidelineID.V.Validate(attrPath)...)
	}

	// Validate narrative
	attrPath = path + ".narrative"
	errs = append(errs, i.Narrative.Validate(attrPath)...)

	// Validate expiry_time
	if i.ExpiryTime.E {
		attrPath = path + ".expiry_time"
		errs = append(errs, i.ExpiryTime.V.Validate(attrPath)...)
	}

	// Validate wf_definition
	if i.WFDefinition.E {
		attrPath = path + ".wf_definition"
		errs = append(errs, i.WFDefinition.V.Validate(attrPath)...)
	}

	// Validate activities
	if i.Activities.E {
		for j := range i.Activities.V {
			attrPath = fmt.Sprintf("%s.activities[%d]", path, j)
			errs = append(errs, i.Activities.V[j].Validate(attrPath)...)
		}
	}

	return errs
}
