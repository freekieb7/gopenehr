package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const OBSERVATION_MODEL_NAME string = "OBSERVATION"

type OBSERVATION struct {
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
	Data                HISTORY                         `json:"data"`
	State               util.Optional[HISTORY]          `json:"state,omitzero"`
}

func (o *OBSERVATION) isContentItemModel() {}

func (o *OBSERVATION) HasModelName() bool {
	return o.Type_.E
}

func (o *OBSERVATION) SetModelName() {
	o.Type_ = util.Some(OBSERVATION_MODEL_NAME)
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

func (o *OBSERVATION) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if o.Type_.E && o.Type_.V != OBSERVATION_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          OBSERVATION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + OBSERVATION_MODEL_NAME,
			Recommendation: "Set _type to " + OBSERVATION_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	errs = append(errs, o.Name.Validate(attrPath)...)

	// Validate uid
	if o.UID.E {
		attrPath = path + ".uid"
		errs = append(errs, o.UID.V.Validate(attrPath)...)
	}

	// Validate links
	if o.Links.E {
		for i := range o.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			errs = append(errs, o.Links.V[i].Validate(attrPath)...)
		}
	}

	// Validate archetype_details
	if o.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errs = append(errs, o.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate feeder_audit
	if o.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errs = append(errs, o.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate language
	attrPath = path + ".language"
	errs = append(errs, o.Language.Validate(attrPath)...)

	// Validate encoding
	attrPath = path + ".encoding"
	errs = append(errs, o.Encoding.Validate(attrPath)...)

	// Validate other_participations
	if o.OtherParticipations.E {
		for i := range o.OtherParticipations.V {
			attrPath = fmt.Sprintf("%s.other_participations[%d]", path, i)
			errs = append(errs, o.OtherParticipations.V[i].Validate(attrPath)...)
		}
	}

	// Validate workflow_id
	if o.WorkflowID.E {
		attrPath = path + ".workflow_id"
		errs = append(errs, o.WorkflowID.V.Validate(attrPath)...)
	}

	// Validate subject
	attrPath = path + ".subject"
	errs = append(errs, o.Subject.Validate(attrPath)...)

	// Validate provider
	if o.Provider.E {
		attrPath = path + ".provider"
		errs = append(errs, o.Provider.V.Validate(attrPath)...)
	}

	// Validate protocol
	if o.Protocol.E {
		attrPath = path + ".protocol"
		errs = append(errs, o.Protocol.V.Validate(attrPath)...)
	}

	// Validate guideline_id
	if o.GuidelineID.E {
		attrPath = path + ".guideline_id"
		errs = append(errs, o.GuidelineID.V.Validate(attrPath)...)
	}

	// Validate data
	attrPath = path + ".data"
	errs = append(errs, o.Data.Validate(attrPath)...)

	// Validate state
	if o.State.E {
		attrPath = path + ".state"
		errs = append(errs, o.State.V.Validate(attrPath)...)
	}

	return errs
}
