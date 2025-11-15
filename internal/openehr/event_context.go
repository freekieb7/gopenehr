package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const EVENT_CONTEXT_MODEL_NAME string = "EVENT_CONTEXT"

type EVENT_CONTEXT struct {
	Type_              util.Optional[string]           `json:"_type,omitzero"`
	StartTime          DV_DATE_TIME                    `json:"start_time"`
	EndTime            util.Optional[DV_DATE_TIME]     `json:"end_time,omitzero"`
	Location           util.Optional[string]           `json:"location,omitzero"`
	Setting            DV_CODED_TEXT                   `json:"setting"`
	OtherContext       util.Optional[X_ITEM_STRUCTURE] `json:"other_context,omitzero"`
	HealthCareFacility util.Optional[PARTY_IDENTIFIED] `json:"health_care_facility,omitzero"`
	Participations     util.Optional[[]PARTICIPATION]  `json:"participations,omitzero"`
}

func (e *EVENT_CONTEXT) SetModelName() {
	e.Type_ = util.Some(EVENT_CONTEXT_MODEL_NAME)
	e.StartTime.SetModelName()
	if e.EndTime.E {
		e.EndTime.V.SetModelName()
	}
	if e.OtherContext.E {
		e.OtherContext.V.SetModelName()
	}
	if e.HealthCareFacility.E {
		e.HealthCareFacility.V.SetModelName()
	}
	if e.Participations.E {
		for i := range e.Participations.V {
			e.Participations.V[i].SetModelName()
		}
	}
}

func (e *EVENT_CONTEXT) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if e.Type_.E && e.Type_.V != EVENT_CONTEXT_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          EVENT_CONTEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + EVENT_CONTEXT_MODEL_NAME,
			Recommendation: "Set _type to " + EVENT_CONTEXT_MODEL_NAME,
		})
	}

	// Validate start_time
	attrPath = path + ".start_time"
	errs = append(errs, e.StartTime.Validate(attrPath)...)

	// Validate end_time
	if e.EndTime.E {
		attrPath = path + ".end_time"
		errs = append(errs, e.EndTime.V.Validate(attrPath)...)
	}

	// Validate setting
	attrPath = path + ".setting"
	errs = append(errs, e.Setting.Validate(attrPath)...)

	// Validate other_context
	if e.OtherContext.E {
		attrPath = path + ".other_context"
		errs = append(errs, e.OtherContext.V.Validate(attrPath)...)
	}

	// Validate health_care_facility
	if e.HealthCareFacility.E {
		attrPath = path + ".health_care_facility"
		errs = append(errs, e.HealthCareFacility.V.Validate(attrPath)...)
	}

	// Validate participations
	if e.Participations.E {
		if len(e.Participations.V) == 0 {
			attrPath = path + ".participations"
			errs = append(errs, util.ValidationError{
				Model:          EVENT_CONTEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        "participations array cannot be empty",
				Recommendation: "Ensure participations array has at least one PARTICIPATION item",
			})
		} else {
			for i := range e.Participations.V {
				attrPath = fmt.Sprintf("%s.participations[%d]", path, i)
				errs = append(errs, e.Participations.V[i].Validate(attrPath)...)
			}
		}
	}

	return errs
}
