package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const EVENT_CONTEXT_MODEL_NAME string = "EVENT_CONTEXT"

type EVENT_CONTEXT struct {
	Type_              utils.Optional[string]           `json:"_type,omitzero"`
	StartTime          DV_DATE_TIME                     `json:"start_time"`
	EndTime            utils.Optional[DV_DATE_TIME]     `json:"end_time,omitzero"`
	Location           utils.Optional[string]           `json:"location,omitzero"`
	Setting            DV_CODED_TEXT                    `json:"setting"`
	OtherContext       utils.Optional[X_ITEM_STRUCTURE] `json:"other_context,omitzero"`
	HealthCareFacility utils.Optional[PARTY_IDENTIFIED] `json:"health_care_facility,omitzero"`
	Participations     utils.Optional[[]PARTICIPATION]  `json:"participations,omitzero"`
}

func (e *EVENT_CONTEXT) SetModelName() {
	e.Type_ = utils.Some(EVENT_CONTEXT_MODEL_NAME)
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

func (e *EVENT_CONTEXT) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if e.Type_.E && e.Type_.V != EVENT_CONTEXT_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          EVENT_CONTEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + EVENT_CONTEXT_MODEL_NAME,
			Recommendation: "Set _type to " + EVENT_CONTEXT_MODEL_NAME,
		})
	}

	// Validate start_time
	attrPath = path + ".start_time"
	validateErr.Errs = append(validateErr.Errs, e.StartTime.Validate(attrPath).Errs...)

	// Validate end_time
	if e.EndTime.E {
		attrPath = path + ".end_time"
		validateErr.Errs = append(validateErr.Errs, e.EndTime.V.Validate(attrPath).Errs...)
	}

	// Validate setting
	attrPath = path + ".setting"
	validateErr.Errs = append(validateErr.Errs, e.Setting.Validate(attrPath).Errs...)

	// Validate other_context
	if e.OtherContext.E {
		attrPath = path + ".other_context"
		validateErr.Errs = append(validateErr.Errs, e.OtherContext.V.Validate(attrPath).Errs...)
	}

	// Validate health_care_facility
	if e.HealthCareFacility.E {
		attrPath = path + ".health_care_facility"
		validateErr.Errs = append(validateErr.Errs, e.HealthCareFacility.V.Validate(attrPath).Errs...)
	}

	// Validate participations
	if e.Participations.E {
		if len(e.Participations.V) == 0 {
			attrPath = path + ".participations"
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          EVENT_CONTEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        "participations array cannot be empty",
				Recommendation: "Ensure participations array has at least one PARTICIPATION item",
			})
		} else {
			for i := range e.Participations.V {
				attrPath = fmt.Sprintf("%s.participations[%d]", path, i)
				validateErr.Errs = append(validateErr.Errs, e.Participations.V[i].Validate(attrPath).Errs...)
			}
		}
	}

	return validateErr
}
