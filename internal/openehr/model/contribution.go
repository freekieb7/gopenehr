package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const CONTRIBUTION_MODEL_NAME string = "CONTRIBUTION"

type CONTRIBUTION struct {
	Type_    utils.Optional[string] `json:"_type,omitzero"`
	UID      HIER_OBJECT_ID         `json:"uid"`
	Versions []OBJECT_REF           `json:"versions"`
	Audit    AUDIT_DETAILS          `json:"audit"`
}

func (c *CONTRIBUTION) SetModelName() {
	c.Type_ = utils.Some(CONTRIBUTION_MODEL_NAME)
	c.UID.SetModelName()
	for i := range c.Versions {
		c.Versions[i].SetModelName()
	}
	c.Audit.SetModelName()
}

func (c *CONTRIBUTION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if c.Type_.E && c.Type_.V != CONTRIBUTION_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          CONTRIBUTION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + CONTRIBUTION_MODEL_NAME,
			Recommendation: "Set _type to " + CONTRIBUTION_MODEL_NAME,
		})
	}

	// Validate uid
	attrPath = path + ".uid"
	validateErr.Errs = append(validateErr.Errs, c.UID.Validate(attrPath).Errs...)

	// Validate versions
	attrPath = path + ".versions"
	if len(c.Versions) == 0 {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          CONTRIBUTION_MODEL_NAME,
			Path:           attrPath,
			Message:        "versions must contain at least one OBJECT_REF",
			Recommendation: "Add at least one OBJECT_REF to versions",
		})
	} else {
		for i := range c.Versions {
			uidPath := fmt.Sprintf("%s[%d]", attrPath, i)
			validateErr.Errs = append(validateErr.Errs, c.Versions[i].Validate(uidPath).Errs...)
		}
	}

	// Validate audit
	attrPath = path + ".audit"
	validateErr.Errs = append(validateErr.Errs, c.Audit.Validate(attrPath).Errs...)

	return validateErr
}
