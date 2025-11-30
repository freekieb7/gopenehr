package rm

import (
	"encoding/json"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const CONTENT_ITEM_TYPE = "CONTENT_ITEM"

type ContentItemKind int

const (
	ContentItemKind_Unknown ContentItemKind = iota
	ContentItemKind_Section
	ContentItemKind_AdminEntry
	ContentItemKind_Observation
	ContentItemKind_Evaluation
	ContentItemKind_Instruction
	ContentItemKind_Activity
	ContentItemKind_Action
	ContentItemKind_GenericEntry
)

type ContentItemUnion struct {
	Kind  ContentItemKind
	Value any
}

func (c *ContentItemUnion) SetModelName() {
	switch c.Kind {
	case ContentItemKind_Section:
		c.Value.(*SECTION).SetModelName()
	case ContentItemKind_AdminEntry:
		c.Value.(*ADMIN_ENTRY).SetModelName()
	case ContentItemKind_Observation:
		c.Value.(*OBSERVATION).SetModelName()
	case ContentItemKind_Evaluation:
		c.Value.(*EVALUATION).SetModelName()
	case ContentItemKind_Instruction:
		c.Value.(*INSTRUCTION).SetModelName()
	case ContentItemKind_Activity:
		c.Value.(*ACTIVITY).SetModelName()
	case ContentItemKind_Action:
		c.Value.(*ACTION).SetModelName()
	case ContentItemKind_GenericEntry:
		c.Value.(*GENERIC_ENTRY).SetModelName()
	}
}

func (c *ContentItemUnion) Validate(path string) util.ValidateError {
	switch c.Kind {
	case ContentItemKind_Section:
		return c.Value.(*SECTION).Validate(path)
	case ContentItemKind_AdminEntry:
		return c.Value.(*ADMIN_ENTRY).Validate(path)
	case ContentItemKind_Observation:
		return c.Value.(*OBSERVATION).Validate(path)
	case ContentItemKind_Evaluation:
		return c.Value.(*EVALUATION).Validate(path)
	case ContentItemKind_Instruction:
		return c.Value.(*INSTRUCTION).Validate(path)
	case ContentItemKind_Activity:
		return c.Value.(*ACTIVITY).Validate(path)
	case ContentItemKind_Action:
		return c.Value.(*ACTION).Validate(path)
	case ContentItemKind_GenericEntry:
		return c.Value.(*GENERIC_ENTRY).Validate(path)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          CONTENT_ITEM_TYPE,
					Path:           path,
					Message:        "Unknown CONTENT_ITEM value",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}
}

func (c ContentItemUnion) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c *ContentItemUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case SECTION_TYPE:
		c.Kind = ContentItemKind_Section
		c.Value = &SECTION{}
	case ADMIN_ENTRY_TYPE:
		c.Kind = ContentItemKind_AdminEntry
		c.Value = &ADMIN_ENTRY{}
	case OBSERVATION_TYPE:
		c.Kind = ContentItemKind_Observation
		c.Value = &OBSERVATION{}
	case EVALUATION_TYPE:
		c.Kind = ContentItemKind_Evaluation
		c.Value = &EVALUATION{}
	case INSTRUCTION_TYPE:
		c.Kind = ContentItemKind_Instruction
		c.Value = &INSTRUCTION{}
	case ACTIVITY_TYPE:
		c.Kind = ContentItemKind_Activity
		c.Value = &ACTIVITY{}
	case ACTION_TYPE:
		c.Kind = ContentItemKind_Action
		c.Value = &ACTION{}
	case GENERIC_ENTRY_TYPE:
		c.Kind = ContentItemKind_GenericEntry
		c.Value = &GENERIC_ENTRY{}
	default:
		c.Kind = ContentItemKind_Unknown
		return nil
	}

	return json.Unmarshal(data, c.Value)
}
