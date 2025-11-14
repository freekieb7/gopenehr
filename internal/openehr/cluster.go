package openehr

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const CLUSTER_MODEL_NAME string = "CLUSTER"

type CLUSTER struct {
	Type_            util.Optional[string]         `json:"_type,omitzero"`
	Name             X_DV_TEXT                     `json:"name"`
	ArchetypeNodeID  string                        `json:"archetype_node_id"`
	UID              util.Optional[X_UID_BASED_ID] `json:"uid,omitzero"`
	Links            util.Optional[[]LINK]         `json:"links,omitzero"`
	ArchetypeDetails util.Optional[ARCHETYPED]     `json:"archetype_details,omitzero"`
	FeederAudit      util.Optional[FEEDER_AUDIT]   `json:"feeder_audit,omitzero"`
	Items            []X_ITEM                      `json:"items"`
}

func (c CLUSTER) isItemModel() {}

func (c CLUSTER) HasModelName() bool {
	return c.Type_.E
}

func (c *CLUSTER) SetModelName() {
	c.Type_ = util.Some(CLUSTER_MODEL_NAME)
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
	for i := range c.Items {
		c.Items[i].SetModelName()
	}
}

func (c CLUSTER) Validate(path string) []util.ValidationError {
	var errors []util.ValidationError
	var attrPath string

	// Validate _type
	if c.Type_.E && c.Type_.V != CLUSTER_MODEL_NAME {
		attrPath = path + "._type"
		errors = append(errors, util.ValidationError{
			Model:   CLUSTER_MODEL_NAME,
			Path:    attrPath,
			Message: "invalid _type value",
		})
	}

	// Validate name
	attrPath = path + ".name"
	errors = append(errors, c.Name.Validate(attrPath)...)

	// Validate uid
	if c.UID.E {
		attrPath = path + ".uid"
		errors = append(errors, c.UID.V.Validate(attrPath)...)
	}

	// Validate links
	if c.Links.E {
		for i := range c.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			errors = append(errors, c.Links.V[i].Validate(attrPath)...)
		}
	}

	// Validate archetype_details
	if c.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		errors = append(errors, c.ArchetypeDetails.V.Validate(attrPath)...)
	}

	// Validate feeder_audit
	if c.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		errors = append(errors, c.FeederAudit.V.Validate(attrPath)...)
	}

	// Validate items
	for i := range c.Items {
		attrPath = fmt.Sprintf("%s.items[%d]", path, i)
		errors = append(errors, c.Items[i].Validate(attrPath)...)
	}

	return errors
}
