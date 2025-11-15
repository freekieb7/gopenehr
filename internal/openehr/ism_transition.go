package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const ISM_TRANSITION_MODEL_NAME string = "ISM_TRANSITION"

type ISM_TRANSITION struct {
	Type_        util.Optional[string]        `json:"_type,omitzero"`
	CurrentState DV_CODED_TEXT                `json:"current_state"`
	Transition   util.Optional[DV_CODED_TEXT] `json:"transition,omitzero"`
	CareflowStep util.Optional[DV_CODED_TEXT] `json:"careflow_step,omitzero"`
	Reason       util.Optional[X_DV_TEXT]     `json:"reason,omitzero"`
}

func (i ISM_TRANSITION) HasModelName() bool {
	return i.Type_.E
}

func (i *ISM_TRANSITION) SetModelName() {
	i.Type_ = util.Some(ISM_TRANSITION_MODEL_NAME)
	i.CurrentState.SetModelName()
	if i.Transition.E {
		i.Transition.V.SetModelName()
	}
	if i.CareflowStep.E {
		i.CareflowStep.V.SetModelName()
	}
	if i.Reason.E {
		i.Reason.V.SetModelName()
	}
}

func (i *ISM_TRANSITION) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if i.Type_.E && i.Type_.V != ISM_TRANSITION_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          ISM_TRANSITION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + ISM_TRANSITION_MODEL_NAME,
			Recommendation: "Set _type to " + ISM_TRANSITION_MODEL_NAME,
		})
	}

	// Validate current_state
	attrPath = path + ".current_state"
	errs = append(errs, i.CurrentState.Validate(attrPath)...)

	// Validate transition
	if i.Transition.E {
		attrPath = path + ".transition"
		errs = append(errs, i.Transition.V.Validate(attrPath)...)
	}

	// Validate careflow_step
	if i.CareflowStep.E {
		attrPath = path + ".careflow_step"
		errs = append(errs, i.CareflowStep.V.Validate(attrPath)...)
	}

	// Validate reason
	if i.Reason.E {
		attrPath = path + ".reason"
		errs = append(errs, i.Reason.V.Validate(attrPath)...)
	}

	return errs
}
