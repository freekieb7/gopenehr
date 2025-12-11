package rm

import (
	"github.com/bytedance/sonic"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DATA_VALUE_TYPE string = "DATA_VALUE"

type DataValueKind int

const (
	DATA_VALUE_kind_unknown DataValueKind = iota
	DATA_VALUE_kind_DV_BOOLEAN
	DATA_VALUE_kind_DV_CODED_TEXT
	DATA_VALUE_kind_DV_IDENTIFIER
	DATA_VALUE_kind_DV_INTERVAL
	DATA_VALUE_kind_DV_ORDINAL
	DATA_VALUE_kind_DV_PARAGRAPH
	DATA_VALUE_kind_DV_PROPORTION
	DATA_VALUE_kind_DV_QUANTITY
	DATA_VALUE_kind_DV_SCALE
	DATA_VALUE_kind_DV_STATE
	DATA_VALUE_kind_DV_TEXT
	DATA_VALUE_kind_DV_TIME
	DATA_VALUE_kind_DV_DATE
	DATA_VALUE_kind_DV_DATE_TIME
	DATA_VALUE_kind_DV_DURATION
	DATA_VALUE_kind_DV_PERIODIC_TIME_SPECIFICATION
	DATA_VALUE_kind_DV_GENERAL_TIME_SPECIFICATION
	DATA_VALUE_kind_DV_MULTIMEDIA
	DATA_VALUE_kind_DV_PARSABLE
	DATA_VALUE_kind_DV_URI
	DATA_VALUE_kind_DV_EHR_URI
	DATA_VALUE_kind_DV_COUNT
)

type DataValueUnion struct {
	Kind  DataValueKind
	Value any
}

func (d *DataValueUnion) SetModelName() {
	switch d.Kind {
	case DATA_VALUE_kind_DV_BOOLEAN:
		d.Value.(*DV_BOOLEAN).SetModelName()
	case DATA_VALUE_kind_DV_CODED_TEXT:
		d.Value.(*DV_CODED_TEXT).SetModelName()
	case DATA_VALUE_kind_DV_IDENTIFIER:
		d.Value.(*DV_IDENTIFIER).SetModelName()
	case DATA_VALUE_kind_DV_INTERVAL:
		d.Value.(*DV_INTERVAL[any]).SetModelName()
	case DATA_VALUE_kind_DV_ORDINAL:
		d.Value.(*DV_ORDINAL).SetModelName()
	case DATA_VALUE_kind_DV_PARAGRAPH:
		d.Value.(*DV_PARAGRAPH).SetModelName()
	case DATA_VALUE_kind_DV_PROPORTION:
		d.Value.(*DV_PROPORTION).SetModelName()
	case DATA_VALUE_kind_DV_QUANTITY:
		d.Value.(*DV_QUANTITY).SetModelName()
	case DATA_VALUE_kind_DV_SCALE:
		d.Value.(*DV_SCALE).SetModelName()
	case DATA_VALUE_kind_DV_STATE:
		d.Value.(*DV_STATE).SetModelName()
	case DATA_VALUE_kind_DV_TEXT:
		d.Value.(*DV_TEXT).SetModelName()
	case DATA_VALUE_kind_DV_TIME:
		d.Value.(*DV_TIME).SetModelName()
	case DATA_VALUE_kind_DV_DATE:
		d.Value.(*DV_DATE).SetModelName()
	case DATA_VALUE_kind_DV_DATE_TIME:
		d.Value.(*DV_DATE_TIME).SetModelName()
	case DATA_VALUE_kind_DV_DURATION:
		d.Value.(*DV_DURATION).SetModelName()
	case DATA_VALUE_kind_DV_PERIODIC_TIME_SPECIFICATION:
		d.Value.(*DV_PERIODIC_TIME_SPECIFICATION).SetModelName()
	case DATA_VALUE_kind_DV_GENERAL_TIME_SPECIFICATION:
		d.Value.(*DV_GENERAL_TIME_SPECIFICATION).SetModelName()
	case DATA_VALUE_kind_DV_MULTIMEDIA:
		d.Value.(*DV_MULTIMEDIA).SetModelName()
	case DATA_VALUE_kind_DV_PARSABLE:
		d.Value.(*DV_PARSABLE).SetModelName()
	case DATA_VALUE_kind_DV_URI:
		d.Value.(*DV_URI).SetModelName()
	case DATA_VALUE_kind_DV_EHR_URI:
		d.Value.(*DV_EHR_URI).SetModelName()
	case DATA_VALUE_kind_DV_COUNT:
		d.Value.(*DV_COUNT).SetModelName()
	}
}

func (d *DataValueUnion) Validate(path string) util.ValidateError {
	switch d.Kind {
	case DATA_VALUE_kind_DV_BOOLEAN:
		return d.Value.(*DV_BOOLEAN).Validate(path)
	case DATA_VALUE_kind_DV_CODED_TEXT:
		return d.Value.(*DV_CODED_TEXT).Validate(path)
	case DATA_VALUE_kind_DV_IDENTIFIER:
		return d.Value.(*DV_IDENTIFIER).Validate(path)
	case DATA_VALUE_kind_DV_INTERVAL:
		return d.Value.(*DV_INTERVAL[any]).Validate(path)
	case DATA_VALUE_kind_DV_ORDINAL:
		return d.Value.(*DV_ORDINAL).Validate(path)
	case DATA_VALUE_kind_DV_PARAGRAPH:
		return d.Value.(*DV_PARAGRAPH).Validate(path)
	case DATA_VALUE_kind_DV_PROPORTION:
		return d.Value.(*DV_PROPORTION).Validate(path)
	case DATA_VALUE_kind_DV_QUANTITY:
		return d.Value.(*DV_QUANTITY).Validate(path)
	case DATA_VALUE_kind_DV_SCALE:
		return d.Value.(*DV_SCALE).Validate(path)
	case DATA_VALUE_kind_DV_STATE:
		return d.Value.(*DV_STATE).Validate(path)
	case DATA_VALUE_kind_DV_TEXT:
		return d.Value.(*DV_TEXT).Validate(path)
	case DATA_VALUE_kind_DV_TIME:
		return d.Value.(*DV_TIME).Validate(path)
	case DATA_VALUE_kind_DV_DATE:
		return d.Value.(*DV_DATE).Validate(path)
	case DATA_VALUE_kind_DV_DATE_TIME:
		return d.Value.(*DV_DATE_TIME).Validate(path)
	case DATA_VALUE_kind_DV_DURATION:
		return d.Value.(*DV_DURATION).Validate(path)
	case DATA_VALUE_kind_DV_PERIODIC_TIME_SPECIFICATION:
		return d.Value.(*DV_PERIODIC_TIME_SPECIFICATION).Validate(path)
	case DATA_VALUE_kind_DV_GENERAL_TIME_SPECIFICATION:
		return d.Value.(*DV_GENERAL_TIME_SPECIFICATION).Validate(path)
	case DATA_VALUE_kind_DV_MULTIMEDIA:
		return d.Value.(*DV_MULTIMEDIA).Validate(path)
	case DATA_VALUE_kind_DV_PARSABLE:
		return d.Value.(*DV_PARSABLE).Validate(path)
	case DATA_VALUE_kind_DV_URI:
		return d.Value.(*DV_URI).Validate(path)
	case DATA_VALUE_kind_DV_EHR_URI:
		return d.Value.(*DV_EHR_URI).Validate(path)
	case DATA_VALUE_kind_DV_COUNT:
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
	return sonic.Marshal(d.Value)
}

func (d *DataValueUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case DV_BOOLEAN_TYPE:
		d.Kind = DATA_VALUE_kind_DV_BOOLEAN
		d.Value = &DV_BOOLEAN{}
	case DV_STATE_TYPE:
		d.Kind = DATA_VALUE_kind_DV_STATE
		d.Value = &DV_STATE{}
	case DV_IDENTIFIER_TYPE:
		d.Kind = DATA_VALUE_kind_DV_IDENTIFIER
		d.Value = &DV_IDENTIFIER{}
	case DV_TEXT_TYPE:
		d.Kind = DATA_VALUE_kind_DV_TEXT
		d.Value = &DV_TEXT{}
	case DV_CODED_TEXT_TYPE:
		d.Kind = DATA_VALUE_kind_DV_CODED_TEXT
		d.Value = &DV_CODED_TEXT{}
	case DV_PARAGRAPH_TYPE:
		d.Kind = DATA_VALUE_kind_DV_PARAGRAPH
		d.Value = &DV_PARAGRAPH{}
	case DV_INTERVAL_TYPE:
		d.Kind = DATA_VALUE_kind_DV_INTERVAL
		d.Value = &DV_INTERVAL[any]{}
	case DV_ORDINAL_TYPE:
		d.Kind = DATA_VALUE_kind_DV_ORDINAL
		d.Value = &DV_ORDINAL{}
	case DV_SCALE_TYPE:
		d.Kind = DATA_VALUE_kind_DV_SCALE
		d.Value = &DV_SCALE{}
	case DV_QUANTITY_TYPE:
		d.Kind = DATA_VALUE_kind_DV_QUANTITY
		d.Value = &DV_QUANTITY{}
	case DV_COUNT_TYPE:
		d.Kind = DATA_VALUE_kind_DV_COUNT
		d.Value = &DV_COUNT{}
	case DV_PROPORTION_TYPE:
		d.Kind = DATA_VALUE_kind_DV_PROPORTION
		d.Value = &DV_PROPORTION{}
	case DV_DATE_TYPE:
		d.Kind = DATA_VALUE_kind_DV_DATE
		d.Value = &DV_DATE{}
	case DV_TIME_TYPE:
		d.Kind = DATA_VALUE_kind_DV_TIME
		d.Value = &DV_TIME{}
	case DV_DATE_TIME_TYPE:
		d.Kind = DATA_VALUE_kind_DV_DATE_TIME
		d.Value = &DV_DATE_TIME{}
	case DV_DURATION_TYPE:
		d.Kind = DATA_VALUE_kind_DV_DURATION
		d.Value = &DV_DURATION{}
	case DV_PERIODIC_TIME_SPECIFICATION_TYPE:
		d.Kind = DATA_VALUE_kind_DV_PERIODIC_TIME_SPECIFICATION
		d.Value = &DV_PERIODIC_TIME_SPECIFICATION{}
	case DV_GENERAL_TIME_SPECIFICATION_TYPE:
		d.Kind = DATA_VALUE_kind_DV_GENERAL_TIME_SPECIFICATION
		d.Value = &DV_GENERAL_TIME_SPECIFICATION{}
	case DV_MULTIMEDIA_TYPE:
		d.Kind = DATA_VALUE_kind_DV_MULTIMEDIA
		d.Value = &DV_MULTIMEDIA{}
	case DV_PARSABLE_TYPE:
		d.Kind = DATA_VALUE_kind_DV_PARSABLE
		d.Value = &DV_PARSABLE{}
	case DV_URI_TYPE:
		d.Kind = DATA_VALUE_kind_DV_URI
		d.Value = &DV_URI{}
	case DV_EHR_URI_TYPE:
		d.Kind = DATA_VALUE_kind_DV_EHR_URI
		d.Value = &DV_EHR_URI{}
	default:
		d.Kind = DATA_VALUE_kind_unknown
		return nil
	}

	return sonic.Unmarshal(data, d.Value)
}

func (d *DataValueUnion) DV_BOOLEAN() *DV_BOOLEAN {
	if d.Kind == DATA_VALUE_kind_DV_BOOLEAN {
		return d.Value.(*DV_BOOLEAN)
	}
	return nil
}

func (d *DataValueUnion) DV_CODED_TEXT() *DV_CODED_TEXT {
	if d.Kind == DATA_VALUE_kind_DV_CODED_TEXT {
		return d.Value.(*DV_CODED_TEXT)
	}
	return nil
}

func (d *DataValueUnion) DV_IDENTIFIER() *DV_IDENTIFIER {
	if d.Kind == DATA_VALUE_kind_DV_IDENTIFIER {
		return d.Value.(*DV_IDENTIFIER)
	}
	return nil
}

func (d *DataValueUnion) DV_INTERVAL() *DV_INTERVAL[any] {
	if d.Kind == DATA_VALUE_kind_DV_INTERVAL {
		return d.Value.(*DV_INTERVAL[any])
	}
	return nil
}

func (d *DataValueUnion) DV_ORDINAL() *DV_ORDINAL {
	if d.Kind == DATA_VALUE_kind_DV_ORDINAL {
		return d.Value.(*DV_ORDINAL)
	}
	return nil
}

func (d *DataValueUnion) DV_PARAGRAPH() *DV_PARAGRAPH {
	if d.Kind == DATA_VALUE_kind_DV_PARAGRAPH {
		return d.Value.(*DV_PARAGRAPH)
	}
	return nil
}

func (d *DataValueUnion) DV_PROPORTION() *DV_PROPORTION {
	if d.Kind == DATA_VALUE_kind_DV_PROPORTION {
		return d.Value.(*DV_PROPORTION)
	}
	return nil
}

func (d *DataValueUnion) DV_QUANTITY() *DV_QUANTITY {
	if d.Kind == DATA_VALUE_kind_DV_QUANTITY {
		return d.Value.(*DV_QUANTITY)
	}
	return nil
}

func (d *DataValueUnion) DV_SCALE() *DV_SCALE {
	if d.Kind == DATA_VALUE_kind_DV_SCALE {
		return d.Value.(*DV_SCALE)
	}
	return nil
}

func (d *DataValueUnion) DV_STATE() *DV_STATE {
	if d.Kind == DATA_VALUE_kind_DV_STATE {
		return d.Value.(*DV_STATE)
	}
	return nil
}

func (d *DataValueUnion) DV_TEXT() *DV_TEXT {
	if d.Kind == DATA_VALUE_kind_DV_TEXT {
		return d.Value.(*DV_TEXT)
	}
	return nil
}

func (d *DataValueUnion) DV_TIME() *DV_TIME {
	if d.Kind == DATA_VALUE_kind_DV_TIME {
		return d.Value.(*DV_TIME)
	}
	return nil
}

func (d *DataValueUnion) DV_DATE() *DV_DATE {
	if d.Kind == DATA_VALUE_kind_DV_DATE {
		return d.Value.(*DV_DATE)
	}
	return nil
}

func (d *DataValueUnion) DV_DATE_TIME() *DV_DATE_TIME {
	if d.Kind == DATA_VALUE_kind_DV_DATE_TIME {
		return d.Value.(*DV_DATE_TIME)
	}
	return nil
}

func (d *DataValueUnion) DV_DURATION() *DV_DURATION {
	if d.Kind == DATA_VALUE_kind_DV_DURATION {
		return d.Value.(*DV_DURATION)
	}
	return nil
}

func (d *DataValueUnion) DV_PERIODIC_TIME_SPECIFICATION() *DV_PERIODIC_TIME_SPECIFICATION {
	if d.Kind == DATA_VALUE_kind_DV_PERIODIC_TIME_SPECIFICATION {
		return d.Value.(*DV_PERIODIC_TIME_SPECIFICATION)
	}
	return nil
}

func (d *DataValueUnion) DV_GENERAL_TIME_SPECIFICATION() *DV_GENERAL_TIME_SPECIFICATION {
	if d.Kind == DATA_VALUE_kind_DV_GENERAL_TIME_SPECIFICATION {
		return d.Value.(*DV_GENERAL_TIME_SPECIFICATION)
	}
	return nil
}

func (d *DataValueUnion) DV_MULTIMEDIA() *DV_MULTIMEDIA {
	if d.Kind == DATA_VALUE_kind_DV_MULTIMEDIA {
		return d.Value.(*DV_MULTIMEDIA)
	}
	return nil
}

func (d *DataValueUnion) DV_PARSABLE() *DV_PARSABLE {
	if d.Kind == DATA_VALUE_kind_DV_PARSABLE {
		return d.Value.(*DV_PARSABLE)
	}
	return nil
}

func (d *DataValueUnion) DV_URI() *DV_URI {
	if d.Kind == DATA_VALUE_kind_DV_URI {
		return d.Value.(*DV_URI)
	}
	return nil
}

func (d *DataValueUnion) DV_EHR_URI() *DV_EHR_URI {
	if d.Kind == DATA_VALUE_kind_DV_EHR_URI {
		return d.Value.(*DV_EHR_URI)
	}
	return nil
}

func (d *DataValueUnion) DV_COUNT() *DV_COUNT {
	if d.Kind == DATA_VALUE_kind_DV_COUNT {
		return d.Value.(*DV_COUNT)
	}
	return nil
}
