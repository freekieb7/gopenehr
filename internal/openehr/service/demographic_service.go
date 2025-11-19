package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/openehr"
	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/google/uuid"
)

var (
	ErrAgentAlreadyExists        = fmt.Errorf("agent with the given UID already exists")
	ErrAgentNotFound             = fmt.Errorf("agent not found")
	ErrGroupAlreadyExists        = fmt.Errorf("group with the given UID already exists")
	ErrGroupNotFound             = fmt.Errorf("group not found")
	ErrOrganisationAlreadyExists = fmt.Errorf("organisation with the given UID already exists")
	ErrOrganisationNotFound      = fmt.Errorf("organisation not found")
	ErrPersonAlreadyExists       = fmt.Errorf("person with the given UID already exists")
	ErrPersonNotFound            = fmt.Errorf("person not found")
	ErrRoleAlreadyExists         = fmt.Errorf("role with the given UID already exists")
	ErrRoleNotFound              = fmt.Errorf("role not found")
	ErrVersionedPartyNotFound    = fmt.Errorf("versioned party not found")
)

type DemographicService struct {
	Logger *slog.Logger
	DB     *database.Database
}

func NewDemographicService(logger *slog.Logger, db *database.Database) DemographicService {
	return DemographicService{
		Logger: logger,
		DB:     db,
	}
}

func (s *DemographicService) CreateAgent(ctx context.Context, agent openehr.AGENT) (openehr.AGENT, error) {
	// If UID is not provided, generate a new one
	if !agent.UID.E {
		agent.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", uuid.NewString(), config.SYSTEM_ID_GOPENEHR),
			},
		})
	} else {
		// Convert HIER_OBJECT_ID to OBJECT_VERSION_ID if necessary
		switch agent.UID.V.Value.(type) {
		case *openehr.HIER_OBJECT_ID:
			agent.UID.V.Value = &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", agent.UID.V.Value.(*openehr.HIER_OBJECT_ID).Value, config.SYSTEM_ID_GOPENEHR),
			}
		case *openehr.OBJECT_VERSION_ID:
		// valid type
		default:
			return openehr.AGENT{}, errors.New("agent UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID")
		}

		// Check if agent with the same UID already exists
		var exists bool
		err := s.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM openehr.tbl_agent WHERE id = $1)", agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value).Scan(&exists)
		if err != nil {
			return openehr.AGENT{}, fmt.Errorf("failed to check existing agent: %w", err)
		}
		if exists {
			return openehr.AGENT{}, ErrAgentAlreadyExists
		}
	}

	// Build Versioned object for the agent
	versionedParty := openehr.VERSIONED_PARTY{
		UID: openehr.HIER_OBJECT_ID{
			Value: agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.AGENT_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Ensures that the versioned party is valid
	versionedParty.SetModelName()
	if errs := versionedParty.Validate("$"); len(errs) > 0 {
		return openehr.AGENT{}, fmt.Errorf("validation errors for versioned party: %v", errs)
	}

	// Build original version of the agent
	agentVersion := openehr.ORIGINAL_VERSION{
		UID: *agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &agent,
	}

	// Ensures that the agent version is valid
	agentVersion.SetModelName()
	if errs := agentVersion.Validate("$"); len(errs) > 0 {
		return openehr.AGENT{}, fmt.Errorf("validation errors for agent version: %v", errs)
	}

	// Create contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.AGENT_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &agentVersion.UID,
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_CREATION),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Agent created",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	// Ensures that the contribution is valid
	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.AGENT{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, data) VALUES ($1, $2, $3)",
		versionedParty.UID.Value, openehr.AGENT_MODEL_NAME, versionedParty)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}

	// Insert agent
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_agent (id, versioned_object_id, data) VALUES ($1, $2, $3)",
		agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, versionedParty.UID.Value, agentVersion)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert agent: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return agent, nil
}

// GetAgent retrieves the agent as the latest version when providing versioned object id, or the specified ID version.
func (s *DemographicService) GetAgent(ctx context.Context, uidBasedID string) (openehr.AGENT, error) {
	var agent openehr.AGENT

	query := "SELECT data FROM openehr.tbl_agent WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM openehr.tbl_agent WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&agent)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.AGENT{}, ErrAgentNotFound
		}
		return openehr.AGENT{}, fmt.Errorf("failed to get agent by ID: %w", err)
	}
	return agent, nil
}

// GetAgentAsJSON retrieves the agent as the latest version when providing versioned object id, or the specified ID version, and returns it as the raw JSON representation.
// This is useful for scenarios where the raw JSON is needed without unmarshalling into Go structs.
func (s *DemographicService) GetAgentAsJSON(ctx context.Context, uidBasedID string) ([]byte, error) {
	var rawAgentJSON []byte

	query := "SELECT data FROM openehr.tbl_agent WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM openehr.tbl_agent WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&rawAgentJSON)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrAgentNotFound
		}
		return nil, fmt.Errorf("failed to get agent by ID: %w", err)
	}
	return rawAgentJSON, nil
}

func (s *DemographicService) UpdateAgent(ctx context.Context, agent openehr.AGENT) (openehr.AGENT, error) {
	// Validate agent
	if !agent.UID.E {
		return openehr.AGENT{}, fmt.Errorf("agent UID is required for update")
	}

	agentVersionID, ok := agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return openehr.AGENT{}, fmt.Errorf("agent UID must be of type OBJECT_VERSION_ID, got %T", agent.UID.V.Value)
	}

	if errs := agent.Validate("$"); len(errs) > 0 {
		return openehr.AGENT{}, fmt.Errorf("validation errors for agent: %v", errs)
	}

	// Prepare contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.AGENT_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: agentVersionID,
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Agent updated",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.AGENT{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Update agent
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_agent (id, versioned_object_id, data) VALUES ($1, $2, $3)",
		agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(), agent)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert agent: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return agent, nil
}

func (s *DemographicService) DeleteAgent(ctx context.Context, versionedObjectID string) error {
	// Create contribution for deletion
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.AGENT_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &openehr.OBJECT_VERSION_ID{
						Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
						Value: versionedObjectID,
					},
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_DELETED),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Agent deleted",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}
	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Start transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	// Delete agent
	if _, err := tx.Exec(ctx, "DELETE FROM openehr.tbl_versioned_object WHERE id = $1 AND type = $2", versionedObjectID, openehr.AGENT_MODEL_NAME); err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DemographicService) CreateGroup(ctx context.Context, group openehr.GROUP) (openehr.GROUP, error) {
	// If UID is not provided, generate a new one
	if !group.UID.E {
		group.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", uuid.NewString(), config.SYSTEM_ID_GOPENEHR),
			},
		})
	} else {
		// Convert HIER_OBJECT_ID to OBJECT_VERSION_ID if necessary
		switch group.UID.V.Value.(type) {
		case *openehr.HIER_OBJECT_ID:
			group.UID.V.Value = &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", group.UID.V.Value.(*openehr.HIER_OBJECT_ID).Value, config.SYSTEM_ID_GOPENEHR),
			}
		case *openehr.OBJECT_VERSION_ID:
			// valid type
		default:
			return openehr.GROUP{}, errors.New("group UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID")
		}

		// Check if group with the same UID already exists
		var exists bool
		err := s.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM openehr.tbl_group WHERE id = $1)", group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value).Scan(&exists)
		if err != nil {
			return openehr.GROUP{}, fmt.Errorf("failed to check existing group: %w", err)
		}
		if exists {
			return openehr.GROUP{}, ErrGroupAlreadyExists
		}
	}

	// Build Versioned object for the group
	versionedParty := openehr.VERSIONED_PARTY{
		UID: openehr.HIER_OBJECT_ID{
			Value: group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.GROUP_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Ensures that the versioned party is valid
	versionedParty.SetModelName()
	if errs := versionedParty.Validate("$"); len(errs) > 0 {
		return openehr.GROUP{}, fmt.Errorf("validation errors for versioned party: %v", errs)
	}

	// Build original version of the group
	groupVersion := openehr.ORIGINAL_VERSION{
		UID: *group.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &group,
	}

	// Ensures that the group version is valid
	groupVersion.SetModelName()
	if errs := groupVersion.Validate("$"); len(errs) > 0 {
		return openehr.GROUP{}, fmt.Errorf("validation errors for group version: %v", errs)
	}

	// Create contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.GROUP_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &groupVersion.UID,
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_CREATION),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Group created",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	// Ensures that the contribution is valid
	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.GROUP{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, data) VALUES ($1, $2, $3)",
		versionedParty.UID.Value, openehr.GROUP_MODEL_NAME, versionedParty)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}

	// Insert group
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_group (id, versioned_object_id, data) VALUES ($1, $2, $3)",
		group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, versionedParty.UID.Value, groupVersion)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert group: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return group, nil
}

func (s *DemographicService) GetGroup(ctx context.Context, uidBasedID string) (openehr.GROUP, error) {
	var group openehr.GROUP

	query := "SELECT data FROM openehr.tbl_group WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM openehr.tbl_group WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&group)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.GROUP{}, ErrGroupNotFound
		}
		return openehr.GROUP{}, fmt.Errorf("failed to get group by ID: %w", err)
	}
	return group, nil
}

func (s *DemographicService) GetGroupAsJSON(ctx context.Context, uidBasedID string) ([]byte, error) {
	var rawGroupJSON []byte

	query := "SELECT data FROM openehr.tbl_group WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM openehr.tbl_group WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&rawGroupJSON)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, fmt.Errorf("failed to get group by ID: %w", err)
	}
	return rawGroupJSON, nil
}

func (s *DemographicService) UpdateGroup(ctx context.Context, group openehr.GROUP) (openehr.GROUP, error) {
	// Validate group
	if !group.UID.E {
		return openehr.GROUP{}, fmt.Errorf("group UID is required for update")
	}

	groupVersionID, ok := group.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return openehr.GROUP{}, fmt.Errorf("group UID must be of type OBJECT_VERSION_ID, got %T", group.UID.V.Value)
	}

	if errs := group.Validate("$"); len(errs) > 0 {
		return openehr.GROUP{}, fmt.Errorf("validation errors for group: %v", errs)
	}

	// Prepare contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.GROUP_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: groupVersionID,
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Group updated",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.GROUP{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Update group
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_group (id, versioned_object_id, data) VALUES ($1, $2, $3)",
		group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(), group)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert group: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return group, nil
}

func (s *DemographicService) DeleteGroup(ctx context.Context, versionedObjectID string) error {
	// Create contribution for deletion
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.GROUP_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &openehr.OBJECT_VERSION_ID{
						Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
						Value: versionedObjectID,
					},
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_DELETED),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Group deleted",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}
	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Start transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	// Delete group
	if _, err := tx.Exec(ctx, "DELETE FROM openehr.tbl_versioned_object WHERE id = $1 AND type = $2", versionedObjectID, openehr.GROUP_MODEL_NAME); err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DemographicService) CreatePerson(ctx context.Context, person openehr.PERSON) (openehr.PERSON, error) {
	// If UID is not provided, generate a new one
	if !person.UID.E {
		person.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", uuid.NewString(), config.SYSTEM_ID_GOPENEHR),
			},
		})
	} else {
		// Convert HIER_OBJECT_ID to OBJECT_VERSION_ID if necessary
		switch person.UID.V.Value.(type) {
		case *openehr.HIER_OBJECT_ID:
			person.UID.V.Value = &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", person.UID.V.Value.(*openehr.HIER_OBJECT_ID).Value, config.SYSTEM_ID_GOPENEHR),
			}
		case *openehr.OBJECT_VERSION_ID:
			// valid type
		default:
			return openehr.PERSON{}, errors.New("person UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID")
		}

		// Check if person with the same UID already exists
		var exists bool
		err := s.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM openehr.tbl_person WHERE id = $1)", person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value).Scan(&exists)
		if err != nil {
			return openehr.PERSON{}, fmt.Errorf("failed to check existing person: %w", err)
		}
		if exists {
			return openehr.PERSON{}, ErrPersonAlreadyExists
		}
	}

	// Build Versioned object for the person
	versionedParty := openehr.VERSIONED_PARTY{
		UID: openehr.HIER_OBJECT_ID{
			Value: person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.PERSON_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Ensures that the versioned party is valid
	versionedParty.SetModelName()
	if errs := versionedParty.Validate("$"); len(errs) > 0 {
		return openehr.PERSON{}, fmt.Errorf("validation errors for versioned party: %v", errs)
	}

	// Build original version of the person
	personVersion := openehr.ORIGINAL_VERSION{
		UID: *person.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &person,
	}

	// Ensures that the person version is valid
	personVersion.SetModelName()
	if errs := personVersion.Validate("$"); len(errs) > 0 {
		return openehr.PERSON{}, fmt.Errorf("validation errors for person version: %v", errs)
	}

	// Create contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.PERSON_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &personVersion.UID,
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_CREATION),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Person created",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	// Ensures that the contribution is valid
	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.PERSON{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, data) VALUES ($1, $2, $3)",
		versionedParty.UID.Value, openehr.PERSON_MODEL_NAME, versionedParty)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}

	// Insert person
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_person (id, versioned_object_id, data) VALUES ($1, $2, $3)",
		person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, versionedParty.UID.Value, personVersion)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert person: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return person, nil
}

func (s *DemographicService) GetPerson(ctx context.Context, uidBasedID string) (openehr.PERSON, error) {
	var person openehr.PERSON

	query := "SELECT data FROM openehr.tbl_person WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM openehr.tbl_person WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&person)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.PERSON{}, ErrPersonNotFound
		}
		return openehr.PERSON{}, fmt.Errorf("failed to get person by ID: %w", err)
	}
	return person, nil
}

func (s *DemographicService) GetPersonAsJSON(ctx context.Context, uidBasedID string) ([]byte, error) {
	var rawPersonJSON []byte

	query := "SELECT data FROM openehr.tbl_person WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM openehr.tbl_person WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&rawPersonJSON)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrPersonNotFound
		}
		return nil, fmt.Errorf("failed to get person by ID: %w", err)
	}
	return rawPersonJSON, nil
}

func (s *DemographicService) UpdatePerson(ctx context.Context, person openehr.PERSON) (openehr.PERSON, error) {
	// Validate person
	if !person.UID.E {
		return openehr.PERSON{}, fmt.Errorf("person UID is required for update")
	}

	personVersionID, ok := person.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return openehr.PERSON{}, fmt.Errorf("person UID must be of type OBJECT_VERSION_ID, got %T", person.UID.V.Value)
	}

	if errs := person.Validate("$"); len(errs) > 0 {
		return openehr.PERSON{}, fmt.Errorf("validation errors for person: %v", errs)
	}

	// Prepare contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.PERSON_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: personVersionID,
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Person updated",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.PERSON{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Update person
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_person (id, versioned_object_id, data) VALUES ($1, $2, $3)",
		person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(), person)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert person: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return person, nil
}

func (s *DemographicService) DeletePerson(ctx context.Context, versionedObjectID string) error {
	// Create contribution for deletion
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.PERSON_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &openehr.OBJECT_VERSION_ID{
						Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
						Value: versionedObjectID,
					},
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_DELETED),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Person deleted",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}
	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Start transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	// Delete person
	if _, err := tx.Exec(ctx, "DELETE FROM openehr.tbl_versioned_object WHERE id = $1 AND type = $2", versionedObjectID, openehr.PERSON_MODEL_NAME); err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DemographicService) CreateOrganisation(ctx context.Context, organisation openehr.ORGANISATION) (openehr.ORGANISATION, error) {
	// If UID is not provided, generate a new one
	if !organisation.UID.E {
		organisation.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", uuid.NewString(), config.SYSTEM_ID_GOPENEHR),
			},
		})
	} else {
		// Convert HIER_OBJECT_ID to OBJECT_VERSION_ID if necessary
		switch organisation.UID.V.Value.(type) {
		case *openehr.HIER_OBJECT_ID:
			organisation.UID.V.Value = &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", organisation.UID.V.Value.(*openehr.HIER_OBJECT_ID).Value, config.SYSTEM_ID_GOPENEHR),
			}
		case *openehr.OBJECT_VERSION_ID:
			// valid type
		default:
			return openehr.ORGANISATION{}, errors.New("organisation UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID")
		}

		// Check if organisation with the same UID already exists
		var exists bool
		err := s.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM openehr.tbl_organisation WHERE id = $1)", organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value).Scan(&exists)
		if err != nil {
			return openehr.ORGANISATION{}, fmt.Errorf("failed to check existing organisation: %w", err)
		}
		if exists {
			return openehr.ORGANISATION{}, ErrOrganisationAlreadyExists
		}
	}

	// Build Versioned object for the organisation
	versionedParty := openehr.VERSIONED_PARTY{
		UID: openehr.HIER_OBJECT_ID{
			Value: organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.ORGANISATION_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Ensures that the versioned party is valid
	versionedParty.SetModelName()
	if errs := versionedParty.Validate("$"); len(errs) > 0 {
		return openehr.ORGANISATION{}, fmt.Errorf("validation errors for versioned party: %v", errs)
	}

	// Build original version of the organisation
	organisationVersion := openehr.ORIGINAL_VERSION{
		UID: *organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &organisation,
	}

	// Ensures that the organisation version is valid
	organisationVersion.SetModelName()
	if errs := organisationVersion.Validate("$"); len(errs) > 0 {
		return openehr.ORGANISATION{}, fmt.Errorf("validation errors for organisation version: %v", errs)
	}

	// Create contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.ORGANISATION_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &organisationVersion.UID,
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_CREATION),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Organisation created",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	// Ensures that the contribution is valid
	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.ORGANISATION{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, data) VALUES ($1, $2, $3)",
		versionedParty.UID.Value, openehr.ORGANISATION_MODEL_NAME, versionedParty)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}

	// Insert organisation
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_organisation (id, versioned_object_id, data) VALUES ($1, $2, $3)",
		organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, versionedParty.UID.Value, organisationVersion)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert organisation: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return organisation, nil
}

func (s *DemographicService) GetOrganisation(ctx context.Context, uidBasedID string) (openehr.ORGANISATION, error) {
	var organisation openehr.ORGANISATION

	query := "SELECT data FROM openehr.tbl_organisation WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM openehr.tbl_organisation WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&organisation)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.ORGANISATION{}, ErrOrganisationNotFound
		}
		return openehr.ORGANISATION{}, fmt.Errorf("failed to get organisation by ID: %w", err)
	}
	return organisation, nil
}

func (s *DemographicService) GetOrganisationAsJSON(ctx context.Context, uidBasedID string) ([]byte, error) {
	var rawOrganisationJSON []byte

	query := "SELECT data FROM openehr.tbl_organisation WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM openehr.tbl_organisation WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&rawOrganisationJSON)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrOrganisationNotFound
		}
		return nil, fmt.Errorf("failed to get organisation by ID: %w", err)
	}
	return rawOrganisationJSON, nil
}

func (s *DemographicService) UpdateOrganisation(ctx context.Context, organisation openehr.ORGANISATION) error {
	// Validate organisation
	if !organisation.UID.E {
		return fmt.Errorf("organisation UID is required for update")
	}

	organisationVersionID, ok := organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return fmt.Errorf("organisation UID must be of type OBJECT_VERSION_ID, got %T", organisation.UID.V.Value)
	}

	if errs := organisation.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for organisation: %v", errs)
	}

	// Prepare contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.ORGANISATION_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: organisationVersionID,
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Organisation updated",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Update organisation
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_organisation (id, versioned_object_id, data) VALUES ($1, $2, $3)",
		organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(), organisation)
	if err != nil {
		return fmt.Errorf("failed to insert organisation: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DemographicService) DeleteOrganisation(ctx context.Context, versionedObjectID string) error {
	// Create contribution for deletion
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.ORGANISATION_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &openehr.OBJECT_VERSION_ID{
						Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
						Value: versionedObjectID,
					},
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_DELETED),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Organisation deleted",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}
	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Start transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	// Delete organisation
	if _, err := tx.Exec(ctx, "DELETE FROM openehr.tbl_versioned_object WHERE id = $1 AND type = $2", versionedObjectID, openehr.ORGANISATION_MODEL_NAME); err != nil {
		return fmt.Errorf("failed to delete organisation: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DemographicService) CreateRole(ctx context.Context, role openehr.ROLE) (openehr.ROLE, error) {
	// If UID is not provided, generate a new one
	if !role.UID.E {
		role.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::v1", uuid.NewString(), config.SYSTEM_ID_GOPENEHR),
			},
		})
	} else {
		// Check if role with the same UID already exists
		var exists bool
		err := s.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM openehr.tbl_role WHERE id = $1)", role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value).Scan(&exists)
		if err != nil {
			return openehr.ROLE{}, fmt.Errorf("failed to check existing role: %w", err)
		}
		if exists {
			return openehr.ROLE{}, ErrRoleAlreadyExists
		}
	}

	// Build Versioned object for the role
	versionedParty := openehr.VERSIONED_PARTY{
		UID: openehr.HIER_OBJECT_ID{
			Value: role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.ROLE_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Ensures that the versioned party is valid
	versionedParty.SetModelName()
	if errs := versionedParty.Validate("$"); len(errs) > 0 {
		return openehr.ROLE{}, fmt.Errorf("validation errors for versioned party: %v", errs)
	}

	// Build original version of the role
	roleVersion := openehr.ORIGINAL_VERSION{
		UID: *role.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &role,
	}

	// Ensures that the role version is valid
	roleVersion.SetModelName()
	if errs := roleVersion.Validate("$"); len(errs) > 0 {
		return openehr.ROLE{}, fmt.Errorf("validation errors for role version: %v", errs)
	}

	// Create contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.ROLE_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &roleVersion.UID,
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_CREATION),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Role created",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	// Ensures that the contribution is valid
	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.ROLE{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, data) VALUES ($1, $2, $3)",
		versionedParty.UID.Value, openehr.ROLE_MODEL_NAME, versionedParty)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}

	// Insert role
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_role (id, versioned_object_id, data) VALUES ($1, $2, $3)",
		roleVersion.UID.Value, versionedParty.UID.Value, role)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert role: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return role, nil
}

func (s *DemographicService) GetRole(ctx context.Context, uidBasedID string) (openehr.ROLE, error) {
	var role openehr.ROLE

	query := "SELECT data FROM openehr.tbl_role WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM openehr.tbl_role WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&role)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.ROLE{}, ErrRoleNotFound
		}

		return openehr.ROLE{}, fmt.Errorf("failed to get role by ID: %w", err)
	}
	return role, nil
}

func (s *DemographicService) GetRoleAsJSON(ctx context.Context, uidBasedID string) ([]byte, error) {
	var rawRoleJSON []byte

	query := "SELECT data FROM openehr.tbl_role WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM openehr.tbl_role WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&rawRoleJSON)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to get role by ID: %w", err)
	}
	return rawRoleJSON, nil
}

func (s *DemographicService) UpdateRole(ctx context.Context, role openehr.ROLE) (openehr.ROLE, error) {
	// Validate role
	if !role.UID.E {
		return openehr.ROLE{}, fmt.Errorf("role UID is required for update")
	}

	roleVersionID, ok := role.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return openehr.ROLE{}, fmt.Errorf("role UID must be of type OBJECT_VERSION_ID, got %T", role.UID.V.Value)
	}

	if errs := role.Validate("$"); len(errs) > 0 {
		return openehr.ROLE{}, fmt.Errorf("validation errors for role: %v", errs)
	}

	// Prepare contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.ROLE_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: roleVersionID,
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Role updated",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.ROLE{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Update role
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_role (id, versioned_object_id, data) VALUES ($1, $2, $3)",
		role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(), role)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert role: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return role, nil
}

func (s *DemographicService) DeleteRole(ctx context.Context, versionedObjectID string) error {
	// Create contribution for deletion
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.ROLE_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &openehr.OBJECT_VERSION_ID{
						Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
						Value: versionedObjectID,
					},
				},
			},
		},
		Audit: openehr.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: openehr.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			ChangeType: openehr.DV_CODED_TEXT{
				Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_DELETED),
				DefiningCode: openehr.CODE_PHRASE{
					CodeString: terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
					TerminologyID: openehr.TERMINOLOGY_ID{
						Value: terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR,
					},
				},
			},
			Description: util.Some(openehr.DV_TEXT{
				Value: "Role deleted",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}
	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Start transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	// Delete role
	if _, err := tx.Exec(ctx, "DELETE FROM openehr.tbl_versioned_object WHERE id = $1 AND type = $2", versionedObjectID, openehr.ROLE_MODEL_NAME); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DemographicService) GetVersionedPartyAsJSON(ctx context.Context, uidBasedID string) ([]byte, error) {
	var rawVersionedPartyJSON []byte

	query := "SELECT data FROM openehr.tbl_versioned_object WHERE id = $1 AND data->>'_type' = $2 ORDER BY created_at DESC LIMIT 1"

	err := s.DB.QueryRow(ctx, query, uidBasedID, openehr.VERSIONED_PARTY_MODEL_NAME).Scan(&rawVersionedPartyJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to get versioned party by ID: %w", err)
	}
	return rawVersionedPartyJSON, nil
}

func (s *DemographicService) GetVersionedPartyRevisionHistoryAsJSON(ctx context.Context, versionedObjectID string) ([]byte, error) {
	// Fetch Revision History
	// Build array of REVISION_HISTORY_ITEM objects
	query := `
		SELECT jsonb_build_object(
			'items', jsonb_agg(
				jsonb_build_object(
					'version_id', jsonb_build_object(
						'value', version_id
					),
					'audits', audits
				)
				ORDER BY version_id
			)
		)
		FROM (
			SELECT 
				version->'id'->>'value' as version_id,
				jsonb_agg(c.data->'audit' ORDER BY c.data->'audit'->'time_committed'->>'value') as audits
			FROM openehr.tbl_contribution c,
				jsonb_array_elements(c.data->'versions') as version
			WHERE c.ehr_id IS NULL
				AND version->>'type' = ANY($1::text[])
				AND version->'id'->>'value' LIKE $2 || '%'
			GROUP BY version->'id'->>'value'
		) grouped
	`
	args := []any{[]string{openehr.AGENT_MODEL_NAME, openehr.PERSON_MODEL_NAME, openehr.GROUP_MODEL_NAME, openehr.ORGANISATION_MODEL_NAME}, versionedObjectID}
	row := s.DB.QueryRow(ctx, query, args...)

	var rawRevisionHistoryDataJSON []byte
	if err := row.Scan(&rawRevisionHistoryDataJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrCompositionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Revision History from database: %w", err)
	}

	return rawRevisionHistoryDataJSON, nil
}

func (s *DemographicService) GetVersionedPartyVersionAsJSON(ctx context.Context, versionedObjectID string, filterAtTime time.Time, filterVersionID string) ([]byte, error) {
	// Fetch EHR Status at given time
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`
		SELECT data 
		FROM (
			SELECT * FROM openehr.tbl_agent
			UNION ALL
			SELECT * FROM openehr.tbl_person
			UNION ALL
			SELECT * FROM openehr.tbl_group
			UNION ALL
			SELECT * FROM openehr.tbl_organisation
		) as allparties
		WHERE versioned_object_id = $1 `)
	args = []any{versionedObjectID}
	argNum++

	if !filterAtTime.IsZero() {
		query.WriteString(fmt.Sprintf(` AND created_at <= $%d`, argNum))
		args = append(args, filterAtTime)
		argNum++
	}

	if filterVersionID != "" {
		query.WriteString(fmt.Sprintf(` AND id = $%d`, argNum))
		args = append(args, filterVersionID)
	}

	query.WriteString(` ORDER BY created_at DESC LIMIT 1`)
	row := s.DB.QueryRow(ctx, query.String(), args...)

	var rawEhrStatusJSON []byte
	if err := row.Scan(&rawEhrStatusJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR Status at time from database: %w", err)
	}

	return rawEhrStatusJSON, nil
}

func (s *DemographicService) CreateContribution(ctx context.Context, contribution openehr.CONTRIBUTION) (openehr.CONTRIBUTION, error) {
	if contribution.UID.Value == "" {
		contribution.UID.Value = uuid.NewString()
	}

	// Validate Contribution
	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.CONTRIBUTION{}, fmt.Errorf("validation errors for Contribution: %v", errs)
	}

	// Insert Contribution
	query := `INSERT INTO openehr.tbl_contribution (id, data) VALUES ($1, $2)`
	args := []any{contribution.UID.Value, contribution}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return openehr.CONTRIBUTION{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	return contribution, nil
}

func (s *DemographicService) GetContributionAsJSON(ctx context.Context, contributionID string) ([]byte, error) {
	query := `
		SELECT c.data 
		FROM openehr.tbl_contribution c
		WHERE c.ehr_id IS NULL AND c.id = $1
		LIMIT 1
	`
	args := []any{contributionID}
	row := s.DB.QueryRow(ctx, query, args...)

	var rawContributionJSON []byte
	if err := row.Scan(&rawContributionJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrContributionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Contribution by ID from database: %w", err)
	}

	return rawContributionJSON, nil
}
