package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const REVISION_HISTORY_ITEM_MODEL_NAME = "REVISION_HISTORY_ITEM"

type REVISION_HISTORY_ITEM struct {
	Type_     util.Optional[string] `json:"_type,omitzero"`
	VersionID OBJECT_VERSION_ID     `json:"version_id"`
	Audits    []AUDIT_DETAILS       `json:"audits"`
}

func (r *REVISION_HISTORY_ITEM) SetModelName() {
	r.Type_ = util.Some(REVISION_HISTORY_ITEM_MODEL_NAME)
	r.VersionID.SetModelName()
	for i := range r.Audits {
		r.Audits[i].SetModelName()
	}
}

func (r *REVISION_HISTORY_ITEM) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if r.Type_.E && r.Type_.V != REVISION_HISTORY_ITEM_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model: REVISION_HISTORY_ITEM_MODEL_NAME,
			Path:  attrPath,
		})
	}

	// Validate version_id
	attrPath = path + ".version_id"
	errs = append(errs, r.VersionID.Validate(attrPath)...)

	// Validate audits
	if len(r.Audits) == 0 {
		attrPath = path + ".audits"
		errs = append(errs, util.ValidationError{
			Model:          REVISION_HISTORY_ITEM_MODEL_NAME,
			Path:           attrPath,
			Message:        "audits array cannot be empty",
			Recommendation: "Ensure audits array has at least one AUDIT_DETAILS item",
		})
	}

	for i := range r.Audits {
		attrPath = fmt.Sprintf("%s.audits[%d]", path, i)
		errs = append(errs, r.Audits[i].Validate(attrPath)...)
	}

	return errs
}
