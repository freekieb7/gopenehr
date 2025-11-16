package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const COMPOSITION_MODEL_NAME string = "COMPOSITION"

type COMPOSITION struct {
	Type_            util.Optional[string]           `json:"_type,omitzero"`
	Name             X_DV_TEXT                       `json:"name"`
	ArchetypeNodeID  string                          `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID]   `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]           `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]       `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]     `json:"feeder_audit,omitzero"`
	Language         CODE_PHRASE                     `json:"language"`
	Territory        CODE_PHRASE                     `json:"territory"`
	Category         DV_CODED_TEXT                   `json:"category"`
	Context          util.Optional[EVENT_CONTEXT]    `json:"context,omitzero"`
	Composer         X_PARTY_PROXY                   `json:"composer"`
	Content          util.Optional[[]X_CONTENT_ITEM] `json:"content,omitzero"`
}

func (c *COMPOSITION) isVersionModel() {}

func (c *COMPOSITION) SetModelName() {
	c.Type_ = util.Some(COMPOSITION_MODEL_NAME)
	c.Name.SetModelName()
	if c.UID.E {
		c.UID.V.SetModelName()
	}
	if c.Links.E {
		for i := range c.Links.V {
			c.Links.V[i].SetModelName()
		}
	}
	if c.ArchetypeDetails.E {
		c.ArchetypeDetails.V.SetModelName()
	}
	if c.FeederAudit.E {
		c.FeederAudit.V.SetModelName()
	}
	c.Language.SetModelName()
	c.Territory.SetModelName()
	c.Category.SetModelName()
	if c.Context.E {
		c.Context.V.SetModelName()
	}
	c.Composer.SetModelName()
	if c.Content.E {
		for i := range c.Content.V {
			c.Content.V[i].SetModelName()
		}
	}
}

func (c *COMPOSITION) Validate(path string) []util.ValidationError {
	var errs []util.ValidationError
	var attrPath string

	// Validate _type
	if c.Type_.E && c.Type_.V != COMPOSITION_MODEL_NAME {
		attrPath = path + "._type"
		errs = append(errs, util.ValidationError{
			Model:          COMPOSITION_MODEL_NAME,
			Path:           attrPath,
			Message:        "_type must be " + COMPOSITION_MODEL_NAME,
			Recommendation: "Set _type to " + COMPOSITION_MODEL_NAME,
		})
	}

	// Validate name
	attrPath = path + ".name"
	errs = append(errs, c.Name.Validate(attrPath)...)

	// Validate uid
	if c.UID.E {
		attrPath = path + ".uid"
		errs = append(errs, c.UID.V.Validate(attrPath)...)
	}

	// Validate links
	if c.Links.E {
		for i := range c.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			errs = append(errs, c.Links.V[i].Validate(attrPath)...)
		}
	}

	// Validate archetype_details
	if c.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errs = append(errs, c.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate feeder_audit
	if c.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errs = append(errs, c.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate language
	attrPath = path + ".language"
	if !terminology.IsValidLanguageTerminologyID(c.Language.TerminologyID.Value) {
		attrPath = path + ".terminology_id.value"
		errs = append(errs, util.ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid language terminology ID: %s", c.Language.TerminologyID.Value),
			Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
		})
	}

	if !terminology.IsValidLanguageCode(c.Language.CodeString) {
		attrPath = path + ".code_string"
		errs = append(errs, util.ValidationError{
			Model:          DV_TEXT_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid language code: %s", c.Language.CodeString),
			Recommendation: "Ensure language field is a known ISO 639-1 or ISO 639-2 language code",
		})
	}
	errs = append(errs, c.Language.Validate(attrPath)...)

	// Validate territory
	attrPath = path + ".territory"
	errs = append(errs, c.Territory.Validate(attrPath)...)

	// Validate category
	attrPath = path + ".category"
	errs = append(errs, c.Category.Validate(attrPath)...)

	// Validate context
	if c.Context.E {
		attrPath = path + ".context"
		errs = append(errs, c.Context.V.Validate(attrPath)...)
	}

	// Validate composer
	attrPath = path + ".composer"
	errs = append(errs, c.Composer.Validate(attrPath)...)

	// Validate content
	if c.Content.E {
		for i := range c.Content.V {
			attrPath = fmt.Sprintf("%s.content[%d]", path, i)
			errs = append(errs, c.Content.V[i].Validate(attrPath)...)
		}
	}

	return errs
}
