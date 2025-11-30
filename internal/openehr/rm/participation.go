package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const PARTICIPATION_TYPE string = "PARTICIPATION"

type PARTICIPATION struct {
	Type_     utils.Optional[string]        `json:"_type,omitzero"`
	Function  DvTextUnion                   `json:"function"`
	Mode      utils.Optional[DV_CODED_TEXT] `json:"mode,omitzero"`
	Performer PartyProxyUnion               `json:"performer"`
	Time      utils.Optional[DV_INTERVAL]   `json:"time,omitzero"`
}

func (p *PARTICIPATION) SetModelName() {
	p.Type_ = utils.Some(PARTICIPATION_TYPE)
	p.Function.SetModelName()
	if p.Mode.E {
		p.Mode.V.SetModelName()
	}
	p.Performer.SetModelName()
	if p.Time.E {
		p.Time.V.SetModelName()
	}
}

func (p *PARTICIPATION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if p.Type_.E && p.Type_.V != PARTICIPATION_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          PARTICIPATION_TYPE,
			Path:           attrPath,
			Message:        "_type must be " + PARTICIPATION_TYPE,
			Recommendation: "Set _type to " + PARTICIPATION_TYPE,
		})
	}

	// Validate function
	attrPath = path + ".function"
	validateErr.Errs = append(validateErr.Errs, p.Function.Validate(attrPath).Errs...)

	// Validate mode
	if p.Mode.E {
		attrPath = path + ".mode"
		validateErr.Errs = append(validateErr.Errs, p.Mode.V.Validate(attrPath).Errs...)
	}

	// Validate performer
	attrPath = path + ".performer"
	validateErr.Errs = append(validateErr.Errs, p.Performer.Validate(attrPath).Errs...)

	// Validate time
	if p.Time.E {
		attrPath = path + ".time"
		validateErr.Errs = append(validateErr.Errs, p.Time.V.Validate(attrPath).Errs...)
	}

	return validateErr
}
