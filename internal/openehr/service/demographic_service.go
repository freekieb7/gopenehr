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
	ErrOrganisationVersionLowerOrEqualToCurrent = fmt.Errorf("organisation version is lower than or equal to the current version")
	ErrRoleVersionLowerOrEqualToCurrent         = fmt.Errorf("role version is lower than or equal to the current version")
	ErrVersionedPartyVersionNotFound            = fmt.Errorf("versioned party version not found")
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
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: uuid.NewString(),
			},
		})
	}

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
	_, err := s.GetAgent(ctx, agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID())
	if err == nil {
		return openehr.AGENT{}, ErrAgentAlreadyExists
	}
	if !errors.Is(err, ErrAgentNotFound) {
		return openehr.AGENT{}, fmt.Errorf("failed to check existing agent: %w", err)
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
				},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert contribution: %w", err)
	}
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id, contribution_id) VALUES ($1, $2, NULL, $3)",
		versionedParty.UID.Value, openehr.VERSIONED_PARTY_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}
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

// GetAgent retrieves the agent as the latest version when providing versioned object id, or the specified ID version, and returns it as the raw JSON representation.
func (s *DemographicService) GetAgent(ctx context.Context, uidBasedID string) (openehr.AGENT, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.type = $1 AND
	`
	args := []any{openehr.AGENT_MODEL_NAME}
	argNum := 2

	if strings.Count(uidBasedID, "::") == 2 {
		// It's an OBJECT_VERSION_ID
		query += fmt.Sprintf("ov.id = $%d ", argNum)
	} else {
		// It's a versioned object ID
		query += fmt.Sprintf("ov.versioned_object_id = $%d ORDER BY ov.created_at DESC", argNum)
	}
	args = append(args, uidBasedID)

	query += " LIMIT 1"

	var agent openehr.AGENT
	err := s.DB.QueryRow(ctx, query, args...).Scan(&agent)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.AGENT{}, ErrAgentNotFound
		}
		return openehr.AGENT{}, fmt.Errorf("failed to get agent by ID: %w", err)
	}

	return agent, nil
}

func (s *DemographicService) UpdateAgent(ctx context.Context, agent openehr.AGENT) (openehr.AGENT, error) {
	// Validate Composition
	if !agent.UID.E {
		return openehr.AGENT{}, fmt.Errorf("agent UID must be provided for update")
	}

	switch agent.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check if version is being updated
		currentAgent, err := s.GetAgent(ctx, agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID())
		if err != nil {
			return openehr.AGENT{}, fmt.Errorf("failed to get current Agent: %w", err)
		}

		currentVersionID := currentAgent.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		newVersionID := agent.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		if newVersionID.VersionTreeID().CompareTo(currentVersionID.VersionTreeID()) <= 0 {
			return openehr.AGENT{}, ErrAgentVersionLowerOrEqualToCurrent
		}
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Add namespace and version to convert to OBJECT_VERSION_ID
		currentAgent, err := s.GetAgent(ctx, agent.UID.V.Value.(*openehr.HIER_OBJECT_ID).Value)
		if err != nil {
			return openehr.AGENT{}, fmt.Errorf("failed to get current Agent: %w", err)
		}

		hierID := agent.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		versionTreeID := currentAgent.UID.V.Value.(*openehr.OBJECT_VERSION_ID).VersionTreeID()
		versionTreeID.Major++

		agent.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", hierID.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.AGENT{}, fmt.Errorf("agent UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", agent.UID.V.Value)
	}

	if errs := agent.Validate("$"); len(errs) > 0 {
		return openehr.AGENT{}, fmt.Errorf("validation errors for agent: %v", errs)
	}

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

	// Validate agent version
	agentVersion.SetModelName()
	if errs := agentVersion.Validate("$"); len(errs) > 0 {
		return openehr.AGENT{}, fmt.Errorf("validation errors for agent version: %v", errs)
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.AGENT{}, fmt.Errorf("failed to insert contribution: %w", err)
	}
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}
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

func (s *DemographicService) CreatePerson(ctx context.Context, person openehr.PERSON) (openehr.PERSON, error) {
	// If UID is not provided, generate a new one
	if !person.UID.E {
		person.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: uuid.NewString(),
			},
		})
	}

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
	_, err := s.GetPerson(ctx, person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID())
	if err == nil {
		return openehr.PERSON{}, ErrPersonAlreadyExists
	}
	if !errors.Is(err, ErrPersonNotFound) {
		return openehr.PERSON{}, fmt.Errorf("failed to check existing person: %w", err)
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
				},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert contribution: %w", err)
	}
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id, contribution_id) VALUES ($1, $2, NULL, $3)",
		versionedParty.UID.Value, openehr.VERSIONED_PARTY_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}
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

// GetPerson retrieves the person as the latest version when providing versioned object id, or the specified ID version, and returns it as the raw JSON representation.
func (s *DemographicService) GetPerson(ctx context.Context, uidBasedID string) (openehr.PERSON, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.type = $1 AND
	`
	args := []any{openehr.PERSON_MODEL_NAME}
	argNum := 2

	if strings.Count(uidBasedID, "::") == 2 {
		// It's an OBJECT_VERSION_ID
		query += fmt.Sprintf("ov.id = $%d ", argNum)
	} else {
		// It's a versioned object ID
		query += fmt.Sprintf("ov.versioned_object_id = $%d ORDER BY ov.created_at DESC", argNum)
	}
	args = append(args, uidBasedID)

	query += " LIMIT 1"

	var person openehr.PERSON
	err := s.DB.QueryRow(ctx, query, args...).Scan(&person)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.PERSON{}, ErrPersonNotFound
		}
		return openehr.PERSON{}, fmt.Errorf("failed to get person by ID: %w", err)
	}

	return person, nil
}

func (s *DemographicService) UpdatePerson(ctx context.Context, person openehr.PERSON) (openehr.PERSON, error) {
	// Validate Composition
	if !person.UID.E {
		return openehr.PERSON{}, fmt.Errorf("person UID must be provided for update")
	}

	switch person.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check if version is being updated
		currentPerson, err := s.GetPerson(ctx, person.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID())
		if err != nil {
			return openehr.PERSON{}, fmt.Errorf("failed to get current Person: %w", err)
		}

		currentVersionID := currentPerson.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		newVersionID := person.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		if newVersionID.VersionTreeID().CompareTo(currentVersionID.VersionTreeID()) <= 0 {
			return openehr.PERSON{}, ErrPersonVersionLowerOrEqualToCurrent
		}
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Add namespace and version to convert to OBJECT_VERSION_ID
		currentPerson, err := s.GetPerson(ctx, person.UID.V.Value.(*openehr.HIER_OBJECT_ID).Value)
		if err != nil {
			return openehr.PERSON{}, fmt.Errorf("failed to get current Person: %w", err)
		}

		hierID := person.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		versionTreeID := currentPerson.UID.V.Value.(*openehr.OBJECT_VERSION_ID).VersionTreeID()
		versionTreeID.Major++

		person.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", hierID.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.PERSON{}, fmt.Errorf("person UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", person.UID.V.Value)
	}

	if errs := person.Validate("$"); len(errs) > 0 {
		return openehr.PERSON{}, fmt.Errorf("validation errors for person: %v", errs)
	}

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

	// Validate person version
	personVersion.SetModelName()
	if errs := personVersion.Validate("$"); len(errs) > 0 {
		return openehr.PERSON{}, fmt.Errorf("validation errors for person version: %v", errs)
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.PERSON{}, fmt.Errorf("failed to insert contribution: %w", err)
	}
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}
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

func (s *DemographicService) CreateGroup(ctx context.Context, group openehr.GROUP) (openehr.GROUP, error) {
	// If UID is not provided, generate a new one
	if !group.UID.E {
		group.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: uuid.NewString(),
			},
		})
	}

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
	_, err := s.GetGroup(ctx, group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID())
	if err == nil {
		return openehr.GROUP{}, ErrGroupAlreadyExists
	}
	if !errors.Is(err, ErrGroupNotFound) {
		return openehr.GROUP{}, fmt.Errorf("failed to check existing group: %w", err)
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
				},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert contribution: %w", err)
	}
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id, contribution_id) VALUES ($1, $2, NULL, $3)",
		versionedParty.UID.Value, openehr.VERSIONED_PARTY_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}
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

// GetGroup retrieves the group as the latest version when providing versioned object id, or the specified ID version, and returns it as the raw JSON representation.
func (s *DemographicService) GetGroup(ctx context.Context, uidBasedID string) (openehr.GROUP, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ehr_id IS NULL AND ov.type = $1 AND
	`
	var args []any
	argNum := 2

	if strings.Count(uidBasedID, "::") == 2 {
		// It's an OBJECT_VERSION_ID
		query += fmt.Sprintf("ov.id = $%d ", argNum)
	} else {
		// It's a versioned object ID
		query += fmt.Sprintf("ov.versioned_object_id = $%d ORDER BY ov.created_at DESC", argNum)
	}
	args = append(args, uidBasedID)

	query += " LIMIT 1"

	var group openehr.GROUP
	err := s.DB.QueryRow(ctx, query, args...).Scan(&group)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.GROUP{}, ErrGroupNotFound
		}
		return openehr.GROUP{}, fmt.Errorf("failed to get group by ID: %w", err)
	}

	return group, nil
}

func (s *DemographicService) UpdateGroup(ctx context.Context, group openehr.GROUP) (openehr.GROUP, error) {
	// Validate Composition
	if !group.UID.E {
		return openehr.GROUP{}, fmt.Errorf("group UID must be provided for update")
	}

	switch group.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check if version is being updated
		currentGroup, err := s.GetGroup(ctx, group.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID())
		if err != nil {
			return openehr.GROUP{}, fmt.Errorf("failed to get current Group: %w", err)
		}

		currentVersionID := currentGroup.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		newVersionID := group.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		if newVersionID.VersionTreeID().CompareTo(currentVersionID.VersionTreeID()) <= 0 {
			return openehr.GROUP{}, ErrGroupVersionLowerOrEqualToCurrent
		}
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Add namespace and version to convert to OBJECT_VERSION_ID
		currentGroup, err := s.GetGroup(ctx, group.UID.V.Value.(*openehr.HIER_OBJECT_ID).Value)
		if err != nil {
			return openehr.GROUP{}, fmt.Errorf("failed to get current Group: %w", err)
		}

		hierID := group.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		versionTreeID := currentGroup.UID.V.Value.(*openehr.OBJECT_VERSION_ID).VersionTreeID()
		versionTreeID.Major++

		group.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", hierID.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.GROUP{}, fmt.Errorf("group UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", group.UID.V.Value)
	}

	if errs := group.Validate("$"); len(errs) > 0 {
		return openehr.GROUP{}, fmt.Errorf("validation errors for group: %v", errs)
	}

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

	// Validate group version
	groupVersion.SetModelName()
	if errs := groupVersion.Validate("$"); len(errs) > 0 {
		return openehr.GROUP{}, fmt.Errorf("validation errors for group version: %v", errs)
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.GROUP{}, fmt.Errorf("failed to insert contribution: %w", err)
	}
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}
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

func (s *DemographicService) CreateOrganisation(ctx context.Context, organisation openehr.ORGANISATION) (openehr.ORGANISATION, error) {
	// If UID is not provided, generate a new one
	if !organisation.UID.E {
		organisation.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: uuid.NewString(),
			},
		})
	}

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
	_, err := s.GetOrganisation(ctx, organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID())
	if err == nil {
		return openehr.ORGANISATION{}, ErrOrganisationAlreadyExists
	}
	if !errors.Is(err, ErrOrganisationNotFound) {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to check existing organisation: %w", err)
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
				},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert contribution: %w", err)
	}
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id, contribution_id) VALUES ($1, $2, NULL, $3)",
		versionedParty.UID.Value, openehr.VERSIONED_PARTY_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}
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

// GetOrganisation retrieves the organisation as the latest version when providing versioned object id, or the specified ID version, and returns it as the raw JSON representation.
func (s *DemographicService) GetOrganisation(ctx context.Context, uidBasedID string) (openehr.ORGANISATION, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ehr_id IS NULL AND ov.type = $1 AND
	`
	args := []any{openehr.ORGANISATION_MODEL_NAME}
	argNum := 2

	if strings.Count(uidBasedID, "::") == 2 {
		// It's an OBJECT_VERSION_ID
		query += fmt.Sprintf("ov.id = $%d ", argNum)
	} else {
		// It's a versioned object ID
		query += fmt.Sprintf("ov.versioned_object_id = $%d ORDER BY ov.created_at DESC", argNum)
	}
	args = append(args, uidBasedID)

	query += " LIMIT 1"

	var organisation openehr.ORGANISATION
	err := s.DB.QueryRow(ctx, query, args...).Scan(&organisation)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.ORGANISATION{}, ErrOrganisationNotFound
		}
		return openehr.ORGANISATION{}, fmt.Errorf("failed to get organisation by ID: %w", err)
	}

	return organisation, nil
}

func (s *DemographicService) UpdateOrganisation(ctx context.Context, organisation openehr.ORGANISATION) (openehr.ORGANISATION, error) {
	// Validate Composition
	if !organisation.UID.E {
		return openehr.ORGANISATION{}, fmt.Errorf("organisation UID must be provided for update")
	}

	switch organisation.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check if version is being updated
		currentOrganisation, err := s.GetOrganisation(ctx, organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID())
		if err != nil {
			return openehr.ORGANISATION{}, fmt.Errorf("failed to get current Organisation: %w", err)
		}

		currentVersionID := currentOrganisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		newVersionID := organisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		if newVersionID.VersionTreeID().CompareTo(currentVersionID.VersionTreeID()) <= 0 {
			return openehr.ORGANISATION{}, ErrOrganisationVersionLowerOrEqualToCurrent
		}
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Add namespace and version to convert to OBJECT_VERSION_ID
		currentOrganisation, err := s.GetOrganisation(ctx, organisation.UID.V.Value.(*openehr.HIER_OBJECT_ID).Value)
		if err != nil {
			return openehr.ORGANISATION{}, fmt.Errorf("failed to get current Organisation: %w", err)
		}

		hierID := organisation.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		versionTreeID := currentOrganisation.UID.V.Value.(*openehr.OBJECT_VERSION_ID).VersionTreeID()
		versionTreeID.Major++

		organisation.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", hierID.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.ORGANISATION{}, fmt.Errorf("organisation UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", organisation.UID.V.Value)
	}

	if errs := organisation.Validate("$"); len(errs) > 0 {
		return openehr.ORGANISATION{}, fmt.Errorf("validation errors for organisation: %v", errs)
	}

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

	// Validate organisation version
	organisationVersion.SetModelName()
	if errs := organisationVersion.Validate("$"); len(errs) > 0 {
		return openehr.ORGANISATION{}, fmt.Errorf("validation errors for organisation version: %v", errs)
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.ORGANISATION{}, fmt.Errorf("failed to insert contribution: %w", err)
	}
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}
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

func (s *DemographicService) CreateRole(ctx context.Context, role openehr.ROLE) (openehr.ROLE, error) {
	// If UID is not provided, generate a new one
	if !role.UID.E {
		role.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: uuid.NewString(),
			},
		})
	}

	// Convert HIER_OBJECT_ID to OBJECT_VERSION_ID if necessary
	switch role.UID.V.Value.(type) {
	case *openehr.HIER_OBJECT_ID:
		role.UID.V.Value = &openehr.OBJECT_VERSION_ID{
			Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
			Value: fmt.Sprintf("%s::%s::1", role.UID.V.Value.(*openehr.HIER_OBJECT_ID).Value, config.SYSTEM_ID_GOPENEHR),
		}
	case *openehr.OBJECT_VERSION_ID:
	// valid type
	default:
		return openehr.ROLE{}, errors.New("role UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID")
	}

	// Check if role with the same UID already exists
	_, err := s.GetRole(ctx, role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID())
	if err == nil {
		return openehr.ROLE{}, ErrRoleAlreadyExists
	}
	if !errors.Is(err, ErrRoleNotFound) {
		return openehr.ROLE{}, fmt.Errorf("failed to check existing role: %w", err)
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
				},
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert contribution: %w", err)
	}
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)",
		contribution.UID.Value, contribution)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert contribution data: %w", err)
	}

	// Insert versioned party
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id, contribution_id) VALUES ($1, $2, NULL, $3)",
		versionedParty.UID.Value, openehr.VERSIONED_PARTY_MODEL_NAME, contribution.UID.Value)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert versioned party: %w", err)
	}
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

// GetRole retrieves the role as the latest version when providing versioned object id, or the specified ID version, and returns it as the raw JSON representation.
func (s *DemographicService) GetRole(ctx context.Context, uidBasedID string) (openehr.ROLE, error) {
	query := `
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ehr_id IS NULL AND ov.type = $1 AND
	`
	args := []any{openehr.ROLE_MODEL_NAME}
	argNum := 2

	if strings.Count(uidBasedID, "::") == 2 {
		// It's an OBJECT_VERSION_ID
		query += fmt.Sprintf("ov.id = $%d ", argNum)
	} else {
		// It's a versioned object ID
		query += fmt.Sprintf("ov.versioned_object_id = $%d ORDER BY ov.created_at DESC", argNum)
	}
	args = append(args, uidBasedID)

	query += " LIMIT 1"

	var role openehr.ROLE
	err := s.DB.QueryRow(ctx, query, args...).Scan(&role)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return openehr.ROLE{}, ErrRoleNotFound
		}
		return openehr.ROLE{}, fmt.Errorf("failed to get role by ID: %w", err)
	}

	return role, nil
}

func (s *DemographicService) UpdateRole(ctx context.Context, role openehr.ROLE) (openehr.ROLE, error) {
	// Validate Composition
	if !role.UID.E {
		return openehr.ROLE{}, fmt.Errorf("role UID must be provided for update")
	}

	switch role.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check if version is being updated
		currentRole, err := s.GetRole(ctx, role.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID())
		if err != nil {
			return openehr.ROLE{}, fmt.Errorf("failed to get current Role: %w", err)
		}

		currentVersionID := currentRole.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		newVersionID := role.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		if newVersionID.VersionTreeID().CompareTo(currentVersionID.VersionTreeID()) <= 0 {
			return openehr.ROLE{}, ErrRoleVersionLowerOrEqualToCurrent
		}
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Add namespace and version to convert to OBJECT_VERSION_ID
		currentRole, err := s.GetRole(ctx, role.UID.V.Value.(*openehr.HIER_OBJECT_ID).Value)
		if err != nil {
			return openehr.ROLE{}, fmt.Errorf("failed to get current Role: %w", err)
		}

		hierID := role.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		versionTreeID := currentRole.UID.V.Value.(*openehr.OBJECT_VERSION_ID).VersionTreeID()
		versionTreeID.Major++

		role.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", hierID.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.ROLE{}, fmt.Errorf("role UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", role.UID.V.Value)
	}

	if errs := role.Validate("$"); len(errs) > 0 {
		return openehr.ROLE{}, fmt.Errorf("validation errors for role: %v", errs)
	}

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

	// Validate role version
	roleVersion.SetModelName()
	if errs := roleVersion.Validate("$"); len(errs) > 0 {
		return openehr.ROLE{}, fmt.Errorf("validation errors for role version: %v", errs)
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

	// Validate contribution
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return openehr.ROLE{}, fmt.Errorf("failed to insert contribution: %w", err)
	}
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

	// Validate contribution
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

	// Insert contribution
	_, err = tx.Exec(ctx, "INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)",
		contribution.UID.Value)
	if err != nil {
		return fmt.Errorf("failed to insert contribution: %w", err)
	}
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
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, NULL)`
	args := []any{contribution.UID.Value}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return openehr.CONTRIBUTION{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}
	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return openehr.CONTRIBUTION{}, fmt.Errorf("failed to insert contribution data into the database: %w", err)
	}

	return contribution, nil
}

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
