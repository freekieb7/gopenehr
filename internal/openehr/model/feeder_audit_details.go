package model

import (
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const FEEDER_AUDIT_DETAILS_MODEL_NAME string = "FEEDER_AUDIT_DETAILS"

type FEEDER_AUDIT_DETAILS struct {
	Type_        utils.Optional[string]           `json:"_type,omitzero"`
	SystemID     string                           `json:"system_id"`
	Location     utils.Optional[PARTY_IDENTIFIED] `json:"location,omitzero"`
	Subject      utils.Optional[X_PARTY_PROXY]    `json:"subject,omitzero"`
	Provider     utils.Optional[PARTY_IDENTIFIED] `json:"provider,omitzero"`
	Time         utils.Optional[DV_DATE_TIME]     `json:"time,omitzero"`
	VersionID    utils.Optional[string]           `json:"version_id,omitzero"`
	OtherDetails utils.Optional[X_ITEM_STRUCTURE] `json:"other_details,omitzero"`
}

func (f *FEEDER_AUDIT_DETAILS) SetModelName() {
	f.Type_ = utils.Some(FEEDER_AUDIT_DETAILS_MODEL_NAME)
	if f.Location.E {
		f.Location.V.SetModelName()
	}
	if f.Subject.E {
		f.Subject.V.SetModelName()
	}
	if f.Provider.E {
		f.Provider.V.SetModelName()
	}
	if f.Time.E {
		f.Time.V.SetModelName()
	}
	if f.OtherDetails.E {
		f.OtherDetails.V.SetModelName()
	}
}

func (f *FEEDER_AUDIT_DETAILS) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if f.Type_.E && f.Type_.V != FEEDER_AUDIT_DETAILS_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          FEEDER_AUDIT_DETAILS_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to FEEDER_AUDIT_DETAILS",
		})
	}

	// Validate system_id
	attrPath = path + ".system_id"
	if f.SystemID == "" {
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          FEEDER_AUDIT_DETAILS_MODEL_NAME,
			Path:           attrPath,
			Message:        "system_id field cannot be empty",
			Recommendation: "Ensure system_id field is not empty",
		})
	}

	// Validate location
	if f.Location.E {
		attrPath = path + ".location"
		validateErr.Errs = append(validateErr.Errs, f.Location.V.Validate(attrPath).Errs...)
	}

	// Validate subject
	if f.Subject.E {
		attrPath = path + ".subject"
		validateErr.Errs = append(validateErr.Errs, f.Subject.V.Validate(attrPath).Errs...)
	}

	// Validate provider
	if f.Provider.E {
		attrPath = path + ".provider"
		validateErr.Errs = append(validateErr.Errs, f.Provider.V.Validate(attrPath).Errs...)
	}

	// Validate time
	if f.Time.E {
		attrPath = path + ".time"
		validateErr.Errs = append(validateErr.Errs, f.Time.V.Validate(attrPath).Errs...)
	}

	// Validate other_details
	if f.OtherDetails.E {
		attrPath = path + ".other_details"
		validateErr.Errs = append(validateErr.Errs, f.OtherDetails.V.Validate(attrPath).Errs...)
	}

	return validateErr
}
