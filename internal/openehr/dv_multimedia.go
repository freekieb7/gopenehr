package openehr

import (
	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const DV_MULTIMEDIA_MODEL_NAME string = "DV_MULTIMEDIA"

type DV_MULTIMEDIA struct {
	Type_                   util.Optional[string]         `json:"_type,omitzero"`
	Charset                 util.Optional[CODE_PHRASE]    `json:"charset,omitzero"`
	Language                util.Optional[CODE_PHRASE]    `json:"language,omitzero"`
	AlternateText           util.Optional[string]         `json:"alternate_text,omitzero"`
	Uri                     util.Optional[DV_URI]         `json:"uri,omitzero"`
	Data                    util.Optional[string]         `json:"data,omitzero"`
	MediaType               CODE_PHRASE                   `json:"media_type"`
	CompressionAlgorithm    util.Optional[CODE_PHRASE]    `json:"compression_algorithm,omitzero"`
	IntegrityCheck          util.Optional[string]         `json:"integrity_check,omitzero"`
	IntegrityCheckAlgorithm util.Optional[CODE_PHRASE]    `json:"integrity_check_algorithm,omitzero"`
	Thumbnail               util.Optional[*DV_MULTIMEDIA] `json:"thumbnail,omitzero"`
	Size                    int64                         `json:"size"`
}

func (d *DV_MULTIMEDIA) isDataValueModel() {}

func (d *DV_MULTIMEDIA) HasModelName() bool {
	return d.Type_.E
}

func (d *DV_MULTIMEDIA) SetModelName() {
	d.Type_ = util.Some(DV_MULTIMEDIA_MODEL_NAME)
	if d.Charset.E {
		d.Charset.V.SetModelName()
	}
	if d.Language.E {
		d.Language.V.SetModelName()
	}
	if d.Uri.E {
		d.Uri.V.SetModelName()
	}
	if d.CompressionAlgorithm.E {
		d.CompressionAlgorithm.V.SetModelName()
	}
	if d.IntegrityCheckAlgorithm.E {
		d.IntegrityCheckAlgorithm.V.SetModelName()
	}
	if d.Thumbnail.E {
		d.Thumbnail.V.SetModelName()
	}
}

func (d *DV_MULTIMEDIA) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_MULTIMEDIA_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          DV_MULTIMEDIA_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + DV_MULTIMEDIA_MODEL_NAME,
		})
	}

	// Validate charset
	if d.Charset.E {
		attrPath = path + ".charset"
		if !terminology.IsValidCharsetTerminologyID(d.Charset.V.TerminologyID.Value) {
			errors = append(errors, util.ValidationError{
				Model:          DV_MULTIMEDIA_MODEL_NAME,
				Path:           attrPath,
				Message:        "invalid charset terminology ID",
				Recommendation: "Ensure charset terminology ID is valid",
			})
		}
		if !terminology.IsValidCharset(d.Charset.V.CodeString) {
			errors = append(errors, util.ValidationError{
				Model:          DV_MULTIMEDIA_MODEL_NAME,
				Path:           attrPath,
				Message:        "invalid charset code string",
				Recommendation: "Ensure charset code string is valid",
			})
		}

		errors = append(errors, d.Charset.V.Validate(attrPath)...)
	}

	// Validate language
	if d.Language.E {
		attrPath = path + ".language"
		if !terminology.IsValidLanguageTerminologyID(d.Language.V.TerminologyID.Value) {
			errors = append(errors, util.ValidationError{
				Model:          DV_MULTIMEDIA_MODEL_NAME,
				Path:           attrPath,
				Message:        "invalid language terminology ID",
				Recommendation: "Ensure language terminology ID is valid",
			})
		}
		if !terminology.IsValidLanguageCode(d.Language.V.CodeString) {
			errors = append(errors, util.ValidationError{
				Model:          DV_MULTIMEDIA_MODEL_NAME,
				Path:           attrPath,
				Message:        "invalid language code string",
				Recommendation: "Ensure language code string is valid",
			})
		}

		errors = append(errors, d.Language.V.Validate(attrPath)...)
	}

	// Validate uri
	if d.Uri.E {
		attrPath = path + ".uri"
		errors = append(errors, d.Uri.V.Validate(attrPath)...)
	}

	// Validate media_type
	attrPath = path + ".media_type"
	errors = append(errors, d.MediaType.Validate(attrPath)...)

	// Validate compression_algorithm
	if d.CompressionAlgorithm.E {
		attrPath = path + ".compression_algorithm"
		errors = append(errors, d.CompressionAlgorithm.V.Validate(attrPath)...)
	}

	// Validate integrity_check_algorithm
	if d.IntegrityCheckAlgorithm.E {
		attrPath = path + ".integrity_check_algorithm"
		errors = append(errors, d.IntegrityCheckAlgorithm.V.Validate(attrPath)...)
	}

	// Validate thumbnail
	if d.Thumbnail.E {
		attrPath = path + ".thumbnail"
		errors = append(errors, d.Thumbnail.V.Validate(attrPath)...)
	}

	// Validate size
	if d.Size < 0 {
		attrPath = path + ".size"
		errors = append(errors, util.ValidationError{
			Model:          DV_MULTIMEDIA_MODEL_NAME,
			Path:           attrPath,
			Message:        "size must be non-negative",
			Recommendation: "Ensure size is zero or positive",
		})
	}

	return errors
}
