package rm

import (
	"github.com/bytedance/sonic"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
)

const CONTENT_ITEM_TYPE = "CONTENT_ITEM"

type ContentItemKind int

const (
	CONTENT_ITEM_kind_unknown ContentItemKind = iota
	CONTENT_ITEM_kind_Section
	CONTENT_ITEM_kind_AdminEntry
	CONTENT_ITEM_kind_Observation
	CONTENT_ITEM_kind_Evaluation
	CONTENT_ITEM_kind_Instruction
	CONTENT_ITEM_kind_Activity
	CONTENT_ITEM_kind_Action
	CONTENT_ITEM_kind_GenericEntry
)

type ContentItemUnion struct {
	Kind  ContentItemKind
	Value any
}

func (c *ContentItemUnion) SetModelName() {
	switch c.Kind {
	case CONTENT_ITEM_kind_Section:
		c.Value.(*SECTION).SetModelName()
	case CONTENT_ITEM_kind_AdminEntry:
		c.Value.(*ADMIN_ENTRY).SetModelName()
	case CONTENT_ITEM_kind_Observation:
		c.Value.(*OBSERVATION).SetModelName()
	case CONTENT_ITEM_kind_Evaluation:
		c.Value.(*EVALUATION).SetModelName()
	case CONTENT_ITEM_kind_Instruction:
		c.Value.(*INSTRUCTION).SetModelName()
	case CONTENT_ITEM_kind_Activity:
		c.Value.(*ACTIVITY).SetModelName()
	case CONTENT_ITEM_kind_Action:
		c.Value.(*ACTION).SetModelName()
	case CONTENT_ITEM_kind_GenericEntry:
		c.Value.(*GENERIC_ENTRY).SetModelName()
	}
}

func (c *ContentItemUnion) Validate(path string) util.ValidateError {
	switch c.Kind {
	case CONTENT_ITEM_kind_Section:
		return c.Value.(*SECTION).Validate(path)
	case CONTENT_ITEM_kind_AdminEntry:
		return c.Value.(*ADMIN_ENTRY).Validate(path)
	case CONTENT_ITEM_kind_Observation:
		return c.Value.(*OBSERVATION).Validate(path)
	case CONTENT_ITEM_kind_Evaluation:
		return c.Value.(*EVALUATION).Validate(path)
	case CONTENT_ITEM_kind_Instruction:
		return c.Value.(*INSTRUCTION).Validate(path)
	case CONTENT_ITEM_kind_Activity:
		return c.Value.(*ACTIVITY).Validate(path)
	case CONTENT_ITEM_kind_Action:
		return c.Value.(*ACTION).Validate(path)
	case CONTENT_ITEM_kind_GenericEntry:
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
	return sonic.Marshal(c.Value)
}

func (c *ContentItemUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case SECTION_TYPE:
		c.Kind = CONTENT_ITEM_kind_Section
		c.Value = &SECTION{}
	case ADMIN_ENTRY_TYPE:
		c.Kind = CONTENT_ITEM_kind_AdminEntry
		c.Value = &ADMIN_ENTRY{}
	case OBSERVATION_TYPE:
		c.Kind = CONTENT_ITEM_kind_Observation
		c.Value = &OBSERVATION{}
	case EVALUATION_TYPE:
		c.Kind = CONTENT_ITEM_kind_Evaluation
		c.Value = &EVALUATION{}
	case INSTRUCTION_TYPE:
		c.Kind = CONTENT_ITEM_kind_Instruction
		c.Value = &INSTRUCTION{}
	case ACTIVITY_TYPE:
		c.Kind = CONTENT_ITEM_kind_Activity
		c.Value = &ACTIVITY{}
	case ACTION_TYPE:
		c.Kind = CONTENT_ITEM_kind_Action
		c.Value = &ACTION{}
	case GENERIC_ENTRY_TYPE:
		c.Kind = CONTENT_ITEM_kind_GenericEntry
		c.Value = &GENERIC_ENTRY{}
	default:
		c.Kind = CONTENT_ITEM_kind_unknown
		return nil
	}

	return sonic.Unmarshal(data, c.Value)
}

func (c *ContentItemUnion) SECTION() *SECTION {
	if c.Kind != CONTENT_ITEM_kind_Section {
		return nil
	}
	return c.Value.(*SECTION)
}

func (c *ContentItemUnion) ADMIN_ENTRY() *ADMIN_ENTRY {
	if c.Kind != CONTENT_ITEM_kind_AdminEntry {
		return nil
	}
	return c.Value.(*ADMIN_ENTRY)
}

func (c *ContentItemUnion) OBSERVATION() *OBSERVATION {
	if c.Kind != CONTENT_ITEM_kind_Observation {
		return nil
	}
	return c.Value.(*OBSERVATION)
}

func (c *ContentItemUnion) EVALUATION() *EVALUATION {
	if c.Kind != CONTENT_ITEM_kind_Evaluation {
		return nil
	}
	return c.Value.(*EVALUATION)
}

func (c *ContentItemUnion) INSTRUCTION() *INSTRUCTION {
	if c.Kind != CONTENT_ITEM_kind_Instruction {
		return nil
	}
	return c.Value.(*INSTRUCTION)
}

func (c *ContentItemUnion) ACTIVITY() *ACTIVITY {
	if c.Kind != CONTENT_ITEM_kind_Activity {
		return nil
	}
	return c.Value.(*ACTIVITY)
}

func (c *ContentItemUnion) ACTION() *ACTION {
	if c.Kind != CONTENT_ITEM_kind_Action {
		return nil
	}
	return c.Value.(*ACTION)
}

func (c *ContentItemUnion) GENERIC_ENTRY() *GENERIC_ENTRY {
	if c.Kind != CONTENT_ITEM_kind_GenericEntry {
		return nil
	}
	return c.Value.(*GENERIC_ENTRY)
}
