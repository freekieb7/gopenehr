package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const SECTION_MODEL_NAME string = "SECTION"

type SECTION struct {
	Type_            util.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                       `json:"name"`
	ArchetypeNodeID  string                          `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Items            util.Optional[[]X_CONTENT_ITEM] `json:"items,omitzero"`
}

func (s *SECTION) isContentItemModel() {}

func (s *SECTION) HasModelName() bool {
	return s.Type_.E
}

func (s *SECTION) SetModelName() {
	s.Type_ = util.Some(SECTION_MODEL_NAME)
	s.Name.SetModelName()
	if s.UID.E {
		s.UID.V.SetModelName()
	}
	if s.Links.E {
		for i := range s.Links.V {
			s.Links.V[i].SetModelName()
		}
	}
	if s.ArchetypeDetails.E {
		s.ArchetypeDetails.V.SetModelName()
	}
	if s.FeederAudit.E {
		s.FeederAudit.V.SetModelName()
	}
	if s.Items.E {
		for i := range s.Items.V {
			s.Items.V[i].SetModelName()
		}
	}
}

func (s *SECTION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if s.Type_.E && s.Type_.V != SECTION_MODEL_NAME {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          SECTION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + SECTION_MODEL_NAME,
			Recommendation: "Set _type to " + SECTION_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, s.Name.Validate(attrPath).Errs...)

	// Validate uid
	if s.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, s.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if s.Links.E {
		for i := range s.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, s.Links.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if s.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, s.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate feeder_audit
	if s.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, s.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate items
	if s.Items.E {
		for i := range s.Items.V {
			attrPath = fmt.Sprintf("%s.items[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, s.Items.V[i].Validate(attrPath).Errs...)
		}
	}

	return validateErr
}
