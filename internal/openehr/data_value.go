package openehr

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DATA_VALUE_MODEL_NAME string = "DATA_VALUE"

// Abstract
type DATA_VALUE struct {
	Type_ util.Optional[string] `json:"_type,omitzero"`
}

func (d DATA_VALUE) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot marshal abstract DATA_VALUE type")
}

func (d *DATA_VALUE) UnmarshalJSON(data []byte) error {
	return fmt.Errorf("cannot unmarshal abstract DATA_VALUE type")
}

// ======== Union of DATA_VALUE subtypes ========

type DataValueModel interface {
	isDataValueModel()
	HasModelName() bool
	SetModelName()
	Validate(path string) util.ValidateError
}

type X_DATA_VALUE struct {
	Value DataValueModel
}

func (x *X_DATA_VALUE) SetModelName() {
	x.Value.SetModelName()
}

func (x *X_DATA_VALUE) Validate(path string) util.ValidateError {
	if x.Value == nil {
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          DATA_VALUE_MODEL_NAME,
					Path:           path,
					Message:        "value is not known DATA_VALUE subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}

	var validateErr util.ValidateError
	var attrPath string

	// Abstract model requires _type to be defined
	if !x.Value.HasModelName() {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DATA_VALUE_MODEL_NAME,
			Path:           attrPath,
			Message:        "empty _type field",
			Recommendation: "Ensure _type field is defined",
		})
	}
	validateErr.Errs = append(validateErr.Errs, x.Value.Validate(path).Errs...)

	return validateErr
}

func (x X_DATA_VALUE) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Value)
}

func (x *X_DATA_VALUE) UnmarshalJSON(data []byte) error {
	var extractor util.TypeExtractor
	if err := json.Unmarshal(data, &extractor); err != nil {
		return err
	}

	t := extractor.Type_
	switch t {
	case DV_BOOLEAN_MODEL_NAME:
		x.Value = new(DV_BOOLEAN)
	case DV_STATE_MODEL_NAME:
		x.Value = new(DV_STATE)
	case DV_IDENTIFIER_MODEL_NAME:
		x.Value = new(DV_IDENTIFIER)
	case DV_TEXT_MODEL_NAME:
		x.Value = new(DV_TEXT)
	case DV_CODED_TEXT_MODEL_NAME:
		x.Value = new(DV_CODED_TEXT)
	case DV_PARAGRAPH_MODEL_NAME:
		x.Value = new(DV_PARAGRAPH)
	case DV_INTERVAL_MODEL_NAME:
		x.Value = new(DV_INTERVAL)
	case DV_ORDINAL_MODEL_NAME:
		x.Value = new(DV_ORDINAL)
	case DV_SCALE_MODEL_NAME:
		x.Value = new(DV_SCALE)
	case DV_QUANTITY_MODEL_NAME:
		x.Value = new(DV_QUANTITY)
	case DV_COUNT_MODEL_NAME:
		x.Value = new(DV_COUNT)
	case DV_PROPORTION_MODEL_NAME:
		x.Value = new(DV_PROPORTION)
	case DV_DATE_MODEL_NAME:
		x.Value = new(DV_DATE)
	case DV_TIME_MODEL_NAME:
		x.Value = new(DV_TIME)
	case DV_DATE_TIME_MODEL_NAME:
		x.Value = new(DV_DATE_TIME)
	case DV_DURATION_MODEL_NAME:
		x.Value = new(DV_DURATION)
	case DV_PERIODIC_TIME_SPECIFICATION_MODEL_NAME:
		x.Value = new(DV_PERIODIC_TIME_SPECIFICATION)
	case DV_GENERAL_TIME_SPECIFICATION_MODEL_NAME:
		x.Value = new(DV_GENERAL_TIME_SPECIFICATION)
	case DV_MULTIMEDIA_MODEL_NAME:
		x.Value = new(DV_MULTIMEDIA)
	case DV_PARSABLE_MODEL_NAME:
		x.Value = new(DV_PARSABLE)
	case DV_URI_MODEL_NAME:
		x.Value = new(DV_URI)
	case DV_EHR_URI_MODEL_NAME:
		x.Value = new(DV_EHR_URI)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          DATA_VALUE_MODEL_NAME,
					Path:           "$.**._type",
					Message:        fmt.Sprintf("unexpected DATA_VALUE _type %s", t),
					Recommendation: "Ensure _type field is one of the known DATA_VALUE subtypes",
				},
			},
		}
	}

	return json.Unmarshal(data, x.Value)
}
