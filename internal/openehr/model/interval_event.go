package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const INTERVAL_EVENT_MODEL_NAME string = "INTERVAL_EVENT"

type INTERVAL_EVENT struct {
	Type_            utils.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                        `json:"name"`
	ArchetypeNodeID  string                           `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Time             DV_DATE_TIME                     `json:"time"`
	State            utils.Optional[X_ITEM_STRUCTURE] `json:"state,omitzero"`
	Data             X_ITEM_STRUCTURE                 `json:"data"`
	Width            DV_DURATION                      `json:"width"`
	SampleCount      utils.Optional[int64]            `json:"sample_count,omitzero"`
	MathFunction     DV_CODED_TEXT                    `json:"math_function"`
}

func (i *INTERVAL_EVENT) isEventModel() {}

func (i *INTERVAL_EVENT) HasModelName() bool {
	return i.Type_.E
}

func (i *INTERVAL_EVENT) SetModelName() {
	i.Type_ = utils.Some(INTERVAL_EVENT_MODEL_NAME)
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
	i.Time.SetModelName()
	if i.State.E {
		i.State.V.SetModelName()
	}
	i.Data.SetModelName()
	i.Width.SetModelName()
	// if i.SampleCount.E {
	// 	// int64 has no SetModelName
	// }
	i.MathFunction.SetModelName()
}

func (i *INTERVAL_EVENT) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != INTERVAL_EVENT_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          INTERVAL_EVENT_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + INTERVAL_EVENT_MODEL_NAME,
			Recommendation: "Set _type to " + INTERVAL_EVENT_MODEL_NAME,
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

	// Validate time
	attrPath = path + ".time"
	validateErr.Errs = append(validateErr.Errs, i.Time.Validate(attrPath).Errs...)

	// Validate state
	if i.State.E {
		attrPath = path + ".state"
		validateErr.Errs = append(validateErr.Errs, i.State.V.Validate(attrPath).Errs...)
	}

	// Validate data
	attrPath = path + ".data"
	validateErr.Errs = append(validateErr.Errs, i.Data.Validate(attrPath).Errs...)

	// Validate width
	attrPath = path + ".width"
	validateErr.Errs = append(validateErr.Errs, i.Width.Validate(attrPath).Errs...)

	// Validate sample_count
	// if i.SampleCount.E {
	// 	// int64 has no Validate
	// }

	// Validate math_function
	attrPath = path + ".math_function"
	validateErr.Errs = append(validateErr.Errs, i.MathFunction.Validate(attrPath).Errs...)

	return validateErr
}
