package rm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const OBJECT_VERSION_ID_TYPE string = "OBJECT_VERSION_ID"

type OBJECT_VERSION_ID struct {
	Type_ utils.Optional[string] `json:"_type,omitzero"`
	Value string                 `json:"value"`
}

func (o *OBJECT_VERSION_ID) SetModelName() {
	o.Type_ = utils.Some(OBJECT_VERSION_ID_TYPE)
}

func (o *OBJECT_VERSION_ID) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if o.Type_.E && o.Type_.V != OBJECT_VERSION_ID_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          OBJECT_VERSION_ID_TYPE,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", OBJECT_VERSION_ID_TYPE, o.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", OBJECT_VERSION_ID_TYPE),
		})
	}

	// Validate value
	attrPath = path + ".value"
	if o.Value == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          OBJECT_VERSION_ID_TYPE,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	} else {
		// lexical form: object_id '::' creating_system_id '::' version_tree_id.
		parts := strings.Split(o.Value, "::")
		if len(parts) != 3 {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          OBJECT_VERSION_ID_TYPE,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid value format: %s", o.Value),
				Recommendation: "Ensure value field follows the lexical form: object_id '::' creating_system_id '::' version_tree_id",
			})
		}

		// First part UID
		uid := parts[0]
		if err := util.ValidateUID(uid); err != nil {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          OBJECT_VERSION_ID_TYPE,
				Path:           attrPath + ".object_id",
				Message:        fmt.Sprintf("invalid object_id format: %s", uid),
				Recommendation: "Ensure object_id follows a valid UID format",
			})
		}

		// Second part creating_system_id
		creatingSystemID := parts[1]
		if creatingSystemID == "" {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          OBJECT_VERSION_ID_TYPE,
				Path:           attrPath + ".creating_system_id",
				Message:        "creating_system_id cannot be empty",
				Recommendation: "Ensure creating_system_id is not empty",
			})
		}

		// Third part version_tree_id
		// Lexical form: trunk_version [ '.' branch_number '.' branch_version ]
		versionTreeID := parts[2]
		if versionTreeID == "" {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          OBJECT_VERSION_ID_TYPE,
				Path:           attrPath + ".version_tree_id",
				Message:        "version_tree_id cannot be empty",
				Recommendation: "Ensure version_tree_id is not empty",
			})
		} else if !util.VersionTreeIDRegex.MatchString(versionTreeID) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          OBJECT_VERSION_ID_TYPE,
				Path:           attrPath + ".version_tree_id",
				Message:        fmt.Sprintf("invalid version_tree_id format: %s", versionTreeID),
				Recommendation: "Ensure version_tree_id follows the lexical form: trunk_version [ '.' branch_number '.' branch_version ]",
			})
		}
	}

	return validateErr
}

func (o OBJECT_VERSION_ID) UID() string {
	return strings.Split(o.Value, "::")[0]
}

func (o OBJECT_VERSION_ID) SystemID() string {
	return strings.Split(o.Value, "::")[1]
}

func (o OBJECT_VERSION_ID) VersionTreeID() VersionTreeID {
	versionTreeIDStr := strings.Split(o.Value, "::")[2]

	parseVersionTreeIDPart := func(versionTreeIDStr string, partIndex int) uint8 {
		parts := strings.Split(versionTreeIDStr, ".")
		if partIndex >= len(parts) {
			return 0
		}
		partValue, err := strconv.ParseUint(parts[partIndex], 10, 8)
		if err != nil {
			return 0
		}
		return uint8(partValue)
	}

	return VersionTreeID{
		Major: parseVersionTreeIDPart(versionTreeIDStr, 0),
		Minor: parseVersionTreeIDPart(versionTreeIDStr, 1),
		Patch: parseVersionTreeIDPart(versionTreeIDStr, 2),
	}
}

type VersionTreeID struct {
	Major uint8
	Minor uint8
	Patch uint8
}

func (v VersionTreeID) String() string {
	if v.Minor == 0 && v.Patch == 0 {
		return fmt.Sprintf("%d", v.Major)
	}

	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v VersionTreeID) Int() int32 {
	return int32(v.Major)<<16 | int32(v.Minor)<<8 | int32(v.Patch)
}

func VersionTreeIDFromString(s string) VersionTreeID {
	parts := strings.Split(s, ".")
	var major, minor, patch uint8

	if len(parts) > 0 {
		if val, err := strconv.ParseUint(parts[0], 10, 8); err == nil {
			major = uint8(val)
		}
	}
	if len(parts) > 1 {
		if val, err := strconv.ParseUint(parts[1], 10, 8); err == nil {
			minor = uint8(val)
		}
	}
	if len(parts) > 2 {
		if val, err := strconv.ParseUint(parts[2], 10, 8); err == nil {
			patch = uint8(val)
		}
	}

	return VersionTreeID{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func VersionTreeIDFromInt(i int32) VersionTreeID {
	major := uint8((i >> 16) & 0xFF)
	minor := uint8((i >> 8) & 0xFF)
	patch := uint8(i & 0xFF)

	return VersionTreeID{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func (v VersionTreeID) CompareTo(other VersionTreeID) int8 {
	if v == other {
		return 0
	}

	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		}
		return -1
	}

	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		}
		return -1
	}

	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		}
		return -1
	}

	return 0
}
