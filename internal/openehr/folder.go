package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const FOLDER_MODEL_NAME = "FOLDER"

type FOLDER struct {
	Type_            util.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                       `json:"name"`
	ArchetypeNodeID  string                          `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Items            util.Optional[[]OBJECT_REF]     `json:"items,omitzero"`
	Folders          util.Optional[[]FOLDER]         `json:"folders,omitzero"`
	Details          util.Optional[X_ITEM_STRUCTURE] `json:"details,omitzero"`
}

func (f *FOLDER) isVersionModel() {}

func (f *FOLDER) SetModelName() {
	f.Type_ = util.Some(FOLDER_MODEL_NAME)
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

func (f *FOLDER) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if f.Type_.E && f.Type_.V != FOLDER_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          FOLDER_MODEL_NAME,
			Path:           attrPath,
			Message:        "invalid _type field",
			Recommendation: "Ensure _type field is set to FOLDER",
		})
	}

	// Validate uid
	if f.UID.E {
		attrPath = path + ".uid"
		errors = append(errors, f.UID.V.Validate(attrPath)...)
	}

	// Validate links
	if f.Links.E {
		for i := range f.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			errors = append(errors, f.Links.V[i].Validate(attrPath)...)
		}
	}

	// Validate archetype_details
	if f.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errors = append(errors, f.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate feeder_audit
	if f.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errors = append(errors, f.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate items
	if f.Items.E {
		for i := range f.Items.V {
			attrPath = fmt.Sprintf("%s.items[%d]", path, i)
			errors = append(errors, f.Items.V[i].Validate(attrPath)...)
		}
	}

	// Validate folders
	if f.Folders.E {
		for i := range f.Folders.V {
			attrPath = fmt.Sprintf("%s.folders[%d]", path, i)
			errors = append(errors, f.Folders.V[i].Validate(attrPath)...)
		}
	}

	// Validate details
	if f.Details.E {
		attrPath = path + ".details"
		errors = append(errors, f.Details.V.Validate(attrPath)...)
	}

	return errors
}
