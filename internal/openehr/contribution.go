package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const CONTRIBUTION_MODEL_NAME string = "CONTRIBUTION"

type CONTRIBUTION struct {
	Type_    util.Optional[string] `json:"_type,omitzero"`
	UID      HIER_OBJECT_ID        `json:"uid"`
	Versions []OBJECT_REF          `json:"versions"`
	Audit    AUDIT_DETAILS         `json:"audit"`
}

func (c *CONTRIBUTION) SetModelName() {
	c.Type_ = util.Some(CONTRIBUTION_MODEL_NAME)
	c.UID.SetModelName()
	for i := range c.Versions {
		c.Versions[i].SetModelName()
	}
	c.Audit.SetModelName()
}

func (c *CONTRIBUTION) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if c.Type_.E && c.Type_.V != CONTRIBUTION_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          CONTRIBUTION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + CONTRIBUTION_MODEL_NAME,
			Recommendation: "Set _type to " + CONTRIBUTION_MODEL_NAME,
		})
	}

	// Validate uid
	attrPath = path + ".uid"
	errs = append(errs, c.UID.Validate(attrPath)...)

	// Validate versions
	attrPath = path + ".versions"
	if len(c.Versions) == 0 {
		errs = append(errs, util.ValidationError{
			Model:          CONTRIBUTION_MODEL_NAME,
			Path:           attrPath,
			Message:        "versions must contain at least one OBJECT_REF",
			Recommendation: "Add at least one OBJECT_REF to versions",
		})
	} else {
		for i := range c.Versions {
			uidPath := fmt.Sprintf("%s[%d]", attrPath, i)
			errs = append(errs, c.Versions[i].Validate(uidPath)...)
		}
	}

	// Validate audit
	attrPath = path + ".audit"
	errs = append(errs, c.Audit.Validate(attrPath)...)

	return errs
}
