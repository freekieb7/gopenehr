package gopenehr

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/freekieb7/gopenehr/encoding/json"
)

var ErrNotFound = errors.New("value not found")
var ErrBadType = errors.New("parse error")

type Option[T any] struct {
	Value T
	Some  bool
}

func Some[T any](v T) Option[T] {
	return Option[T]{Value: v, Some: true}
}

func None[T any]() Option[T] {
	var zero T
	return Option[T]{Value: zero, Some: false}
}

func (o Option[T]) IsSome() bool { return o.Some }
func (o Option[T]) Unwrap() T    { return o.Value }

func (o *Option[T]) Unmarshal(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		var zero T
		o.Some = false
		o.Value = zero
	} else {
		var v T
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		o.Value = v
		o.Some = true
	}

	return nil
}

func (o Option[T]) Marshal() ([]byte, error) {
	if o.Some {
		return json.Marshal(o.Unwrap())
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
	Name             DV_TEXT                `json:"name"`
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
	EndTime            Option[DV_DATE_TIME]     `json:"end_time"`
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

func (c *CONTENT_ITEM) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := ContentItemType(typeData)
	switch t {
	case CONTENT_ITEM_TYPE_SECTION:
		{
			c.Value = new(SECTION)
		}
	case CONTENT_ITEM_TYPE_ADMIN_ENTRY:
		{
			c.Value = new(ADMIN_ENTRY)
		}
	case CONTENT_ITEM_TYPE_OBSERVATION:
		{
			c.Value = new(OBSERVATION)
		}
	case CONTENT_ITEM_TYPE_EVALUATION:
		{
			c.Value = new(EVALUATION)
		}
	case CONTENT_ITEM_TYPE_INSTRUCTION:
		{
			c.Value = new(INSTRUCTION)
		}
	case CONTENT_ITEM_TYPE_ACTIVITY:
		{
			c.Value = new(ACTIVITY)
		}
	case CONTENT_ITEM_TYPE_ACTION:
		{
			c.Value = new(ACTION)
		}
	default:
		{
			return fmt.Errorf("CONTENT_ITEM unexpected _type %s", t)
		}
	}

	c.Type = t
	return json.Unmarshal(data, c.Value)
}

func (c CONTENT_ITEM) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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

func (c *ENTRY) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := EntryType(typeData)
	switch t {
	case ENTRY_TYPE_ADMIN_ENTRY:
		{
			c.Value = new(ADMIN_ENTRY)
		}
	case ENTRY_TYPE_OBSERVATION:
		{
			c.Value = new(OBSERVATION)
		}
	case ENTRY_TYPE_EVALUATION:
		{
			c.Value = new(EVALUATION)
		}
	case ENTRY_TYPE_INSTRUCTION:
		{
			c.Value = new(INSTRUCTION)
		}
	case ENTRY_TYPE_ACTIVITY:
		{
			c.Value = new(ACTIVITY)
		}
	case ENTRY_TYPE_ACTION:
		{
			c.Value = new(ACTION)
		}
	default:
		{
			return fmt.Errorf("ENTRY unexpected _type %s", t)
		}
	}

	c.Type = t
	return json.Unmarshal(data, c.Value)
}

func (c ENTRY) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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

func (c *CARE_ENTRY) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := CareEntryType(typeData)
	switch t {
	case CARE_ENTRY_TYPE_OBSERVATION:
		{
			c.Value = new(OBSERVATION)
		}
	case CARE_ENTRY_TYPE_EVALUATION:
		{
			c.Value = new(EVALUATION)
		}
	case CARE_ENTRY_TYPE_INSTRUCTION:
		{
			c.Value = new(INSTRUCTION)
		}
	case CARE_ENTRY_TYPE_ACTIVITY:
		{
			c.Value = new(ACTIVITY)
		}
	case CARE_ENTRY_TYPE_ACTION:
		{
			c.Value = new(ACTION)
		}
	default:
		{
			return fmt.Errorf("CARE_ENTRY unexpected _type %s", t)
		}
	}

	c.Type = t
	return json.Unmarshal(data, c.Value)
}

func (c CARE_ENTRY) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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
	Narrative           DV_TEXT                 `json:"narrative"`
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
	Type_                    Option[string]               `json:"_type"`
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

func (p *PARTY_PROXY) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := PartyProxyType(typeData)
	switch t {
	case PARTY_PROXY_TYPE_PARTY_SELF:
		{
			p.Value = new(PARTY_SELF)
		}
	case PARTY_PROXY_TYPE_PARTY_IDENTIFIED:
		{
			p.Value = new(PARTY_IDENTIFIED)
		}
	case PARTY_PROXY_TYPE_PARTY_RELATED:
		{
			p.Value = new(PARTY_RELATED)
		}
	default:
		{
			return fmt.Errorf("PARTY_PROXY unexpected _type %s", t)
		}
	}

	p.Type = t
	return json.Unmarshal(data, p.Value)
}

func (c PARTY_PROXY) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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

func (i *ITEM_STRUCTURE) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := ItemStructureType(typeData)
	switch t {
	case ITEM_STRUCTURE_TYPE_ITEM_SINGLE:
		{
			i.Value = new(ITEM_SINGLE)
		}
	case ITEM_STRUCTURE_TYPE_ITEM_LIST:
		{
			i.Value = new(ITEM_LIST)
		}
	case ITEM_STRUCTURE_TYPE_ITEM_TABLE:
		{
			i.Value = new(ITEM_TABLE)
		}
	case ITEM_STRUCTURE_TYPE_ITEM_TREE:
		{
			i.Value = new(ITEM_TREE)
		}
	default:
		{
			return fmt.Errorf("ITEM_STRUCTURE unexpected _type %s", t)
		}
	}

	i.Type = t
	return json.Unmarshal(data, i.Value)
}

func (c ITEM_STRUCTURE) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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
	Type_            Option[string]       `json:"_type"`
	Name             DV_TEXT              `json:"name"`
	ArchetypeNodeId  string               `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `json:"uid"`
	Links            Option[[]LINK]       `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `json:"feeder_audit"`
	Items            Option[[]ELEMENT]    `json:"items"`
}

type ITEM_TABLE struct {
	Type_            Option[string]       `json:"_type"`
	Name             DV_TEXT              `json:"name"`
	ArchetypeNodeId  string               `json:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `json:"uid"`
	Links            Option[[]LINK]       `json:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `json:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `json:"feeder_audit"`
	Rows             Option[[]CLUSTER]    `json:"rows"`
}

type ITEM_TREE struct {
	Type_            Option[string]       `json:"_type"`
	Name             DV_TEXT              `json:"name"`
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

func (i *ITEM) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := ItemType(typeData)
	switch t {
	case ITEM_TYPE_CLUSTER:
		{
			i.Value = new(CLUSTER)
		}
	case ITEM_TYPE_ELEMENT:
		{
			i.Value = new(ELEMENT)
		}
	default:
		{
			return fmt.Errorf("ITEM unexpected _type %s", t)
		}
	}

	i.Type = t
	return json.Unmarshal(data, i.Value)
}

func (c ITEM) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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

func (e *EVENT[T]) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := EventType(typeData)
	switch t {
	case EVENT_TYPE_POINT_EVENT:
		{
			e.Value = new(POINT_EVENT[T])
		}
	case EVENT_TYPE_INTERVAL_EVENT:
		{
			e.Value = new(INTERVAL_EVENT[T])
		}
	default:
		{
			return fmt.Errorf("EVENT unexpected _type %s", t)
		}
	}

	e.Type = t
	return json.Unmarshal(data, e.Value)
}

func (c EVENT[T]) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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

func (d *DATA_VALUE) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := DataValueType(typeData)
	switch t {
	case DATA_VALUE_TYPE_DV_BOOLEAN:
		{
			d.Value = new(DV_BOOLEAN)
		}
	case DATA_VALUE_TYPE_DV_STATE:
		{
			d.Value = new(DV_STATE)
		}
	case DATA_VALUE_TYPE_DV_IDENTIFIER:
		{
			d.Value = new(DV_IDENTIFIER)
		}
	case DATA_VALUE_TYPE_DV_TEXT:
		{
			d.Value = new(DV_TEXT)
		}
	case DATA_VALUE_TYPE_DV_CODED_TEXT:
		{
			d.Value = new(DV_CODED_TEXT)
		}
	case DATA_VALUE_TYPE_DV_PARAGRAPH:
		{
			d.Value = new(DV_PARAGRAPH)
		}
	case DATA_VALUE_TYPE_DV_INTERVAL:
		{
			d.Value = new(DV_INTERVAL)
		}
	case DATA_VALUE_TYPE_DV_ORDINAL:
		{
			d.Value = new(DV_ORDINAL)
		}
	case DATA_VALUE_TYPE_DV_SCALE:
		{
			d.Value = new(DV_SCALE)
		}
	case DATA_VALUE_TYPE_DV_QUANTITY:
		{
			d.Value = new(DV_QUANTITY)
		}
	case DATA_VALUE_TYPE_DV_COUNT:
		{
			d.Value = new(DV_COUNT)
		}
	case DATA_VALUE_TYPE_DV_PROPORTION:
		{
			d.Value = new(DV_PROPORTION)
		}
	case DATA_VALUE_TYPE_DV_DATE:
		{
			d.Value = new(DV_DATE)
		}
	case DATA_VALUE_TYPE_DV_TIME:
		{
			d.Value = new(DV_TIME)
		}
	case DATA_VALUE_TYPE_DV_DATE_TIME:
		{
			d.Value = new(DV_DATE_TIME)
		}
	case DATA_VALUE_TYPE_DV_DURATION:
		{
			d.Value = new(DV_DURATION)
		}
	case DATA_VALUE_TYPE_DV_PERIODIC_TIME_SPECIFICATION:
		{
			d.Value = new(DV_PERIODIC_TIME_SPECIFICATION)
		}
	case DATA_VALUE_TYPE_DV_GENERAL_TIME_SPECIFICATION:
		{
			d.Value = new(DV_GENERAL_TIME_SPECIFICATION)
		}
	case DATA_VALUE_TYPE_DV_MULTIMEDIA:
		{
			d.Value = new(DV_MULTIMEDIA)
		}
	case DATA_VALUE_TYPE_DV_PARSABLE:
		{
			d.Value = new(DV_PARSABLE)
		}
	case DATA_VALUE_TYPE_DV_URI:
		{
			d.Value = new(DV_URI)
		}
	case DATA_VALUE_TYPE_DV_EHR_URI:
		{
			d.Value = new(DV_EHR_URI)
		}
	default:
		{
			return fmt.Errorf("DATA_VALUE unexpected _type %s", t)
		}
	}

	d.Type = t
	return json.Unmarshal(data, d.Value)
}

func (c DATA_VALUE) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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

func (d *DV_ORDERED) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := DvOrderedType(typeData)
	switch t {
	case DV_ORDERED_TYPE_DV_ORDINAL:
		{
			d.Value = new(DV_ORDINAL)
		}
	case DV_ORDERED_TYPE_DV_SCALE:
		{
			d.Value = new(DV_SCALE)
		}
	case DV_ORDERED_TYPE_DV_QUANTITY:
		{
			d.Value = new(DV_QUANTITY)
		}
	case DV_ORDERED_TYPE_DV_COUNT:
		{
			d.Value = new(DV_COUNT)
		}
	case DV_ORDERED_TYPE_DV_PROPORTION:
		{
			d.Value = new(DV_PROPORTION)
		}
	case DV_ORDERED_TYPE_DV_DATE:
		{
			d.Value = new(DV_DATE)
		}
	case DV_ORDERED_TYPE_DV_TIME:
		{
			d.Value = new(DV_TIME)
		}
	case DV_ORDERED_TYPE_DV_DATE_TIME:
		{
			d.Value = new(DV_DATE_TIME)
		}
	case DV_ORDERED_TYPE_DV_DURATION:
		{
			d.Value = new(DV_DURATION)
		}
	default:
		{
			return fmt.Errorf("DV_ORDERED unexpected _type %s", t)
		}
	}

	d.Type = t
	return json.Unmarshal(data, d.Value)
}

func (c DV_ORDERED) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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
	Value                string                    `json:"value"`
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

func (d *DV_ENCAPSULATED) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := DvEncapsulatedType(typeData)
	switch t {
	case DV_ENCAPSULATED_TYPE_DV_MULTIMEDIA:
		{
			d.Value = new(DV_MULTIMEDIA)
		}
	case DV_ENCAPSULATED_TYPE_DV_PARSABLE:
		{
			d.Value = new(DV_PARSABLE)
		}
	default:
		{
			return fmt.Errorf("DV_ENCAPSULATED unexpected _type %s", t)
		}
	}

	d.Type = t
	return json.Unmarshal(data, d.Value)
}

func (c DV_ENCAPSULATED) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
}

type DV_MULTIMEDIA struct {
	Type_                   Option[string]         `json:"_type"`
	Charset                 Option[CODE_PHRASE]    `json:"charset"`
	Language                Option[CODE_PHRASE]    `json:"language"`
	AlternateText           Option[string]         `json:"alternate_text"`
	Uri                     Option[DV_URI]         `json:"uri"`
	Data                    Option[string]         `json:"data"`
	MediaType               CODE_PHRASE            `json:"media_type"`
	CompressionAlgorithm    Option[CODE_PHRASE]    `json:"compression_algorithm"`
	IntegrityCheck          Option[string]         `json:"integrity_check"`
	IntegrityCheckAlgorithm Option[CODE_PHRASE]    `json:"integrity_check_algorithm"`
	Thumbnail               Option[*DV_MULTIMEDIA] `json:"thumbnail"`
	Size                    int64                  `json:"size"`
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

func (u *UID) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := UidType(typeData)
	switch t {
	case UID_TYPE_ISO_OID:
		{
			u.Value = new(ISO_OID)
		}
	case UID_TYPE_UUID:
		{
			u.Value = new(UUID)
		}
	case UID_TYPE_INTERNET_ID:
		{
			u.Value = new(INTERNET_ID)
		}
	default:
		{
			return fmt.Errorf("UID unexpected _type %s", t)
		}
	}

	u.Type = t
	return json.Unmarshal(data, u.Value)
}

func (c UID) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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

func (o *OBJECT_ID) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := ObjectIdType(typeData)
	switch t {
	case OBJECT_ID_TYPE_HIER_OBJECT_ID:
		{
			o.Value = new(HIER_OBJECT_ID)
		}
	case OBJECT_ID_TYPE_OBJECT_VERSION_ID:
		{
			o.Value = new(OBJECT_VERSION_ID)
		}
	case OBJECT_ID_TYPE_ARCHETYPE_ID:
		{
			o.Value = new(ARCHETYPE_ID)
		}
	case OBJECT_ID_TYPE_TEMPLATE_ID:
		{
			o.Value = new(TEMPLATE_ID)
		}
	case OBJECT_ID_TYPE_GENERIC_ID:
		{
			o.Value = new(GENERIC_ID)
		}
	default:
		{
			return fmt.Errorf("OBJECT_ID unexpected _type %s", t)
		}
	}

	o.Type = t
	return json.Unmarshal(data, o.Value)
}

func (c OBJECT_ID) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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

func (u *UID_BASED_ID) Unmarshal(data []byte) error {
	typeData, err := json.Search(data, "_type")
	if err != nil {
		return err
	}
	typeData = typeData[1 : len(typeData)-1]

	t := UidBasedIdType(typeData)
	switch t {
	case UID_BASED_ID_TYPE_HIER_OBJECT_ID:
		{
			u.Value = new(HIER_OBJECT_ID)
		}
	case UID_BASED_ID_TYPE_OBJECT_VERSION_ID:
		{
			u.Value = new(OBJECT_VERSION_ID)
		}

	default:
		{
			return fmt.Errorf("UID_BASED_ID unexpected _type %s", t)
		}
	}

	u.Type = t
	return json.Unmarshal(data, u.Value)
}

func (c UID_BASED_ID) Marshal() ([]byte, error) {
	return json.Marshal(c.Value)
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
