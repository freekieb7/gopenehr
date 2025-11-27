package model

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const FOLDER_MODEL_NAME = "FOLDER"

type FOLDER struct {
	Type_            utils.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                        `json:"name"`
	ArchetypeNodeID  string                           `json:"archetype_node_id"`
	UID              utils.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Items            utils.Optional[[]OBJECT_REF]     `json:"items,omitzero"`
	Folders          utils.Optional[[]FOLDER]         `json:"folders,omitzero"`
	Details          utils.Optional[X_ITEM_STRUCTURE] `json:"details,omitzero"`
}

func (f *FOLDER) isVersionModel() {}

func (f *FOLDER) SetModelName() {
	f.Type_ = utils.Some(FOLDER_MODEL_NAME)
	if f.UID.E {
		f.UID.V.SetModelName()
	}
	if f.Links.E {
		for i := range f.Links.V {
			f.Links.V[i].SetModelName()
		}
	}
	if f.ArchetypeDetails.E {
		f.ArchetypeDetails.V.SetModelName()
	}
	if f.FeederAudit.E {
		f.FeederAudit.V.SetModelName()
	}
	if f.Items.E {
		for i := range f.Items.V {
			f.Items.V[i].SetModelName()
		}
	}
	if f.Folders.E {
		for i := range f.Folders.V {
			f.Folders.V[i].SetModelName()
		}
	}
	if f.Details.E {
		f.Details.V.SetModelName()
	}
}

func (f *FOLDER) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if f.Type_.E && f.Type_.V != FOLDER_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          FOLDER_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to FOLDER",
		})
	}

	// Validate uid
	if f.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, f.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if f.Links.E {
		for i := range f.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, f.Links.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if f.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, f.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate feeder_audit
	if f.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, f.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate items
	if f.Items.E {
		for i := range f.Items.V {
			attrPath = fmt.Sprintf("%s.items[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, f.Items.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate folders
	if f.Folders.E {
		for i := range f.Folders.V {
			attrPath = fmt.Sprintf("%s.folders[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, f.Folders.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate details
	if f.Details.E {
		attrPath = path + ".details"
		validateErr.Errs = append(validateErr.Errs, f.Details.V.Validate(attrPath).Errs...)
	}

	return validateErr
}

func (f *FOLDER) ObjectVersionID() OBJECT_VERSION_ID {
	if f.UID.E {
		if objVerID, ok := f.UID.V.Value.(*OBJECT_VERSION_ID); ok {
			return *objVerID
		}
	}
	return OBJECT_VERSION_ID{}
}
