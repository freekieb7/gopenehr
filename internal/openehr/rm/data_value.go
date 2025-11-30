package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DATA_VALUE_TYPE string = "DATA_VALUE"

type DataValueKind int

const (
	DataValueKind_Unknown DataValueKind = iota
	DataValueKind_DV_BOOLEAN
	DataValueKind_DV_CODED_TEXT
	DataValueKind_DV_IDENTIFIER
	DataValueKind_DV_INTERVAL
	DataValueKind_DV_ORDINAL
	DataValueKind_DV_PARAGRAPH
	DataValueKind_DV_PROPORTION
	DataValueKind_DV_QUANTITY
	DataValueKind_DV_SCALE
	DataValueKind_DV_STATE
	DataValueKind_DV_TEXT
	DataValueKind_DV_TIME
	DataValueKind_DV_DATE
	DataValueKind_DV_DATE_TIME
	DataValueKind_DV_DURATION
	DataValueKind_DV_PERIODIC_TIME_SPECIFICATION
	DataValueKind_DV_GENERAL_TIME_SPECIFICATION
	DataValueKind_DV_MULTIMEDIA
	DataValueKind_DV_PARSABLE
	DataValueKind_DV_URI
	DataValueKind_DV_EHR_URI
	DataValueKind_DV_COUNT
)

type DataValueUnion struct {
	Kind  DataValueKind
	Value any
}

func (d *DataValueUnion) SetModelName() {
	switch d.Kind {
	case DataValueKind_DV_BOOLEAN:
		d.Value.(*DV_BOOLEAN).SetModelName()
	case DataValueKind_DV_CODED_TEXT:
		d.Value.(*DV_CODED_TEXT).SetModelName()
	case DataValueKind_DV_IDENTIFIER:
		d.Value.(*DV_IDENTIFIER).SetModelName()
	case DataValueKind_DV_INTERVAL:
		d.Value.(*DV_INTERVAL).SetModelName()
	case DataValueKind_DV_ORDINAL:
		d.Value.(*DV_ORDINAL).SetModelName()
	case DataValueKind_DV_PARAGRAPH:
		d.Value.(*DV_PARAGRAPH).SetModelName()
	case DataValueKind_DV_PROPORTION:
		d.Value.(*DV_PROPORTION).SetModelName()
	case DataValueKind_DV_QUANTITY:
		d.Value.(*DV_QUANTITY).SetModelName()
	case DataValueKind_DV_SCALE:
		d.Value.(*DV_SCALE).SetModelName()
	case DataValueKind_DV_STATE:
		d.Value.(*DV_STATE).SetModelName()
	case DataValueKind_DV_TEXT:
		d.Value.(*DV_TEXT).SetModelName()
	case DataValueKind_DV_TIME:
		d.Value.(*DV_TIME).SetModelName()
	case DataValueKind_DV_DATE:
		d.Value.(*DV_DATE).SetModelName()
	case DataValueKind_DV_DATE_TIME:
		d.Value.(*DV_DATE_TIME).SetModelName()
	case DataValueKind_DV_DURATION:
		d.Value.(*DV_DURATION).SetModelName()
	case DataValueKind_DV_PERIODIC_TIME_SPECIFICATION:
		d.Value.(*DV_PERIODIC_TIME_SPECIFICATION).SetModelName()
	case DataValueKind_DV_GENERAL_TIME_SPECIFICATION:
		d.Value.(*DV_GENERAL_TIME_SPECIFICATION).SetModelName()
	case DataValueKind_DV_MULTIMEDIA:
		d.Value.(*DV_MULTIMEDIA).SetModelName()
	case DataValueKind_DV_PARSABLE:
		d.Value.(*DV_PARSABLE).SetModelName()
	case DataValueKind_DV_URI:
		d.Value.(*DV_URI).SetModelName()
	case DataValueKind_DV_EHR_URI:
		d.Value.(*DV_EHR_URI).SetModelName()
	case DataValueKind_DV_COUNT:
		d.Value.(*DV_COUNT).SetModelName()
	}
}

func (d *DataValueUnion) Validate(path string) util.ValidateError {
	switch d.Kind {
	case DataValueKind_DV_BOOLEAN:
		return d.Value.(*DV_BOOLEAN).Validate(path)
	case DataValueKind_DV_CODED_TEXT:
		return d.Value.(*DV_CODED_TEXT).Validate(path)
	case DataValueKind_DV_IDENTIFIER:
		return d.Value.(*DV_IDENTIFIER).Validate(path)
	case DataValueKind_DV_INTERVAL:
		return d.Value.(*DV_INTERVAL).Validate(path)
	case DataValueKind_DV_ORDINAL:
		return d.Value.(*DV_ORDINAL).Validate(path)
	case DataValueKind_DV_PARAGRAPH:
		return d.Value.(*DV_PARAGRAPH).Validate(path)
	case DataValueKind_DV_PROPORTION:
		return d.Value.(*DV_PROPORTION).Validate(path)
	case DataValueKind_DV_QUANTITY:
		return d.Value.(*DV_QUANTITY).Validate(path)
	case DataValueKind_DV_SCALE:
		return d.Value.(*DV_SCALE).Validate(path)
	case DataValueKind_DV_STATE:
		return d.Value.(*DV_STATE).Validate(path)
	case DataValueKind_DV_TEXT:
		return d.Value.(*DV_TEXT).Validate(path)
	case DataValueKind_DV_TIME:
		return d.Value.(*DV_TIME).Validate(path)
	case DataValueKind_DV_DATE:
		return d.Value.(*DV_DATE).Validate(path)
	case DataValueKind_DV_DATE_TIME:
		return d.Value.(*DV_DATE_TIME).Validate(path)
	case DataValueKind_DV_DURATION:
		return d.Value.(*DV_DURATION).Validate(path)
	case DataValueKind_DV_PERIODIC_TIME_SPECIFICATION:
		return d.Value.(*DV_PERIODIC_TIME_SPECIFICATION).Validate(path)
	case DataValueKind_DV_GENERAL_TIME_SPECIFICATION:
		return d.Value.(*DV_GENERAL_TIME_SPECIFICATION).Validate(path)
	case DataValueKind_DV_MULTIMEDIA:
		return d.Value.(*DV_MULTIMEDIA).Validate(path)
	case DataValueKind_DV_PARSABLE:
		return d.Value.(*DV_PARSABLE).Validate(path)
	case DataValueKind_DV_URI:
		return d.Value.(*DV_URI).Validate(path)
	case DataValueKind_DV_EHR_URI:
		return d.Value.(*DV_EHR_URI).Validate(path)
	case DataValueKind_DV_COUNT:
		return d.Value.(*DV_COUNT).Validate(path)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          DATA_VALUE_TYPE,
					Path:           path,
					Message:        "value is not known DATA_VALUE subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}
}

func (d DataValueUnion) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value)
}

func (d *DataValueUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case DV_BOOLEAN_TYPE:
		d.Kind = DataValueKind_DV_BOOLEAN
		d.Value = &DV_BOOLEAN{}
	case DV_STATE_TYPE:
		d.Kind = DataValueKind_DV_STATE
		d.Value = &DV_STATE{}
	case DV_IDENTIFIER_TYPE:
		d.Kind = DataValueKind_DV_IDENTIFIER
		d.Value = &DV_IDENTIFIER{}
	case DV_TEXT_TYPE:
		d.Kind = DataValueKind_DV_TEXT
		d.Value = &DV_TEXT{}
	case DV_CODED_TEXT_TYPE:
		d.Kind = DataValueKind_DV_CODED_TEXT
		d.Value = &DV_CODED_TEXT{}
	case DV_PARAGRAPH_TYPE:
		d.Kind = DataValueKind_DV_PARAGRAPH
		d.Value = &DV_PARAGRAPH{}
	case DV_INTERVAL_TYPE:
		d.Kind = DataValueKind_DV_INTERVAL
		d.Value = &DV_INTERVAL{}
	case DV_ORDINAL_TYPE:
		d.Kind = DataValueKind_DV_ORDINAL
		d.Value = &DV_ORDINAL{}
	case DV_SCALE_TYPE:
		d.Kind = DataValueKind_DV_SCALE
		d.Value = &DV_SCALE{}
	case DV_QUANTITY_TYPE:
		d.Kind = DataValueKind_DV_QUANTITY
		d.Value = &DV_QUANTITY{}
	case DV_COUNT_TYPE:
		d.Kind = DataValueKind_DV_COUNT
		d.Value = &DV_COUNT{}
	case DV_PROPORTION_TYPE:
		d.Kind = DataValueKind_DV_PROPORTION
		d.Value = &DV_PROPORTION{}
	case DV_DATE_TYPE:
		d.Kind = DataValueKind_DV_DATE
		d.Value = &DV_DATE{}
	case DV_TIME_TYPE:
		d.Kind = DataValueKind_DV_TIME
		d.Value = &DV_TIME{}
	case DV_DATE_TIME_TYPE:
		d.Kind = DataValueKind_DV_DATE_TIME
		d.Value = &DV_DATE_TIME{}
	case DV_DURATION_TYPE:
		d.Kind = DataValueKind_DV_DURATION
		d.Value = &DV_DURATION{}
	case DV_PERIODIC_TIME_SPECIFICATION_TYPE:
		d.Kind = DataValueKind_DV_PERIODIC_TIME_SPECIFICATION
		d.Value = &DV_PERIODIC_TIME_SPECIFICATION{}
	case DV_GENERAL_TIME_SPECIFICATION_TYPE:
		d.Kind = DataValueKind_DV_GENERAL_TIME_SPECIFICATION
		d.Value = &DV_GENERAL_TIME_SPECIFICATION{}
	case DV_MULTIMEDIA_TYPE:
		d.Kind = DataValueKind_DV_MULTIMEDIA
		d.Value = &DV_MULTIMEDIA{}
	case DV_PARSABLE_TYPE:
		d.Kind = DataValueKind_DV_PARSABLE
		d.Value = &DV_PARSABLE{}
	case DV_URI_TYPE:
		d.Kind = DataValueKind_DV_URI
		d.Value = &DV_URI{}
	case DV_EHR_URI_TYPE:
		d.Kind = DataValueKind_DV_EHR_URI
		d.Value = &DV_EHR_URI{}
	default:
		d.Kind = DataValueKind_Unknown
		return nil
	}

	return json.Unmarshal(data, d.Value)
}
