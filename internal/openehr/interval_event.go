package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const INTERVAL_EVENT_MODEL_NAME string = "INTERVAL_EVENT"

type INTERVAL_EVENT struct {
	Type_            util.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                       `json:"name"`
	ArchetypeNodeID  string                          `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Time             DV_DATE_TIME                    `json:"time"`
	State            util.Optional[X_ITEM_STRUCTURE] `json:"state,omitzero"`
	Data             X_ITEM_STRUCTURE                `json:"data"`
	Width            DV_DURATION                     `json:"width"`
	SampleCount      util.Optional[int64]            `json:"sample_count,omitzero"`
	MathFunction     DV_CODED_TEXT                   `json:"math_function"`
}

func (i *INTERVAL_EVENT) isEventModel() {}

func (i *INTERVAL_EVENT) HasModelName() bool {
	return i.Type_.E
}

func (i *INTERVAL_EVENT) SetModelName() {
	i.Type_ = util.Some(INTERVAL_EVENT_MODEL_NAME)
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

func (i *INTERVAL_EVENT) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != INTERVAL_EVENT_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          INTERVAL_EVENT_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + INTERVAL_EVENT_MODEL_NAME,
			Recommendation: "Set _type to " + INTERVAL_EVENT_MODEL_NAME,
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

	// Validate time
	attrPath = path + ".time"
	errs = append(errs, i.Time.Validate(attrPath)...)

	// Validate state
	if i.State.E {
		attrPath = path + ".state"
		errs = append(errs, i.State.V.Validate(attrPath)...)
	}

	// Validate data
	attrPath = path + ".data"
	errs = append(errs, i.Data.Validate(attrPath)...)

	// Validate width
	attrPath = path + ".width"
	errs = append(errs, i.Width.Validate(attrPath)...)

	// Validate sample_count
	// if i.SampleCount.E {
	// 	// int64 has no Validate
	// }

	// Validate math_function
	attrPath = path + ".math_function"
	errs = append(errs, i.MathFunction.Validate(attrPath)...)

	return errs
}
