package openehr

import (
	"errors"
	"reflect"
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

func (o *Option[T]) SetSome(v T) {
	o.Value = &v
	o.Some = true
}

func (o *Option[T]) SetNone() {
	o.Value = nil
	o.Some = false
}

func (o *Option[T]) Unmarshal(data any) {
	if data == nil {
		o.SetNone()
	}

	v := reflect.ValueOf(*o.Value)
	PocUnmarshal(v, data)
}

type OPENEHR interface {
	MetaType() string
}

type Marshaller interface {
	Unmarshal(data any) error
}

// -----------------------------------
// EHR
// -----------------------------------

type EHR struct {
	Type_         Option[string]         `openehr:"type_"`
	SystemId      Option[HIER_OBJECT_ID] `openehr:"system_id"`
	EhrId         HIER_OBJECT_ID         `openehr:"ehr_id"`
	Contributions Option[[]OBJECT_REF]   `openehr:"contributions"`
	EhrStatus     OBJECT_REF             `openehr:"ehr_status"`
	EhrAccess     OBJECT_REF             `openehr:"ehr_access"`
	Compositions  Option[[]OBJECT_REF]   `openehr:"compositions"`
	Directory     Option[OBJECT_REF]     `openehr:"directory"`
	TimeCreated   DV_DATE_TIME           `openehr:"time_created"`
	Folders       Option[[]OBJECT_REF]   `openehr:"folders"`
}

func (EHR) MetaType() string {
	return "EHR"
}

// func (ehr *EHR) Unmarshal(data any) error {
// 	if value, found := data["type_"]; found {
// 		if str, ok := value.(string); ok {
// 			ehr.Type_.SetSome(str)
// 		}
// 		return ErrBadType
// 	} else {
// 		ehr.Type_.SetNone()
// 	}

// 	if value, found := data["system_id"]; found {
// 		if s, ok := value.(string); ok {
// 			ehr.Type_.SetSome(s)
// 		} else {
// 			return ErrBadType
// 		}
// 	} else {
// 		ehr.Type_.SetNone()
// 	}

// 	if value, found := data["ehr_id"]; found {
// 		if m, ok := value.(map[string]any); ok {
// 			ehr.EhrId.Unmarshal(m)
// 		} else {
// 			return ErrNotFound
// 		}
// 	} else {
// 		return ErrNotFound
// 	}

// 	if value, found := data["contributions"]; found {
// 		if slice, ok := value.([]any); ok {
// 			t := make([]OBJECT_REF, 0)
// 			for _, el := range slice {
// 				if m, ok := el.(map[string]any); ok {
// 					var a OBJECT_REF
// 					a.Unmarshal(m)
// 					t = append(t, a)
// 				} else {
// 					return ErrBadType
// 				}
// 			}
// 			ehr.Contributions.SetSome(t)
// 		} else {
// 			return ErrBadType
// 		}
// 	} else {
// 		ehr.Contributions.SetNone()
// 	}

// 	if value, found := data[""]

// 	return nil
// }

type VERSIONED_EHR_ACCESS struct {
	Type_       Option[string] `openehr:"type_"`
	Uid         HIER_OBJECT_ID `openehr:"uid"`
	OwnerId     OBJECT_REF     `openehr:"owner_id"`
	TimeCreated DV_DATE_TIME   `openehr:"time_created"`
}

func (VERSIONED_EHR_ACCESS) MetaType() string {
	return "VERSIONED_EHR_ACCESS"
}

// func (v *VERSIONED_EHR_ACCESS) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type EHR_ACCESS struct {
	Type_            Option[string]       `openehr:"type_"`
	Name             DV_TEXT              `openehr:"name"`
	ArchetypeNodeId  string               `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `openehr:"uid"`
	Links            Option[[]LINK]       `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `openehr:"archetype_node_id"`
	FeederAudit      Option[FEEDER_AUDIT] `openehr:"feeder_audit"`
}

func (EHR_ACCESS) MetaType() string {
	return "EHR_ACCESS"
}

// func (v *EHR_ACCESS) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type VERSIONED_EHR_STATUS struct {
	Type_       Option[string] `openehr:"type_"`
	Uid         HIER_OBJECT_ID `openehr:"uid"`
	OwnerId     OBJECT_REF     `openehr:"owner_id"`
	TimeCreated DV_DATE_TIME   `openehr:"time_created"`
}

func (VERSIONED_EHR_STATUS) MetaType() string {
	return "VERSIONED_EHR_STATUS"
}

// func (v *VERSIONED_EHR_STATUS) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type EHR_STATUS struct {
	Type_            Option[string]         `openehr:"type_"`
	Name             DV_TEXT                `openehr:"name"`
	ArchetypeNodeId  string                 `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `openehr:"uid"`
	Links            Option[[]LINK]         `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `openehr:"feeder_audit"`
	Subject          PARTY_SELF             `openehr:"subject"`
	IsQueryable      bool                   `openehr:"is_queryable"`
	IsModifiable     bool                   `openehr:"is_modifiable"`
	OtherDetails     Option[ITEM_STRUCTURE] `openehr:"other_details"`
}

// func (EHR_STATUS) MetaType() string {
// 	return "EHR_STATUS"
// }

// func (ehrStatus *EHR_STATUS) Unmarshal(data any) error {
// 	if value, found := data["type_"]; found {
// 		v, ok := value.(string)
// 		if !ok {
// 			return ErrBadType
// 		}
// 		ehrStatus.Type_.SetSome(v)
// 	} else {
// 		ehrStatus.Type_.SetNone()
// 	}

// 	if value, found := data["name"]; found {
// 		v, ok := value.(map[string]any)
// 		if !ok {
// 			return ErrBadType
// 		}
// 		ehrStatus.Name.Unmarshal(v)
// 	} else {
// 		return ErrNotFound
// 	}

// 	return nil
// }

type VERSIONED_COMPOSITION struct {
	Type_       Option[string] `openehr:"type_"`
	Uid         HIER_OBJECT_ID `openehr:"uid"`
	OwnerId     OBJECT_REF     `openehr:"owner_id"`
	TimeCreated DV_DATE_TIME   `openehr:"time_created"`
}

func (VERSIONED_COMPOSITION) MetaType() string {
	return "VERSIONED_COMPOSITION"
}

// func (v *VERSIONED_COMPOSITION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type COMPOSITION struct {
	Type_            Option[string]         `openehr:"type_"`
	ArchetypeNodeId  string                 `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `openehr:"uid"`
	Links            Option[[]LINK]         `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `openehr:"feeder_audit"`
	Language         CODE_PHRASE            `openehr:"language"`
	Territory        CODE_PHRASE            `openehr:"territory"`
	Category         DV_CODED_TEXT          `openehr:"category"`
	Context          Option[EVENT_CONTEXT]  `openehr:"context"`
	Composer         PARTY_PROXY            `openehr:"composer"`
	Content          Option[[]CONTENT_ITEM] `openehr:"content"`
}

func (COMPOSITION) MetaType() string {
	return "COMPOSITION"
}

// func (v *COMPOSITION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type EVENT_CONTEXT struct {
	Type_              Option[string]           `openehr:"type_"`
	StartTime          DV_DATE_TIME             `openehr:"start_time"`
	EndTime            DV_DATE_TIME             `openehr:"end_time"`
	Location           Option[string]           `openehr:"location"`
	Setting            DV_CODED_TEXT            `openehr:"setting"`
	OtherContext       Option[ITEM_STRUCTURE]   `openehr:"other_context"`
	HealthCareFacility Option[PARTY_IDENTIFIED] `openehr:"health_care_facility"`
	Participations     Option[[]PARTICIPATION]  `openehr:"participations"`
}

func (EVENT_CONTEXT) MetaType() string {
	return "EVENT_CONTEXT"
}

// func (v *EVENT_CONTEXT) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type CONTENT_ITEM interface {
	_CONTENT_ITEM()
}

func (SECTION) _CONTENT_ITEM()     {}
func (ADMIN_ENTRY) _CONTENT_ITEM() {}
func (OBSERVATION) _CONTENT_ITEM() {}
func (EVALUATION) _CONTENT_ITEM()  {}
func (INSTRUCTION) _CONTENT_ITEM() {}
func (ACTIVITY) _CONTENT_ITEM()    {}
func (ACTION) _CONTENT_ITEM()      {}

type SECTION struct {
	Type_            Option[string]         `openehr:"type_"`
	Name             DV_TEXT                `openehr:"name"`
	ArchetypeNodeId  string                 `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `openehr:"uid"`
	Links            Option[[]LINK]         `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `openehr:"feeder_audit"`
	Items            Option[[]CONTENT_ITEM] `openehr:"items"`
}

func (SECTION) MetaType() string {
	return "SECTION"
}

// func (v *SECTION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type ENTRY interface {
	_ENTRY()
}

func (ADMIN_ENTRY) _ENTRY() {}
func (OBSERVATION) _ENTRY() {}
func (EVALUATION) _ENTRY()  {}
func (INSTRUCTION) _ENTRY() {}
func (ACTIVITY) _ENTRY()    {}
func (ACTION) _ENTRY()      {}

type ADMIN_ENTRY struct {
	Type_               Option[string]          `openehr:"type_"`
	Name                DV_TEXT                 `openehr:"name"`
	ArchetypeNodeId     string                  `openehr:"archetype_node_id"`
	Uid                 Option[UID_BASED_ID]    `openehr:"uid"`
	Links               Option[[]LINK]          `openehr:"links"`
	ArchetypeDetails    Option[ARCHETYPED]      `openehr:"archetype_details"`
	FeederAudit         Option[FEEDER_AUDIT]    `openehr:"feeder_audit"`
	Language            CODE_PHRASE             `openehr:"language"`
	Encoding            CODE_PHRASE             `openehr:"encoding"`
	OtherParticipations Option[[]PARTICIPATION] `openehr:"other_participations"`
	WorkflowId          Option[OBJECT_REF]      `openehr:"workflow_id"`
	Subject             PARTY_PROXY             `openehr:"subject"`
	Provider            Option[PARTY_PROXY]     `openehr:"provider"`
	Data                ITEM_STRUCTURE          `openehr:"data"`
}

func (ADMIN_ENTRY) MetaType() string {
	return "ADMIN_ENTRY"
}

// func (v *ADMIN_ENTRY) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type CARE_ENTRY interface {
	_CARE_ENTRY()
}

func (OBSERVATION) _CARE_ENTRY() {}
func (EVALUATION) _CARE_ENTRY()  {}
func (INSTRUCTION) _CARE_ENTRY() {}
func (ACTIVITY) _CARE_ENTRY()    {}
func (ACTION) _CARE_ENTRY()      {}

type OBSERVATION struct {
	Type_               Option[string]                  `openehr:"type_"`
	Name                DV_TEXT                         `openehr:"name"`
	ArchetypeNodeId     string                          `openehr:"archetype_node_id"`
	Uid                 Option[UID_BASED_ID]            `openehr:"uid"`
	Links               Option[[]LINK]                  `openehr:"links"`
	ArchetypeDetails    Option[ARCHETYPED]              `openehr:"archetype_details"`
	FeederAudit         Option[FEEDER_AUDIT]            `openehr:"feeder_audit"`
	Language            CODE_PHRASE                     `openehr:"language"`
	Encoding            CODE_PHRASE                     `openehr:"encoding"`
	OtherParticipations Option[[]PARTICIPATION]         `openehr:"other_participations"`
	WorkflowId          Option[OBJECT_REF]              `openehr:"workflow_id"`
	Subject             PARTY_PROXY                     `openehr:"subject"`
	Provider            Option[PARTY_PROXY]             `openehr:"provider"`
	Protocol            Option[ITEM_STRUCTURE]          `openehr:"protocol"`
	GuidelineId         Option[OBJECT_REF]              `openehr:"guideline_id"`
	Data                HISTORY[ITEM_STRUCTURE]         `openehr:"data"`
	State               Option[HISTORY[ITEM_STRUCTURE]] `openehr:"state"`
}

func (OBSERVATION) MetaType() string {
	return "OBSERVATION"
}

// func (v *OBSERVATION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type EVALUATION struct {
	Type_               Option[string]          `openehr:"type_"`
	Name                DV_TEXT                 `openehr:"name"`
	ArchetypeNodeId     string                  `openehr:"archetype_node_id"`
	Uid                 Option[UID_BASED_ID]    `openehr:"uid"`
	Links               Option[[]LINK]          `openehr:"links"`
	ArchetypeDetails    Option[ARCHETYPED]      `openehr:"archetype_details"`
	FeederAudit         Option[FEEDER_AUDIT]    `openehr:"feeder_audit"`
	Language            CODE_PHRASE             `openehr:"language"`
	Encoding            CODE_PHRASE             `openehr:"encoding"`
	OtherParticipations Option[[]PARTICIPATION] `openehr:"other_participations"`
	WorkflowId          Option[OBJECT_REF]      `openehr:"workflow_id"`
	Subject             PARTY_PROXY             `openehr:"subject"`
	Provider            Option[PARTY_PROXY]     `openehr:"provider"`
	Protocol            Option[ITEM_STRUCTURE]  `openehr:"protocol"`
	GuidelineId         Option[OBJECT_REF]      `openehr:"guideline_id"`
	Data                ITEM_STRUCTURE          `openehr:"data"`
}

func (EVALUATION) MetaType() string {
	return "EVALUATION"
}

// func (v *EVALUATION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type INSTRUCTION struct {
	Type_               Option[string]          `openehr:"type_"`
	Name                DV_TEXT                 `openehr:"name"`
	ArchetypeNodeId     string                  `openehr:"archetype_node_id"`
	Uid                 Option[UID_BASED_ID]    `openehr:"uid"`
	Links               Option[[]LINK]          `openehr:"links"`
	ArchetypeDetails    Option[ARCHETYPED]      `openehr:"archetype_details"`
	FeederAudit         Option[FEEDER_AUDIT]    `openehr:"feeder_audit"`
	Language            CODE_PHRASE             `openehr:"language"`
	Encoding            CODE_PHRASE             `openehr:"encoding"`
	OtherParticipations Option[[]PARTICIPATION] `openehr:"other_participations"`
	WorkflowId          Option[OBJECT_REF]      `openehr:"workflow_id"`
	Subject             PARTY_PROXY             `openehr:"subject"`
	Provider            Option[PARTY_PROXY]     `openehr:"provider"`
	Protocol            Option[ITEM_STRUCTURE]  `openehr:"protocol"`
	GuidelineId         Option[OBJECT_REF]      `openehr:"guideline_id"`
	Narative            DV_TEXT                 `openehr:"narative"`
	ExpiryTime          Option[DV_DATE_TIME]    `openehr:"expiry_time"`
	WfDefinition        Option[DV_PARSABLE]     `openehr:"wf_definition"`
	Activities          Option[[]ACTIVITY]      `openehr:"activities"`
}

func (INSTRUCTION) MetaType() string {
	return "INSTRUCTION"
}

// func (v *INSTRUCTION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type ACTIVITY struct {
	Type_             Option[string]       `openehr:"type_"`
	Name              DV_TEXT              `openehr:"name"`
	ArchetypeNodeId   string               `openehr:"archetype_node_id"`
	Uid               Option[UID_BASED_ID] `openehr:"uid"`
	Links             Option[[]LINK]       `openehr:"links"`
	ArchetypeDetails  Option[ARCHETYPED]   `openehr:"archetype_details"`
	FeederAudit       Option[FEEDER_AUDIT] `openehr:"feeder_audit"`
	Timing            Option[DV_PARSABLE]  `openehr:"timing"`
	ActionArchetypeId string               `openehr:"action_archetype_id"`
	Description       ITEM_STRUCTURE       `openehr:"description"`
}

func (ACTIVITY) MetaType() string {
	return "ACTIVITY"
}

// func (v *ACTIVITY) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type ACTION struct {
	Type_               Option[string]              `openehr:"type_"`
	Name                DV_TEXT                     `openehr:"name"`
	ArchetypeNodeId     string                      `openehr:"archetype_node_id"`
	Uid                 Option[UID_BASED_ID]        `openehr:"uid"`
	Links               Option[[]LINK]              `openehr:"links"`
	ArchetypeDetails    Option[ARCHETYPED]          `openehr:"archetype_details"`
	FeederAudit         Option[FEEDER_AUDIT]        `openehr:"feeder_audit"`
	Language            CODE_PHRASE                 `openehr:"language"`
	Encoding            CODE_PHRASE                 `openehr:"encoding"`
	OtherParticipations Option[[]PARTICIPATION]     `openehr:"other_participations"`
	WorkflowId          Option[OBJECT_REF]          `openehr:"workflow_id"`
	Subject             PARTY_PROXY                 `openehr:"subject"`
	Provider            Option[PARTY_PROXY]         `openehr:"provider"`
	Protocol            Option[ITEM_STRUCTURE]      `openehr:"protocol"`
	GuidelineId         Option[OBJECT_REF]          `openehr:"guideline_id"`
	Time                DV_DATE_TIME                `openehr:"time"`
	IsmTransition       ISM_TRANSITION              `openehr:"ism_transition"`
	InstructionDetails  Option[INSTRUCTION_DETAILS] `openehr:"instruction_details"`
	Description         ITEM_STRUCTURE              `openehr:"description"`
}

func (ACTION) MetaType() string {
	return "ACTION"
}

// func (v *ACTION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type INSTRUCTION_DETAILS struct {
	Type_         Option[string]         `openehr:"type_"`
	InstructionId LOCATABLE_REF          `openehr:"instruction_id"`
	ActivityId    string                 `openehr:"activity"`
	WfDetails     Option[ITEM_STRUCTURE] `openehr:"wf_details"`
}

func (INSTRUCTION_DETAILS) MetaType() string {
	return "INSTRUCTION_DETAILS"
}

// func (v *INSTRUCTION_DETAILS) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type ISM_TRANSITION struct {
	Type_        Option[string]        `openehr:"type_"`
	CurrentState DV_CODED_TEXT         `openehr:"current_state"`
	Transition   Option[DV_CODED_TEXT] `openehr:"transition"`
	CareflowStep Option[DV_CODED_TEXT] `openehr:"cateflow_step"`
	Reason       Option[DV_TEXT]       `openehr:"reason"`
}

func (ISM_TRANSITION) MetaType() string {
	return "ISM_TRANSITION"
}

// func (v *ISM_TRANSITION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

// -----------------------------------
// COMMON
// -----------------------------------

type ARCHETYPED struct {
	Type_       Option[string]      `openehr:"type_"`
	ArchetypeId ARCHETYPE_ID        `openehr:"archetype_id"`
	TemplateId  Option[TEMPLATE_ID] `openehr:"template_id"`
	RmVersion   string              `openehr:"rm_version"`
}

func (ARCHETYPED) MetaType() string {
	return "ARCHETYPED"
}

// func (v *ARCHETYPED) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type LINK struct {
	Type_   Option[string] `openehr:"type_"`
	Meaning DV_TEXT        `openehr:"meaning"`
	Type    DV_TEXT        `openehr:"type"`
	Target  DV_EHR_URI     `openehr:"target"`
}

func (LINK) MetaType() string {
	return "LINK"
}

// func (v *LINK) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type FEEDER_AUDIT struct {
	Type_                    Option[FEEDER_AUDIT]         `openehr:"type_"`
	OriginatingSystemItemIds Option[[]DV_IDENTIFIER]      `openehr:"originating_system_item_ids"`
	FeederSystemItemIds      Option[[]DV_IDENTIFIER]      `openehr:"feeder_system_item_ids"`
	OriginalContent          Option[DV_ENCAPSULATED]      `openehr:"original_content"`
	OriginatingSystemAudit   FEEDER_AUDIT_DETAILS         `openehr:"originating_system_audit"`
	FeederSystemAudit        Option[FEEDER_AUDIT_DETAILS] `openehr:"feeder_system_audit"`
}

func (FEEDER_AUDIT) MetaType() string {
	return "FEEDER_AUDIT"
}

// func (v *FEEDER_AUDIT) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type FEEDER_AUDIT_DETAILS struct {
	Type_        Option[string]           `openehr:"type_"`
	SystemId     string                   `openehr:"system_id"`
	Location     Option[PARTY_IDENTIFIED] `openehr:"location"`
	Subject      Option[PARTY_PROXY]      `openehr:"subject"`
	Provider     Option[PARTY_IDENTIFIED] `openehr:"provider"`
	Time         Option[DV_DATE_TIME]     `openehr:"time"`
	VersionId    Option[string]           `openehr:"version_id"`
	OtherDetails Option[ITEM_STRUCTURE]   `openehr:"other_details"`
}

func (FEEDER_AUDIT_DETAILS) MetaType() string {
	return "FEEDER_AUDIT_DETAILS"
}

// func (v *FEEDER_AUDIT_DETAILS) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type PARTY_PROXY interface {
	_PARTY_PROXY()
}

func (PARTY_SELF) _PARTY_PROXY()       {}
func (PARTY_IDENTIFIED) _PARTY_PROXY() {}
func (PARTY_RELATED) _PARTY_PROXY()    {}

type PARTY_SELF struct {
	Type_       Option[string]
	ExternalRef Option[PARTY_REF]
}

func (PARTY_SELF) MetaType() string {
	return "PARTY_SELF"
}

// func (v *PARTY_SELF) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type PARTY_IDENTIFIED struct {
	Type_       Option[string]          `openehr:"type_"`
	ExternalRef Option[PARTY_REF]       `openehr:"external_ref"`
	Name        Option[string]          `openehr:"name"`
	Identifiers Option[[]DV_IDENTIFIER] `openehr:"identifiers"`
}

func (PARTY_IDENTIFIED) MetaType() string {
	return "PARTY_IDENTIFIED"
}

// func (v *PARTY_IDENTIFIED) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type PARTY_RELATED struct {
	Type_        Option[string]          `openehr:"type_"`
	ExternalRef  Option[PARTY_REF]       `openehr:"external_ref"`
	Name         Option[string]          `openehr:"name"`
	Identifiers  Option[[]DV_IDENTIFIER] `openehr:"identifiers"`
	Relationship DV_CODED_TEXT           `openehr:"relationship"`
}

func (PARTY_RELATED) MetaType() string {
	return "PARTY_RELATED"
}

// func (v *PARTY_RELATED) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type PARTICIPATION struct {
	Type_     Option[string]        `openehr:"type_"`
	Function  DV_TEXT               `openehr:"function"`
	Mode      Option[DV_CODED_TEXT] `openehr:"mode"`
	Performer PARTY_PROXY           `openehr:"performer"`
	Time      Option[DV_INTERVAL]   `openehr:"time"`
}

func (PARTICIPATION) MetaType() string {
	return "PARTICIPATION"
}

// func (v *PARTICIPATION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

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

type ITEM_STRUCTURE interface {
	_ITEM_STRUCTURE()
}

func (ITEM_SINGLE) _ITEM_STRUCTURE() {}
func (ITEM_LIST) _ITEM_STRUCTURE()   {}
func (ITEM_TABLE) _ITEM_STRUCTURE()  {}
func (ITEM_TREE) _ITEM_STRUCTURE()   {}

type ITEM_SINGLE struct {
	Type_            Option[string]       `openehr:"type_"`
	Name             DV_TEXT              `openehr:"name"`
	ArchetypeNodeId  string               `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `openehr:"uid"`
	Links            Option[LINK]         `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `openehr:"feeder_audit"`
	Item             ELEMENT              `openehr:"item"`
}

func (ITEM_SINGLE) MetaType() string {
	return "ITEM_SINGLE"
}

type ITEM_LIST struct {
	Type_            Option[DV_TEXT]      `openehr:"type_"`
	ArchetypeNodeId  string               `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `openehr:"uid"`
	Links            Option[[]LINK]       `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `openehr:"feeder_audit"`
	Items            Option[[]ELEMENT]    `openehr:"items"`
}

func (ITEM_LIST) MetaType() string {
	return "ITEM_LIST"
}

// func (v *ITEM_LIST) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type ITEM_TABLE struct {
	Type_            Option[DV_TEXT]      `openehr:"type_"`
	ArchetypeNodeId  string               `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `openehr:"uid"`
	Links            Option[[]LINK]       `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `openehr:"feeder_audit"`
	Rows             Option[[]CLUSTER]    `openehr:"rows"`
}

func (ITEM_TABLE) MetaType() string {
	return "ITEM_TABLE"
}

// func (v *ITEM_TABLE) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type ITEM_TREE struct {
	Type_            Option[DV_TEXT]      `openehr:"type_"`
	ArchetypeNodeId  string               `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `openehr:"uid"`
	Links            Option[[]LINK]       `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `openehr:"feeder_audit"`
	Items            Option[[]ITEM]       `openehr:"items"`
}

func (ITEM_TREE) MetaType() string {
	return "ITEM_TREE"
}

// func (v *ITEM_TREE) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type ITEM interface {
	_ITEM()
}

func (CLUSTER) _ITEM() {}
func (ELEMENT) _ITEM() {}

type CLUSTER struct {
	Type_            Option[string]       `openehr:"type_"`
	Name             DV_TEXT              `openehr:"name"`
	ArchetypeNodeId  string               `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID] `openehr:"uid"`
	Links            Option[[]LINK]       `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]   `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT] `openehr:"feeder_audit"`
	Items            []ITEM               `openehr:"items"`
}

func (CLUSTER) MetaType() string {
	return "CLUSTER"
}

// func (v *CLUSTER) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type ELEMENT struct {
	Type_            Option[string]        `openehr:"type_"`
	Name             DV_TEXT               `openehr:"name"`
	ArchetypeNodeId  string                `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]  `openehr:"uid"`
	Links            Option[[]LINK]        `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]    `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]  `openehr:"feeder_audit"`
	NullFlavour      Option[DV_CODED_TEXT] `openehr:"null_flavour"`
	Value            Option[DATA_VALUE]    `openehr:"value"`
	NullReason       Option[DV_TEXT]       `openehr:"null_reason"`
}

func (ELEMENT) MetaType() string {
	return "ELEMENT"
}

// func (v *ELEMENT) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type HISTORY[T any] struct {
	Type_            Option[string]         `openehr:"type_"`
	Name             DV_TEXT                `openehr:"name"`
	ArchetypeNodeId  string                 `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `openehr:"uid"`
	Links            Option[[]LINK]         `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `openehr:"feeder_audit"`
	Origin           DV_DATE_TIME           `openehr:"origin"`
	Period           Option[DV_DURATION]    `openehr:"period"`
	Duration         Option[DV_DURATION]    `openehr:"duration"`
	Summary          Option[ITEM_STRUCTURE] `openehr:"summary"`
	Events           Option[[]EVENT]        `openehr:"events"`
}

func (HISTORY[T]) MetaType() string {
	return "HISTORY"
}

// func (v *HISTORY[T]) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type EVENT interface {
	_EVENT()
}

func (POINT_EVENT[T]) _EVENT()    {}
func (INTERVAL_EVENT[T]) _EVENT() {}

type POINT_EVENT[T any] struct {
	Type_            string                 `openehr:"type_"`
	Name             DV_TEXT                `openehr:"name"`
	ArchetypeNodeId  string                 `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `openehr:"uid"`
	Links            Option[[]LINK]         `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `openehr:"feeder_audit"`
	Time             DV_DATE_TIME           `openehr:"time"`
	State            Option[ITEM_STRUCTURE] `openehr:"state"`
	Data             T                      `openehr:"data"`
}

func (POINT_EVENT[T]) MetaType() string {
	return "POINT_EVENT"
}

// func (v *POINT_EVENT[T]) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type INTERVAL_EVENT[T any] struct {
	Type_            string                 `openehr:"type_"`
	Name             DV_TEXT                `openehr:"name"`
	ArchetypeNodeId  string                 `openehr:"archetype_node_id"`
	Uid              Option[UID_BASED_ID]   `openehr:"uid"`
	Links            Option[[]LINK]         `openehr:"links"`
	ArchetypeDetails Option[ARCHETYPED]     `openehr:"archetype_details"`
	FeederAudit      Option[FEEDER_AUDIT]   `openehr:"feeder_audit"`
	Time             DV_DATE_TIME           `openehr:"time"`
	State            Option[ITEM_STRUCTURE] `openehr:"state"`
	Data             T                      `openehr:"data"`
	Width            DV_DURATION            `openehr:"width"`
	SampleCount      Option[int64]          `openehr:"sample_count"`
	MathFunction     DV_CODED_TEXT          `openehr:"math_function"`
}

func (INTERVAL_EVENT[T]) MetaType() string {
	return "INTERVAL_EVENT"
}

// func (v *INTERVAL_EVENT[T]) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

// -----------------------------------
// DATA_STRUCTURES
// -----------------------------------

type DATA_VALUE interface {
	_DATA_VALUE()
}

func (DV_BOOLEAN) _DATA_VALUE()                     {}
func (DV_STATE) _DATA_VALUE()                       {}
func (DV_IDENTIFIER) _DATA_VALUE()                  {}
func (DV_TEXT) _DATA_VALUE()                        {}
func (DV_CODED_TEXT) _DATA_VALUE()                  {}
func (DV_PARAGRAPH) _DATA_VALUE()                   {}
func (DV_INTERVAL) _DATA_VALUE()                    {}
func (DV_ORDINAL) _DATA_VALUE()                     {}
func (DV_SCALE) _DATA_VALUE()                       {}
func (DV_QUANTITY) _DATA_VALUE()                    {}
func (DV_COUNT) _DATA_VALUE()                       {}
func (DV_PROPORTION) _DATA_VALUE()                  {}
func (DV_DATE) _DATA_VALUE()                        {}
func (DV_TIME) _DATA_VALUE()                        {}
func (DV_DATE_TIME) _DATA_VALUE()                   {}
func (DV_DURATION) _DATA_VALUE()                    {}
func (DV_PERIODIC_TIME_SPECIFICATION) _DATA_VALUE() {}
func (DV_GENERAL_TIME_SPECIFICATION) _DATA_VALUE()  {}
func (DV_MULTIMEDIA) _DATA_VALUE()                  {}
func (DV_PARSABLE) _DATA_VALUE()                    {}
func (DV_URI) _DATA_VALUE()                         {}
func (DV_EHR_URI) _DATA_VALUE()                     {}

type DV_BOOLEAN struct {
	Type_ Option[string] `openehr:"type_"`
	Value bool           `openehr:"value"`
}

func (DV_BOOLEAN) MetaType() string {
	return "DV_BOOLEAN"
}

// func (v *DV_BOOLEAN) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_STATE struct {
	Type_      Option[string] `openehr:"type_"`
	Value      DV_CODED_TEXT  `openehr:"value"`
	IsTerminal bool           `openehr:"is_terminal"`
}

func (DV_STATE) MetaType() string {
	return "DV_STATE"
}

// func (v *DV_STATE) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_IDENTIFIER struct {
	Type_    Option[string] `openehr:"type_"`
	Issuer   Option[string] `openehr:"issuer"`
	Assigner Option[string] `openehr:"assigner"`
	Id       string         `openehr:"id"`
	Type     Option[string] `openehr:"type"`
}

func (DV_IDENTIFIER) MetaType() string {
	return "DV_IDENTIFIER"
}

// func (v *DV_IDENTIFIER) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_TEXT struct {
	Type_      Option[string]         `openehr:"type_"`
	Value      string                 `openehr:"value"`
	Hyperlink  Option[DV_URI]         `openehr:"hyperlink"`
	Formatting Option[string]         `openehr:"formatting"`
	Mappings   Option[[]TERM_MAPPING] `openehr:"mappings"`
	Language   Option[CODE_PHRASE]    `openehr:"language"`
	Encoding   Option[CODE_PHRASE]    `openehr:"encoding"`
}

func (DV_TEXT) MetaType() string {
	return "DV_TEXT"
}

// func (dvText *DV_TEXT) Unmarshal(data any) error {
// 	if value, found := data["type_"]; found {
// 		v, ok := value.(string)
// 		if !ok {
// 			return ErrBadType
// 		}
// 		dvText.Type_.SetSome(v)
// 	} else {
// 		dvText.Type_.SetNone()
// 	}

// 	if value, found := data["value"]; found {
// 		v, ok := value.(string)
// 		if !ok {
// 			return ErrBadType
// 		}
// 		dvText.Value = v
// 	} else {
// 		return ErrNotFound
// 	}

// 	return nil
// }

type TERM_MAPPING struct {
	Type_   Option[string]        `openehr:"type_"`
	Match   byte                  `openehr:"match"`
	Purpose Option[DV_CODED_TEXT] `openehr:"purpose"`
	Target  CODE_PHRASE           `openehr:"target"`
}

func (TERM_MAPPING) MetaType() string {
	return "TERM_MAPPING"
}

// func (v *TERM_MAPPING) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type CODE_PHRASE struct {
	Type_         Option[string] `openehr:"type_"`
	TerminologyId TERMINOLOGY_ID `openehr:"terminology_id"`
	CodeString    string         `openehr:"code_string"`
	PreferredTerm Option[string] `openehr:"preferred_term"`
}

func (CODE_PHRASE) MetaType() string {
	return "CODE_PHRASE"
}

// func (v *CODE_PHRASE) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_CODED_TEXT struct {
	Type_        Option[string]         `openehr:"type_"`
	Value        string                 `openehr:"value"`
	Hyperlink    Option[DV_URI]         `openehr:"hyperlink"`
	Formatting   Option[string]         `openehr:"formatting"`
	Mappings     Option[[]TERM_MAPPING] `openehr:"mappings"`
	Language     Option[CODE_PHRASE]    `openehr:"language"`
	Encoding     Option[CODE_PHRASE]    `openehr:"encoding"`
	DefiningCode CODE_PHRASE            `openehr:"defining_code"`
}

func (DV_CODED_TEXT) MetaType() string {
	return "DV_CODED_TEXT"
}

// func (v *DV_CODED_TEXT) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_PARAGRAPH struct {
	Type_ Option[string] `openehr:"type_"`
	Items []DV_TEXT      `openehr:"items"`
}

func (DV_PARAGRAPH) MetaType() string {
	return "DV_PARAGRAPH"
}

// func (v *DV_PARAGRAPH) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_ORDERED interface {
	_DV_ORDERED()
}

func (DV_ORDINAL) _DV_ORDERED()    {}
func (DV_SCALE) _DV_ORDERED()      {}
func (DV_QUANTITY) _DV_ORDERED()   {}
func (DV_COUNT) _DV_ORDERED()      {}
func (DV_PROPORTION) _DV_ORDERED() {}
func (DV_DATE) _DV_ORDERED()       {}
func (DV_TIME) _DV_ORDERED()       {}
func (DV_DATE_TIME) _DV_ORDERED()  {}
func (DV_DURATION) _DV_ORDERED()   {}

type DV_INTERVAL struct {
	Type_          Option[string] `openehr:"type_"`
	Lower          DV_ORDERED     `openehr:"lower"`
	Upper          DV_ORDERED     `openehr:"upper"`
	LowerUnbounded bool           `openehr:"lower_unbounded"`
	UpperUnbounded bool           `openehr:"upper_unbounded"`
	LowerIncluded  bool           `openehr:"lower_included"`
	UpperIncluded  bool           `openehr:"upper_included"`
}

func (DV_INTERVAL) MetaType() string {
	return "DV_INTERVAL"
}

// func (v *DV_INTERVAL) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type REFERENCE_RANGE struct {
	Type_   Option[string] `openehr:"type_"`
	Meaning DV_TEXT        `openehr:"meaning"`
	Range   DV_INTERVAL    `openehr:"range"`
}

func (REFERENCE_RANGE) MetaType() string {
	return "REFERENCE_RANGE"
}

// func (v *REFERENCE_RANGE) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_ORDINAL struct {
	Type_                Option[string]          `openehr:"type_"`
	NormalStatus         Option[CODE_PHRASE]     `openehr:"normal_status"`
	NormalRange          Option[DV_INTERVAL]     `openehr:"normal_range"`
	OtherReferenceRanges Option[REFERENCE_RANGE] `openehr:"other_reference_ranges"`
	Symbol               DV_CODED_TEXT           `openehr:"symbol"`
	Value                int64                   `openehr:"value"`
}

func (DV_ORDINAL) MetaType() string {
	return "DV_ORDINAL"
}

// func (v *DV_ORDINAL) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_SCALE struct {
	Type_                Option[string]            `openehr:"type_"`
	NormalStatus         Option[CODE_PHRASE]       `openehr:"normal_status"`
	NormalRange          Option[DV_INTERVAL]       `openehr:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `openehr:"other_reference_ranges"`
	Symbol               DV_CODED_TEXT             `openehr:"symbol"`
	Value                float64                   `openehr:"value"`
}

func (DV_SCALE) MetaType() string {
	return "DV_SCALE"
}

// func (v *DV_SCALE) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_QUANTITY struct {
	Type_                Option[string]            `openehr:"type_"`
	NormalStatus         Option[CODE_PHRASE]       `openehr:"normal_status"`
	MagnitudeStatus      Option[string]            `openehr:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `openehr:"accuracy_is_percent"`
	Accuracy             Option[float64]           `openehr:"accuracy"`
	Magnitude            float64                   `openehr:"magnitude"`
	Precision            Option[int64]             `openehr:"precision"`
	Units                string                    `openehr:"units"`
	UnitsSystem          Option[string]            `openehr:"units_system"`
	UnitsDisplayName     Option[string]            `openehr:"units_display_name"`
	NormalRange          Option[DV_INTERVAL]       `openehr:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `openehr:"other_reference_ranges"`
}

func (DV_QUANTITY) MetaType() string {
	return "DV_QUANTITY"
}

// func (v *DV_QUANTITY) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_COUNT struct {
	Type_                Option[string]            `openehr:"type_"`
	NormalStatus         Option[CODE_PHRASE]       `openehr:"normal_status"`
	MagnitudeStatus      Option[string]            `openehr:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `openehr:"accuracy_is_percent"`
	Accuracy             Option[float64]           `openehr:"accuracy"`
	Magnitude            int64                     `openehr:"magnitude"`
	NormalRange          Option[DV_INTERVAL]       `openehr:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `openehr:"other_reference_ranges"`
}

func (DV_COUNT) MetaType() string {
	return "DV_COUNT"
}

// func (v *DV_COUNT) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_PROPORTION struct {
	Type_                Option[string]            `openehr:"type_"`
	NormalStatus         Option[CODE_PHRASE]       `openehr:"normal_status"`
	MagnitudeStatus      Option[string]            `openehr:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `openehr:"accuracy_is_percent"`
	Accuracy             Option[float64]           `openehr:"accuracy"`
	Numerator            float64                   `openehr:"numerator"`
	Denominator          float64                   `openehr:"denominator"`
	Type                 int64                     `openehr:"type"`
	Precision            Option[int64]             `openehr:"precision"`
	NormalRange          Option[DV_INTERVAL]       `openehr:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `openehr:"other_reference_ranges"`
}

func (DV_PROPORTION) MetaType() string {
	return "DV_PROPORTION"
}

// func (v *DV_PROPORTION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_DATE struct {
	Type_                Option[string]            `openehr:"type_"`
	NormalStatus         Option[CODE_PHRASE]       `openehr:"normal_status"`
	NormalRange          Option[DV_INTERVAL]       `openehr:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `openehr:"other_reference_ranges"`
	MagnitudeStatus      Option[string]            `openehr:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `openehr:"accuracy_is_percent"`
	Accuracy             Option[DV_DURATION]       `openehr:"accuracy"`
	Value                string                    `openehr:"value"`
}

func (DV_DATE) MetaType() string {
	return "DV_DATE"
}

// func (v *DV_DATE) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_TIME struct {
	Type_                Option[string]            `openehr:"type_"`
	NormalStatus         Option[CODE_PHRASE]       `openehr:"normal_status"`
	NormalRange          Option[DV_INTERVAL]       `openehr:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `openehr:"other_reference_ranges"`
	MagnitudeStatus      Option[string]            `openehr:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `openehr:"accuracy_is_percent"`
	Accuracy             Option[DV_DURATION]       `openehr:"accuracy"`
	Value                string                    `openehr:"value"`
}

func (DV_TIME) MetaType() string {
	return "DV_TIME"
}

// func (v *DV_TIME) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_DATE_TIME struct {
	Type_                Option[string]            `openehr:"type_"`
	NormalStatus         Option[CODE_PHRASE]       `openehr:"normal_status"`
	NormalRange          Option[DV_INTERVAL]       `openehr:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `openehr:"other_reference_ranges"`
	MagnitudeStatus      Option[string]            `openehr:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `openehr:"accuracy_is_percent"`
	Accuracy             Option[DV_DURATION]       `openehr:"accuracy"`
	Value                string                    `openehr:"value"`
}

func (DV_DATE_TIME) MetaType() string {
	return "DV_DATE_TIME"
}

// func (v *DV_DATE_TIME) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_DURATION struct {
	Type_                string                    `openehr:"type_"`
	NormalStatus         Option[CODE_PHRASE]       `openehr:"normal_status"`
	NormalRange          Option[DV_INTERVAL]       `openehr:"normal_range"`
	OtherReferenceRanges Option[[]REFERENCE_RANGE] `openehr:"other_reference_ranges"`
	MagnitudeStatus      Option[string]            `openehr:"magnitude_status"`
	AccuracyIsPercent    Option[bool]              `openehr:"accuracy_is_percent"`
	Accuracy             Option[bool]              `openehr:"accuracy"`
	Value                float64                   `openehr:"value"`
}

func (DV_DURATION) MetaType() string {
	return "DV_DURATION"
}

// func (v *DV_DURATION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_PERIODIC_TIME_SPECIFICATION struct {
	Type_ Option[string] `openehr:"type_"`
	Value DV_PARSABLE    `openehr:"value"`
}

func (DV_PERIODIC_TIME_SPECIFICATION) MetaType() string {
	return "DV_PERIODIC_TIME_SPECIFICATION"
}

// func (v *DV_PERIODIC_TIME_SPECIFICATION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_GENERAL_TIME_SPECIFICATION struct {
	Type_ Option[string] `openehr:"type_"`
	Value DV_PARSABLE    `openehr:"value"`
}

func (DV_GENERAL_TIME_SPECIFICATION) MetaType() string {
	return "DV_GENERAL_TIME_SPECIFICATION"
}

// func (v *DV_GENERAL_TIME_SPECIFICATION) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_ENCAPSULATED interface {
	_DV_ENCAPSULATED()
}

func (DV_MULTIMEDIA) _DV_ENCAPSULATED() {}
func (DV_PARSABLE) _DV_ENCAPSULATED()   {}

type DV_MULTIMEDIA struct {
	Type_                   Option[string]        `openehr:"type_"`
	Charset                 Option[CODE_PHRASE]   `openehr:"charset"`
	Language                Option[CODE_PHRASE]   `openehr:"language"`
	AlternateText           Option[string]        `openehr:"alternate_text"`
	Uri                     Option[DV_URI]        `openehr:"uri"`
	Data                    Option[string]        `openehr:"data"`
	MediaType               CODE_PHRASE           `openehr:"media_type"`
	CompressionAlgorithm    Option[CODE_PHRASE]   `openehr:"compression_algorithm"`
	IntegrityCheck          Option[string]        `openehr:"integrity_check"`
	IntegrityCheckAlgorithm Option[DV_MULTIMEDIA] `openehr:"integrity_check_algorithm"`
	Thumbnail               Option[DV_MULTIMEDIA] `openehr:"thumbnail"`
	Size                    int64                 `openehr:"size"`
}

func (DV_MULTIMEDIA) MetaType() string {
	return "DV_MULTIMEDIA"
}

// func (v *DV_MULTIMEDIA) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_PARSABLE struct {
	Type_     Option[string]      `openehr:"type_"`
	Charset   Option[CODE_PHRASE] `openehr:"charset"`
	Language  Option[CODE_PHRASE] `openehr:"language"`
	Value     string              `openehr:"value"`
	Formalism string              `openehr:"formalism"`
}

func (DV_PARSABLE) MetaType() string {
	return "DV_PARSABLE"
}

// func (v *DV_PARSABLE) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_URI struct {
	Type_ Option[string] `openehr:"type_"`
	Value string         `openehr:"value"`
}

func (DV_URI) MetaType() string {
	return "DV_URI"
}

// func (v *DV_URI) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type DV_EHR_URI struct {
	Type_              Option[string] `openehr:"type_"`
	Value              string         `openehr:"value"`
	LocalTerminologyId string         `openehr:"local_terminology_id"`
}

func (DV_EHR_URI) MetaType() string {
	return "DV_EHR_URI"
}

// func (v *DV_EHR_URI) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

// -----------------------------------
// BASE_TYPES
// -----------------------------------

type UID interface {
	_UID()
}

func (ISO_OID) _UID()     {}
func (UUID) _UID()        {}
func (INTERNET_ID) _UID() {}

type ISO_OID struct {
	Type_ Option[string] `openehr:"type_"`
	Value string         `openehr:"value"`
}

func (ISO_OID) MetaType() string {
	return "ISO_OID"
}

// func (v *ISO_OID) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type UUID struct {
	Type_ Option[string] `openehr:"type_"`
	Value string         `openehr:"value"`
}

func (UUID) MetaType() string {
	return "UUID"
}

// func (v *UUID) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type INTERNET_ID struct {
	Type_ Option[string] `openehr:"type_"`
	Value string         `openehr:"value"`
}

func (INTERNET_ID) MetaType() string {
	return "INTERNET_ID"
}

// func (v *INTERNET_ID) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type OBJECT_ID interface {
	_OBJECT_ID()
}

func (HIER_OBJECT_ID) _OBJECT_ID()    {}
func (OBJECT_VERSION_ID) _OBJECT_ID() {}
func (ARCHETYPE_ID) _OBJECT_ID()      {}
func (TEMPLATE_ID) _OBJECT_ID()       {}
func (GENERIC_ID) _OBJECT_ID()        {}

type UID_BASED_ID interface {
	_UID_BASED_ID()
}

func (HIER_OBJECT_ID) _UID_BASED_ID()    {}
func (OBJECT_VERSION_ID) _UID_BASED_ID() {}

type HIER_OBJECT_ID struct {
	Type_ Option[string] `openehr:"type_"`
	Value string         `openehr:"value"`
}

func (HIER_OBJECT_ID) MetaType() string {
	return "HIER_OBJECT_ID"
}

// func (v *HIER_OBJECT_ID) Unmarshal(data any) error {
// 	return nil
// }

type OBJECT_VERSION_ID struct {
	Type_ Option[string] `openehr:"type_"`
	Value string         `openehr:"value"`
}

func (OBJECT_VERSION_ID) MetaType() string {
	return "OBJECT_VERSION_ID"
}

// func (v *OBJECT_VERSION_ID) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type ARCHETYPE_ID struct {
	Type_ Option[string] `openehr:"type_"`
	Value string         `openehr:"value"`
}

func (ARCHETYPE_ID) MetaType() string {
	return "ARCHETYPE_ID"
}

// func (v *ARCHETYPE_ID) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type TEMPLATE_ID struct {
	Type_ Option[string] `openehr:"type_"`
	Value string         `openehr:"value"`
}

func (TEMPLATE_ID) MetaType() string {
	return "TEMPLATE_ID"
}

// func (v *TEMPLATE_ID) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type TERMINOLOGY_ID struct {
	Type_ Option[string] `openehr:"type_"`
	Value string         `openehr:"value"`
}

func (TERMINOLOGY_ID) MetaType() string {
	return "TERMINOLOGY_ID"
}

// func (v *TERMINOLOGY_ID) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type GENERIC_ID struct {
	Type_  Option[string] `openehr:"type_"`
	Value  string         `openehr:"value"`
	Schema string         `openehr:"schema"`
}

func (GENERIC_ID) MetaType() string {
	return "GENERIC_ID"
}

// func (v *GENERIC_ID) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type OBJECT_REF struct {
	Type_     Option[string] `openehr:"object_ref"`
	Namespace string         `openehr:"namespace"`
	Type      string         `openehr:"type"`
	Id        OBJECT_ID      `openehr:"id"`
}

func (OBJECT_REF) MetaType() string {
	return "OBJECT_REF"
}

// func (v *OBJECT_REF) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type PARTY_REF struct {
	Type_     Option[string] `openehr:"type_"`
	Namespace string         `openehr:"namespace"`
	Type      string         `openehr:"type"`
	Path      Option[string] `openehr:"path"`
	Id        UID_BASED_ID   `openehr:"id"`
}

func (PARTY_REF) MetaType() string {
	return "PARTY_REF"
}

// func (v *PARTY_REF) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }

type LOCATABLE_REF struct {
	Type_     Option[string] `openehr:"type_"`
	Namespace string         `openehr:"namespace"`
	Type      string         `openehr:"type"`
	Path      Option[string] `openehr:"path"`
	Id        UID_BASED_ID   `openehr:"id"`
}

func (LOCATABLE_REF) MetaType() string {
	return "LOCATABLE_REF"
}

// func (v *LOCATABLE_REF) Unmarshal(data any) error {
// 	return errors.New("not implemented")
// }
