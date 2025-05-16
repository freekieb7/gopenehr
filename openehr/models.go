package openehr

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("value not found")
var ErrBadType = errors.New("parse error")

type Option[T any] struct {
	Value *T
	Some  bool
}

func Some[T any](v T) Option[T] {
	return Option[T]{Value: &v, Some: true}
}

func None[T any]() Option[T] {
	return Option[T]{Value: nil, Some: false}
}

func (o Option[T]) IsSome() bool { return o.Some }
func (o Option[T]) Unwrap() T    { return *o.Value }

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		o.Some = false
		o.Value = nil
	} else {
		var v T
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		o.Value = &v
		o.Some = true
	}

	return nil
}

func (o *Option[T]) MarshalJSON() ([]byte, error) {
	if o.Some {
		return json.Marshal(o.Unwrap())
	}

	return []byte("null"), nil
}

func (o Option[T]) Marshal() ([]byte, error) {
	if o.Some {
		return Marshal(o.Unwrap())
	}

	return []byte{}, nil
}

// -----------------------------------
// EHR
// -----------------------------------

type EHR struct {
	Type_         Option[string]         `json:"_type"`
	SystemId      Option[HIER_OBJECT_ID] `json:"system_id"`
	EhrId         HIER_OBJECT_ID         `json:"ehr_id"`
	Contributions Option[[]OBJECT_REF]   `json:"contributions"`
	EhrStatus     OBJECT_REF             `json:"ehr_status"`
	EhrAccess     OBJECT_REF             `json:"ehr_access"`
	Compositions  Option[[]OBJECT_REF]   `json:"compositions"`
	Directory     Option[OBJECT_REF]     `json:"directory"`
	TimeCreated   DV_DATE_TIME           `json:"time_created"`
	Folders       Option[[]OBJECT_REF]   `json:"folders"`
}

type VERSIONED_EHR_ACCESS struct {
	Type_       Option[string] `json:"_type"`
	Uid         HIER_OBJECT_ID `json:"uid"`
	OwnerId     OBJECT_REF     `json:"owner_id"`
	TimeCreated DV_DATE_TIME   `json:"time_created"`
}

type EHR_ACCESS struct {
	Type_            Option[string]       `json:"_type"`
	Name             DV_TEXT              `json:"name"`
	ArchetypeNodeId  string               `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `json:"uid"`
	Links            Option[[]LINK]       `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `json:"feeder_audit"`
}

type VERSIONED_EHR_STATUS struct {
	Type_       Option[string] `json:"_type"`
	Uid         HIER_OBJECT_ID `json:"uid"`
	OwnerId     OBJECT_REF     `json:"owner_id"`
	TimeCreated DV_DATE_TIME   `json:"time_created"`
}

type EHR_STATUS struct {
	Type_            Option[string]         `json:"_type"`
	Name             DV_TEXT                `json:"name"`
	ArchetypeNodeId  string                 `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `json:"uid"`
	Links            Option[[]LINK]         `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `json:"feeder_audit"`
	Subject          PARTY_SELF             `json:"subject"`
	IsQueryable      bool                   `json:"is_queryable"`
	IsModifiable     bool                   `json:"is_modifiable"`
	OtherDetails     Option[ITEM_STRUCTURE] `json:"other_details"`
}

type VERSIONED_COMPOSITION struct {
	Type_       Option[string] `json:"_type"`
	Uid         HIER_OBJECT_ID `json:"uid"`
	OwnerId     OBJECT_REF     `json:"owner_id"`
	TimeCreated DV_DATE_TIME   `json:"time_created"`
}

type COMPOSITION struct {
	Type_            Option[string]         `json:"_type"`
	ArchetypeNodeId  string                 `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `json:"uid"`
	Links            Option[[]LINK]         `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `json:"feeder_audit"`
	Language         CODE_PHRASE            `json:"language"`
	Territory        CODE_PHRASE            `json:"territory"`
	Category         DV_CODED_TEXT          `json:"category"`
	Context          Option[EVENT_CONTEXT]  `json:"context"`
	Composer         PARTY_PROXY            `json:"composer"`
	Content          Option[[]CONTENT_ITEM] `json:"content"`
}

type EVENT_CONTEXT struct {
	Type_              Option[string]           `json:"_type"`
	StartTime          DV_DATE_TIME             `json:"start_time"`
	EndTime            DV_DATE_TIME             `json:"end_time"`
	Location           Option[string]           `json:"location"`
	Setting            DV_CODED_TEXT            `json:"setting"`
	OtherContext       Option[ITEM_STRUCTURE]   `json:"other_context"`
	HealthCareFacility Option[PARTY_IDENTIFIED] `json:"health_care_facility"`
	Participations     Option[[]PARTICIPATION]  `json:"participations"`
}

type ContentItemType string

const (
	CONTENT_ITEM_TYPE_SECTION     ContentItemType = "SECTION"
	CONTENT_ITEM_TYPE_ADMIN_ENTRY ContentItemType = "ADMIN_ENTRY"
	CONTENT_ITEM_TYPE_OBSERVATION ContentItemType = "OBSERVATION"
	CONTENT_ITEM_TYPE_EVALUATION  ContentItemType = "EVALUATION"
	CONTENT_ITEM_TYPE_INSTRUCTION ContentItemType = "INSTRUCTION"
	CONTENT_ITEM_TYPE_ACTIVITY    ContentItemType = "ACTIVITY"
	CONTENT_ITEM_TYPE_ACTION      ContentItemType = "ACTION"
)

type CONTENT_ITEM struct {
	Type  ContentItemType
	Value any
}

func (c *CONTENT_ITEM) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t ContentItemType
	switch ContentItemType(typStr) {
	case CONTENT_ITEM_TYPE_SECTION:
		{
			v = new(SECTION)
			t = CONTENT_ITEM_TYPE_SECTION
		}
	case CONTENT_ITEM_TYPE_ADMIN_ENTRY:
		{
			v = new(ADMIN_ENTRY)
			t = CONTENT_ITEM_TYPE_ADMIN_ENTRY
		}
	case CONTENT_ITEM_TYPE_OBSERVATION:
		{
			v = new(OBSERVATION)
			t = CONTENT_ITEM_TYPE_OBSERVATION
		}
	case CONTENT_ITEM_TYPE_EVALUATION:
		{
			v = new(EVALUATION)
			t = CONTENT_ITEM_TYPE_EVALUATION
		}
	case CONTENT_ITEM_TYPE_INSTRUCTION:
		{
			v = new(INSTRUCTION)
			t = CONTENT_ITEM_TYPE_INSTRUCTION
		}
	case CONTENT_ITEM_TYPE_ACTIVITY:
		{
			v = new(ACTIVITY)
			t = CONTENT_ITEM_TYPE_ACTIVITY
		}
	case CONTENT_ITEM_TYPE_ACTION:
		{
			v = new(ACTION)
			t = CONTENT_ITEM_TYPE_ACTION
		}
	default:
		{
			return fmt.Errorf("CONTENT_ITEM unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *CONTENT_ITEM) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c CONTENT_ITEM) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type SECTION struct {
	Type_            Option[string]         `json:"_type"`
	Name             DV_TEXT                `json:"name"`
	ArchetypeNodeId  string                 `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `json:"uid"`
	Links            Option[[]LINK]         `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `json:"feeder_audit"`
	Items            Option[[]CONTENT_ITEM] `json:"items"`
}

type EntryType string

const (
	ENTRY_TYPE_ADMIN_ENTRY EntryType = "ADMIN_ENTRY"
	ENTRY_TYPE_OBSERVATION EntryType = "OBSERVATION"
	ENTRY_TYPE_EVALUATION  EntryType = "EVALUATION"
	ENTRY_TYPE_INSTRUCTION EntryType = "INSTRUCTION"
	ENTRY_TYPE_ACTIVITY    EntryType = "ACTIVITY"
	ENTRY_TYPE_ACTION      EntryType = "ACTION"
)

type ENTRY struct {
	Type  EntryType
	Value any
}

func (c *ENTRY) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t EntryType
	switch EntryType(typStr) {
	case ENTRY_TYPE_ADMIN_ENTRY:
		{
			v = new(ADMIN_ENTRY)
			t = ENTRY_TYPE_ADMIN_ENTRY
		}
	case ENTRY_TYPE_OBSERVATION:
		{
			v = new(OBSERVATION)
			t = ENTRY_TYPE_OBSERVATION
		}
	case ENTRY_TYPE_EVALUATION:
		{
			v = new(EVALUATION)
			t = ENTRY_TYPE_EVALUATION
		}
	case ENTRY_TYPE_INSTRUCTION:
		{
			v = new(INSTRUCTION)
			t = ENTRY_TYPE_INSTRUCTION
		}
	case ENTRY_TYPE_ACTIVITY:
		{
			v = new(ACTIVITY)
			t = ENTRY_TYPE_ACTIVITY
		}
	case ENTRY_TYPE_ACTION:
		{
			v = new(ACTION)
			t = ENTRY_TYPE_ACTION
		}
	default:
		{
			return fmt.Errorf("ENTRY unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *ENTRY) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c ENTRY) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type ADMIN_ENTRY struct {
	Type_               Option[string]          `json:"_type"`
	Name                DV_TEXT                 `json:"name"`
	ArchetypeNodeId     string                  `json:"archetype_node_id"`
	Uid                 Option[UID_BASED_ID]    `json:"uid"`
	Links               Option[[]LINK]          `json:"links"`
	ArchetypeDetails    Option[ARCHETYPED]      `json:"archetype_details"`
	FeederAudit         Option[FEEDER_AUDIT]    `json:"feeder_audit"`
	Language            CODE_PHRASE             `json:"language"`
	Encoding            CODE_PHRASE             `json:"encoding"`
	OtherParticipations Option[[]PARTICIPATION] `json:"other_participations"`
	WorkflowId          Option[OBJECT_REF]      `json:"workflow_id"`
	Subject             PARTY_PROXY             `json:"subject"`
	Provider            Option[PARTY_PROXY]     `json:"provider"`
	Data                ITEM_STRUCTURE          `json:"data"`
}

type CareEntryType string

const (
	CARE_ENTRY_TYPE_OBSERVATION CareEntryType = "OBSERVATION"
	CARE_ENTRY_TYPE_EVALUATION  CareEntryType = "EVALUATION"
	CARE_ENTRY_TYPE_INSTRUCTION CareEntryType = "INSTRUCTION"
	CARE_ENTRY_TYPE_ACTIVITY    CareEntryType = "ACTIVITY"
	CARE_ENTRY_TYPE_ACTION      CareEntryType = "ACTION"
)

type CARE_ENTRY struct {
	Type  CareEntryType
	Value any
}

func (c *CARE_ENTRY) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t CareEntryType
	switch CareEntryType(typStr) {
	case CARE_ENTRY_TYPE_OBSERVATION:
		{
			v = new(OBSERVATION)
			t = CARE_ENTRY_TYPE_OBSERVATION
		}
	case CARE_ENTRY_TYPE_EVALUATION:
		{
			v = new(EVALUATION)
			t = CARE_ENTRY_TYPE_EVALUATION
		}
	case CARE_ENTRY_TYPE_INSTRUCTION:
		{
			v = new(INSTRUCTION)
			t = CARE_ENTRY_TYPE_INSTRUCTION
		}
	case CARE_ENTRY_TYPE_ACTIVITY:
		{
			v = new(ACTIVITY)
			t = CARE_ENTRY_TYPE_ACTIVITY
		}
	case CARE_ENTRY_TYPE_ACTION:
		{
			v = new(ACTION)
			t = CARE_ENTRY_TYPE_ACTION
		}
	default:
		{
			return fmt.Errorf("CARE_ENTRY unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *CARE_ENTRY) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c CARE_ENTRY) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type OBSERVATION struct {
	Type_               Option[string]                  `json:"_type"`
	Name                DV_TEXT                         `json:"name"`
	ArchetypeNodeId     string                          `json:"archetype_node_id"`
	Uid                 Option[UID_BASED_ID]            `json:"uid"`
	Links               Option[[]LINK]                  `json:"links"`
	ArchetypeDetails    Option[ARCHETYPED]              `json:"archetype_details"`
	FeederAudit         Option[FEEDER_AUDIT]            `json:"feeder_audit"`
	Language            CODE_PHRASE                     `json:"language"`
	Encoding            CODE_PHRASE                     `json:"encoding"`
	OtherParticipations Option[[]PARTICIPATION]         `json:"other_participations"`
	WorkflowId          Option[OBJECT_REF]              `json:"workflow_id"`
	Subject             PARTY_PROXY                     `json:"subject"`
	Provider            Option[PARTY_PROXY]             `json:"provider"`
	Protocol            Option[ITEM_STRUCTURE]          `json:"protocol"`
	GuidelineId         Option[OBJECT_REF]              `json:"guideline_id"`
	Data                HISTORY[ITEM_STRUCTURE]         `json:"data"`
	State               Option[HISTORY[ITEM_STRUCTURE]] `json:"state"`
}

type EVALUATION struct {
	Type_               Option[string]          `json:"_type"`
	Name                DV_TEXT                 `json:"name"`
	ArchetypeNodeId     string                  `json:"archetype_node_id"`
	Uid                 Option[UID_BASED_ID]    `json:"uid"`
	Links               Option[[]LINK]          `json:"links"`
	ArchetypeDetails    Option[ARCHETYPED]      `json:"archetype_details"`
	FeederAudit         Option[FEEDER_AUDIT]    `json:"feeder_audit"`
	Language            CODE_PHRASE             `json:"language"`
	Encoding            CODE_PHRASE             `json:"encoding"`
	OtherParticipations Option[[]PARTICIPATION] `json:"other_participations"`
	WorkflowId          Option[OBJECT_REF]      `json:"workflow_id"`
	Subject             PARTY_PROXY             `json:"subject"`
	Provider            Option[PARTY_PROXY]     `json:"provider"`
	Protocol            Option[ITEM_STRUCTURE]  `json:"protocol"`
	GuidelineId         Option[OBJECT_REF]      `json:"guideline_id"`
	Data                ITEM_STRUCTURE          `json:"data"`
}

type INSTRUCTION struct {
	Type_               Option[string]          `json:"_type"`
	Name                DV_TEXT                 `json:"name"`
	ArchetypeNodeId     string                  `json:"archetype_node_id"`
	Uid                 Option[UID_BASED_ID]    `json:"uid"`
	Links               Option[[]LINK]          `json:"links"`
	ArchetypeDetails    Option[ARCHETYPED]      `json:"archetype_details"`
	FeederAudit         Option[FEEDER_AUDIT]    `json:"feeder_audit"`
	Language            CODE_PHRASE             `json:"language"`
	Encoding            CODE_PHRASE             `json:"encoding"`
	OtherParticipations Option[[]PARTICIPATION] `json:"other_participations"`
	WorkflowId          Option[OBJECT_REF]      `json:"workflow_id"`
	Subject             PARTY_PROXY             `json:"subject"`
	Provider            Option[PARTY_PROXY]     `json:"provider"`
	Protocol            Option[ITEM_STRUCTURE]  `json:"protocol"`
	GuidelineId         Option[OBJECT_REF]      `json:"guideline_id"`
	Narative            DV_TEXT                 `json:"narative"`
	ExpiryTime          Option[DV_DATE_TIME]    `json:"expiry_time"`
	WfDefinition        Option[DV_PARSABLE]     `json:"wf_definition"`
	Activities          Option[[]ACTIVITY]      `json:"activities"`
}

type ACTIVITY struct {
	Type_             Option[string]       `json:"_type"`
	Name              DV_TEXT              `json:"name"`
	ArchetypeNodeId   string               `json:"archetype_node_id"`
	Uid               Option[UID_BASED_ID] `json:"uid"`
	Links             Option[[]LINK]       `json:"links"`
	ArchetypeDetails  Option[ARCHETYPED]   `json:"archetype_details"`
	FeederAudit       Option[FEEDER_AUDIT] `json:"feeder_audit"`
	Timing            Option[DV_PARSABLE]  `json:"timing"`
	ActionArchetypeId string               `json:"action_archetype_id"`
	Description       ITEM_STRUCTURE       `json:"description"`
}

type ACTION struct {
	Type_               Option[string]              `json:"_type"`
	Name                DV_TEXT                     `json:"name"`
	ArchetypeNodeId     string                      `json:"archetype_node_id"`
	Uid                 Option[UID_BASED_ID]        `json:"uid"`
	Links               Option[[]LINK]              `json:"links"`
	ArchetypeDetails    Option[ARCHETYPED]          `json:"archetype_details"`
	FeederAudit         Option[FEEDER_AUDIT]        `json:"feeder_audit"`
	Language            CODE_PHRASE                 `json:"language"`
	Encoding            CODE_PHRASE                 `json:"encoding"`
	OtherParticipations Option[[]PARTICIPATION]     `json:"other_participations"`
	WorkflowId          Option[OBJECT_REF]          `json:"workflow_id"`
	Subject             PARTY_PROXY                 `json:"subject"`
	Provider            Option[PARTY_PROXY]         `json:"provider"`
	Protocol            Option[ITEM_STRUCTURE]      `json:"protocol"`
	GuidelineId         Option[OBJECT_REF]          `json:"guideline_id"`
	Time                DV_DATE_TIME                `json:"time"`
	IsmTransition       ISM_TRANSITION              `json:"ism_transition"`
	InstructionDetails  Option[INSTRUCTION_DETAILS] `json:"instruction_details"`
	Description         ITEM_STRUCTURE              `json:"description"`
}

type INSTRUCTION_DETAILS struct {
	Type_         Option[string]         `json:"_type"`
	InstructionId LOCATABLE_REF          `json:"instruction_id"`
	ActivityId    string                 `json:"activity"`
	WfDetails     Option[ITEM_STRUCTURE] `json:"wf_details"`
}

type ISM_TRANSITION struct {
	Type_        Option[string]        `json:"_type"`
	CurrentState DV_CODED_TEXT         `json:"current_state"`
	Transition   Option[DV_CODED_TEXT] `json:"transition"`
	CareflowStep Option[DV_CODED_TEXT] `json:"cateflow_step"`
	Reason       Option[DV_TEXT]       `json:"reason"`
}

// -----------------------------------
// COMMON
// -----------------------------------

type ARCHETYPED struct {
	Type_       Option[string]      `json:"_type"`
	ArchetypeId ARCHETYPE_ID        `json:"archetype_id"`
	TemplateId  Option[TEMPLATE_ID] `json:"template_id"`
	RmVersion   string              `json:"rm_version"`
}

type LINK struct {
	Type_   Option[string] `json:"_type"`
	Meaning DV_TEXT        `json:"meaning"`
	Type    DV_TEXT        `json:"type"`
	Target  DV_EHR_URI     `json:"target"`
}

type FEEDER_AUDIT struct {
	Type_                    Option[FEEDER_AUDIT]         `json:"_type"`
	OriginatingSystemItemIds Option[[]DV_IDENTIFIER]      `json:"originating_system_item_ids"`
	FeederSystemItemIds      Option[[]DV_IDENTIFIER]      `json:"feeder_system_item_ids"`
	OriginalContent          Option[DV_ENCAPSULATED]      `json:"original_content"`
	OriginatingSystemAudit   FEEDER_AUDIT_DETAILS         `json:"originating_system_audit"`
	FeederSystemAudit        Option[FEEDER_AUDIT_DETAILS] `json:"feeder_system_audit"`
}

type FEEDER_AUDIT_DETAILS struct {
	Type_        Option[string]           `json:"_type"`
	SystemId     string                   `json:"system_id"`
	Location     Option[PARTY_IDENTIFIED] `json:"location"`
	Subject      Option[PARTY_PROXY]      `json:"subject"`
	Provider     Option[PARTY_IDENTIFIED] `json:"provider"`
	Time         Option[DV_DATE_TIME]     `json:"time"`
	VersionId    Option[string]           `json:"version_id"`
	OtherDetails Option[ITEM_STRUCTURE]   `json:"other_details"`
}

type PartyProxyType string

const (
	PARTY_PROXY_TYPE_PARTY_SELF       PartyProxyType = "PARTY_SELF"
	PARTY_PROXY_TYPE_PARTY_IDENTIFIED PartyProxyType = "PARTY_IDENTIFIED"
	PARTY_PROXY_TYPE_PARTY_RELATED    PartyProxyType = "PARTY_RELATED"
)

type PARTY_PROXY struct {
	Type  PartyProxyType
	Value any
}

func (c *PARTY_PROXY) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t PartyProxyType
	switch PartyProxyType(typStr) {
	case PARTY_PROXY_TYPE_PARTY_SELF:
		{
			v = new(PARTY_SELF)
			t = PARTY_PROXY_TYPE_PARTY_SELF
		}
	case PARTY_PROXY_TYPE_PARTY_IDENTIFIED:
		{
			v = new(PARTY_IDENTIFIED)
			t = PARTY_PROXY_TYPE_PARTY_IDENTIFIED
		}
	case PARTY_PROXY_TYPE_PARTY_RELATED:
		{
			v = new(PARTY_RELATED)
			t = PARTY_PROXY_TYPE_PARTY_RELATED
		}
	default:
		{
			return fmt.Errorf("PARTY_PROXY unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *PARTY_PROXY) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c PARTY_PROXY) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type PARTY_SELF struct {
	Type_       Option[string]    `json:"_type"`
	ExternalRef Option[PARTY_REF] `json:"external_ref"`
}

type PARTY_IDENTIFIED struct {
	Type_       Option[string]          `json:"_type"`
	ExternalRef Option[PARTY_REF]       `json:"external_ref"`
	Name        Option[string]          `json:"name"`
	Identifiers Option[[]DV_IDENTIFIER] `json:"identifiers"`
}

type PARTY_RELATED struct {
	Type_        Option[string]          `json:"_type"`
	ExternalRef  Option[PARTY_REF]       `json:"external_ref"`
	Name         Option[string]          `json:"name"`
	Identifiers  Option[[]DV_IDENTIFIER] `json:"identifiers"`
	Relationship DV_CODED_TEXT           `json:"relationship"`
}

type PARTICIPATION struct {
	Type_     Option[string]        `json:"_type"`
	Function  DV_TEXT               `json:"function"`
	Mode      Option[DV_CODED_TEXT] `json:"mode"`
	Performer PARTY_PROXY           `json:"performer"`
	Time      Option[DV_INTERVAL]   `json:"time"`
}

// pub const AUDIT_DETAILS = struct {};
// pub const ATTESTATION = struct {};
// pub const REVISION_HISTORY = struct {};
// pub const REVISION_HISTORY_ITEM = struct {};
// pub const VERSIONED_FOLDER = struct {};
// pub const FOLDER = struct {};
// pub const VERSIONED_OBJECT = struct {};
// pub const VERSION = struct {};
// pub const ORIGINAL_VERSION = struct {};
// pub const IMPORTED_VERSION = struct {};
// pub const CONTRIBUTION = struct {};
// pub const AUTHORED_RESOURCE = struct {};
// pub const TRANSLATION_DETAILS = struct {};
// pub const RESOURCE_DESCRIPTION = struct {};
// pub const RESOURCE_DESCRIPTION_ITEM = struct {};

// -----------------------------------
// DATA_STRUCTURES
// -----------------------------------

type ItemStructureType string

const (
	ITEM_STRUCTURE_TYPE_ITEM_SINGLE ItemStructureType = "ITEM_SINGLE"
	ITEM_STRUCTURE_TYPE_ITEM_LIST   ItemStructureType = "ITEM_LIST"
	ITEM_STRUCTURE_TYPE_ITEM_TABLE  ItemStructureType = "ITEM_TABLE"
	ITEM_STRUCTURE_TYPE_ITEM_TREE   ItemStructureType = "ITEM_TREE"
)

type ITEM_STRUCTURE struct {
	Type  ItemStructureType
	Value any
}

func (c *ITEM_STRUCTURE) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t ItemStructureType
	switch ItemStructureType(typStr) {
	case ITEM_STRUCTURE_TYPE_ITEM_SINGLE:
		{
			v = new(ITEM_SINGLE)
			t = ITEM_STRUCTURE_TYPE_ITEM_SINGLE
		}
	case ITEM_STRUCTURE_TYPE_ITEM_LIST:
		{
			v = new(ITEM_LIST)
			t = ITEM_STRUCTURE_TYPE_ITEM_LIST
		}
	case ITEM_STRUCTURE_TYPE_ITEM_TABLE:
		{
			v = new(ITEM_TABLE)
			t = ITEM_STRUCTURE_TYPE_ITEM_TABLE
		}
	case ITEM_STRUCTURE_TYPE_ITEM_TREE:
		{
			v = new(ITEM_TREE)
			t = ITEM_STRUCTURE_TYPE_ITEM_TREE
		}
	default:
		{
			return fmt.Errorf("ITEM_STRUCTURE unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *ITEM_STRUCTURE) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c ITEM_STRUCTURE) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type ITEM_SINGLE struct {
	Type_            Option[string]       `json:"_type"`
	Name             DV_TEXT              `json:"name"`
	ArchetypeNodeId  string               `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `json:"uid"`
	Links            Option[LINK]         `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `json:"feeder_audit"`
	Item             ELEMENT              `json:"item"`
}

type ITEM_LIST struct {
	Type_            Option[DV_TEXT]      `json:"_type"`
	ArchetypeNodeId  string               `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `json:"uid"`
	Links            Option[[]LINK]       `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `json:"feeder_audit"`
	Items            Option[[]ELEMENT]    `json:"items"`
}

type ITEM_TABLE struct {
	Type_            Option[DV_TEXT]      `json:"_type"`
	ArchetypeNodeId  string               `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `json:"uid"`
	Links            Option[[]LINK]       `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `json:"feeder_audit"`
	Rows             Option[[]CLUSTER]    `json:"rows"`
}

type ITEM_TREE struct {
	Type_            Option[DV_TEXT]      `json:"_type"`
	ArchetypeNodeId  string               `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `json:"uid"`
	Links            Option[[]LINK]       `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `json:"feeder_audit"`
	Items            Option[[]ITEM]       `json:"items"`
}

type ItemType string

const (
	ITEM_TYPE_CLUSTER ItemType = "CLUSTER"
	ITEM_TYPE_ELEMENT ItemType = "ELEMENT"
)

type ITEM struct {
	Type  ItemType
	Value any
}

func (c *ITEM) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t ItemType
	switch ItemType(typStr) {
	case ITEM_TYPE_CLUSTER:
		{
			v = new(CLUSTER)
			t = ITEM_TYPE_CLUSTER
		}
	case ITEM_TYPE_ELEMENT:
		{
			v = new(ELEMENT)
			t = ITEM_TYPE_ELEMENT
		}
	default:
		{
			return fmt.Errorf("ITEM unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *ITEM) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c ITEM) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type CLUSTER struct {
	Type_            Option[string]       `json:"_type"`
	Name             DV_TEXT              `json:"name"`
	ArchetypeNodeId  string               `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `json:"uid"`
	Links            Option[[]LINK]       `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `json:"feeder_audit"`
	Items            []ITEM               `json:"items"`
}

type ELEMENT struct {
	Type_            Option[string]        `json:"_type"`
	Name             DV_TEXT               `json:"name"`
	ArchetypeNodeId  string                `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]  `json:"uid"`
	Links            Option[[]LINK]        `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]    `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]  `json:"feeder_audit"`
	NullFlavour      Option[DV_CODED_TEXT] `json:"null_flavour"`
	Value            Option[DATA_VALUE]    `json:"value"`
	NullReason       Option[DV_TEXT]       `json:"null_reason"`
}

type HISTORY[T any] struct {
	Type_            Option[string]         `json:"_type"`
	Name             DV_TEXT                `json:"name"`
	ArchetypeNodeId  string                 `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `json:"uid"`
	Links            Option[[]LINK]         `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `json:"feeder_audit"`
	Origin           DV_DATE_TIME           `json:"origin"`
	Period           Option[DV_DURATION]    `json:"period"`
	Duration         Option[DV_DURATION]    `json:"duration"`
	Summary          Option[ITEM_STRUCTURE] `json:"summary"`
	Events           Option[[]EVENT[T]]     `json:"events"`
}

type EventType string

const (
	EVENT_TYPE_POINT_EVENT    EventType = "POINT_EVENT"
	EVENT_TYPE_INTERVAL_EVENT EventType = "INTERVAL_EVENT"
)

type EVENT[T any] struct {
	Type  EventType
	Value any
}

func (c *EVENT[T]) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t EventType
	switch EventType(typStr) {
	case EVENT_TYPE_POINT_EVENT:
		{
			v = new(POINT_EVENT[T])
			t = EVENT_TYPE_POINT_EVENT
		}
	case EVENT_TYPE_INTERVAL_EVENT:
		{
			v = new(INTERVAL_EVENT[T])
			t = EVENT_TYPE_INTERVAL_EVENT
		}
	default:
		{
			return fmt.Errorf("EVENT unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *EVENT[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c EVENT[T]) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type POINT_EVENT[T any] struct {
	Type_            string                 `json:"_type"`
	Name             DV_TEXT                `json:"name"`
	ArchetypeNodeId  string                 `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `json:"uid"`
	Links            Option[[]LINK]         `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `json:"feeder_audit"`
	Time             DV_DATE_TIME           `json:"time"`
	State            Option[ITEM_STRUCTURE] `json:"state"`
	Data             T                      `json:"data"`
}

type INTERVAL_EVENT[T any] struct {
	Type_            string                 `json:"_type"`
	Name             DV_TEXT                `json:"name"`
	ArchetypeNodeId  string                 `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `json:"uid"`
	Links            Option[[]LINK]         `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `json:"feeder_audit"`
	Time             DV_DATE_TIME           `json:"time"`
	State            Option[ITEM_STRUCTURE] `json:"state"`
	Data             T                      `json:"data"`
	Width            DV_DURATION            `json:"width"`
	SampleCount      Option[int64]          `json:"sample_count"`
	MathFunction     DV_CODED_TEXT          `json:"math_function"`
}

// -----------------------------------
// DATA_STRUCTURES
// -----------------------------------

type DataValueType string

const (
	DATA_VALUE_TYPE_DV_BOOLEAN                     DataValueType = "DV_BOOLEAN"
	DATA_VALUE_TYPE_DV_STATE                       DataValueType = "DV_STATE"
	DATA_VALUE_TYPE_DV_IDENTIFIER                  DataValueType = "DV_IDENTIFIER"
	DATA_VALUE_TYPE_DV_TEXT                        DataValueType = "DV_TEXT"
	DATA_VALUE_TYPE_DV_CODED_TEXT                  DataValueType = "DV_CODED_TEXT"
	DATA_VALUE_TYPE_DV_PARAGRAPH                   DataValueType = "DV_PARAGRAPH"
	DATA_VALUE_TYPE_DV_INTERVAL                    DataValueType = "DV_INTERVAL"
	DATA_VALUE_TYPE_DV_ORDINAL                     DataValueType = "DV_ORDINAL"
	DATA_VALUE_TYPE_DV_SCALE                       DataValueType = "DV_SCALE"
	DATA_VALUE_TYPE_DV_QUANTITY                    DataValueType = "DV_QUANTITY"
	DATA_VALUE_TYPE_DV_COUNT                       DataValueType = "DV_COUNT"
	DATA_VALUE_TYPE_DV_PROPORTION                  DataValueType = "DV_PROPORTION"
	DATA_VALUE_TYPE_DV_DATE                        DataValueType = "DV_DATE"
	DATA_VALUE_TYPE_DV_TIME                        DataValueType = "DV_TIME"
	DATA_VALUE_TYPE_DV_DATE_TIME                   DataValueType = "DV_DATE_TIME"
	DATA_VALUE_TYPE_DV_DURATION                    DataValueType = "DV_DURATION"
	DATA_VALUE_TYPE_DV_PERIODIC_TIME_SPECIFICATION DataValueType = "DV_PERIODIC_TIME_SPECIFICATION"
	DATA_VALUE_TYPE_DV_GENERAL_TIME_SPECIFICATION  DataValueType = "DV_GENERAL_TIME_SPECIFICATION"
	DATA_VALUE_TYPE_DV_MULTIMEDIA                  DataValueType = "DV_MULTIMEDIA"
	DATA_VALUE_TYPE_DV_PARSABLE                    DataValueType = "DV_PARSABLE"
	DATA_VALUE_TYPE_DV_URI                         DataValueType = "DV_URI"
	DATA_VALUE_TYPE_DV_EHR_URI                     DataValueType = "DV_EHR_URI"
)

type DATA_VALUE struct {
	Type  DataValueType
	Value any
}

func (c *DATA_VALUE) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t DataValueType
	switch DataValueType(typStr) {
	case DATA_VALUE_TYPE_DV_BOOLEAN:
		{
			v = new(DV_BOOLEAN)
			t = DATA_VALUE_TYPE_DV_BOOLEAN
		}
	case DATA_VALUE_TYPE_DV_STATE:
		{
			v = new(DV_STATE)
			t = DATA_VALUE_TYPE_DV_STATE
		}
	case DATA_VALUE_TYPE_DV_IDENTIFIER:
		{
			v = new(DV_IDENTIFIER)
			t = DATA_VALUE_TYPE_DV_IDENTIFIER
		}
	case DATA_VALUE_TYPE_DV_TEXT:
		{
			v = new(DV_TEXT)
			t = DATA_VALUE_TYPE_DV_TEXT
		}
	case DATA_VALUE_TYPE_DV_CODED_TEXT:
		{
			v = new(DV_CODED_TEXT)
			t = DATA_VALUE_TYPE_DV_CODED_TEXT
		}
	case DATA_VALUE_TYPE_DV_PARAGRAPH:
		{
			v = new(DV_PARAGRAPH)
			t = DATA_VALUE_TYPE_DV_PARAGRAPH
		}
	case DATA_VALUE_TYPE_DV_INTERVAL:
		{
			v = new(DV_INTERVAL)
			t = DATA_VALUE_TYPE_DV_INTERVAL
		}
	case DATA_VALUE_TYPE_DV_ORDINAL:
		{
			v = new(DV_ORDINAL)
			t = DATA_VALUE_TYPE_DV_ORDINAL
		}
	case DATA_VALUE_TYPE_DV_SCALE:
		{
			v = new(DV_SCALE)
			t = DATA_VALUE_TYPE_DV_SCALE
		}
	case DATA_VALUE_TYPE_DV_QUANTITY:
		{
			v = new(DV_QUANTITY)
			t = DATA_VALUE_TYPE_DV_QUANTITY
		}
	case DATA_VALUE_TYPE_DV_COUNT:
		{
			v = new(DV_COUNT)
			t = DATA_VALUE_TYPE_DV_COUNT
		}
	case DATA_VALUE_TYPE_DV_PROPORTION:
		{
			v = new(DV_PROPORTION)
			t = DATA_VALUE_TYPE_DV_PROPORTION
		}
	case DATA_VALUE_TYPE_DV_DATE:
		{
			v = new(DV_DATE)
			t = DATA_VALUE_TYPE_DV_DATE
		}
	case DATA_VALUE_TYPE_DV_TIME:
		{
			v = new(DV_TIME)
			t = DATA_VALUE_TYPE_DV_TIME
		}
	case DATA_VALUE_TYPE_DV_DATE_TIME:
		{
			v = new(DV_DATE_TIME)
			t = DATA_VALUE_TYPE_DV_DATE_TIME
		}
	case DATA_VALUE_TYPE_DV_DURATION:
		{
			v = new(DV_DURATION)
			t = DATA_VALUE_TYPE_DV_DURATION
		}
	case DATA_VALUE_TYPE_DV_PERIODIC_TIME_SPECIFICATION:
		{
			v = new(DV_PERIODIC_TIME_SPECIFICATION)
			t = DATA_VALUE_TYPE_DV_PERIODIC_TIME_SPECIFICATION
		}
	case DATA_VALUE_TYPE_DV_GENERAL_TIME_SPECIFICATION:
		{
			v = new(DV_GENERAL_TIME_SPECIFICATION)
			t = DATA_VALUE_TYPE_DV_GENERAL_TIME_SPECIFICATION
		}
	case DATA_VALUE_TYPE_DV_MULTIMEDIA:
		{
			v = new(DV_MULTIMEDIA)
			t = DATA_VALUE_TYPE_DV_MULTIMEDIA
		}
	case DATA_VALUE_TYPE_DV_PARSABLE:
		{
			v = new(DV_PARSABLE)
			t = DATA_VALUE_TYPE_DV_PARSABLE
		}
	case DATA_VALUE_TYPE_DV_URI:
		{
			v = new(DV_URI)
			t = DATA_VALUE_TYPE_DV_URI
		}
	case DATA_VALUE_TYPE_DV_EHR_URI:
		{
			v = new(DV_EHR_URI)
			t = DATA_VALUE_TYPE_DV_EHR_URI
		}
	default:
		{
			return fmt.Errorf("DATA_VALUE unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *DATA_VALUE) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c DATA_VALUE) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type DV_BOOLEAN struct {
	Type_ Option[string] `json:"_type"`
	Value bool           `json:"value"`
}

type DV_STATE struct {
	Type_      Option[string] `json:"_type"`
	Value      DV_CODED_TEXT  `json:"value"`
	IsTerminal bool           `json:"is_terminal"`
}

type DV_IDENTIFIER struct {
	Type_    Option[string] `json:"_type"`
	Issuer   Option[string] `json:"issuer"`
	Assigner Option[string] `json:"assigner"`
	Id       string         `json:"id"`
	Type     Option[string] `json:"type"`
}

type DV_TEXT struct {
	Type_      Option[string]         `json:"_type"`
	Value      string                 `json:"value"`
	Hyperlink  Option[DV_URI]         `json:"hyperlink"`
	Formatting Option[string]         `json:"formatting"`
	Mappings   Option[[]TERM_MAPPING] `json:"mappings"`
	Language   Option[CODE_PHRASE]    `json:"language"`
	Encoding   Option[CODE_PHRASE]    `json:"encoding"`
}

type TERM_MAPPING struct {
	Type_   Option[string]        `json:"_type"`
	Match   byte                  `json:"match"`
	Purpose Option[DV_CODED_TEXT] `json:"purpose"`
	Target  CODE_PHRASE           `json:"target"`
}

type CODE_PHRASE struct {
	Type_         Option[string] `json:"_type"`
	TerminologyId TERMINOLOGY_ID `json:"terminology_id"`
	CodeString    string         `json:"code_string"`
	PreferredTerm Option[string] `json:"preferred_term"`
}

type DV_CODED_TEXT struct {
	Type_        Option[string]         `json:"_type"`
	Value        string                 `json:"value"`
	Hyperlink    Option[DV_URI]         `json:"hyperlink"`
	Formatting   Option[string]         `json:"formatting"`
	Mappings     Option[[]TERM_MAPPING] `json:"mappings"`
	Language     Option[CODE_PHRASE]    `json:"language"`
	Encoding     Option[CODE_PHRASE]    `json:"encoding"`
	DefiningCode CODE_PHRASE            `json:"defining_code"`
}

type DV_PARAGRAPH struct {
	Type_ Option[string] `json:"_type"`
	Items []DV_TEXT      `json:"items"`
}

type DvOrderedType string

const (
	DV_ORDERED_TYPE_DV_ORDINAL    DvOrderedType = "DV_ORDINAL"
	DV_ORDERED_TYPE_DV_SCALE      DvOrderedType = "DV_SCALE"
	DV_ORDERED_TYPE_DV_QUANTITY   DvOrderedType = "DV_QUANTITY"
	DV_ORDERED_TYPE_DV_COUNT      DvOrderedType = "DV_COUNT"
	DV_ORDERED_TYPE_DV_PROPORTION DvOrderedType = "DV_PROPORTION"
	DV_ORDERED_TYPE_DV_DATE       DvOrderedType = "DV_DATE"
	DV_ORDERED_TYPE_DV_TIME       DvOrderedType = "DV_TIME"
	DV_ORDERED_TYPE_DV_DATE_TIME  DvOrderedType = "DV_DATE_TIME"
	DV_ORDERED_TYPE_DV_DURATION   DvOrderedType = "DV_DURATION"
)

type DV_ORDERED struct {
	Type  DvOrderedType
	Value any
}

func (c *DV_ORDERED) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t DvOrderedType
	switch DvOrderedType(typStr) {
	case DV_ORDERED_TYPE_DV_ORDINAL:
		{
			v = new(DV_ORDINAL)
			t = DV_ORDERED_TYPE_DV_ORDINAL
		}
	case DV_ORDERED_TYPE_DV_SCALE:
		{
			v = new(DV_SCALE)
			t = DV_ORDERED_TYPE_DV_SCALE
		}
	case DV_ORDERED_TYPE_DV_QUANTITY:
		{
			v = new(DV_QUANTITY)
			t = DV_ORDERED_TYPE_DV_QUANTITY
		}
	case DV_ORDERED_TYPE_DV_COUNT:
		{
			v = new(DV_COUNT)
			t = DV_ORDERED_TYPE_DV_COUNT
		}
	case DV_ORDERED_TYPE_DV_PROPORTION:
		{
			v = new(DV_PROPORTION)
			t = DV_ORDERED_TYPE_DV_PROPORTION
		}
	case DV_ORDERED_TYPE_DV_DATE:
		{
			v = new(DV_DATE)
			t = DV_ORDERED_TYPE_DV_DATE
		}
	case DV_ORDERED_TYPE_DV_TIME:
		{
			v = new(DV_TIME)
			t = DV_ORDERED_TYPE_DV_TIME
		}
	case DV_ORDERED_TYPE_DV_DATE_TIME:
		{
			v = new(DV_DATE_TIME)
			t = DV_ORDERED_TYPE_DV_DATE_TIME
		}
	case DV_ORDERED_TYPE_DV_DURATION:
		{
			v = new(DV_DURATION)
			t = DV_ORDERED_TYPE_DV_DURATION
		}
	default:
		{
			return fmt.Errorf("DV_ORDERED unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *DV_ORDERED) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c DV_ORDERED) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type DV_INTERVAL struct {
	Type_          Option[string] `json:"_type"`
	Lower          DV_ORDERED     `json:"lower"`
	Upper          DV_ORDERED     `json:"upper"`
	LowerUnbounded bool           `json:"lower_unbounded"`
	UpperUnbounded bool           `json:"upper_unbounded"`
	LowerIncluded  bool           `json:"lower_included"`
	UpperIncluded  bool           `json:"upper_included"`
}

type REFERENCE_RANGE struct {
	Type_   Option[string] `json:"_type"`
	Meaning DV_TEXT        `json:"meaning"`
	Range   DV_INTERVAL    `json:"range"`
}

type DV_ORDINAL struct {
	Type_                Option[string]          `json:"_type"`
	NormalStatus         Option[CODE_PHRASE]     `json:"normal_status"`
	NormalRange          Option[DV_INTERVAL]     `json:"normal_range"`
	OtherReferenceRanges Option[REFERENCE_RANGE] `json:"other_reference_ranges"`
	Symbol               DV_CODED_TEXT           `json:"symbol"`
	Value                int64                   `json:"value"`
}

type DV_SCALE struct {
	Type_                Option[string]            `json:"_type"`
	NormalStatus         Option[CODE_PHRASE]       `json:"normal_status"`
	NormalRange          Option[DV_INTERVAL]       `json:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `json:"other_reference_ranges"`
	Symbol               DV_CODED_TEXT             `json:"symbol"`
	Value                float64                   `json:"value"`
}

type DV_QUANTITY struct {
	Type_                Option[string]            `json:"_type"`
	NormalStatus         Option[CODE_PHRASE]       `json:"normal_status"`
	MagnitudeStatus      Option[string]            `json:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `json:"accuracy_is_percent"`
	Accuracy             Option[float64]           `json:"accuracy"`
	Magnitude            float64                   `json:"magnitude"`
	Precision            Option[int64]             `json:"precision"`
	Units                string                    `json:"units"`
	UnitsSystem          Option[string]            `json:"units_system"`
	UnitsDisplayName     Option[string]            `json:"units_display_name"`
	NormalRange          Option[DV_INTERVAL]       `json:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `json:"other_reference_ranges"`
}

type DV_COUNT struct {
	Type_                Option[string]            `json:"_type"`
	NormalStatus         Option[CODE_PHRASE]       `json:"normal_status"`
	MagnitudeStatus      Option[string]            `json:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `json:"accuracy_is_percent"`
	Accuracy             Option[float64]           `json:"accuracy"`
	Magnitude            int64                     `json:"magnitude"`
	NormalRange          Option[DV_INTERVAL]       `json:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `json:"other_reference_ranges"`
}

type DV_PROPORTION struct {
	Type_                Option[string]            `json:"_type"`
	NormalStatus         Option[CODE_PHRASE]       `json:"normal_status"`
	MagnitudeStatus      Option[string]            `json:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `json:"accuracy_is_percent"`
	Accuracy             Option[float64]           `json:"accuracy"`
	Numerator            float64                   `json:"numerator"`
	Denominator          float64                   `json:"denominator"`
	Type                 int64                     `json:"type"`
	Precision            Option[int64]             `json:"precision"`
	NormalRange          Option[DV_INTERVAL]       `json:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `json:"other_reference_ranges"`
}

type DV_DATE struct {
	Type_                Option[string]            `json:"_type"`
	NormalStatus         Option[CODE_PHRASE]       `json:"normal_status"`
	NormalRange          Option[DV_INTERVAL]       `json:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `json:"other_reference_ranges"`
	MagnitudeStatus      Option[string]            `json:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `json:"accuracy_is_percent"`
	Accuracy             Option[DV_DURATION]       `json:"accuracy"`
	Value                string                    `json:"value"`
}

type DV_TIME struct {
	Type_                Option[string]            `json:"_type"`
	NormalStatus         Option[CODE_PHRASE]       `json:"normal_status"`
	NormalRange          Option[DV_INTERVAL]       `json:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `json:"other_reference_ranges"`
	MagnitudeStatus      Option[string]            `json:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `json:"accuracy_is_percent"`
	Accuracy             Option[DV_DURATION]       `json:"accuracy"`
	Value                string                    `json:"value"`
}

type DV_DATE_TIME struct {
	Type_                Option[string]            `json:"_type"`
	NormalStatus         Option[CODE_PHRASE]       `json:"normal_status"`
	NormalRange          Option[DV_INTERVAL]       `json:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `json:"other_reference_ranges"`
	MagnitudeStatus      Option[string]            `json:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `json:"accuracy_is_percent"`
	Accuracy             Option[DV_DURATION]       `json:"accuracy"`
	Value                string                    `json:"value"`
}

type DV_DURATION struct {
	Type_                string                    `json:"_type"`
	NormalStatus         Option[CODE_PHRASE]       `json:"normal_status"`
	NormalRange          Option[DV_INTERVAL]       `json:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `json:"other_reference_ranges"`
	MagnitudeStatus      Option[string]            `json:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `json:"accuracy_is_percent"`
	Accuracy             Option[bool]              `json:"accuracy"`
	Value                float64                   `json:"value"`
}

type DV_PERIODIC_TIME_SPECIFICATION struct {
	Type_ Option[string] `json:"_type"`
	Value DV_PARSABLE    `json:"value"`
}

type DV_GENERAL_TIME_SPECIFICATION struct {
	Type_ Option[string] `json:"_type"`
	Value DV_PARSABLE    `json:"value"`
}

type DvEncapsulatedType string

const (
	DV_ENCAPSULATED_TYPE_DV_MULTIMEDIA DvEncapsulatedType = "DV_MULTIMEDIA"
	DV_ENCAPSULATED_TYPE_DV_PARSABLE   DvEncapsulatedType = "DV_PARSABLE"
)

type DV_ENCAPSULATED struct {
	Type  DvEncapsulatedType
	Value any
}

func (c *DV_ENCAPSULATED) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t DvEncapsulatedType
	switch DvEncapsulatedType(typStr) {
	case DV_ENCAPSULATED_TYPE_DV_MULTIMEDIA:
		{
			v = new(DV_MULTIMEDIA)
			t = DV_ENCAPSULATED_TYPE_DV_MULTIMEDIA
		}
	case DV_ENCAPSULATED_TYPE_DV_PARSABLE:
		{
			v = new(DV_PARSABLE)
			t = DV_ENCAPSULATED_TYPE_DV_PARSABLE
		}
	default:
		{
			return fmt.Errorf("DV_ENCAPSULATED unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *DV_ENCAPSULATED) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c DV_ENCAPSULATED) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type DV_MULTIMEDIA struct {
	Type_                   Option[string]        `json:"_type"`
	Charset                 Option[CODE_PHRASE]   `json:"charset"`
	Language                Option[CODE_PHRASE]   `json:"language"`
	AlternateText           Option[string]        `json:"alternate_text"`
	Uri                     Option[DV_URI]        `json:"uri"`
	Data                    Option[string]        `json:"data"`
	MediaType               CODE_PHRASE           `json:"media_type"`
	CompressionAlgorithm    Option[CODE_PHRASE]   `json:"compression_algorithm"`
	IntegrityCheck          Option[string]        `json:"integrity_check"`
	IntegrityCheckAlgorithm Option[DV_MULTIMEDIA] `json:"integrity_check_algorithm"`
	Thumbnail               Option[DV_MULTIMEDIA] `json:"thumbnail"`
	Size                    int64                 `json:"size"`
}

type DV_PARSABLE struct {
	Type_     Option[string]      `json:"_type"`
	Charset   Option[CODE_PHRASE] `json:"charset"`
	Language  Option[CODE_PHRASE] `json:"language"`
	Value     string              `json:"value"`
	Formalism string              `json:"formalism"`
}

type DV_URI struct {
	Type_ Option[string] `json:"_type"`
	Value string         `json:"value"`
}

type DV_EHR_URI struct {
	Type_              Option[string] `json:"_type"`
	Value              string         `json:"value"`
	LocalTerminologyId string         `json:"local_terminology_id"`
}

// -----------------------------------
// BASE_TYPES
// -----------------------------------

type UidType string

const (
	UID_TYPE_ISO_OID     UidType = "ISO_OID"
	UID_TYPE_UUID        UidType = "UUID"
	UID_TYPE_INTERNET_ID UidType = "INTERNET_ID"
)

type UID struct {
	Type  UidType
	Value any
}

func (c *UID) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t UidType
	switch UidType(typStr) {
	case UID_TYPE_ISO_OID:
		{
			v = new(ISO_OID)
			t = UID_TYPE_ISO_OID
		}
	case UID_TYPE_UUID:
		{
			v = new(UUID)
			t = UID_TYPE_UUID
		}
	case UID_TYPE_INTERNET_ID:
		{
			v = new(INTERNET_ID)
			t = UID_TYPE_INTERNET_ID
		}
	default:
		{
			return fmt.Errorf("UID unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *UID) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c UID) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type ISO_OID struct {
	Type_ Option[string] `json:"_type"`
	Value string         `json:"value"`
}

type UUID struct {
	Type_ Option[string] `json:"_type"`
	Value string         `json:"value"`
}

type INTERNET_ID struct {
	Type_ Option[string] `json:"_type"`
	Value string         `json:"value"`
}

type ObjectIdType string

const (
	OBJECT_ID_TYPE_HIER_OBJECT_ID    ObjectIdType = "HIER_OBJECT_ID"
	OBJECT_ID_TYPE_OBJECT_VERSION_ID ObjectIdType = "OBJECT_VERSION_ID"
	OBJECT_ID_TYPE_ARCHETYPE_ID      ObjectIdType = "ARCHETYPE_ID"
	OBJECT_ID_TYPE_TEMPLATE_ID       ObjectIdType = "TEMPLATE_ID"
	OBJECT_ID_TYPE_GENERIC_ID        ObjectIdType = "GENERIC_ID"
)

type OBJECT_ID struct {
	Type  ObjectIdType
	Value any
}

func (c *OBJECT_ID) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t ObjectIdType
	switch ObjectIdType(typStr) {
	case OBJECT_ID_TYPE_HIER_OBJECT_ID:
		{
			v = new(HIER_OBJECT_ID)
			t = OBJECT_ID_TYPE_HIER_OBJECT_ID
		}
	case OBJECT_ID_TYPE_OBJECT_VERSION_ID:
		{
			v = new(OBJECT_VERSION_ID)
			t = OBJECT_ID_TYPE_OBJECT_VERSION_ID
		}
	case OBJECT_ID_TYPE_ARCHETYPE_ID:
		{
			v = new(ARCHETYPE_ID)
			t = OBJECT_ID_TYPE_ARCHETYPE_ID
		}
	case OBJECT_ID_TYPE_TEMPLATE_ID:
		{
			v = new(TEMPLATE_ID)
			t = OBJECT_ID_TYPE_TEMPLATE_ID
		}
	case OBJECT_ID_TYPE_GENERIC_ID:
		{
			v = new(GENERIC_ID)
			t = OBJECT_ID_TYPE_GENERIC_ID
		}
	default:
		{
			return fmt.Errorf("OBJECT_ID unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *OBJECT_ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c OBJECT_ID) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type UidBasedIdType string

const (
	UID_BASED_ID_TYPE_HIER_OBJECT_ID    UidBasedIdType = "HIER_OBJECT_ID"
	UID_BASED_ID_TYPE_OBJECT_VERSION_ID UidBasedIdType = "OBJECT_VERSION_ID"
)

type UID_BASED_ID struct {
	Type  UidBasedIdType
	Value any
}

func (c *UID_BASED_ID) UnmarshalJSON(data []byte) error {
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	typ, found := value["_type"]
	if !found {
		return errors.New("required _type for abstract")
	}

	typStr, ok := typ.(string)
	if !ok {
		return errors.New("expected _type to be a string")
	}

	var v any
	var t UidBasedIdType
	switch UidBasedIdType(typStr) {
	case UID_BASED_ID_TYPE_HIER_OBJECT_ID:
		{
			v = new(HIER_OBJECT_ID)
			t = UID_BASED_ID_TYPE_HIER_OBJECT_ID
		}
	case UID_BASED_ID_TYPE_OBJECT_VERSION_ID:
		{
			v = new(OBJECT_VERSION_ID)
			t = UID_BASED_ID_TYPE_OBJECT_VERSION_ID
		}

	default:
		{
			return fmt.Errorf("UID_BASED_ID unexpected _type %s", typStr)
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = t
	c.Value = v
	return nil
}

func (c *UID_BASED_ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

func (c UID_BASED_ID) Marshal() ([]byte, error) {
	return Marshal(c.Value)
}

type HIER_OBJECT_ID struct {
	Type_ Option[string] `json:"_type"`
	Value string         `json:"value"`
}

type OBJECT_VERSION_ID struct {
	Type_ Option[string] `json:"_type"`
	Value string         `json:"value"`
}

type ARCHETYPE_ID struct {
	Type_ Option[string] `json:"_type"`
	Value string         `json:"value"`
}

type TEMPLATE_ID struct {
	Type_ Option[string] `json:"_type"`
	Value string         `json:"value"`
}

type TERMINOLOGY_ID struct {
	Type_ Option[string] `json:"_type"`
	Value string         `json:"value"`
}

type GENERIC_ID struct {
	Type_  Option[string] `json:"_type"`
	Value  string         `json:"value"`
	Scheme string         `json:"scheme"`
}

type OBJECT_REF struct {
	Type_     Option[string] `json:"object_ref"`
	Namespace string         `json:"namespace"`
	Type      string         `json:"type"`
	Id        OBJECT_ID      `json:"id"`
}

type PARTY_REF struct {
	Type_     Option[string] `json:"_type"`
	Namespace string         `json:"namespace"`
	Type      string         `json:"type"`
	Id        OBJECT_ID      `json:"id"`
}

type LOCATABLE_REF struct {
	Type_     Option[string] `json:"_type"`
	Namespace string         `json:"namespace"`
	Type      string         `json:"type"`
	Path      Option[string] `json:"path"`
	Id        UID_BASED_ID   `json:"id"`
}
