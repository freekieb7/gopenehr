package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const PARTICIPATION_MODEL_NAME string = "PARTICIPATION"

type PARTICIPATION struct {
	Type_     util.Optional[string]        `json:"_type,omitzero"`
	Function  X_DV_TEXT                    `json:"function"`
	Mode      util.Optional[DV_CODED_TEXT] `json:"mode,omitzero"`
	Performer X_PARTY_PROXY                `json:"performer"`
	Time      util.Optional[DV_INTERVAL]   `json:"time,omitzero"`
}

func (p *PARTICIPATION) SetModelName() {
	p.Type_ = util.Some(PARTICIPATION_MODEL_NAME)
	p.Function.SetModelName()
	if p.Mode.E {
		p.Mode.V.SetModelName()
	}
	p.Performer.SetModelName()
	if p.Time.E {
		p.Time.V.SetModelName()
	}
}

func (p *PARTICIPATION) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if p.Type_.E && p.Type_.V != PARTICIPATION_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          PARTICIPATION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + PARTICIPATION_MODEL_NAME,
			Recommendation: "Set _type to " + PARTICIPATION_MODEL_NAME,
		})
	}

	// Validate function
	attrPath = path + ".function"
	errs = append(errs, p.Function.Validate(attrPath)...)

	// Validate mode
	if p.Mode.E {
		attrPath = path + ".mode"
		errs = append(errs, p.Mode.V.Validate(attrPath)...)
	}

	// Validate performer
	attrPath = path + ".performer"
	errs = append(errs, p.Performer.Validate(attrPath)...)

	// Validate time
	if p.Time.E {
		attrPath = path + ".time"
		errs = append(errs, p.Time.V.Validate(attrPath)...)
	}

	return errs
}
