package audit

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID      `json:"id"`
	ActorID   uuid.UUID      `json:"actor_id,omitempty"`   // References Account.ID or SYSTEM_USER_ID for system actions
	ActorType string         `json:"actor_type,omitempty"` // "user", "system", "client_app", etc.
	Resource  string         `json:"resource"`             // What was accessed/modified
	Action    string         `json:"action"`
	Success   bool           `json:"success"`
	IPAddress net.IP         `json:"ip_address,omitempty"`
	UserAgent string         `json:"user_agent,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

type Resource string

const (
	ResourceEHR                         Resource = "ehr"
	ResourceEHRStatus                   Resource = "ehr_status"
	ResourceVersionedEHRStatus          Resource = "versioned_ehr_status"
	ResourceVersionedEHRStatusVersion   Resource = "versioned_ehr_status_version"
	ResourceComposition                 Resource = "composition"
	ResourceVersionedComposition        Resource = "versioned_composition"
	ResourceVersionedCompositionVersion Resource = "versioned_composition_version"
	ResourceDirectory                   Resource = "directory"
	ResourceFolder                      Resource = "folder"
	ResourceContribution                Resource = "contribution"
	ResourceAgent                       Resource = "agent"
	ResourceGroup                       Resource = "group"
	ResourcePerson                      Resource = "person"
	ResourceOrganisation                Resource = "organisation"
	ResourceRole                        Resource = "role"
	ResourceVersionedParty              Resource = "versioned_party"
	ResourceVersionedPartyVersion       Resource = "versioned_party_version"
	ResourceTemplate                    Resource = "template"
	ResourceQuery                       Resource = "query"
	ResourceWebhook                     Resource = "webhook"
)

type Action string

const (
	ActionCreate  Action = "create"
	ActionRead    Action = "read"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionExecute Action = "execute"
)
