package rm

import (
	"fmt"
	"strings"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const HIER_OBJECT_ID_TYPE string = "HIER_OBJECT_ID"

type HIER_OBJECT_ID struct {
	Type_ utils.Optional[string] `json:"_type,omitzero"`
	Value string                 `json:"value"`
}

func (h *HIER_OBJECT_ID) SetModelName() {
	h.Type_ = utils.Some(HIER_OBJECT_ID_TYPE)
}

func (h *HIER_OBJECT_ID) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if h.Type_.E && h.Type_.V != HIER_OBJECT_ID_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          HIER_OBJECT_ID_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", HIER_OBJECT_ID_TYPE, h.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", HIER_OBJECT_ID_TYPE),
		})
	}

	// Validate UID-based identifier format: root '::' extension (extension is optional)
	attrPath = path + ".value"
	if h.Value == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          HIER_OBJECT_ID_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("%s value cannot be empty", HIER_OBJECT_ID_TYPE),
			Recommendation: fmt.Sprintf("Ensure %s value is set", HIER_OBJECT_ID_TYPE),
		})
	} else {
		// Split by '::' separator
		parts := strings.Split(h.Value, "::")

		// Must have 1 (root only) or 2 parts (root + extension)
		if len(parts) > 2 {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          HIER_OBJECT_ID_TYPE,
				Path:           attrPath,
				Message:        fmt.Sprintf("%s invalid format: too many '::'", HIER_OBJECT_ID_TYPE),
				Recommendation: fmt.Sprintf("Ensure %s value is in the format 'root::extension'", HIER_OBJECT_ID_TYPE),
			})
			return validateErr
		}

		// Validate root part (first part)
		root := parts[0]
		if root == "" {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          HIER_OBJECT_ID_TYPE,
				Path:           attrPath,
				Message:        fmt.Sprintf("%s root part cannot be empty in '%s'", HIER_OBJECT_ID_TYPE, h.Value),
				Recommendation: fmt.Sprintf("Ensure %s value has a non-empty root part", HIER_OBJECT_ID_TYPE),
			})
			return validateErr
		}

		// Root should be a valid UID (UUID, ISO_OID, or INTERNET_ID format)
		if err := util.ValidateUID(root); err != nil {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          HIER_OBJECT_ID_TYPE,
				Path:           attrPath,
				Message:        fmt.Sprintf("%s invalid root UID '%s': %v", HIER_OBJECT_ID_TYPE, root, err),
				Recommendation: fmt.Sprintf("Ensure %s root part is a valid UUID, ISO_OID, or INTERNET_ID", HIER_OBJECT_ID_TYPE),
			})
		}

		// If extension exists, validate it's not empty
		if len(parts) == 2 {
			extension := parts[1]
			if extension == "" {
				validateErr.Errs = append(validateErr.Errs, util.ValidationError{
					Model:          HIER_OBJECT_ID_TYPE,
					Path:           attrPath,
					Message:        fmt.Sprintf("%s extension cannot be empty when '::' is present", HIER_OBJECT_ID_TYPE),
					Recommendation: fmt.Sprintf("Ensure %s value has a non-empty extension part", HIER_OBJECT_ID_TYPE),
				})
			}
		}
	}

	return validateErr
}
