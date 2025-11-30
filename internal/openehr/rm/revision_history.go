package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const REVISION_HISTORY_TYPE = "REVISION_HISTORY"

type REVISION_HISTORY struct {
	Type_ utils.Optional[string]  `json:"_type,omitzero"`
	Items []REVISION_HISTORY_ITEM `json:"items"`
}

func (r *REVISION_HISTORY) SetModelName() {
	r.Type_ = utils.Some(REVISION_HISTORY_TYPE)
	for i := range r.Items {
		r.Items[i].SetModelName()
	}
}

func (r *REVISION_HISTORY) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if r.Type_.E && r.Type_.V != REVISION_HISTORY_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          REVISION_HISTORY_TYPE,
			Path:           attrPath,
			Message:        "_type must be " + REVISION_HISTORY_TYPE,
			Recommendation: "Set _type to " + REVISION_HISTORY_TYPE,
		})
	}

	// Validate items
	for i := range r.Items {
		attrPath = fmt.Sprintf("%s.items[%d]", path, i)
		validateErr.Errs = append(validateErr.Errs, r.Items[i].Validate(attrPath).Errs...)
	}

	return validateErr
}
