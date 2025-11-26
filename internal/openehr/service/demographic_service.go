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
	ErrAgentAlreadyExists                       = fmt.Errorf("agent with the given UID already exists")
	ErrAgentNotFound                            = fmt.Errorf("agent not found")
	ErrGroupAlreadyExists                       = fmt.Errorf("group with the given UID already exists")
	ErrGroupNotFound                            = fmt.Errorf("group not found")
	ErrOrganisationAlreadyExists                = fmt.Errorf("organisation with the given UID already exists")
	ErrOrganisationNotFound                     = fmt.Errorf("organisation not found")
	ErrPersonAlreadyExists                      = fmt.Errorf("person with the given UID already exists")
	ErrPersonNotFound                           = fmt.Errorf("person not found")
	ErrRoleAlreadyExists                        = fmt.Errorf("role with the given UID already exists")
	ErrRoleNotFound                             = fmt.Errorf("role not found")
	ErrVersionedPartyNotFound                   = fmt.Errorf("versioned party not found")
	ErrRevisionHistoryNotFound                  = fmt.Errorf("revision history not found")
	ErrAgentVersionLowerOrEqualToCurrent        = fmt.Errorf("agent version is lower than or equal to the current version")
	ErrPersonVersionLowerOrEqualToCurrent       = fmt.Errorf("person version is lower than or equal to the current version")
	ErrGroupVersionLowerOrEqualToCurrent        = fmt.Errorf("group version is lower than or equal to the current version")
	ErrInvalidGroupUIDMismatch                  = fmt.Errorf("group UID does not match current group UID")
	ErrOrganisationVersionLowerOrEqualToCurrent = fmt.Errorf("organisation version is lower than or equal to the current version")
	ErrInvalidOrganisationUIDMismatch           = fmt.Errorf("organisation UID does not match current organisation UID")
	ErrRoleVersionLowerOrEqualToCurrent         = fmt.Errorf("role version is lower than or equal to the current version")
	ErrInvalidRoleUIDMismatch                   = fmt.Errorf("role UID does not match current role UID")
	ErrVersionedPartyVersionNotFound            = fmt.Errorf("versioned party version not found")
	ErrInvalidAgentUIDMismatch                  = fmt.Errorf("agent UID does not match current agent UID")
	ErrInvalidPersonUIDMismatch                 = fmt.Errorf("person UID does not match current person UID")
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

func (s *DemographicService) ValidateAgent(ctx context.Context, agent openehr.AGENT) error {
	validateErr := agent.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Agent validation can be added here

	return nil
}

func (s *DemographicService) ExistsAgent(ctx context.Context, versionID string) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.type = $1 AND ov.id = $2 LIMIT 1`

	var exists int
	err := s.DB.QueryRow(ctx, query, openehr.AGENT_MODEL_NAME, versionID).Scan(&exists)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if agent exists: %w", err)
	}
	return true, nil
}

func (s *DemographicService) CreateAgent(ctx context.Context, agent openehr.AGENT) (openehr.AGENT, error) {
	// Validate agent
	if err := s.ValidateAgent(ctx, agent); err != nil {
		return openehr.AGENT{}, err
	}

	// Provide ID when Agent does not have one
	if !agent.UID.E {
		agent.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: uuid.NewString(),
			},
		})
	}

	switch agent.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Convert to OBJECT_VERSION_ID
		hierID := agent.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		agent.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", hierID.Value, config.SYSTEM_ID_GOPENEHR),
			},
		})
	default:
		return openehr.AGENT{}, fmt.Errorf("agent UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", agent.UID.V.Value)
	}

	// Extract UID type
	agentID, ok := agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return openehr.AGENT{}, fmt.Errorf("agent UID must be of type OBJECT_VERSION_ID, got %T", agent.UID.V.Value)
	}

	// Check if agent with the same UID already exists
	exists, err := s.ExistsAgent(ctx, agentID.Value)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to check existing agent: %w", err)
	}
	if exists {
		return openehr.AGENT{}, ErrAgentAlreadyExists
	}

	// Build Versioned object for the agent
	versionedParty := openehr.VERSIONED_PARTY{
		UID: openehr.HIER_OBJECT_ID{
			Value: agentID.UID(),
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.AGENT_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: agentID.UID(),
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
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
				},
			},
		},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, object_type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		versionedParty.UID.Value, openehr.VERSIONED_PARTY_MODEL_NAME, openehr.AGENT_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}

	versionedParty.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object_data (id, data) VALUES ($1, $2)",
		versionedParty.UID.Value, versionedParty)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert versioned party data: %w", err)
	}

	// Insert agent
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, versionedParty.UID.Value, openehr.AGENT_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert agent: %w", err)
	}

	agentVersion.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)",
		agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, agentVersion)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert agent data: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return agent, nil
}

func (s *DemographicService) GetCurrentAgentVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (openehr.AGENT, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var agent openehr.AGENT
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), openehr.AGENT_MODEL_NAME).Scan(&agent)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.AGENT{}, ErrAgentNotFound
		}
		return openehr.AGENT{}, fmt.Errorf("failed to get latest agent by versioned party ID: %w", err)
	}

	return agent, nil
}

func (s *DemographicService) GetAgentAtVersion(ctx context.Context, versionID string) (openehr.AGENT, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var agent openehr.AGENT
	err := s.DB.QueryRow(ctx, query, versionID, openehr.AGENT_MODEL_NAME).Scan(&agent)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.AGENT{}, ErrAgentNotFound
		}
		return openehr.AGENT{}, fmt.Errorf("failed to get agent at version: %w", err)
	}

	return agent, nil
}

func (s *DemographicService) UpdateAgent(ctx context.Context, versionedPartyID uuid.UUID, agent openehr.AGENT) (openehr.AGENT, error) {
	// Validate Agent
	if err := s.ValidateAgent(ctx, agent); err != nil {
		return openehr.AGENT{}, err
	}

	// Get current Agent
	currentAgent, err := s.GetCurrentAgentVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if errors.Is(err, ErrAgentNotFound) {
			return openehr.AGENT{}, ErrAgentNotFound
		}
		return openehr.AGENT{}, fmt.Errorf("failed to get current Agent: %w", err)
	}
	currentAgentID := currentAgent.UID.V.Value.(*openehr.OBJECT_VERSION_ID)

	// Provide ID when Agent does not have one
	if !agent.UID.E {
		agent.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: currentAgentID.UID(),
			},
		})
	}

	// Handle Agent UID types
	switch v := agent.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check version is incremented
		if v.VersionTreeID().CompareTo(currentAgentID.VersionTreeID()) <= 0 {
			return openehr.AGENT{}, ErrAgentVersionLowerOrEqualToCurrent
		}
	case *openehr.HIER_OBJECT_ID:
		// Check HIER_OBJECT_ID matches current Agent UID
		if currentAgentID.UID() != v.Value {
			return openehr.AGENT{}, ErrInvalidAgentUIDMismatch
		}

		// Increment version
		versionTreeID := currentAgentID.VersionTreeID()
		versionTreeID.Major++

		// Convert to OBJECT_VERSION_ID
		agent.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", v.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.AGENT{}, fmt.Errorf("agent UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", agent.UID.V.Value)
	}

	// Build Agent version
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
					Value: agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
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
				},
			},
		},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Update agent
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(), openehr.AGENT_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert agent: %w", err)
	}

	agentVersion.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)",
		agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, agentVersion)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert agent data: %w", err)
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Delete agent
	var deleted uint8
	row := tx.QueryRow(ctx, "DELETE FROM openehr.tbl_versioned_object WHERE ehr_id IS NULL AND id = $1 AND type = $2 RETURNING 1", strings.Split(versionedObjectID, "::")[0], openehr.VERSIONED_PARTY_MODEL_NAME)
	err = row.Scan(&deleted)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return ErrAgentNotFound
		}
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DemographicService) ValidatePerson(ctx context.Context, person openehr.PERSON) error {
	validateErr := person.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Person validation can be added here

	return nil
}

func (s *DemographicService) ExistsPerson(ctx context.Context, versionID string) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.type = $1 AND ov.id = $2 LIMIT 1`

	var exists int
	err := s.DB.QueryRow(ctx, query, openehr.PERSON_MODEL_NAME, versionID).Scan(&exists)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if person exists: %w", err)
	}
	return true, nil
}

func (s *DemographicService) CreatePerson(ctx context.Context, person openehr.PERSON) (openehr.PERSON, error) {
	// Validate person
	if err := s.ValidatePerson(ctx, person); err != nil {
		return openehr.PERSON{}, err
	}

	// Provide ID when Person does not have one
	if !person.UID.E {
		person.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: uuid.NewString(),
			},
		})
	}

	switch person.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Convert to OBJECT_VERSION_ID
		hierID := person.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		person.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", hierID.Value, config.SYSTEM_ID_GOPENEHR),
			},
		})
	default:
		return openehr.PERSON{}, fmt.Errorf("person UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", person.UID.V.Value)
	}

	// Extract UID type
	personID, ok := person.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return openehr.PERSON{}, fmt.Errorf("person UID must be of type OBJECT_VERSION_ID, got %T", person.UID.V.Value)
	}

	// Check if person with the same UID already exists
	exists, err := s.ExistsPerson(ctx, personID.Value)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to check existing person: %w", err)
	}
	if exists {
		return openehr.PERSON{}, ErrPersonAlreadyExists
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
				},
			},
		},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, object_type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		versionedParty.UID.Value, openehr.VERSIONED_PARTY_MODEL_NAME, openehr.PERSON_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}

	versionedParty.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object_data (id, data) VALUES ($1, $2)",
		versionedParty.UID.Value, versionedParty)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert versioned party data: %w", err)
	}

	// Insert person
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, versionedParty.UID.Value, openehr.PERSON_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert person: %w", err)
	}

	personVersion.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)",
		person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, personVersion)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert person data: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return person, nil
}

func (s *DemographicService) GetCurrentPersonVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (openehr.PERSON, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var person openehr.PERSON
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), openehr.PERSON_MODEL_NAME).Scan(&person)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.PERSON{}, ErrPersonNotFound
		}
		return openehr.PERSON{}, fmt.Errorf("failed to get latest person by versioned party ID: %w", err)
	}

	return person, nil
}

func (s *DemographicService) GetPersonAtVersion(ctx context.Context, versionID string) (openehr.PERSON, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var person openehr.PERSON
	err := s.DB.QueryRow(ctx, query, versionID, openehr.GROUP_MODEL_NAME).Scan(&person)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.PERSON{}, ErrPersonNotFound
		}
		return openehr.PERSON{}, fmt.Errorf("failed to get person at version: %w", err)
	}

	return person, nil
}

func (s *DemographicService) UpdatePerson(ctx context.Context, versionedPartyID uuid.UUID, person openehr.PERSON) (openehr.PERSON, error) {
	// Validate Person
	if err := s.ValidatePerson(ctx, person); err != nil {
		return openehr.PERSON{}, err
	}

	// Get current Person
	currentPerson, err := s.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to get current Person: %w", err)
	}
	currentPersonID := currentPerson.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	// Provide ID when Person does not have one
	if !person.UID.E {
		person.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: currentPersonID.UID(),
			},
		})
	}

	// Handle Person UID types
	switch v := person.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check version is incremented
		if v.VersionTreeID().CompareTo(currentPersonID.VersionTreeID()) <= 0 {
			return openehr.PERSON{}, ErrPersonVersionLowerOrEqualToCurrent
		}
	case *openehr.HIER_OBJECT_ID:
		// Check HIER_OBJECT_ID matches current Person UID
		if currentPersonID.UID() != v.Value {
			return openehr.PERSON{}, ErrInvalidPersonUIDMismatch
		}

		// Increment version
		versionTreeID := currentPersonID.VersionTreeID()
		versionTreeID.Major++

		// Convert to OBJECT_VERSION_ID
		person.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", v.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.PERSON{}, fmt.Errorf("person UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", person.UID.V.Value)
	}

	// Build Person version
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
					Value: person.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
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
				},
			},
		},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Update person
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(), openehr.PERSON_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert person: %w", err)
	}

	personVersion.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)",
		person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, personVersion)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert person data: %w", err)
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Delete person
	var deleted uint8
	row := tx.QueryRow(ctx, "DELETE FROM openehr.tbl_versioned_object WHERE ehr_id IS NULL AND type = $1 AND id = $2 RETURNING 1", openehr.VERSIONED_PARTY_MODEL_NAME, strings.Split(versionedObjectID, "::")[0])
	err = row.Scan(&deleted)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return ErrPersonNotFound
		}
		return fmt.Errorf("failed to delete person: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DemographicService) ValidateGroup(ctx context.Context, group openehr.GROUP) error {
	validateErr := group.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Group validation can be added here

	return nil
}

func (s *DemographicService) ExistsGroup(ctx context.Context, versionID string) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.type = $1 AND ov.id = $2 LIMIT 1`

	var exists int
	err := s.DB.QueryRow(ctx, query, openehr.GROUP_MODEL_NAME, versionID).Scan(&exists)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if group exists: %w", err)
	}
	return true, nil
}

func (s *DemographicService) CreateGroup(ctx context.Context, group openehr.GROUP) (openehr.GROUP, error) {
	// Validate group
	if err := s.ValidateGroup(ctx, group); err != nil {
		return openehr.GROUP{}, err
	}

	// Provide ID when Group does not have one
	if !group.UID.E {
		group.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: uuid.NewString(),
			},
		})
	}

	switch group.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Convert to OBJECT_VERSION_ID
		hierID := group.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		group.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", hierID.Value, config.SYSTEM_ID_GOPENEHR),
			},
		})
	default:
		return openehr.GROUP{}, fmt.Errorf("group UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", group.UID.V.Value)
	}

	// Extract UID type
	groupID, ok := group.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return openehr.GROUP{}, fmt.Errorf("group UID must be of type OBJECT_VERSION_ID, got %T", group.UID.V.Value)
	}

	// Check if group with the same UID already exists
	exists, err := s.ExistsGroup(ctx, groupID.Value)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to check existing group: %w", err)
	}
	if exists {
		return openehr.GROUP{}, ErrGroupAlreadyExists
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
				},
			},
		},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, object_type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		versionedParty.UID.Value, openehr.VERSIONED_PARTY_MODEL_NAME, openehr.GROUP_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}

	versionedParty.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object_data (id, data) VALUES ($1, $2)",
		versionedParty.UID.Value, versionedParty)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert versioned party data: %w", err)
	}

	// Insert group
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		groupVersion.UID.Value, versionedParty.UID.Value, openehr.GROUP_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert group: %w", err)
	}

	groupVersion.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)",
		groupVersion.UID.Value, groupVersion)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert group data: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return group, nil
}

func (s *DemographicService) GetCurrentGroupVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (openehr.GROUP, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var group openehr.GROUP
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), openehr.GROUP_MODEL_NAME).Scan(&group)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.GROUP{}, ErrGroupNotFound
		}
		return openehr.GROUP{}, fmt.Errorf("failed to get latest group by versioned party ID: %w", err)
	}

	return group, nil
}

func (s *DemographicService) GetGroupAtVersion(ctx context.Context, versionID string) (openehr.GROUP, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var group openehr.GROUP
	err := s.DB.QueryRow(ctx, query, versionID, openehr.GROUP_MODEL_NAME).Scan(&group)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.GROUP{}, ErrGroupNotFound
		}
		return openehr.GROUP{}, fmt.Errorf("failed to get group at version: %w", err)
	}

	return group, nil
}

func (s *DemographicService) UpdateGroup(ctx context.Context, versionedPartyID uuid.UUID, group openehr.GROUP) (openehr.GROUP, error) {
	// Validate Group
	if err := s.ValidateGroup(ctx, group); err != nil {
		return openehr.GROUP{}, err
	}

	// Get current Group
	currentGroup, err := s.GetCurrentGroupVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to get current Group: %w", err)
	}
	currentGroupID := currentGroup.UID.V.Value.(*openehr.OBJECT_VERSION_ID)

	// Provide ID when Group does not have one
	if !group.UID.E {
		group.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: currentGroupID.UID(),
			},
		})
	}

	// Handle Group UID types
	switch v := group.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check version is incremented
		if v.VersionTreeID().CompareTo(currentGroupID.VersionTreeID()) <= 0 {
			return openehr.GROUP{}, ErrGroupVersionLowerOrEqualToCurrent
		}
	case *openehr.HIER_OBJECT_ID:
		// Check HIER_OBJECT_ID matches current Group UID
		if currentGroupID.UID() != v.Value {
			return openehr.GROUP{}, ErrInvalidGroupUIDMismatch
		}

		// Increment version
		versionTreeID := currentGroupID.VersionTreeID()
		versionTreeID.Major++

		// Convert to OBJECT_VERSION_ID
		group.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", v.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.GROUP{}, fmt.Errorf("group UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", group.UID.V.Value)
	}

	// Build group version
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
					Value: group.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
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
				},
			},
		},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Update group
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(), openehr.GROUP_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert group: %w", err)
	}

	groupVersion.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)",
		group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, groupVersion)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert group data: %w", err)
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Delete group
	var deleted uint8
	row := tx.QueryRow(ctx, "DELETE FROM openehr.tbl_versioned_object WHERE ehr_id IS NULL AND id = $1 AND type = $2 RETURNING 1", strings.Split(versionedObjectID, "::")[0], openehr.VERSIONED_PARTY_MODEL_NAME)
	err = row.Scan(&deleted)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return ErrGroupNotFound
		}
		return fmt.Errorf("failed to delete group: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DemographicService) ValidateOrganisation(ctx context.Context, organisation openehr.ORGANISATION) error {
	validateErr := organisation.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Organisation validation can be added here

	return nil
}

func (s *DemographicService) ExistsOrganisation(ctx context.Context, versionID string) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.type = $1 AND ov.id = $2 LIMIT 1`

	var exists int
	err := s.DB.QueryRow(ctx, query, openehr.ORGANISATION_MODEL_NAME, versionID).Scan(&exists)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if organisation exists: %w", err)
	}
	return true, nil
}

func (s *DemographicService) CreateOrganisation(ctx context.Context, organisation openehr.ORGANISATION) (openehr.ORGANISATION, error) {
	// Validate organisation
	if err := s.ValidateOrganisation(ctx, organisation); err != nil {
		return openehr.ORGANISATION{}, err
	}

	// Provide ID when Organisation does not have one
	if !organisation.UID.E {
		organisation.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: uuid.NewString(),
			},
		})
	}

	switch organisation.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Convert to OBJECT_VERSION_ID
		hierID := organisation.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		organisation.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", hierID.Value, config.SYSTEM_ID_GOPENEHR),
			},
		})
	default:
		return openehr.ORGANISATION{}, fmt.Errorf("organisation UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", organisation.UID.V.Value)
	}

	// Extract UID type
	organisationID, ok := organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return openehr.ORGANISATION{}, fmt.Errorf("organisation UID must be of type OBJECT_VERSION_ID, got %T", organisation.UID.V.Value)
	}

	// Check if organisation with the same UID already exists
	exists, err := s.ExistsOrganisation(ctx, organisationID.Value)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to check existing organisation: %w", err)
	}
	if exists {
		return openehr.ORGANISATION{}, ErrOrganisationAlreadyExists
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
				},
			},
		},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, object_type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		versionedParty.UID.Value, openehr.VERSIONED_PARTY_MODEL_NAME, openehr.ORGANISATION_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}

	versionedParty.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object_data (id, data) VALUES ($1, $2)",
		versionedParty.UID.Value, versionedParty)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert versioned party data: %w", err)
	}

	// Insert organisation
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		organisationVersion.UID.Value, versionedParty.UID.Value, openehr.ORGANISATION_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert organisation: %w", err)
	}

	organisationVersion.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)",
		organisationVersion.UID.Value, organisationVersion)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert organisation data: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return organisation, nil
}

func (s *DemographicService) GetCurrentOrganisationVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (openehr.ORGANISATION, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var organisation openehr.ORGANISATION
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), openehr.AGENT_MODEL_NAME).Scan(&organisation)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.ORGANISATION{}, ErrOrganisationNotFound
		}
		return openehr.ORGANISATION{}, fmt.Errorf("failed to get latest organisation by versioned party ID: %w", err)
	}

	return organisation, nil
}

func (s *DemographicService) GetOrganisationAtVersion(ctx context.Context, versionID string) (openehr.ORGANISATION, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var organisation openehr.ORGANISATION
	err := s.DB.QueryRow(ctx, query, versionID, openehr.ORGANISATION_MODEL_NAME).Scan(&organisation)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.ORGANISATION{}, ErrOrganisationNotFound
		}
		return openehr.ORGANISATION{}, fmt.Errorf("failed to get organisation at version: %w", err)
	}

	return organisation, nil
}

func (s *DemographicService) UpdateOrganisation(ctx context.Context, versionedPartyID uuid.UUID, organisation openehr.ORGANISATION) (openehr.ORGANISATION, error) {
	// Validate Organisation
	if err := s.ValidateOrganisation(ctx, organisation); err != nil {
		return openehr.ORGANISATION{}, err
	}

	// Get current Organisation
	currentOrganisation, err := s.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to get current Organisation: %w", err)
	}
	currentOrganisationID := currentOrganisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	// Provide ID when Organisation does not have one
	if !organisation.UID.E {
		organisation.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: currentOrganisationID.UID(),
			},
		})
	}

	// Handle Organisation UID types
	switch v := organisation.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check version is incremented
		if v.VersionTreeID().CompareTo(currentOrganisationID.VersionTreeID()) <= 0 {
			return openehr.ORGANISATION{}, ErrOrganisationVersionLowerOrEqualToCurrent
		}
	case *openehr.HIER_OBJECT_ID:
		// Check HIER_OBJECT_ID matches current Organisation UID
		if currentOrganisationID.UID() != v.Value {
			return openehr.ORGANISATION{}, ErrInvalidOrganisationUIDMismatch
		}

		// Increment version
		versionTreeID := currentOrganisationID.VersionTreeID()
		versionTreeID.Major++

		// Convert to OBJECT_VERSION_ID
		organisation.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", v.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.ORGANISATION{}, fmt.Errorf("organisation UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", organisation.UID.V.Value)
	}

	// Build organisation version
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
					Value: organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
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
				},
			},
		},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Update organisation
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(), openehr.ORGANISATION_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert organisation: %w", err)
	}

	organisationVersion.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)",
		organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, organisationVersion)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert organisation data: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return organisation, nil
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Delete organisation
	var deleted uint8
	row := tx.QueryRow(ctx, "DELETE FROM openehr.tbl_versioned_object WHERE ehr_id IS NULL AND id = $1 AND type = $2 RETURNING 1", strings.Split(versionedObjectID, "::")[0], openehr.VERSIONED_PARTY_MODEL_NAME)
	err = row.Scan(&deleted)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return ErrOrganisationNotFound
		}
		return fmt.Errorf("failed to delete organisation: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DemographicService) ValidateRole(ctx context.Context, role openehr.ROLE) error {
	validateErr := role.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Organisation validation can be added here

	return nil
}

func (s *DemographicService) ExistsRole(ctx context.Context, versionID string) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.type = $1 AND ov.id = $2 LIMIT 1`

	var exists int
	err := s.DB.QueryRow(ctx, query, openehr.ROLE_MODEL_NAME, versionID).Scan(&exists)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if role exists: %w", err)
	}
	return true, nil
}

func (s *DemographicService) CreateRole(ctx context.Context, role openehr.ROLE) (openehr.ROLE, error) {
	// Validate rol
	if err := s.ValidateRole(ctx, role); err != nil {
		return openehr.ROLE{}, err
	}

	// Provide ID when Role does not have one
	if !role.UID.E {
		role.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: uuid.NewString(),
			},
		})
	}

	switch role.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Convert to OBJECT_VERSION_ID
		hierID := role.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		role.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", hierID.Value, config.SYSTEM_ID_GOPENEHR),
			},
		})
	default:
		return openehr.ROLE{}, fmt.Errorf("role UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", role.UID.V.Value)
	}

	// Extract UID type
	roleID, ok := role.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return openehr.ROLE{}, fmt.Errorf("role UID must be of type OBJECT_VERSION_ID, got %T", role.UID.V.Value)
	}

	// Check if role with the same UID already exists
	exists, err := s.ExistsRole(ctx, roleID.Value)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to check existing role: %w", err)
	}
	if exists {
		return openehr.ROLE{}, ErrRoleAlreadyExists
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
				},
			},
		},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, object_type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		versionedParty.UID.Value, openehr.VERSIONED_PARTY_MODEL_NAME, openehr.ROLE_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}

	versionedParty.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object_data (id, data) VALUES ($1, $2)",
		versionedParty.UID.Value, versionedParty)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert versioned party data: %w", err)
	}

	// Insert role
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, versionedParty.UID.Value, openehr.ROLE_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert role: %w", err)
	}

	roleVersion.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)",
		role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, roleVersion)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert role data: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return role, nil
}

func (s *DemographicService) GetCurrentRoleVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (openehr.ROLE, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var role openehr.ROLE
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), openehr.ROLE_MODEL_NAME).Scan(&role)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.ROLE{}, ErrRoleNotFound
		}
		return openehr.ROLE{}, fmt.Errorf("failed to get latest role by versioned party ID: %w", err)
	}

	return role, nil
}

func (s *DemographicService) GetRoleAtVersion(ctx context.Context, versionID string) (openehr.ROLE, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var role openehr.ROLE
	err := s.DB.QueryRow(ctx, query, versionID, openehr.ROLE_MODEL_NAME).Scan(&role)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.ROLE{}, ErrRoleNotFound
		}
		return openehr.ROLE{}, fmt.Errorf("failed to get role at version: %w", err)
	}

	return role, nil
}

func (s *DemographicService) UpdateRole(ctx context.Context, versionedPartyID uuid.UUID, role openehr.ROLE) (openehr.ROLE, error) {
	// Validate Role
	if err := s.ValidateRole(ctx, role); err != nil {
		return openehr.ROLE{}, err
	}

	// Get current Role
	currentRole, err := s.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to get current Role: %w", err)
	}
	currentRoleID := currentRole.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	// Provide ID when Role does not have one
	if !role.UID.E {
		role.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: currentRoleID.UID(),
			},
		})
	}

	// Handle Role UID types
	switch v := role.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check version is incremented
		if v.VersionTreeID().CompareTo(currentRoleID.VersionTreeID()) <= 0 {
			return openehr.ROLE{}, ErrRoleVersionLowerOrEqualToCurrent
		}
	case *openehr.HIER_OBJECT_ID:
		// Check HIER_OBJECT_ID matches current Role UID
		if currentRoleID.UID() != v.Value {
			return openehr.ROLE{}, ErrInvalidRoleUIDMismatch
		}

		// Increment version
		versionTreeID := currentRoleID.VersionTreeID()
		versionTreeID.Major++

		// Convert to OBJECT_VERSION_ID
		role.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", v.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.ROLE{}, fmt.Errorf("role UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", role.UID.V.Value)
	}

	// Build role version
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
		PrecedingVersionUID: util.Some(*currentRoleID),
		Data:                &role,
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
					Value: role.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
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
				},
			},
		},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Update role
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, NULL, $4)",
		role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(), openehr.ROLE_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert role: %w", err)
	}

	roleVersion.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)",
		role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, roleVersion)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert role data: %w", err)
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}

	contribution.SetModelName()
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Delete role
	var deleted uint8
	row := tx.QueryRow(ctx, "DELETE FROM openehr.tbl_versioned_object WHERE ehr_id IS NULL AND id = $1 AND type = $2 RETURNING 1", strings.Split(versionedObjectID, "::")[0], openehr.VERSIONED_PARTY_MODEL_NAME)
	err = row.Scan(&deleted)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return ErrRoleNotFound
		}
		return fmt.Errorf("failed to delete role: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DemographicService) GetVersionedParty(ctx context.Context, versionedObjectID uuid.UUID) (openehr.VERSIONED_PARTY, error) {
	query := `
		SELECT vod.data 
		FROM openehr.tbl_versioned_object vo
		JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id
		WHERE vo.ehr_id IS NULL AND vo.type = $1 AND vo.id = $2 
		LIMIT 1
	`

	var versionedParty openehr.VERSIONED_PARTY
	err := s.DB.QueryRow(ctx, query, openehr.VERSIONED_PARTY_MODEL_NAME, versionedObjectID).Scan(&versionedParty)
	if err != nil {
		return openehr.VERSIONED_PARTY{}, fmt.Errorf("failed to get versioned party by ID: %w", err)
	}
	return versionedParty, nil
}

func (s *DemographicService) GetVersionedPartyRevisionHistory(ctx context.Context, versionedObjectID uuid.UUID) (openehr.REVISION_HISTORY, error) {
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
				jsonb_agg(cd.data->'audit' ORDER BY cd.data->'audit'->'time_committed'->>'value') as audits
			FROM openehr.tbl_contribution c
			JOIN openehr.tbl_contribution_data cd ON c.id = cd.id,
				jsonb_array_elements(cd.data->'versions') as version
			WHERE c.ehr_id IS NULL
				AND version->>'type' = ANY($1::text[])
				AND version->'id'->>'value' LIKE $2 || '%'
			GROUP BY version->'id'->>'value'
		) grouped
	`
	args := []any{[]string{openehr.AGENT_MODEL_NAME, openehr.PERSON_MODEL_NAME, openehr.GROUP_MODEL_NAME, openehr.ORGANISATION_MODEL_NAME, openehr.ROLE_MODEL_NAME}, versionedObjectID}
	row := s.DB.QueryRow(ctx, query, args...)

	var revisionHistory openehr.REVISION_HISTORY
	if err := row.Scan(&revisionHistory); err != nil {
		if err == database.ErrNoRows {
			return openehr.REVISION_HISTORY{}, ErrRevisionHistoryNotFound
		}
		return openehr.REVISION_HISTORY{}, fmt.Errorf("failed to fetch Revision History from database: %w", err)
	}

	return revisionHistory, nil
}

func (s *DemographicService) GetVersionedPartyVersionJSON(ctx context.Context, versionedObjectID uuid.UUID, filterAtTime time.Time, filterVersionID string) ([]byte, error) {
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`
		SELECT data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id
		WHERE ov.ehr_id IS NULL AND ov.versioned_object_id = $1
	`)
	args = []any{versionedObjectID}
	argNum++

	if !filterAtTime.IsZero() {
		query.WriteString(fmt.Sprintf(`AND ov.created_at <= $%d `, argNum))
		args = append(args, filterAtTime)
		argNum++
	}

	if filterVersionID != "" {
		query.WriteString(fmt.Sprintf(`AND ov.id = $%d `, argNum))
		args = append(args, filterVersionID)
	}

	query.WriteString(`ORDER BY ov.created_at DESC LIMIT 1`)
	row := s.DB.QueryRow(ctx, query.String(), args...)

	var partyVersionJSON []byte
	if err := row.Scan(&partyVersionJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrVersionedPartyVersionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Party version at time from database: %w", err)
	}

	return partyVersionJSON, nil
}

// func (s *DemographicService) CreateContribution(ctx context.Context, contribution openehr.CONTRIBUTION) (openehr.CONTRIBUTION, error) {
// 	if contribution.UID.Value == "" {
// 		contribution.UID.Value = uuid.NewString()
// 	}

// 	// Insert Contribution
// 	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)`
// 	args := []any{contribution.UID.Value}
// 	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
// 		return openehr.CONTRIBUTION{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
// 	}
// 	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
// 	args = []any{contribution.UID.Value, contribution}
// 	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
// 		return openehr.CONTRIBUTION{}, fmt.Errorf("failed to insert contribution data into the database: %w", err)
// 	}

// 	return contribution, nil
// }

func (s *DemographicService) GetContribution(ctx context.Context, contributionID string) (openehr.CONTRIBUTION, error) {
	query := `
		SELECT c.data 
		FROM openehr.tbl_contribution c
		WHERE c.ehr_id IS NULL AND c.id = $1
		LIMIT 1
	`
	args := []any{contributionID}
	row := s.DB.QueryRow(ctx, query, args...)

	var contribution openehr.CONTRIBUTION
	if err := row.Scan(&contribution); err != nil {
		if err == database.ErrNoRows {
			return openehr.CONTRIBUTION{}, ErrContributionNotFound
		}
		return openehr.CONTRIBUTION{}, fmt.Errorf("failed to fetch Contribution from database: %w", err)
	}

	return contribution, nil
}
