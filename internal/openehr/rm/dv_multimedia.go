package rm

import (
	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const DV_MULTIMEDIA_TYPE string = "DV_MULTIMEDIA"

type DV_MULTIMEDIA struct {
	Type_                   utils.Optional[string]         `json:"_type,omitzero"`
	Charset                 utils.Optional[CODE_PHRASE]    `json:"charset,omitzero"`
	Language                utils.Optional[CODE_PHRASE]    `json:"language,omitzero"`
	AlternateText           utils.Optional[string]         `json:"alternate_text,omitzero"`
	Uri                     utils.Optional[DV_URI]         `json:"uri,omitzero"`
	Data                    utils.Optional[string]         `json:"data,omitzero"`
	MediaType               CODE_PHRASE                    `json:"media_type"`
	CompressionAlgorithm    utils.Optional[CODE_PHRASE]    `json:"compression_algorithm,omitzero"`
	IntegrityCheck          utils.Optional[string]         `json:"integrity_check,omitzero"`
	IntegrityCheckAlgorithm utils.Optional[CODE_PHRASE]    `json:"integrity_check_algorithm,omitzero"`
	Thumbnail               utils.Optional[*DV_MULTIMEDIA] `json:"thumbnail,omitzero"`
	Size                    int64                          `json:"size"`
}

func (d *DV_MULTIMEDIA) SetModelName() {
	d.Type_ = utils.Some(DV_MULTIMEDIA_TYPE)
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

func (d *DV_MULTIMEDIA) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if d.Type_.E && d.Type_.V != DV_MULTIMEDIA_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_MULTIMEDIA_TYPE,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to " + DV_MULTIMEDIA_TYPE,
		})
	}

	// Validate charset
	if d.Charset.E {
		attrPath = path + ".charset"
		if !terminology.IsValidCharsetTerminologyID(d.Charset.V.TerminologyID.Value) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_MULTIMEDIA_TYPE,
				Path:           attrPath,
				Message:        "invalid charset terminology ID",
				Recommendation: "Ensure charset terminology ID is valid",
			})
		}
		if !terminology.IsValidCharset(d.Charset.V.CodeString) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_MULTIMEDIA_TYPE,
				Path:           attrPath,
				Message:        "invalid charset code string",
				Recommendation: "Ensure charset code string is valid",
			})
		}

		validateErr.Errs = append(validateErr.Errs, d.Charset.V.Validate(attrPath).Errs...)
	}

	// Validate language
	if d.Language.E {
		attrPath = path + ".language"
		if !terminology.IsValidLanguageTerminologyID(d.Language.V.TerminologyID.Value) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_MULTIMEDIA_TYPE,
				Path:           attrPath,
				Message:        "invalid language terminology ID",
				Recommendation: "Ensure language terminology ID is valid",
			})
		}
		if !terminology.IsValidLanguageCode(d.Language.V.CodeString) {
			validateErr.Errs = append(validateErr.Errs, util.ValidationError{
				Model:          DV_MULTIMEDIA_TYPE,
				Path:           attrPath,
				Message:        "invalid language code string",
				Recommendation: "Ensure language code string is valid",
			})
		}

		validateErr.Errs = append(validateErr.Errs, d.Language.V.Validate(attrPath).Errs...)
	}

	// Validate uri
	if d.Uri.E {
		attrPath = path + ".uri"
		validateErr.Errs = append(validateErr.Errs, d.Uri.V.Validate(attrPath).Errs...)
	}

	// Validate media_type
	attrPath = path + ".media_type"
	validateErr.Errs = append(validateErr.Errs, d.MediaType.Validate(attrPath).Errs...)

	// Validate compression_algorithm
	if d.CompressionAlgorithm.E {
		attrPath = path + ".compression_algorithm"
		validateErr.Errs = append(validateErr.Errs, d.CompressionAlgorithm.V.Validate(attrPath).Errs...)
	}

	// Validate integrity_check_algorithm
	if d.IntegrityCheckAlgorithm.E {
		attrPath = path + ".integrity_check_algorithm"
		validateErr.Errs = append(validateErr.Errs, d.IntegrityCheckAlgorithm.V.Validate(attrPath).Errs...)
	}

	// Validate thumbnail
	if d.Thumbnail.E {
		attrPath = path + ".thumbnail"
		validateErr.Errs = append(validateErr.Errs, d.Thumbnail.V.Validate(attrPath).Errs...)
	}

	// Validate size
	if d.Size < 0 {
		attrPath = path + ".size"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          DV_MULTIMEDIA_TYPE,
			Path:           attrPath,
			Message:        "size must be non-negative",
			Recommendation: "Ensure size is zero or positive",
		})
	}

	return validateErr
}
