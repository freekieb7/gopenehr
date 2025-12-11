package webhook

type Event struct {
	Version string         `json:"version"`
	Type    string         `json:"type"`
	Source  string         `json:"source"`
	ID      string         `json:"id"`
	Time    string         `json:"time"`
	Data    map[string]any `json:"data"`
}

type EventType string

const (
	EventTypeEHRCreated EventType = "ehr.created"
	EventTypeEHRDeleted EventType = "ehr.deleted"

	EventTypeEHRStatusUpdated EventType = "ehr_status.updated"

	EventTypeCompositionCreated EventType = "composition.created"
	EventTypeCompositionUpdated EventType = "composition.updated"
	EventTypeCompositionDeleted EventType = "composition.deleted"

	EventTypeDirectoryCreated EventType = "directory.created"
	EventTypeDirectoryUpdated EventType = "directory.updated"
	EventTypeDirectoryDeleted EventType = "directory.deleted"

	EventTypePersonCreated EventType = "person.created"
	EventTypePersonUpdated EventType = "person.updated"
	EventTypePersonDeleted EventType = "person.deleted"

	EventTypeAgentCreated EventType = "agent.created"
	EventTypeAgentUpdated EventType = "agent.updated"
	EventTypeAgentDeleted EventType = "agent.deleted"

	EventTypeGroupCreated EventType = "group.created"
	EventTypeGroupUpdated EventType = "group.updated"
	EventTypeGroupDeleted EventType = "group.deleted"

	EventTypeOrganisationCreated EventType = "organisation.created"
	EventTypeOrganisationUpdated EventType = "organisation.updated"
	EventTypeOrganisationDeleted EventType = "organisation.deleted"

	EventTypeRoleCreated EventType = "role.created"
	EventTypeRoleUpdated EventType = "role.updated"
	EventTypeRoleDeleted EventType = "role.deleted"

	EventTypeQueryExecuted EventType = "query.executed"
	EventTypeQueryStored   EventType = "query.stored"
)

var EventTypes = map[EventType]string{
	EventTypeEHRCreated:          "EHR Created",
	EventTypeEHRDeleted:          "EHR Deleted",
	EventTypeEHRStatusUpdated:    "EHR Status Updated",
	EventTypeCompositionCreated:  "Composition Created",
	EventTypeCompositionDeleted:  "Composition Deleted",
	EventTypePersonCreated:       "Person Created",
	EventTypePersonUpdated:       "Person Updated",
	EventTypePersonDeleted:       "Person Deleted",
	EventTypeAgentCreated:        "Agent Created",
	EventTypeAgentUpdated:        "Agent Updated",
	EventTypeAgentDeleted:        "Agent Deleted",
	EventTypeGroupCreated:        "Group Created",
	EventTypeGroupUpdated:        "Group Updated",
	EventTypeGroupDeleted:        "Group Deleted",
	EventTypeOrganisationCreated: "Organisation Created",
	EventTypeOrganisationUpdated: "Organisation Updated",
	EventTypeOrganisationDeleted: "Organisation Deleted",
	EventTypeRoleCreated:         "Role Created",
	EventTypeRoleUpdated:         "Role Updated",
	EventTypeRoleDeleted:         "Role Deleted",
	EventTypeQueryExecuted:       "Query Executed",
	EventTypeQueryStored:         "Query Stored",
}

func IsValidEventType(event EventType) bool {
	_, exists := EventTypes[event]
	return exists
}
