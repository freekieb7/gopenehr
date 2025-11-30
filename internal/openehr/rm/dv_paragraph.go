package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_PARAGRAPH_TYPE string = "DV_PARAGRAPH"

type DV_PARAGRAPH struct {
	Type_ utils.Optional[string] `json:"_type,omitzero"`
	Items []DvTextUnion          `json:"items"`
}

func (d *DV_PARAGRAPH) SetModelName() {
	d.Type_ = utils.Some(DV_PARAGRAPH_TYPE)
	for i := range d.Items {
		d.Items[i].SetModelName()
	}
}

func (d *DV_PARAGRAPH) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_PARAGRAPH_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_PARAGRAPH_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_PARAGRAPH",
		})
	}

	// Validate items
	for i := range d.Items {
		itemPath := fmt.Sprintf("%s.items[%d]", path, i)
		validateErr.Errs = append(validateErr.Errs, d.Items[i].Validate(itemPath).Errs...)
	}

	return validateErr
}
