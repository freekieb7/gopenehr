package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const REVISION_HISTORY_MODEL_NAME = "REVISION_HISTORY"

type REVISION_HISTORY struct {
	Type_ util.Optional[string]   `json:"_type,omitzero"`
	Items []REVISION_HISTORY_ITEM `json:"items"`
}

func (r *REVISION_HISTORY) SetModelName() {
	r.Type_ = util.Some(REVISION_HISTORY_MODEL_NAME)
	for i := range r.Items {
		r.Items[i].SetModelName()
	}
}

func (r *REVISION_HISTORY) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if r.Type_.E && r.Type_.V != REVISION_HISTORY_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          REVISION_HISTORY_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + REVISION_HISTORY_MODEL_NAME,
			Recommendation: "Set _type to " + REVISION_HISTORY_MODEL_NAME,
		})
	}

	// Validate items
	for i := range r.Items {
		attrPath = fmt.Sprintf("%s.items[%d]", path, i)
		errs = append(errs, r.Items[i].Validate(attrPath)...)
	}

	return errs
}
