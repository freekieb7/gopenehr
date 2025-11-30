package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const EVALUATION_MODEL_NAME string = "EVALUATION"

type EVALUATION struct {
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
	Data                X_ITEM_STRUCTURE                 `json:"data"`
}

func (e *EVALUATION) isContentItemModel() {}

func (e *EVALUATION) HasModelName() bool {
	return e.Type_.E
}

func (e *EVALUATION) SetModelName() {
	e.Type_ = utils.Some(EVALUATION_MODEL_NAME)
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
	if e.FeederAudit.E {
		e.FeederAudit.V.SetModelName()
	}
	e.Language.SetModelName()
	e.Encoding.SetModelName()
	if e.OtherParticipations.E {
		for i := range e.OtherParticipations.V {
			e.OtherParticipations.V[i].SetModelName()
		}
	}
	if e.WorkflowID.E {
		e.WorkflowID.V.SetModelName()
	}
	e.Subject.SetModelName()
	if e.Provider.E {
		e.Provider.V.SetModelName()
	}
	if e.Protocol.E {
		e.Protocol.V.SetModelName()
	}
	if e.GuidelineID.E {
		e.GuidelineID.V.SetModelName()
	}
	e.Data.SetModelName()
}

func (e *EVALUATION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if e.Type_.E && e.Type_.V != EVALUATION_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          EVALUATION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + EVALUATION_MODEL_NAME,
			Recommendation: "Set _type to " + EVALUATION_MODEL_NAME,
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

	// Validate feeder_audit
	if e.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, e.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate language
	attrPath = path + ".language"
	validateErr.Errs = append(validateErr.Errs, e.Language.Validate(attrPath).Errs...)

	// Validate encoding
	attrPath = path + ".encoding"
	validateErr.Errs = append(validateErr.Errs, e.Encoding.Validate(attrPath).Errs...)
	// Validate other_participations
	if e.OtherParticipations.E {
		for i := range e.OtherParticipations.V {
			attrPath = fmt.Sprintf("%s.other_participations[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, e.OtherParticipations.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate workflow_id
	if e.WorkflowID.E {
		attrPath = path + ".workflow_id"
		validateErr.Errs = append(validateErr.Errs, e.WorkflowID.V.Validate(attrPath).Errs...)
	}

	// Validate subject
	attrPath = path + ".subject"
	validateErr.Errs = append(validateErr.Errs, e.Subject.Validate(attrPath).Errs...)

	// Validate provider
	if e.Provider.E {
		attrPath = path + ".provider"
		validateErr.Errs = append(validateErr.Errs, e.Provider.V.Validate(attrPath).Errs...)
	}

	// Validate protocol
	if e.Protocol.E {
		attrPath = path + ".protocol"
		validateErr.Errs = append(validateErr.Errs, e.Protocol.V.Validate(attrPath).Errs...)
	}

	// Validate guideline_id
	if e.GuidelineID.E {
		attrPath = path + ".guideline_id"
		validateErr.Errs = append(validateErr.Errs, e.GuidelineID.V.Validate(attrPath).Errs...)
	}

	// Validate data
	attrPath = path + ".data"
	validateErr.Errs = append(validateErr.Errs, e.Data.Validate(attrPath).Errs...)

	return validateErr
}
