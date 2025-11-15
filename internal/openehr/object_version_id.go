package openehr

import (
	"fmt"
	"strings"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const OBJECT_VERSION_ID_MODEL_NAME string = "OBJECT_VERSION_ID"

type OBJECT_VERSION_ID struct {
	Type_ util.Optional[string] `json:"_type,omitzero"`
	Value string                `json:"value"`
}

func (o OBJECT_VERSION_ID) isUidBasedIDModel() {}

func (o OBJECT_VERSION_ID) isObjectIDModel() {}

func (o OBJECT_VERSION_ID) HasModelName() bool {
	return o.Type_.E
}

func (o *OBJECT_VERSION_ID) SetModelName() {
	o.Type_ = util.Some(OBJECT_VERSION_ID_MODEL_NAME)
}

func (o OBJECT_VERSION_ID) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if o.Type_.E && o.Type_.V != OBJECT_VERSION_ID_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          OBJECT_VERSION_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", OBJECT_VERSION_ID_MODEL_NAME, o.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", OBJECT_VERSION_ID_MODEL_NAME),
		})
	}

	// Validate value
	attrPath = path + ".value"
	if o.Value == "" {
		errors = append(errors, util.ValidationError{
			Model:          OBJECT_VERSION_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	} else {
		// lexical form: object_id '::' creating_system_id '::' version_tree_id.
		parts := strings.Split(o.Value, "::")
		if len(parts) != 3 {
			errors = append(errors, util.ValidationError{
				Model:          OBJECT_VERSION_ID_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid value format: %s", o.Value),
				Recommendation: "Ensure value field follows the lexical form: object_id '::' creating_system_id '::' version_tree_id",
			})
		}

		// First part UID
		uid := parts[0]
		if err := util.ValidateUID(uid); err != nil {
			errors = append(errors, util.ValidationError{
				Model:          OBJECT_VERSION_ID_MODEL_NAME,
				Path:           attrPath + ".object_id",
				Message:        fmt.Sprintf("invalid object_id format: %s", uid),
				Recommendation: "Ensure object_id follows a valid UID format",
			})
		}

		// Second part creating_system_id
		creatingSystemID := parts[1]
		if creatingSystemID == "" {
			errors = append(errors, util.ValidationError{
				Model:          OBJECT_VERSION_ID_MODEL_NAME,
				Path:           attrPath + ".creating_system_id",
				Message:        "creating_system_id cannot be empty",
				Recommendation: "Ensure creating_system_id is not empty",
			})
		}

		// Third part version_tree_id
		// Lexical form: trunk_version [ '.' branch_number '.' branch_version ]
		versionTreeID := parts[2]
		if versionTreeID == "" {
			errors = append(errors, util.ValidationError{
				Model:          OBJECT_VERSION_ID_MODEL_NAME,
				Path:           attrPath + ".version_tree_id",
				Message:        "version_tree_id cannot be empty",
				Recommendation: "Ensure version_tree_id is not empty",
			})
		} else if !util.VersionTreeIDRegex.MatchString(versionTreeID) {
			errors = append(errors, util.ValidationError{
				Model:          OBJECT_VERSION_ID_MODEL_NAME,
				Path:           attrPath + ".version_tree_id",
				Message:        fmt.Sprintf("invalid version_tree_id format: %s", versionTreeID),
				Recommendation: "Ensure version_tree_id follows the lexical form: trunk_version [ '.' branch_number '.' branch_version ]",
			})
		}
	}

	return errors
}

func (o OBJECT_VERSION_ID) UID() string {
	return strings.Split(o.Value, "::")[0]
}
