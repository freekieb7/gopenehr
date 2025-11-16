package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const CONTACT_MODEL_NAME string = "CONTACT"

type CONTACT struct {
	Type_            util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Addresses        []ADDRESS                     `json:"addresses"`
	TimeValidity     util.Optional[DV_INTERVAL]    `json:"time_validity,omitzero"`
}

func (c *CONTACT) SetModelName() {
	c.Type_ = util.Some(CONTACT_MODEL_NAME)
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
	for i := range c.Addresses {
		c.Addresses[i].SetModelName()
	}
	if c.TimeValidity.E {
		c.TimeValidity.V.SetModelName()
	}
}

func (c *CONTACT) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if c.Type_.E && c.Type_.V != CONTACT_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:          CONTACT_MODEL_NAME,
			Path:           attrPath,
			Message:        fmt.Sprintf("invalid %s _type field: %s", CONTACT_MODEL_NAME, c.Type_.V),
			Recommendation: fmt.Sprintf("Ensure _type field is set to '%s'", CONTACT_MODEL_NAME),
		})
	}

	// Validate Name
	attrPath = path + ".name"
	errors = append(errors, c.Name.Validate(attrPath)...)

	// Validate ArchetypeNodeID
	attrPath = path + ".archetype_node_id"
	if c.ArchetypeNodeID == "" {
		errors = append(errors, util.ValidationError{
			Model:          CONTACT_MODEL_NAME,
			Path:           attrPath,
			Message:        "archetype_node_id is required",
			Recommendation: "Ensure archetype_node_id is not empty",
		})
	}

	// Validate UID
	if c.UID.E {
		attrPath = path + ".uid"
		errors = append(errors, c.UID.V.Validate(attrPath)...)
	}

	// Validate Links
	if c.Links.E {
		for i := range c.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			errors = append(errors, c.Links.V[i].Validate(attrPath)...)
		}
	}

	// Validate ArchetypeDetails
	if c.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errors = append(errors, c.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate FeederAudit
	if c.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errors = append(errors, c.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate Addresses
	for i := range c.Addresses {
		attrPath = fmt.Sprintf("%s.addresses[%d]", path, i)
		errors = append(errors, c.Addresses[i].Validate(attrPath)...)
	}

	// Validate TimeValidity
	if c.TimeValidity.E {
		attrPath = path + ".time_validity"
		errors = append(errors, c.TimeValidity.V.Validate(attrPath)...)
	}

	return errors
}
