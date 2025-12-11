package rm

import (
	"github.com/bytedance/sonic"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_ENCAPSULATED_TYPE string = "DV_ENCAPSULATED"

type DvEncapsulatedKind int

const (
	DV_ENCAPSULATED_kind_unknown DvEncapsulatedKind = iota
	DV_ENCAPSULATED_kind_DV_MULTIMEDIA
	DV_ENCAPSULATED_kind_DV_EHR_URI
)

type DvEncapsulatedUnion struct {
	Kind  DvEncapsulatedKind
	Value any
}

func (x *DvEncapsulatedUnion) SetModelName() {
	switch x.Kind {
	case DV_ENCAPSULATED_kind_DV_MULTIMEDIA:
		x.Value.(*DV_MULTIMEDIA).SetModelName()
	case DV_ENCAPSULATED_kind_DV_EHR_URI:
		x.Value.(*DV_EHR_URI).SetModelName()
	}
}

func (x *DvEncapsulatedUnion) Validate(path string) util.ValidateError {
	switch x.Kind {
	case DV_ENCAPSULATED_kind_DV_MULTIMEDIA:
		return x.Value.(*DV_MULTIMEDIA).Validate(path)
	case DV_ENCAPSULATED_kind_DV_EHR_URI:
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
	return sonic.Marshal(d.Value)
}

func (d *DvEncapsulatedUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case DV_MULTIMEDIA_TYPE:
		d.Kind = DV_ENCAPSULATED_kind_DV_MULTIMEDIA
		d.Value = &DV_MULTIMEDIA{}
	case DV_EHR_URI_TYPE:
		d.Kind = DV_ENCAPSULATED_kind_DV_EHR_URI
		d.Value = &DV_EHR_URI{}
	default:
		d.Kind = DV_ENCAPSULATED_kind_unknown
		return nil
	}

	return sonic.Unmarshal(data, d.Value)
}

func (d *DvEncapsulatedUnion) DV_MULTIMEDIA() *DV_MULTIMEDIA {
	if d.Kind != DV_ENCAPSULATED_kind_DV_MULTIMEDIA {
		return nil
	}
	return d.Value.(*DV_MULTIMEDIA)
}

func (d *DvEncapsulatedUnion) DV_EHR_URI() *DV_EHR_URI {
	if d.Kind != DV_ENCAPSULATED_kind_DV_EHR_URI {
		return nil
	}
	return d.Value.(*DV_EHR_URI)
}
