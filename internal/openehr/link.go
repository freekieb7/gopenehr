package openehr

import "github.com/freekieb7/gopenehr/internal/openehr/util"

const LINK_MODEL_NAME string = "LINK"

type LINK struct {
	Type_   util.Optional[string] `json:"_type,omitzero"`
	Meaning X_DV_TEXT             `json:"meaning"`
	Type    X_DV_TEXT             `json:"type"`
	Target  DV_EHR_URI            `json:"target"`
}

func (l *LINK) SetModelName() {
	l.Type_ = util.Some(LINK_MODEL_NAME)
	l.Meaning.SetModelName()
	l.Type.SetModelName()
	l.Target.SetModelName()
}

func (l *LINK) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if l.Type_.E && l.Type_.V != LINK_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          LINK_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to LINK",
		})
	}

	return errors
}
