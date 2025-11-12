package openehr

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
)

type ReferenceModel interface {
	HasModelName() bool
	Validate(path string) []ValidationError
}

const EHR_MODEL_NAME string = "EHR"

var _ ReferenceModel = (*EHR)(nil)

type EHR struct {
	Type_         Optional[string]         `json:"_type,omitzero"`
	SystemID      Optional[HIER_OBJECT_ID] `json:"system_id,omitzero"`
	EHRID         HIER_OBJECT_ID           `json:"ehr_id"`
	Contributions Optional[[]OBJECT_REF]   `json:"contributions,omitzero"`
	EHRStatus     OBJECT_REF               `json:"ehr_status"`
	EHRAccess     OBJECT_REF               `json:"ehr_access"`
	Compositions  Optional[[]OBJECT_REF]   `json:"compositions,omitzero"`
	Directory     Optional[OBJECT_REF]     `json:"directory,omitzero"`
	TimeCreated   DV_DATE_TIME             `json:"time_created"`
	Folders       Optional[[]OBJECT_REF]   `json:"folders,omitzero"`
	Tags          Optional[[]OBJECT_REF]   `json:"tags,omitzero"`
}

func (e EHR) HasModelName() bool {
	return e.Type_.IsSet()
}

func (e EHR) Validate(path string) []ValidationError {
	var errs []ValidationError
	var attrPath string

	// Validate _type
	if e.Type_.IsSet() && e.Type_.Unwrap() != EHR_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, ValidationError{
			Model:          EHR_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid EHR _type field: %s", e.Type_.Unwrap()),
			Recommendation: "Ensure the _type field is set to 'EHR'",
		})
	}

	// Validate system_id
	if e.SystemID.IsSet() {
		attrPath = path + ".system_id"
		errs = append(errs, e.SystemID.Unwrap().Validate(attrPath)...)
	}

	// Validate ehr_id
	attrPath = path + ".ehr_id"
	errs = append(errs, e.EHRID.Validate(attrPath)...)

	// Validate contributions
	if e.Contributions.IsSet() {
		for i, contribRef := range e.Contributions.Unwrap() {
			attrPath = path + fmt.Sprintf(".contributions[%d]", i)
			if contribRef.Type != CONTRIBUTION_MODEL_NAME {
				errs = append(errs, ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid contribution type: %s", contribRef.Type),
					Recommendation: fmt.Sprintf("Ensure contributions[%d] _type field is set to '%s'", i, CONTRIBUTION_MODEL_NAME),
				})
			}
			errs = append(errs, contribRef.Validate(attrPath)...)
		}
	}

	// Validate ehr_status
	attrPath = path + ".ehr_status"
	if e.EHRStatus.Type != VERSIONED_EHR_STATUS_MODEL_NAME {
		errs = append(errs, ValidationError{
			Model:          EHR_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid EHR status type: %s", e.EHRStatus.Type),
			Recommendation: fmt.Sprintf("Ensure ehr_status _type field is set to '%s'", VERSIONED_EHR_STATUS_MODEL_NAME),
		})

	}
	errs = append(errs, e.EHRStatus.Validate(attrPath)...)

	// Validate ehr_access
	attrPath = path + ".ehr_access"
	if e.EHRAccess.Type != VERSIONED_EHR_ACCESS_MODEL_NAME {
		errs = append(errs, ValidationError{
			Model:          EHR_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid EHR access type: %s", e.EHRAccess.Type),
			Recommendation: fmt.Sprintf("Ensure ehr_access _type field is set to '%s'", VERSIONED_EHR_ACCESS_MODEL_NAME),
		})
	}
	errs = append(errs, e.EHRAccess.Validate(attrPath)...)

	// Validate compositions
	if e.Compositions.IsSet() {
		for i, compRef := range e.Compositions.Unwrap() {
			attrPath = path + fmt.Sprintf(".compositions[%d]", i)
			if compRef.Type != VERSIONED_COMPOSITION_MODEL_NAME {
				errs = append(errs, ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid composition type: %s", compRef.Type),
					Recommendation: fmt.Sprintf("Ensure compositions[%d] _type field is set to '%s'", i, VERSIONED_COMPOSITION_MODEL_NAME),
				})
			}
			errs = append(errs, compRef.Validate(attrPath)...)
		}
	}

	// Validate directory
	if e.Directory.IsSet() {
		directory := e.Directory.Unwrap()
		attrPath = path + ".directory"
		if directory.Type != VERSIONED_FOLDER_MODEL_NAME {
			errs = append(errs, ValidationError{
				Model:          EHR_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid folder type: %s", directory.Type),
				Recommendation: fmt.Sprintf("Ensure directory _type field is set to '%s'", VERSIONED_FOLDER_MODEL_NAME),
			})
		}

		errs = append(errs, directory.Validate(attrPath)...)
	}

	// Validate time_created
	attrPath = path + ".time_created"
	errs = append(errs, e.TimeCreated.Validate(attrPath)...)

	// Validate folders
	if e.Folders.IsSet() {
		for i, folderRef := range e.Folders.Unwrap() {
			attrPath = path + fmt.Sprintf(".folders[%d]", i)
			if folderRef.Type != VERSIONED_FOLDER_MODEL_NAME {
				errs = append(errs, ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid folder type: %s", folderRef.Type),
					Recommendation: fmt.Sprintf("Ensure folders[%d] _type field is set to '%s'", i, VERSIONED_FOLDER_MODEL_NAME),
				})
			}
			errs = append(errs, folderRef.Validate(attrPath)...)
		}
	}

	// Validate tags
	if e.Tags.IsSet() {
		for i, tagRef := range e.Tags.Unwrap() {
			attrPath = path + fmt.Sprintf(".tags[%d]", i)
			if tagRef.Type != VERSIONED_TAG_MODEL_NAME {
				errs = append(errs, ValidationError{
					Model:          EHR_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("invalid tag type: %s", tagRef.Type),
					Recommendation: fmt.Sprintf("Ensure tags[%d] _type field is set to '%s'", i, VERSIONED_TAG_MODEL_NAME),
				})
			}
			errs = append(errs, tagRef.Validate(attrPath)...)
		}
	}

	return errs
}

const VERSIONED_EHR_STATUS_MODEL_NAME string = "VERSIONED_EHR_STATUS"

// todo

const VERSIONED_EHR_ACCESS_MODEL_NAME string = "VERSIONED_EHR_ACCESS"

// todo

const VERSIONED_COMPOSITION_MODEL_NAME string = "VERSIONED_COMPOSITION"

// todo

const VERSIONED_FOLDER_MODEL_NAME string = "VERSIONED_FOLDER"

// todo

const VERSIONED_TAG_MODEL_NAME string = "VERSIONED_TAG"

const CONTRIBUTION_MODEL_NAME string = "CONTRIBUTION"

// todo

const HIER_OBJECT_ID_MODEL_NAME string = "HIER_OBJECT_ID"

var _ ReferenceModel = (*HIER_OBJECT_ID)(nil)

type HIER_OBJECT_ID struct {
	Type_ Optional[string] `json:"_type,omitzero"`
	Value string           `json:"value"`
}

func (h HIER_OBJECT_ID) HasModelName() bool {
	return h.Type_.IsSet()
}

func (h HIER_OBJECT_ID) Validate(path string) []ValidationError {
	var errs []ValidationError
	var attrPath string

	// Validate _type
	if h.Type_.IsSet() && h.Type_.Unwrap() != HIER_OBJECT_ID_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, ValidationError{
			Model:          HIER_OBJECT_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", HIER_OBJECT_ID_MODEL_NAME, h.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", HIER_OBJECT_ID_MODEL_NAME),
		})
	}

	// Validate UID-based identifier format: root '::' extension (extension is optional)
	attrPath = ".value"
	if h.Value == "" {
		errs = append(errs, ValidationError{
			Model:          HIER_OBJECT_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("%s value cannot be empty", HIER_OBJECT_ID_MODEL_NAME),
			Recommendation: fmt.Sprintf("Ensure %s value is set", HIER_OBJECT_ID_MODEL_NAME),
		})
	} else {
		// Split by '::' separator
		parts := strings.Split(h.Value, "::")

		// Must have 1 (root only) or 2 parts (root + extension)
		if len(parts) > 2 {
			errs = append(errs, ValidationError{
				Model:          HIER_OBJECT_ID_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("%s invalid format: too many '::'", HIER_OBJECT_ID_MODEL_NAME),
				Recommendation: fmt.Sprintf("Ensure %s value is in the format 'root::extension'", HIER_OBJECT_ID_MODEL_NAME),
			})
			return errs
		}

		// Validate root part (first part)
		root := parts[0]
		if root == "" {
			errs = append(errs, ValidationError{
				Model:          HIER_OBJECT_ID_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("%s root part cannot be empty in '%s'", HIER_OBJECT_ID_MODEL_NAME, h.Value),
				Recommendation: fmt.Sprintf("Ensure %s value has a non-empty root part", HIER_OBJECT_ID_MODEL_NAME),
			})
			return errs
		}

		// Root should be a valid UID (UUID, ISO_OID, or INTERNET_ID format)
		if err := ValidateUID(root); err != nil {
			errs = append(errs, ValidationError{
				Model:          HIER_OBJECT_ID_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("%s invalid root UID '%s': %v", HIER_OBJECT_ID_MODEL_NAME, root, err),
				Recommendation: fmt.Sprintf("Ensure %s root part is a valid UUID, ISO_OID, or INTERNET_ID", HIER_OBJECT_ID_MODEL_NAME),
			})
		}

		// If extension exists, validate it's not empty
		if len(parts) == 2 {
			extension := parts[1]
			if extension == "" {
				errs = append(errs, ValidationError{
					Model:          HIER_OBJECT_ID_MODEL_NAME,
					Path:           attrPath,
					Message:        fmt.Sprintf("%s extension cannot be empty when '::' is present", HIER_OBJECT_ID_MODEL_NAME),
					Recommendation: fmt.Sprintf("Ensure %s value has a non-empty extension part", HIER_OBJECT_ID_MODEL_NAME),
				})
			}
		}
	}

	return errs
}

const OBJECT_REF_MODEL_NAME = "OBJECT_REF"

var _ ReferenceModel = (*OBJECT_REF)(nil)

type OBJECT_REF struct {
	Type_     Optional[string] `json:"_type,omitzero"`
	Namespace string           `json:"namespace"`
	Type      string           `json:"type"`
	ID        X_OBJECT_ID      `json:"id"`
}

func (o OBJECT_REF) HasModelName() bool {
	return o.Type_.IsSet()
}

func (o OBJECT_REF) Validate(path string) []ValidationError {
	var errs []ValidationError
	var attrPath string

	// Validate _type
	if o.Type_.IsSet() && o.Type_.Unwrap() != OBJECT_REF_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, ValidationError{
			Model:          OBJECT_REF_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", OBJECT_REF_MODEL_NAME, o.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", OBJECT_REF_MODEL_NAME),
		})
	}

	// Validate namespace
	attrPath = path + ".namespace"
	if o.Namespace == "" {
		errs = append(errs, ValidationError{
			Model:          "String",
			Path:           attrPath,
			Message:        "invalid namespace: cannot be empty",
			Recommendation: "Fill in a value matching the regex standard documented in the specifications",
		})
	} else {
		if !NamespaceRegex.MatchString(o.Namespace) {
			errs = append(errs, ValidationError{
				Model:          "String",
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid namespace: %s", o.Namespace),
				Recommendation: "Fill in a value matching the regex standard documented in the specifications",
			})
		}
	}

	// Validate type
	attrPath = path + ".type"
	if o.Type == "" {
		errs = append(errs, ValidationError{
			Model:          "String",
			Path:           attrPath,
			Message:        "invalid type: cannot be empty",
			Recommendation: "Fill in a value matching the regex standard documented in the specifications",
		})
	}

	// Validate id
	attrPath = path + ".id"
	errs = append(errs, o.ID.Validate(attrPath)...)

	return errs
}

var _ ReferenceModel = (*X_OBJECT_ID)(nil)

// Abstract
type X_OBJECT_ID struct {
	Value ReferenceModel
}

func (x X_OBJECT_ID) HasModelName() bool {
	return x.Value.HasModelName()
}

func (x X_OBJECT_ID) Validate(path string) []ValidationError {
	var errs []ValidationError
	var valuePath string

	// Abstract model requires _type to be defined
	if !x.HasModelName() {
		valuePath = path + "._type"
		errs = append(errs, ValidationError{
			Model:          OBJECT_REF_MODEL_NAME,
			Path:           valuePath,
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}

	return append(errs, x.Value.Validate(path)...)
}

func (o X_OBJECT_ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

func (o *X_OBJECT_ID) UnmarshalJSON(data []byte) error {
	var extractor TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.MetaType
	switch t {
	case HIER_OBJECT_ID_MODEL_NAME:
		o.Value = new(HIER_OBJECT_ID)
	case OBJECT_VERSION_ID_MODEL_NAME:
		o.Value = new(OBJECT_VERSION_ID)
	case ARCHETYPE_ID_MODEL_NAME:
		o.Value = new(ARCHETYPE_ID)
	case TEMPLATE_ID_MODEL_NAME:
		o.Value = new(TEMPLATE_ID)
	case GENERIC_ID_MODEL_NAME:
		o.Value = new(GENERIC_ID)
	case "":
		return fmt.Errorf("missing OBJECT_ID _type field")
	default:
		return fmt.Errorf("OBJECT_ID unexpected _type %s", t)
	}

	return json.Unmarshal(data, o.Value)
}

const DV_DATE_TIME_MODEL_NAME string = "DV_DATE_TIME"

var _ ReferenceModel = (*DV_DATE_TIME)(nil)

type DV_DATE_TIME struct {
	Type_                Optional[string]            `json:"_type"`
	NormalStatus         Optional[CODE_PHRASE]       `json:"normal_status"`
	NormalRange          Optional[DV_INTERVAL]       `json:"normal_range"`
	OtherReferenceRanges Optional[[]REFERENCE_RANGE] `json:"other_reference_ranges"`
	MagnitudeStatus      Optional[string]            `json:"magnitude_status"`
	Accuracy             Optional[DV_DURATION]       `json:"accuracy"`
	Value                string                      `json:"value"`
}

// HasModelName implements ReferenceModel.
func (d DV_DATE_TIME) HasModelName() bool {
	return d.Type_.IsSet()
}

// Validate implements ReferenceModel.
func (d DV_DATE_TIME) Validate(path string) []ValidationError {
	var errs []ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.IsSet() && d.Type_.Unwrap() != DV_DATE_TIME_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, ValidationError{
			Model:          DV_DATE_TIME_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_DATE_TIME_MODEL_NAME, d.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_DATE_TIME_MODEL_NAME),
		})
	}

	// Validate normal_status
	if d.NormalStatus.IsSet() {
		attrPath = path + ".normal_status"
		errs = append(errs, d.NormalStatus.Unwrap().Validate(attrPath)...)
	}

	// Validate normal_range
	if d.NormalRange.IsSet() {
		attrPath = path + ".normal_range"
		errs = append(errs, d.NormalRange.Unwrap().Validate(attrPath)...)
	}

	// Validate other_reference_ranges
	if d.OtherReferenceRanges.IsSet() {
		attrPath = path + ".other_reference_ranges"
		for i, v := range d.OtherReferenceRanges.Unwrap() {
			errs = append(errs, v.Validate(fmt.Sprintf("%s[%d]", attrPath, i))...)
		}
	}

	// Validate magnitude_status
	if d.MagnitudeStatus.IsSet() {
		attrPath = path + ".magnitude_status"
		if !slices.Contains([]string{"<", ">", "<=", ">=", "=", "~"}, d.MagnitudeStatus.Unwrap()) {
			errs = append(errs, ValidationError{
				Model:          DV_DATE_TIME_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s magnitude_status field: %s", DV_DATE_TIME_MODEL_NAME, d.MagnitudeStatus.Unwrap()),
				Recommendation: "Ensure magnitude_status field is one of '<', '>', '<=', '>=', '=', '~'",
			})
		}
	}

	// Validate accuracy
	if d.Accuracy.IsSet() {
		attrPath = path + ".accuracy"
		errs = append(errs, d.Accuracy.Unwrap().Validate(attrPath)...)
	}

	// Validate value
	attrPath = path + ".value"
	if d.Value == "" {
		errs = append(errs, ValidationError{
			Model:   DV_DATE_TIME_MODEL_NAME,
			Path:    attrPath,
			Message: fmt.Sprintf("%s value field is required", DV_DATE_TIME_MODEL_NAME),
		})
	} else if !strings.HasSuffix(d.Value, "Z") {
		errs = append(errs, ValidationError{
			Model:          DV_DATE_TIME_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s value field: %s", DV_DATE_TIME_MODEL_NAME, d.Value),
			Recommendation: "Ensure value field is of format YYYY-MM-DDTHH:MM:SSZ",
		})
	} else {
		if _, err := time.Parse(time.RFC3339Nano, d.Value); err != nil {
			errs = append(errs, ValidationError{
				Model:          DV_DATE_TIME_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s value field: %s", DV_DATE_TIME_MODEL_NAME, d.Value),
				Recommendation: "Ensure value field is of format YYYY-MM-DDTHH:MM:SSZ",
			})
		}
	}

	return errs
}

const DV_INTERVAL_MODEL_NAME string = "DV_INTERVAL"

var _ ReferenceModel = (*DV_INTERVAL)(nil)

type DV_INTERVAL struct {
	Type_          Optional[string] `json:"_type,omitzero"`
	Lower          any              `json:"lower"`
	Upper          any              `json:"upper"`
	LowerUnbounded bool             `json:"lower_unbounded"`
	UpperUnbounded bool             `json:"upper_unbounded"`
	LowerIncluded  bool             `json:"lower_included"`
	UpperIncluded  bool             `json:"upper_included"`
}

func (d DV_INTERVAL) HasModelName() bool {
	return d.Type_.IsSet()
}

func (d DV_INTERVAL) Validate(path string) []ValidationError {
	var errors []ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.IsSet() && d.Type_.Unwrap() != DV_INTERVAL_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, ValidationError{
			Model:          DV_INTERVAL_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_INTERVAL_MODEL_NAME, d.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_INTERVAL_MODEL_NAME),
		})
	}

	return errors
}

const DV_DURATION_MODEL_NAME string = "DV_DURATION"

var _ ReferenceModel = (*DV_DURATION)(nil)

type DV_DURATION struct {
	Type_                Optional[string]            `json:"_type,omitzero"`
	NormalStatus         Optional[CODE_PHRASE]       `json:"normal_status,omitzero"`
	NormalRange          Optional[DV_INTERVAL]       `json:"normal_range,omitzero"`
	OtherReferenceRanges Optional[[]REFERENCE_RANGE] `json:"other_reference_ranges,omitzero"`
	MagnitudeStatus      Optional[string]            `json:"magnitude_status,omitzero"`
	AccuracyIsPercent    Optional[bool]              `json:"accuracy_is_percent,omitzero"`
	Accuracy             Optional[float64]           `json:"accuracy,omitzero"`
	Value                string                      `json:"value"`
}

func (d DV_DURATION) HasModelName() bool {
	return d.Type_.IsSet()
}

func (d DV_DURATION) Validate(path string) []ValidationError {
	var errors []ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.IsSet() && d.Type_.Unwrap() != DV_DURATION_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, ValidationError{
			Model:          DV_DURATION_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_DURATION_MODEL_NAME, d.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_DURATION_MODEL_NAME),
		})
	}

	// Validate normal_status
	if d.NormalStatus.IsSet() {
		attrPath = path + ".normal_status"
		errors = append(errors, d.NormalStatus.Unwrap().Validate(attrPath)...)
	}

	// Validate normal_range
	if d.NormalRange.IsSet() {
		attrPath = path + ".normal_range"
		errors = append(errors, d.NormalRange.Unwrap().Validate(attrPath)...)
	}

	// Validate other_reference_ranges
	if d.OtherReferenceRanges.IsSet() {
		attrPath = path + ".other_reference_ranges"
		for i, v := range d.OtherReferenceRanges.Unwrap() {
			errors = append(errors, v.Validate(fmt.Sprintf("%s[%d]", attrPath, i))...)
		}
	}

	// Validate magnitude_status
	if d.MagnitudeStatus.IsSet() {
		attrPath = path + ".magnitude_status"
		validValues := []string{"<", ">", "<=", ">=", "=", "~"}
		value := d.MagnitudeStatus.Unwrap()
		isValid := slices.Contains(validValues, d.MagnitudeStatus.Unwrap())
		if !isValid {
			errors = append(errors, ValidationError{
				Model:          DV_DURATION_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s magnitude_status field: %s", DV_DURATION_MODEL_NAME, value),
				Recommendation: "Ensure magnitude_status field is one of '<', '>', '<=', '>=', '=', '~'",
			})
		}
	}

	// Validate accuracy
	if d.Accuracy.IsSet() {
		attrPath = path + ".accuracy"
		value := d.Accuracy.Unwrap()
		if value < 0 {
			errors = append(errors, ValidationError{
				Model:          DV_DURATION_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid %s accuracy field: %f", DV_DURATION_MODEL_NAME, value),
				Recommendation: "Ensure accuracy field is a non-negative number",
			})
		}
	}

	return errors
}

const CODE_PHRASE_MODEL_NAME string = "CODE_PHRASE"

var _ ReferenceModel = (*CODE_PHRASE)(nil)

type CODE_PHRASE struct {
	Type_         Optional[string] `json:"_type,omitzero"`
	TerminologyId TERMINOLOGY_ID   `json:"terminology_id"`
	CodeString    string           `json:"code_string"`
	PreferredTerm Optional[string] `json:"preferred_term,omitzero"`
}

func (c CODE_PHRASE) HasModelName() bool {
	return c.Type_.IsSet()
}

func (c CODE_PHRASE) Validate(path string) []ValidationError {
	var errors []ValidationError

	// Validate _type
	if c.Type_.IsSet() && c.Type_.Unwrap() != CODE_PHRASE_MODEL_NAME {
		errors = append(errors, ValidationError{
			Model:          CODE_PHRASE_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", CODE_PHRASE_MODEL_NAME, c.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", CODE_PHRASE_MODEL_NAME),
		})
	}

	// Validate terminology_id
	attrPath := path + ".terminology_id"
	errors = append(errors, c.TerminologyId.Validate(attrPath)...)

	// Validate code_string
	attrPath = path + ".code_string"
	if c.CodeString == "" {
		errors = append(errors, ValidationError{
			Model:          CODE_PHRASE_MODEL_NAME,
			Path:           attrPath,
			Message:        "code_string field is required",
			Recommendation: "Ensure code_string field is not empty",
		})
	}

	return errors
}

const TERMINOLOGY_ID_MODEL_NAME string = "TERMINOLOGY_ID"

var _ ReferenceModel = (*TERMINOLOGY_ID)(nil)

type TERMINOLOGY_ID struct {
	Type_ Optional[string] `json:"_type,omitzero"`
	Value string           `json:"value"`
}

func (t TERMINOLOGY_ID) HasModelName() bool {
	return t.Type_.IsSet()
}

func (t TERMINOLOGY_ID) Validate(path string) []ValidationError {
	var errors []ValidationError

	// Validate _type
	if t.Type_.IsSet() && t.Type_.Unwrap() != TERMINOLOGY_ID_MODEL_NAME {
		errors = append(errors, ValidationError{
			Model:          TERMINOLOGY_ID_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", TERMINOLOGY_ID_MODEL_NAME, t.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", TERMINOLOGY_ID_MODEL_NAME),
		})
	}

	// Validate value
	attrPath := path + ".value"
	if t.Value == "" {
		errors = append(errors, ValidationError{
			Model:          TERMINOLOGY_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field is required",
			Recommendation: "Ensure value field is not empty",
		})
	}

	return errors
}

const REFERENCE_RANGE_MODEL_NAME string = "REFERENCE_RANGE"

var _ ReferenceModel = (*REFERENCE_RANGE)(nil)

type REFERENCE_RANGE struct {
	Type_   Optional[string] `json:"_type,omitzero"`
	Meaning X_DV_TEXT        `json:"meaning"`
	Range   DV_INTERVAL      `json:"range"`
}

func (r REFERENCE_RANGE) HasModelName() bool {
	return r.Type_.IsSet()
}

func (r REFERENCE_RANGE) Validate(path string) []ValidationError {
	var errors []ValidationError

	// Validate _type
	if r.Type_.IsSet() && r.Type_.Unwrap() != REFERENCE_RANGE_MODEL_NAME {
		errors = append(errors, ValidationError{
			Model:          REFERENCE_RANGE_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", REFERENCE_RANGE_MODEL_NAME, r.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", REFERENCE_RANGE_MODEL_NAME),
		})
	}

	// Validate meaning
	attrPath := path + ".meaning"
	errors = append(errors, r.Meaning.Validate(attrPath)...)

	// Validate range
	attrPath = path + ".range"
	errors = append(errors, r.Range.Validate(attrPath)...)

	return errors
}

var _ ReferenceModel = (*X_DV_TEXT)(nil)

type X_DV_TEXT struct {
	Value ReferenceModel
}

func (x X_DV_TEXT) HasModelName() bool {
	return x.Value.HasModelName()
}

func (x X_DV_TEXT) Validate(path string) []ValidationError {
	var errs []ValidationError

	// Abstract model requires _type to be defined
	if !x.HasModelName() {
		errs = append(errs, ValidationError{
			Model:          OBJECT_REF_MODEL_NAME,
			Path:           "._type",
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}

	return append(errs, x.Value.Validate(path)...)
}

func (a X_DV_TEXT) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Value)
}

func (a *X_DV_TEXT) UnmarshalJSON(data []byte) error {
	var extractor TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.MetaType
	switch t {
	case DV_TEXT_MODEL_NAME:
		a.Value = new(DV_TEXT)
	case DV_CODED_TEXT_MODEL_NAME:
		a.Value = new(DV_CODED_TEXT)
	case "":
		return fmt.Errorf("missing DV_TEXT _type field")
	default:
		return fmt.Errorf("DV_TEXT unexpected _type %s", t)
	}

	return json.Unmarshal(data, a.Value)
}

const DV_TEXT_MODEL_NAME string = "DV_TEXT"

var _ ReferenceModel = (*DV_TEXT)(nil)

type DV_TEXT struct {
	Type_      Optional[string]         `json:"_type,omitzero"`
	Value      string                   `json:"value"`
	Formatting Optional[string]         `json:"formatting,omitzero"`
	Mappings   Optional[[]TERM_MAPPING] `json:"mappings,omitzero"`
	Language   Optional[CODE_PHRASE]    `json:"language,omitzero"`
	Encoding   Optional[CODE_PHRASE]    `json:"encoding,omitzero"`
}

func (d DV_TEXT) HasModelName() bool {
	return d.Type_.IsSet()
}

func (d DV_TEXT) Validate(path string) []ValidationError {
	var errors []ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.IsSet() && d.Type_.Unwrap() != DV_TEXT_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_TEXT_MODEL_NAME, d.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_TEXT_MODEL_NAME),
		})
	}

	// Validate formatting
	if d.Formatting.IsSet() {
		attrPath = path + ".formatting"
		validFormats := []string{"plain", "plain_no_newlines", "markdown"}
		if !slices.Contains(validFormats, d.Formatting.Unwrap()) {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid formatting field: %s", d.Formatting.Unwrap()),
				Recommendation: "Ensure formatting field is one of 'plain', 'plain_no_newlines', 'markdown'",
			})
		}
	}

	// Validate value
	attrPath = path + ".value"
	if d.Value == "" {
		errors = append(errors, ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field is required",
			Recommendation: "Ensure value field is not empty",
		})
	} else if len(d.Value) > 10000 {
		errors = append(errors, ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field exceeds maximum length of 10000 characters",
			Recommendation: "Ensure value field does not exceed 10000 characters",
		})
	}

	if d.Formatting.IsSet() && d.Formatting.Unwrap() == "plain_no_newlines" {
		if strings.ContainsAny(d.Value, "\n\r") {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        "value field contains newlines but formatting is 'plain_no_newlines'",
				Recommendation: "Ensure value field does not contain newlines when formatting is 'plain_no_newlines'",
			})
		}
	}

	// Validate mappings
	if d.Mappings.IsSet() {
		attrPath = path + ".mappings"
		for i, v := range d.Mappings.Unwrap() {
			errors = append(errors, v.Validate(fmt.Sprintf("%s[%d]", attrPath, i))...)
		}
	}

	// Validate language
	if d.Language.IsSet() {
		attrPath = path + ".language"
		v := d.Language.Unwrap()
		if !terminology.IsValidLanguageTerminologyID(v.TerminologyId.Value) {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language field: %s", d.Language.Unwrap()),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}

		if !terminology.IsValidLanguageCode(v.CodeString) {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language field: %s", d.Language.Unwrap()),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}
		errors = append(errors, d.Language.Unwrap().Validate(attrPath)...)
	}

	// Validate encoding
	if d.Encoding.IsSet() {
		attrPath = path + ".encoding"
		v := d.Encoding.Unwrap()
		if !terminology.IsValidCharsetTerminologyID(v.TerminologyId.Value) {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding field: %s", d.Encoding.Unwrap()),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}

		if !terminology.IsValidCharset(v.CodeString) {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding field: %s", d.Encoding.Unwrap()),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}
		errors = append(errors, d.Encoding.Unwrap().Validate(attrPath)...)
	}

	return errors
}

const DV_CODED_TEXT_MODEL_NAME string = "DV_CODED_TEXT"

var _ ReferenceModel = (*DV_CODED_TEXT)(nil)

type DV_CODED_TEXT struct {
	Type_        Optional[string]         `json:"_type,omitzero"`
	Value        string                   `json:"value"`
	Hyperlink    Optional[DV_URI]         `json:"hyperlink,omitzero"`
	Formatting   Optional[string]         `json:"formatting,omitzero"`
	Mappings     Optional[[]TERM_MAPPING] `json:"mappings,omitzero"`
	Language     Optional[CODE_PHRASE]    `json:"language,omitzero"`
	Encoding     Optional[CODE_PHRASE]    `json:"encoding,omitzero"`
	DefiningCode CODE_PHRASE              `json:"defining_code"`
}

func (d DV_CODED_TEXT) HasModelName() bool {
	return d.Type_.IsSet()
}

func (d DV_CODED_TEXT) Validate(path string) []ValidationError {
	var errors []ValidationError

	// Validate _type
	if d.Type_.IsSet() && d.Type_.Unwrap() != DV_CODED_TEXT_MODEL_NAME {
		errors = append(errors, ValidationError{
			Model:          DV_CODED_TEXT_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_CODED_TEXT_MODEL_NAME, d.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_CODED_TEXT_MODEL_NAME),
		})
	}

	// Validate defining_code
	attrPath := path + ".defining_code"
	errors = append(errors, d.DefiningCode.Validate(attrPath)...)

	// Validate formatting
	if d.Formatting.IsSet() {
		attrPath = path + ".formatting"
		validFormats := []string{"plain", "plain_no_newlines", "markdown"}
		if !slices.Contains(validFormats, d.Formatting.Unwrap()) {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid formatting field: %s", d.Formatting.Unwrap()),
				Recommendation: "Ensure formatting field is one of 'plain', 'plain_no_newlines', 'markdown'",
			})
		}
	}

	// Validate value
	attrPath = path + ".value"
	if d.Value == "" {
		errors = append(errors, ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field is required",
			Recommendation: "Ensure value field is not empty",
		})
	} else if len(d.Value) > 10000 {
		errors = append(errors, ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field exceeds maximum length of 10000 characters",
			Recommendation: "Ensure value field does not exceed 10000 characters",
		})
	}

	if d.Formatting.IsSet() && d.Formatting.Unwrap() == "plain_no_newlines" {
		if strings.ContainsAny(d.Value, "\n\r") {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        "value field contains newlines but formatting is 'plain_no_newlines'",
				Recommendation: "Ensure value field does not contain newlines when formatting is 'plain_no_newlines'",
			})
		}
	}

	// Validate mappings
	if d.Mappings.IsSet() {
		attrPath = path + ".mappings"
		for i, v := range d.Mappings.Unwrap() {
			errors = append(errors, v.Validate(fmt.Sprintf("%s[%d]", attrPath, i))...)
		}
	}

	// Validate language
	if d.Language.IsSet() {
		attrPath = path + ".language"
		v := d.Language.Unwrap()
		if !terminology.IsValidLanguageTerminologyID(v.TerminologyId.Value) {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language field: %s", d.Language.Unwrap()),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}

		if !terminology.IsValidLanguageCode(v.CodeString) {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid language field: %s", d.Language.Unwrap()),
				Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
			})
		}
		errors = append(errors, d.Language.Unwrap().Validate(attrPath)...)
	}

	// Validate encoding
	if d.Encoding.IsSet() {
		attrPath = path + ".encoding"
		v := d.Encoding.Unwrap()
		if !terminology.IsValidCharsetTerminologyID(v.TerminologyId.Value) {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding field: %s", d.Encoding.Unwrap()),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}

		if !terminology.IsValidCharset(v.CodeString) {
			errors = append(errors, ValidationError{
				Model:          DV_TEXT_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid encoding field: %s", d.Encoding.Unwrap()),
				Recommendation: "Ensure encoding field is a known IANA character set",
			})
		}
		errors = append(errors, d.Encoding.Unwrap().Validate(attrPath)...)
	}

	return errors
}

const DV_URI_MODEL_NAME string = "DV_URI"

var _ ReferenceModel = (*DV_URI)(nil)

type DV_URI struct {
	Type_ Optional[string] `json:"_type,omitzero"`
	Value string           `json:"value"`
}

func (d DV_URI) HasModelName() bool {
	return d.Type_.IsSet()
}

func (d DV_URI) Validate(path string) []ValidationError {
	var errors []ValidationError

	// Validate _type
	if d.Type_.IsSet() && d.Type_.Unwrap() != DV_URI_MODEL_NAME {
		errors = append(errors, ValidationError{
			Model:          DV_URI_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", DV_URI_MODEL_NAME, d.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", DV_URI_MODEL_NAME),
		})
	}

	// Validate value
	attrPath := path + ".value"
	if d.Value == "" {
		errors = append(errors, ValidationError{
			Model:          DV_URI_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	} else if !URIRegex.MatchString(d.Value) {
		errors = append(errors, ValidationError{
			Model:          DV_URI_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid URI value: %s", d.Value),
			Recommendation: "Ensure value field is a valid URI according to RFC 3986",
		})
	}

	return errors
}

const TERM_MAPPING_MODEL_NAME string = "TERM_MAPPING"

var _ ReferenceModel = (*TERM_MAPPING)(nil)

type TERM_MAPPING struct {
	Type_   Optional[string]        `json:"_type,omitzero"`
	Match   byte                    `json:"match"`
	Purpose Optional[DV_CODED_TEXT] `json:"purpose,omitzero"`
	Target  CODE_PHRASE             `json:"target"`
}

func (t TERM_MAPPING) HasModelName() bool {
	return t.Type_.IsSet()
}

func (t TERM_MAPPING) Validate(path string) []ValidationError {
	var errors []ValidationError

	// Validate _type
	if t.Type_.IsSet() && t.Type_.Unwrap() != TERM_MAPPING_MODEL_NAME {
		errors = append(errors, ValidationError{
			Model:          TERM_MAPPING_MODEL_NAME,
			Path:           "._type",
			Message:        fmt.Sprintf("invalid %s _type field: %s", TERM_MAPPING_MODEL_NAME, t.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", TERM_MAPPING_MODEL_NAME),
		})
	}

	// Validate purpose
	if t.Purpose.IsSet() {
		attrPath := path + ".purpose"
		errors = append(errors, t.Purpose.Unwrap().Validate(attrPath)...)
	}

	// Validate target
	attrPath := path + ".target"
	errors = append(errors, t.Target.Validate(attrPath)...)

	// Validate match
	validMatches := []byte{'=', '>', '<', '?'}
	if !slices.Contains(validMatches, t.Match) {
		errors = append(errors, ValidationError{
			Model:          TERM_MAPPING_MODEL_NAME,
			Path:           path + ".match",
			Message:        fmt.Sprintf("invalid match value: %c", t.Match),
			Recommendation: "Ensure match field is one of '=', '>', '<', '?'",
		})
	}

	return errors
}

const OBJECT_VERSION_ID_MODEL_NAME string = "OBJECT_VERSION_ID"

var _ ReferenceModel = (*OBJECT_VERSION_ID)(nil)

type OBJECT_VERSION_ID struct {
	Type_ Optional[string] `json:"_type,omitzero"`
	Value string           `json:"value"`
}

func (o OBJECT_VERSION_ID) HasModelName() bool {
	return o.Type_.IsSet()
}

func (o OBJECT_VERSION_ID) Validate(path string) []ValidationError {
	var errors []ValidationError
	var attrPath string

	// Validate _type
	if o.Type_.IsSet() && o.Type_.Unwrap() != OBJECT_VERSION_ID_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, ValidationError{
			Model:          OBJECT_VERSION_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", OBJECT_VERSION_ID_MODEL_NAME, o.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", OBJECT_VERSION_ID_MODEL_NAME),
		})
	}

	// Validate value
	attrPath = path + ".value"
	if o.Value == "" {
		errors = append(errors, ValidationError{
			Model:          OBJECT_VERSION_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	} else {
		// lexical form: object_id '::' creating_system_id '::' version_tree_id.
		parts := strings.Split(o.Value, "::")
		if len(parts) != 3 {
			errors = append(errors, ValidationError{
				Model:          OBJECT_VERSION_ID_MODEL_NAME,
				Path:           attrPath,
				Message:        fmt.Sprintf("invalid value format: %s", o.Value),
				Recommendation: "Ensure value field follows the lexical form: object_id '::' creating_system_id '::' version_tree_id",
			})
		}

		// First part UID
		uid := parts[0]
		if err := ValidateUID(uid); err != nil {
			errors = append(errors, ValidationError{
				Model:          OBJECT_VERSION_ID_MODEL_NAME,
				Path:           attrPath + ".object_id",
				Message:        fmt.Sprintf("invalid object_id format: %s", uid),
				Recommendation: "Ensure object_id follows a valid UID format",
			})
		}

		// Second part creating_system_id
		creatingSystemID := parts[1]
		if creatingSystemID == "" {
			errors = append(errors, ValidationError{
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
			errors = append(errors, ValidationError{
				Model:          OBJECT_VERSION_ID_MODEL_NAME,
				Path:           attrPath + ".version_tree_id",
				Message:        "version_tree_id cannot be empty",
				Recommendation: "Ensure version_tree_id is not empty",
			})
		} else if !VersionTreeIDRegex.MatchString(versionTreeID) {
			errors = append(errors, ValidationError{
				Model:          OBJECT_VERSION_ID_MODEL_NAME,
				Path:           attrPath + ".version_tree_id",
				Message:        fmt.Sprintf("invalid version_tree_id format: %s", versionTreeID),
				Recommendation: "Ensure version_tree_id follows the lexical form: trunk_version [ '.' branch_number '.' branch_version ]",
			})
		}
	}

	return errors
}

const ARCHETYPE_ID_MODEL_NAME string = "ARCHETYPE_ID"

var _ ReferenceModel = (*ARCHETYPE_ID)(nil)

type ARCHETYPE_ID struct {
	Type_ Optional[string] `json:"_type,omitzero"`
	Value string           `json:"value"`
}

func (a ARCHETYPE_ID) HasModelName() bool {
	return a.Type_.IsSet()
}

func (a ARCHETYPE_ID) Validate(path string) []ValidationError {
	var errors []ValidationError
	var attrPath string

	// Validate _type
	if a.Type_.IsSet() && a.Type_.Unwrap() != ARCHETYPE_ID_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, ValidationError{
			Model:          ARCHETYPE_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", ARCHETYPE_ID_MODEL_NAME, a.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", ARCHETYPE_ID_MODEL_NAME),
		})
	}

	// Validate value
	attrPath = path + ".value"
	if a.Value == "" {
		errors = append(errors, ValidationError{
			Model:          ARCHETYPE_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	} else if !ArchetypeIDRegex.MatchString(a.Value) {
		errors = append(errors, ValidationError{
			Model:          ARCHETYPE_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid value format: %s", a.Value),
			Recommendation: "Ensure value field follows the lexical form: rm_originator '-' rm_name '-' rm_entity '.' concept_name { '-' specialisation }* '.v' number.",
		})
	}

	return errors
}

const TEMPLATE_ID_MODEL_NAME string = "TEMPLATE_ID"

var _ ReferenceModel = (*TEMPLATE_ID)(nil)

type TEMPLATE_ID struct {
	Type_ Optional[string] `json:"_type,omitzero"`
	Value string           `json:"value"`
}

func (t TEMPLATE_ID) HasModelName() bool {
	return t.Type_.IsSet()
}

func (t TEMPLATE_ID) Validate(path string) []ValidationError {
	var errors []ValidationError
	var attrPath string

	// Validate _type
	if t.Type_.IsSet() && t.Type_.Unwrap() != TEMPLATE_ID_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, ValidationError{
			Model:          TEMPLATE_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", TEMPLATE_ID_MODEL_NAME, t.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", TEMPLATE_ID_MODEL_NAME),
		})
	}

	// Validate value
	attrPath = path + ".value"
	if t.Value == "" {
		errors = append(errors, ValidationError{
			Model:          TEMPLATE_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	}

	return errors
}

const GENERIC_ID_MODEL_NAME string = "GENERIC_ID"

var _ ReferenceModel = (*GENERIC_ID)(nil)

type GENERIC_ID struct {
	Type_  Optional[string] `json:"_type,omitzero"`
	Value  string           `json:"value"`
	Scheme string           `json:"scheme"`
}

func (g GENERIC_ID) HasModelName() bool {
	return g.Type_.IsSet()
}

func (g GENERIC_ID) Validate(path string) []ValidationError {
	var errors []ValidationError
	var attrPath string

	// Validate _type
	if g.Type_.IsSet() && g.Type_.Unwrap() != GENERIC_ID_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, ValidationError{
			Model:          GENERIC_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", GENERIC_ID_MODEL_NAME, g.Type_.Unwrap()),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", GENERIC_ID_MODEL_NAME),
		})
	}

	// Validate value
	attrPath = path + ".value"
	if g.Value == "" {
		errors = append(errors, ValidationError{
			Model:          GENERIC_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "value field cannot be empty",
			Recommendation: "Ensure value field is not empty",
		})
	}

	// Validate scheme
	attrPath = path + ".scheme"
	if g.Scheme == "" {
		errors = append(errors, ValidationError{
			Model:          GENERIC_ID_MODEL_NAME,
			Path:           attrPath,
			Message:        "scheme field cannot be empty",
			Recommendation: "Ensure scheme field is not empty",
		})
	}

	return errors
}
