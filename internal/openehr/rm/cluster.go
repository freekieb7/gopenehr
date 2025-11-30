package rm

import (
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const CLUSTER_TYPE string = "CLUSTER"

type CLUSTER struct {
	Type_            utils.Optional[string]          `json:"_type,omitzero"`
	Name             DvTextUnion                     `json:"name"`
	ArchetypeNodeID  string                          `json:"archetype_node_id"`
	UID              utils.Optional[UIDBasedIDUnion] `json:"uid,omitzero"`
	Links            utils.Optional[[]LINK]          `json:"links,omitzero"`
	ArchetypeDetails utils.Optional[ARCHETYPED]      `json:"archetype_details,omitzero"`
	FeederAudit      utils.Optional[FEEDER_AUDIT]    `json:"feeder_audit,omitzero"`
	Items            []ItemUnion                     `json:"items"`
}

func (c *CLUSTER) SetModelName() {
	c.Type_ = utils.Some(CLUSTER_TYPE)
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

func (c *CLUSTER) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if c.Type_.E && c.Type_.V != CLUSTER_TYPE {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:   CLUSTER_TYPE,
			Path:    attrPath,
			Message: "invalid _type value",
		})
	}

	// Validate name
	attrPath = path + ".name"
	validateErr.Errs = append(validateErr.Errs, c.Name.Validate(attrPath).Errs...)

	// Validate uid
	if c.UID.E {
		attrPath = path + ".uid"
		validateErr.Errs = append(validateErr.Errs, c.UID.V.Validate(attrPath).Errs...)
	}

	// Validate links
	if c.Links.E {
		for i := range c.Links.V {
			attrPath = fmt.Sprintf("%s.links[%d]", path, i)
			validateErr.Errs = append(validateErr.Errs, c.Links.V[i].Validate(attrPath).Errs...)
		}
	}

	// Validate archetype_details
	if c.ArchetypeDetails.E {
		attrPath = path + ".archetype_details"
		validateErr.Errs = append(validateErr.Errs, c.ArchetypeDetails.V.Validate(attrPath).Errs...)
	}

	// Validate feeder_audit
	if c.FeederAudit.E {
		attrPath = path + ".feeder_audit"
		validateErr.Errs = append(validateErr.Errs, c.FeederAudit.V.Validate(attrPath).Errs...)
	}

	// Validate items
	for i := range c.Items {
		attrPath = fmt.Sprintf("%s.items[%d]", path, i)
		validateErr.Errs = append(validateErr.Errs, c.Items[i].Validate(attrPath).Errs...)
	}

	return validateErr
}
