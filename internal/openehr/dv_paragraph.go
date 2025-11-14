package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_PARAGRAPH_MODEL_NAME string = "DV_PARAGRAPH"

type DV_PARAGRAPH struct {
	Type_ util.Optional[string] `json:"_type,omitzero"`
	Items []X_DV_TEXT           `json:"items"`
}

func (d DV_PARAGRAPH) isDataValueModel() {}

func (d DV_PARAGRAPH) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_PARAGRAPH) SetModelName() {
	d.Type_ = util.Some(DV_PARAGRAPH_MODEL_NAME)
	for i := range d.Items {
		d.Items[i].SetModelName()
	}
}

func (d DV_PARAGRAPH) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_PARAGRAPH_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          DV_PARAGRAPH_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to DV_PARAGRAPH",
		})
	}

	// Validate items
	for i := range d.Items {
		itemPath := fmt.Sprintf("%s.items[%d]", path, i)
		errors = append(errors, d.Items[i].Validate(itemPath)...)
	}

	return errors
}
