package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_ENCAPSULATED_TYPE string = "DV_ENCAPSULATED"

type DvEncapsulatedKind int

const (
	DvEncapsulatedKind_Unknown DvEncapsulatedKind = iota
	DvEncapsulatedKind_DV_MULTIMEDIA
	DvEncapsulatedKind_DV_EHR_URI
)

type DvEncapsulatedUnion struct {
	Kind  DvEncapsulatedKind
	Value any
}

func (x *DvEncapsulatedUnion) SetModelName() {
	switch x.Kind {
	case DvEncapsulatedKind_DV_MULTIMEDIA:
		x.Value.(*DV_MULTIMEDIA).SetModelName()
	case DvEncapsulatedKind_DV_EHR_URI:
		x.Value.(*DV_EHR_URI).SetModelName()
	}
}

func (x *DvEncapsulatedUnion) Validate(path string) util.ValidateError {
	switch x.Kind {
	case DvEncapsulatedKind_DV_MULTIMEDIA:
		return x.Value.(*DV_MULTIMEDIA).Validate(path)
	case DvEncapsulatedKind_DV_EHR_URI:
		return x.Value.(*DV_EHR_URI).Validate(path)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          DV_ENCAPSULATED_TYPE,
					Path:           path,
					Message:        "value is not known DV_ENCAPSULATED subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}
}

func (d DvEncapsulatedUnion) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value)
}

func (d *DvEncapsulatedUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case DV_MULTIMEDIA_TYPE:
		d.Kind = DvEncapsulatedKind_DV_MULTIMEDIA
		d.Value = &DV_MULTIMEDIA{}
	case DV_EHR_URI_TYPE:
		d.Kind = DvEncapsulatedKind_DV_EHR_URI
		d.Value = &DV_EHR_URI{}
	default:
		d.Kind = DvEncapsulatedKind_Unknown
		return nil
	}

	return json.Unmarshal(data, d.Value)
}
