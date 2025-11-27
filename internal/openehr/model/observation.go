package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const OBSERVATION_MODEL_NAME string = "OBSERVATION"

type OBSERVATION struct {
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
	Data                HISTORY                          `json:"data"`
	State               utils.Optional[HISTORY]          `json:"state,omitzero"`
}

func (o *OBSERVATION) isContentItemModel() {}

func (o *OBSERVATION) HasModelName() bool {
	return o.Type_.E
}

func (o *OBSERVATION) SetModelName() {
	o.Type_ = utils.Some(OBSERVATION_MODEL_NAME)
	o.Name.SetModelName()
	if o.UID.E {
		o.UID.V.SetModelName()
	}
	if o.Links.E {
		for i := range o.Links.V {
			o.Links.V[i].SetModelName()
		}
	}
	if o.ArchetypeDetails.E {
		o.ArchetypeDetails.V.SetModelName()
	}
	if o.FeederAudit.E {
		o.FeederAudit.V.SetModelName()
	}
	o.Language.SetModelName()
	o.Encoding.SetModelName()
	if o.OtherParticipations.E {
		for i := range o.OtherParticipations.V {
			o.OtherParticipations.V[i].SetModelName()
		}
	}
	if o.WorkflowID.E {
		o.WorkflowID.V.SetModelName()
	}
	o.Subject.SetModelName()
	if o.Provider.E {
		o.Provider.V.SetModelName()
	}
	if o.Protocol.E {
		o.Protocol.V.SetModelName()
	}
	if o.GuidelineID.E {
		o.GuidelineID.V.SetModelName()
	}
	o.Data.SetModelName()
	if o.State.E {
		o.State.V.SetModelName()
	}
}

func (o *OBSERVATION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if o.Type_.E && o.Type_.V != OBSERVATION_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          OBSERVATION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + OBSERVATION_MODEL_NAME,
			Recommendation: "Set _type to " + OBSERVATION_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, o.Name.Validate(attrPath).Errs...)

	// Validate uid
	if o.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, o.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if o.Links.E {
		for i := range o.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, o.Links.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if o.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, o.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate feeder_audit
	if o.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, o.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate language
	attrPath = path + ".language"
	validateErr.Errs = append(validateErr.Errs, o.Language.Validate(attrPath).Errs...)

	// Validate encoding
	attrPath = path + ".encoding"
	validateErr.Errs = append(validateErr.Errs, o.Encoding.Validate(attrPath).Errs...)

	// Validate other_participations
	if o.OtherParticipations.E {
		for i := range o.OtherParticipations.V {
			attrPath = fmt.Sprintf("%s.other_participations[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, o.OtherParticipations.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate workflow_id
	if o.WorkflowID.E {
		attrPath = path + ".workflow_id"
		validateErr.Errs = append(validateErr.Errs, o.WorkflowID.V.Validate(attrPath).Errs...)
	}

	// Validate subject
	attrPath = path + ".subject"
	validateErr.Errs = append(validateErr.Errs, o.Subject.Validate(attrPath).Errs...)

	// Validate provider
	if o.Provider.E {
		attrPath = path + ".provider"
		validateErr.Errs = append(validateErr.Errs, o.Provider.V.Validate(attrPath).Errs...)
	}

	// Validate protocol
	if o.Protocol.E {
		attrPath = path + ".protocol"
		validateErr.Errs = append(validateErr.Errs, o.Protocol.V.Validate(attrPath).Errs...)
	}

	// Validate guideline_id
	if o.GuidelineID.E {
		attrPath = path + ".guideline_id"
		validateErr.Errs = append(validateErr.Errs, o.GuidelineID.V.Validate(attrPath).Errs...)
	}

	// Validate data
	attrPath = path + ".data"
	validateErr.Errs = append(validateErr.Errs, o.Data.Validate(attrPath).Errs...)

	// Validate state
	if o.State.E {
		attrPath = path + ".state"
		validateErr.Errs = append(validateErr.Errs, o.State.V.Validate(attrPath).Errs...)
	}

	return validateErr
}
