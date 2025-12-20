package openehr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/openehr/aql"
	"github.com/freekieb7/gopenehr/internal/openehr/rm"
	"github.com/freekieb7/gopenehr/internal/openehr/terminology"
	outil "github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrEHRNotFound                             = fmt.Errorf("EHR not found")
	ErrEHRAlreadyExists                        = fmt.Errorf("EHR already exists")
	ErrEHRStatusNotFound                       = fmt.Errorf("EHR Status not found")
	ErrEHRStatusAlreadyExists                  = fmt.Errorf("EHR Status already exists")
	ErrCompositionAlreadyExists                = fmt.Errorf("composition already exists")
	ErrCompositionNotFound                     = fmt.Errorf("composition not found")
	ErrCompositionVersionLowerOrEqualToCurrent = fmt.Errorf("composition version must be incremented")
	ErrInvalidCompositionUIDMismatch           = fmt.Errorf("composition UID HIER_OBJECT_ID does not match current composition UID")
	ErrDirectoryAlreadyExists                  = fmt.Errorf("directory already exists")
	ErrDirectoryNotFound                       = fmt.Errorf("directory not found")
	ErrContributionAlreadyExists               = fmt.Errorf("contribution already exists")
	ErrContributionNotFound                    = fmt.Errorf("contribution not found")
	ErrEHRStatusVersionLowerOrEqualToCurrent   = fmt.Errorf("EHR Status version must be incremented")
	ErrInvalidEHRStatusUIDMismatch             = fmt.Errorf("EHR Status UID HIER_OBJECT_ID does not match current EHR Status UID")
	ErrFolderNotFoundInDirectory               = fmt.Errorf("folder not found in directory")
	ErrDirectoryVersionLowerOrEqualToCurrent   = fmt.Errorf("directory version must be incremented")
	ErrInvalidDirectoryUIDMismatch             = fmt.Errorf("directory UID HIER_OBJECT_ID does not match current directory UID")
	ErrCompositionUIDNotProvided               = fmt.Errorf("composition UID must be provided")

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
	ErrVersionedObjectNotFound                  = fmt.Errorf("versioned object not found")

	ErrVersionLowerOrEqualToCurrent = fmt.Errorf("object version must be incremented")

	ErrQueryNotFound      = fmt.Errorf("AQL query not found")
	ErrQueryAlreadyExists = fmt.Errorf("AQL query with the given name already exists")
)

type StoredQuery struct {
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Query   string    `json:"q"`
	Type    string    `json:"type"`
	Saved   time.Time `json:"saved"`
}

type Service struct {
	Logger *telemetry.Logger
	DB     *database.Database
}

func NewService(logger *telemetry.Logger, db *database.Database) *Service {
	return &Service{
		Logger: logger,
		DB:     db,
	}
}

func (s *Service) CreateEHR(ctx context.Context, ehrID uuid.UUID, ehrStatus rm.EHR_STATUS) (rm.EHR, error) {
	err := s.ValidateEHRStatus(ctx, ehrStatus)
	if err != nil {
		return rm.EHR{}, err
	}

	// Ensure versioned object id does not already exist
	if ehrStatus.UID.E {
		var versionedEHRStatusIDstr string
		switch ehrStatus.UID.V.Kind {
		case rm.UID_BASED_ID_kind_OBJECT_VERSION_ID:
			versionedEHRStatusIDstr = ehrStatus.UID.V.OBJECT_VERSION_ID().UID()
		case rm.UID_BASED_ID_kind_HIER_OBJECT_ID:
			versionedEHRStatusIDstr = ehrStatus.UID.V.HIER_OBJECT_ID().Value
		default:
			return rm.EHR{}, fmt.Errorf("invalid EHR Status UID kind")
		}

		versionedEHRStatusID, err := uuid.Parse(versionedEHRStatusIDstr)
		if err != nil {
			return rm.EHR{}, fmt.Errorf("invalid UID format: %w", err)
		}

		exists, err := s.ExistsVersionedObject(ctx, versionedEHRStatusID)
		if err != nil {
			return rm.EHR{}, fmt.Errorf("failed to check if versioned object exists: %w", err)
		}
		if exists {
			return rm.EHR{}, ErrEHRStatusAlreadyExists
		}
	} else {
		// Provide new versioned object id
		ehrStatus.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&rm.OBJECT_VERSION_ID{
			Value: fmt.Sprintf("%s::%s::1", uuid.NewString(), config.SYSTEM_ID_GOPENEHR),
		}))
	}

	versionedEHRStatus := NewVersionedEHRStatus(uuid.MustParse(ehrStatus.UID.V.OBJECT_VERSION_ID().UID()), ehrID)
	versionedEHRAccess := NewVersionedEHRAccess(uuid.New(), ehrID)
	ehrAccess := NewEHRAccess(uuid.MustParse(versionedEHRAccess.UID.Value))
	ehrStatusVersion := NewOriginalVersion(*ehrStatus.UID.V.Value.(*rm.OBJECT_VERSION_ID), rm.ORIGINAL_VERSION_DATA_from_EHR_STATUS(ehrStatus), utils.None[rm.OBJECT_VERSION_ID]())
	ehrAccessVersion := NewOriginalVersion(*ehrAccess.UID.V.Value.(*rm.OBJECT_VERSION_ID), rm.ORIGINAL_VERSION_DATA_from_EHR_ACCESS(ehrAccess), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("EHR created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_EHR_STATUS_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(ehrStatus.UID.V.OBJECT_VERSION_ID()),
			},
			{
				Type:      rm.VERSIONED_EHR_ACCESS_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(ehrAccess.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)
	ehr := rm.EHR{
		SystemID: rm.HIER_OBJECT_ID{
			Value: config.SYSTEM_ID_GOPENEHR,
		},
		EHRID: rm.HIER_OBJECT_ID{
			Value: ehrID.String(),
		},
		EHRStatus: rm.OBJECT_REF{
			Type:      rm.VERSIONED_EHR_STATUS_TYPE,
			Namespace: rm.Namespace_local,
			ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
				Value: ehrStatus.UID.V.OBJECT_VERSION_ID().UID(),
			}),
		},
		EHRAccess: rm.OBJECT_REF{
			Type:      rm.VERSIONED_EHR_ACCESS_TYPE,
			Namespace: rm.Namespace_local,
			ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
				Value: ehrAccess.UID.V.OBJECT_VERSION_ID().UID(),
			}),
		},
		Contributions: utils.Some([]rm.OBJECT_REF{
			{
				Type:      rm.CONTRIBUTION_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(contribution.UID),
			},
		}),
		Compositions: utils.Some([]rm.OBJECT_REF{}),
		Directory:    utils.None[rm.OBJECT_REF](),
		Folders:      utils.Some([]rm.OBJECT_REF{}),
		Tags:         utils.Some([]rm.OBJECT_REF{}),
		TimeCreated: rm.DV_DATE_TIME{
			Value: time.Now().Format(time.RFC3339),
		},
	}

	// Only 'local' VERSIONED_PARTY external refs are supported
	localRefVersionedParty := utils.None[uuid.UUID]()
	if ehrStatus.Subject.ExternalRef.E && ehrStatus.Subject.ExternalRef.V.Namespace == rm.Namespace_local && ehrStatus.Subject.ExternalRef.V.Type == rm.VERSIONED_PARTY_TYPE {
		localRefVersionedPartyID := uuid.MustParse(ehrStatus.Subject.ExternalRef.V.ID.Value.(*rm.HIER_OBJECT_ID).Value)
		localRefVersionedParty = utils.Some(localRefVersionedPartyID)
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Increment user's registered EHR count
	var ehrCount int
	row := tx.QueryRow(ctx, `UPDATE subscription.tbl_organisation SET ehr_count = ehr_count + 1, updated_at = NOW() WHERE id = $1 AND ehr_count < ehr_limit RETURNING ehr_count;`)
	err = row.Scan(&ehrCount)
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to increment organisation's registered EHR count: %w", err)
	}

	batch := &pgx.Batch{}

	// Insert EHR
	batch.Queue(`INSERT INTO openehr.tbl_ehr (id) VALUES ($1)`, ehr.EHRID.Value)
	ehr.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_ehr_data (id, data) VALUES ($1, $2)`, ehr.EHRID.Value, ehr)
	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`, contribution.UID.Value, ehrID)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)
	// Insert VERSIONED_EHR_STATUS
	batch.Queue(`INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id) VALUES ($1, $2, $3)`, versionedEHRStatus.UID.Value, rm.VERSIONED_EHR_STATUS_TYPE, ehrID)
	versionedEHRStatus.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_versioned_ehr_status_data (id, data) VALUES ($1, $2)`, versionedEHRStatus.UID.Value, versionedEHRStatus)
	// Insert EHR_STATUS
	batch.Queue(`INSERT INTO openehr.tbl_ehr_status (id, version_int, versioned_ehr_status_id, ehr_id, contribution_id, local_ref_versioned_party_id) VALUES ($1, $2, $3, $4, $5, $6)`, ehrStatusVersion.UID.Value, ehrStatusVersion.UID.VersionTreeID().Int(), ehrStatusVersion.UID.UID(), ehrID, contribution.UID.Value, localRefVersionedParty)
	ehrStatusVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_ehr_status_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, ehrStatusVersion.UID.Value, ehrStatusVersion)
	// Insert VERSIONED_EHR_ACCESS
	batch.Queue(`INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id) VALUES ($1, $2, $3)`, versionedEHRAccess.UID.Value, rm.VERSIONED_EHR_ACCESS_TYPE, ehrID)
	versionedEHRAccess.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_versioned_ehr_access_data (id, data) VALUES ($1, $2)`, versionedEHRAccess.UID.Value, versionedEHRAccess)
	// Insert EHR_ACCESS
	batch.Queue(`INSERT INTO openehr.tbl_ehr_access (id, version_int, versioned_ehr_access_id, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`, ehrAccessVersion.UID.Value, ehrAccessVersion.UID.VersionTreeID().Int(), ehrAccessVersion.UID.UID(), ehrID, contribution.UID.Value)
	ehrAccessVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_ehr_access_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, ehrAccessVersion.UID.Value, ehrAccessVersion)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to execute batch insert for EHR creation: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to close batch result for EHR creation: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ehr, nil
}

func (s *Service) ExistsEHR(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_ehr WHERE id = $1 LIMIT 1`
	args := []any{id}

	row := s.DB.QueryRow(ctx, query, args...)

	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if err == database.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check EHR existence in database: %w", err)
	}

	return true, nil
}

func (s *Service) GetEHRRawJSON(ctx context.Context, id uuid.UUID) ([]byte, error) {
	query := `SELECT ed.data FROM openehr.tbl_ehr e JOIN tbl_ehr_data ed ON e.id = ed.id WHERE e.id = $1 LIMIT 1`
	args := []any{id}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetEHRBySubjectRawJSON(ctx context.Context, subjectID, subjectNamespace string) ([]byte, error) {
	query := `
        SELECT ed.data
        FROM openehr.tbl_ehr e
		JOIN openehr.tbl_ehr_data ed ON e.id = ed.id
		JOIN openehr.tbl_ehr_status es ON e.id = es.ehr_id
        WHERE es.data->'subject'->'external_ref'->>'namespace' = $1
		  AND es.data->'subject'->'external_ref'->'id'->>'value' = $2
        LIMIT 1
    `
	args := []any{subjectNamespace, subjectID}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR by subject from database: %w", err)
	}

	return data, nil
}

func (s *Service) DeleteEHR(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM openehr.tbl_ehr WHERE id = $1 RETURNING 1`
	args := []any{id}

	row := s.DB.QueryRow(ctx, query, args...)

	var deleted uint8
	err := row.Scan(&deleted)
	if err != nil {
		if err == database.ErrNoRows {
			return ErrEHRNotFound
		}

		return fmt.Errorf("failed to delete EHR from database: %w", err)
	}
	return nil
}

func (s *Service) DeleteEHRBulk(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(`DELETE FROM openehr.tbl_ehr WHERE id IN (%s)`, strings.Join(placeholders, ", "))

	_, err := s.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete multiple EHRs from database: %w", err)
	}

	return nil
}

func (s *Service) ValidateEHRStatus(ctx context.Context, ehrStatus rm.EHR_STATUS) error {
	validateErr := ehrStatus.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional EHR Status validation
	if ehrStatus.Subject.ExternalRef.E {
		externalRef := ehrStatus.Subject.ExternalRef.V
		if externalRef.Namespace == rm.Namespace_local {

			attrPath := "$.subject.external_ref"

			switch v := externalRef.ID.Value.(type) {
			case *rm.HIER_OBJECT_ID:
				// Must be a valid type
				if externalRef.Type != rm.VERSIONED_PARTY_TYPE {
					validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
						Model:          rm.EHR_STATUS_TYPE,
						Path:           attrPath + ".type",
						Message:        fmt.Sprintf("invalid subject external_ref type: %s", externalRef.Type),
						Recommendation: "Ensure external ref type is VERSIONED_PARTY",
					})
				}

				// Must be a valid UUID
				if err := uuid.Validate(v.Value); err != nil {
					validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
						Model:          rm.EHR_STATUS_TYPE,
						Path:           attrPath + ".id.value",
						Message:        fmt.Sprintf("invalid subject external_ref id: %v", err),
						Recommendation: "Ensure external ref id is a valid UUID",
					})
				}

				if len(validateErr.Errs) > 0 {
					// Stop validation, there are already errors
					break
				}

				// Check if VERSIONED_PARTY exists
				err := s.DB.QueryRow(ctx, "SELECT 1 FROM openehr.tbl_versioned_object WHERE id = $1 AND type = $2", v.Value, externalRef.Type).Scan(new(int))
				if err != nil {
					if err == database.ErrNoRows {
						validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
							Model:          rm.EHR_STATUS_TYPE,
							Path:           attrPath + ".id.value",
							Message:        "Subject external ref id " + v.Value + " with type " + externalRef.Type + " does not exist in tbl_versioned_party",
							Recommendation: "Ensure external ref id exists in tbl_versioned_party",
						})
					} else {
						return err
					}
				}
			default:
				validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
					Model:          rm.EHR_STATUS_TYPE,
					Path:           attrPath + ".id",
					Message:        "Unsupported subject external_ref id type",
					Recommendation: "Ensure external ref is of type HIER_OBJECT_ID and type is VERSIONED_PARTY",
				})
			}
		}
	}

	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	return nil
}

func (s *Service) GetEHRStatusID(ctx context.Context, ehrID uuid.UUID) (rm.OBJECT_VERSION_ID, error) {
	query := `SELECT id FROM openehr.tbl_ehr_status WHERE ehr_id = $1 ORDER BY version_int DESC LIMIT 1`
	args := []any{ehrID}

	row := s.DB.QueryRow(ctx, query, args...)

	var id string
	err := row.Scan(&id)
	if err != nil {
		if err == database.ErrNoRows {
			return rm.OBJECT_VERSION_ID{}, ErrEHRStatusNotFound
		}
		return rm.OBJECT_VERSION_ID{}, fmt.Errorf("failed to fetch EHR Status ID from database: %w", err)
	}

	return rm.OBJECT_VERSION_ID{Value: id}, nil
}

func (s *Service) GetEHRStatusByVersionedEHRStatusIDRawJSON(ctx context.Context, ehrID uuid.UUID, versionedEHRStatusID uuid.UUID) ([]byte, error) {
	query := `
		SELECT esd.data 
		FROM openehr.tbl_ehr_status es 
		JOIN openehr.tbl_ehr_status_data esd ON es.id = esd.id 
		WHERE es.ehr_id = $1 
		  AND es.versioned_ehr_status_id = $2 
		ORDER BY es.version_int DESC
		LIMIT 1
	`
	args := []any{ehrID, versionedEHRStatusID}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRStatusNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR Status from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetEHRStatusByIDRawJSON(ctx context.Context, ehrID uuid.UUID, ehrStatusID string) ([]byte, error) {
	query := `
		SELECT esd.data 
		FROM openehr.tbl_ehr_status es 
		JOIN openehr.tbl_ehr_status_data esd ON es.id = esd.id 
		WHERE es.ehr_id = $1 
		  AND es.id = $2 
		LIMIT 1
	`
	args := []any{ehrID, ehrStatusID}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRStatusNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR Status from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetEHRStatusAtTimeRawJSON(ctx context.Context, ehrID uuid.UUID, filterOnTime time.Time) ([]byte, error) {
	query := `
		SELECT esd.data 
		FROM openehr.tbl_ehr_status es 
		JOIN openehr.tbl_ehr_status_data esd ON es.id = esd.id 
		WHERE es.ehr_id = $1 
	`
	args := []any{ehrID}
	if !filterOnTime.IsZero() {
		query += `AND es.created_at <= $2 `
		args = append(args, filterOnTime)
	}
	query += `ORDER BY es.created_at DESC LIMIT 1`

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRStatusNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR Status from database: %w", err)
	}

	return data, nil
}

func (s *Service) UpdateEHRStatus(ctx context.Context, ehrID uuid.UUID, currentEHRStatusID rm.OBJECT_VERSION_ID, nextEHRStatus rm.EHR_STATUS) (rm.EHR_STATUS, error) {
	if err := s.ValidateEHRStatus(ctx, nextEHRStatus); err != nil {
		return rm.EHR_STATUS{}, err
	}

	// Ensure EHR Status contains a UID to check/upgrade
	if !nextEHRStatus.UID.E {
		nextEHRStatus.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: currentEHRStatusID.UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(nextEHRStatus.UID.V, currentEHRStatusID)
	if err != nil {
		return rm.EHR_STATUS{}, fmt.Errorf("failed to upgrade current EHR Status UID: %w", err)
	}
	nextEHRStatus.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	ehrStatusVersion := rm.ORIGINAL_VERSION{
		UID:                 nextEHRStatus.UID.V.OBJECT_VERSION_ID(),
		PrecedingVersionUID: utils.Some(currentEHRStatusID),
		LifecycleState: rm.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: rm.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: rm.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: rm.ORIGINAL_VERSION_DATA_from_EHR_STATUS(nextEHRStatus),
	}

	contribution := rm.CONTRIBUTION{
		UID: rm.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_EHR_STATUS_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(nextEHRStatus.UID.V.OBJECT_VERSION_ID()),
			},
		},
		Audit: rm.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: rm.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			Description:   utils.Some(rm.DV_TEXT{Value: "EHR Status updated"}),
			Committer: rm.PARTY_PROXY_from_PARTY_SELF(rm.PARTY_SELF{
				Type_: utils.Some(rm.PARTY_SELF_TYPE),
			}),
		},
	}

	// Start transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.EHR_STATUS{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Only 'local' VERSIONED_PARTY external refs are supported
	localRefVersionedParty := utils.None[uuid.UUID]()
	if nextEHRStatus.Subject.ExternalRef.E && nextEHRStatus.Subject.ExternalRef.V.Namespace == rm.Namespace_local && nextEHRStatus.Subject.ExternalRef.V.Type == rm.VERSIONED_PARTY_TYPE {
		localRefVersionedPartyID := uuid.MustParse(nextEHRStatus.Subject.ExternalRef.V.ID.Value.(*rm.HIER_OBJECT_ID).Value)
		localRefVersionedParty = utils.Some(localRefVersionedPartyID)
	}

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`, contribution.UID.Value, ehrID)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert EHR_STATUS
	batch.Queue(`INSERT INTO openehr.tbl_ehr_status (id, version_int, versioned_ehr_status_id, ehr_id, contribution_id, local_ref_versioned_party_id) VALUES ($1, $2, $3, $4, $5, $6)`, nextEHRStatus.UID.V.OBJECT_VERSION_ID().Value, nextEHRStatus.UID.V.OBJECT_VERSION_ID().VersionTreeID().Int(), nextEHRStatus.UID.V.OBJECT_VERSION_ID().UID(), ehrID, contribution.UID.Value, localRefVersionedParty)
	nextEHRStatus.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_ehr_status_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, nextEHRStatus.UID.V.OBJECT_VERSION_ID().Value, ehrStatusVersion)

	// Update EHR with new contribution reference
	batch.Queue(`
		UPDATE openehr.tbl_ehr_data
		SET data = jsonb_insert(data, '{contributions,-1}', $1::jsonb, true)
		WHERE id = $2
	`, rm.OBJECT_REF{
		Type:      rm.CONTRIBUTION_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(contribution.UID),
	}, ehrID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.EHR_STATUS{}, fmt.Errorf("failed to execute batch insert for EHR Status update: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.EHR_STATUS{}, fmt.Errorf("failed to close batch result for EHR Status update: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return rm.EHR_STATUS{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextEHRStatus, nil
}

func (s *Service) GetVersionedEHRStatusRawJSON(ctx context.Context, ehrID uuid.UUID) ([]byte, error) {
	query := `
		SELECT vod.data 
		FROM openehr.tbl_versioned_object vo
		JOIN openehr.tbl_versioned_ehr_status vod ON vo.id = vod.id
		WHERE vo.type = $1 
		  AND vo.ehr_id = $2
		LIMIT 1
	`
	args := []any{rm.VERSIONED_EHR_STATUS_TYPE, ehrID}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch Versioned EHR Status from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetVersionedEHRStatusRevisionHistoryRawJSON(ctx context.Context, ehrID uuid.UUID) ([]byte, error) {
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
			WHERE c.ehr_id = $1
				AND version->>'type' = 'EHR_STATUS'
			GROUP BY version->'id'->>'value'
		) grouped
    `
	args := []any{ehrID}

	var data []byte
	row := s.DB.QueryRow(ctx, query, args...)
	err := row.Scan(&data)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch Revision History from database: %w", err)
	}

	return data, nil
}

func (s *Service) ValidateComposition(ctx context.Context, composition rm.COMPOSITION) error {
	validateErr := composition.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Composition validation can be added here

	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	return nil
}

func (s *Service) ExistsComposition(ctx context.Context, ehrID uuid.UUID, versionID string) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.ehr_id = $1 AND ov.type = $2 AND ov.id = $3 LIMIT 1`
	args := []any{ehrID, rm.COMPOSITION_TYPE, versionID}

	var exists int
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		if err == database.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if Composition exists in database: %w", err)
	}

	return true, nil
}

func (s *Service) CreateComposition(ctx context.Context, ehrID uuid.UUID, composition rm.COMPOSITION) (rm.COMPOSITION, error) {
	err := s.ValidateComposition(ctx, composition)
	if err != nil {
		return rm.COMPOSITION{}, err
	}

	// Check if versioned object ID is already used
	if composition.UID.E {
		var versionedCompositionIDstr string
		switch composition.UID.V.Kind {
		case rm.UID_BASED_ID_kind_HIER_OBJECT_ID:
			versionedCompositionIDstr = composition.UID.V.HIER_OBJECT_ID().Value
		case rm.UID_BASED_ID_kind_OBJECT_VERSION_ID:
			versionedCompositionIDstr = composition.UID.V.OBJECT_VERSION_ID().UID()
		default:
			return rm.COMPOSITION{}, fmt.Errorf("unsupported UID kind: %d", composition.UID.V.Kind)
		}

		versionedCompositionID, err := uuid.Parse(versionedCompositionIDstr)
		if err != nil {
			return rm.COMPOSITION{}, fmt.Errorf("invalid composition UID: %w", err)
		}

		exists, err := s.ExistsVersionedObject(ctx, versionedCompositionID)
		if err != nil {
			return rm.COMPOSITION{}, fmt.Errorf("failed to check if versioned composition exists: %w", err)
		}
		if exists {
			return rm.COMPOSITION{}, ErrCompositionAlreadyExists
		}
	} else {
		// Provide new versioned object id
		composition.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&rm.OBJECT_VERSION_ID{
			Value: fmt.Sprintf("%s::%s::1", uuid.NewString(), config.SYSTEM_ID_GOPENEHR),
		}))
	}

	versionedComposition := NewVersionedComposition(uuid.MustParse(composition.UID.V.OBJECT_VERSION_ID().UID()), ehrID)
	compositionVersion := NewOriginalVersion(composition.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_COMPOSITION(composition), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Composition created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_COMPOSITION_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(composition.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	contributionData, err := sonic.Marshal(contribution)
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to marshal contribution data: %w", err)
	}
	versionedCompositionData, err := sonic.Marshal(versionedComposition)
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to marshal versioned composition data: %w", err)
	}
	compositionVersionData, err := sonic.Marshal(compositionVersion)
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to marshal composition version data: %w", err)
	}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`, contribution.UID.Value, ehrID)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contributionData)

	// Insert VERSIONED_COMPOSITION
	batch.Queue(`INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id) VALUES ($1, $2, $3)`, versionedComposition.UID.Value, rm.VERSIONED_COMPOSITION_TYPE, ehrID)
	versionedComposition.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_versioned_composition_data (id, data) VALUES ($1, $2)`, versionedComposition.UID.Value, versionedCompositionData)

	// Insert COMPOSITION
	batch.Queue(`INSERT INTO openehr.tbl_composition (id, version_int, versioned_composition_id, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`, composition.UID.V.OBJECT_VERSION_ID().Value, composition.UID.V.OBJECT_VERSION_ID().VersionTreeID().Int(), composition.UID.V.OBJECT_VERSION_ID().UID(), ehrID, contribution.UID.Value)
	compositionVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_composition_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, composition.UID.V.OBJECT_VERSION_ID().Value, compositionVersionData)

	// Update EHR, add contribution ref to list
	batch.Queue(`
		UPDATE openehr.tbl_ehr_data
		SET data = jsonb_insert(
			jsonb_insert(data, '{compositions, -1}', $2::jsonb, true)
			, '{contributions, -1}', $1::jsonb, true)
		WHERE id = $3
	`, rm.OBJECT_REF{
		Type:      rm.CONTRIBUTION_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(contribution.UID),
	}, rm.OBJECT_REF{
		Type:      rm.VERSIONED_COMPOSITION_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(versionedComposition.UID),
	}, ehrID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to execute batch insert for Composition creation: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to close batch result for Composition creation: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return composition, nil
}

func (s *Service) GetCompositionID(ctx context.Context, ehrID uuid.UUID, versionedCompositionID uuid.UUID) (rm.OBJECT_VERSION_ID, error) {
	query := `SELECT id FROM openehr.tbl_composition WHERE ehr_id = $1 AND versioned_composition_id = $2 ORDER BY version_int DESC LIMIT 1`
	args := []any{ehrID, versionedCompositionID}

	var id string
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		if err == database.ErrNoRows {
			return rm.OBJECT_VERSION_ID{}, ErrCompositionNotFound
		}
		return rm.OBJECT_VERSION_ID{}, fmt.Errorf("failed to fetch Composition ID from database: %w", err)
	}

	return rm.OBJECT_VERSION_ID{Value: id}, nil
}

func (s *Service) GetCompositionRawJSON(ctx context.Context, ehrID uuid.UUID, objectVersionID string) ([]byte, error) {
	query := `
		SELECT cd.data 
		FROM openehr.tbl_composition c 
		JOIN openehr.tbl_composition_data cd ON c.id = cd.id 
		WHERE c.ehr_id = $1 AND c.id = $2 LIMIT 1
	`
	args := []any{ehrID, objectVersionID}

	var data []byte
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&data); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrCompositionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Composition by ID from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetCompositionByVersionedCompositionIDRawJSON(ctx context.Context, ehrID uuid.UUID, versionedCompositionID uuid.UUID) ([]byte, error) {
	query := `
		SELECT cd.data
		FROM openehr.tbl_composition c 
		JOIN openehr.tbl_composition_data cd ON cd.id = c.id 
		WHERE c.ehr_id = $1 AND c.versioned_composition_id = $2
		ORDER BY c.version_int DESC
		LIMIT 1
	`
	args := []any{ehrID, versionedCompositionID}

	var data []byte
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&data); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrCompositionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Composition by ID from database: %w", err)
	}

	return data, nil
}

func (s *Service) UpdateComposition(ctx context.Context, ehrID uuid.UUID, currentCompositionID rm.OBJECT_VERSION_ID, nextComposition rm.COMPOSITION) (rm.COMPOSITION, error) {
	err := s.ValidateComposition(ctx, nextComposition)
	if err != nil {
		return rm.COMPOSITION{}, err
	}

	// Ensure Composition contains a UID to check/upgrade
	if !nextComposition.UID.E {
		nextComposition.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: currentCompositionID.UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(nextComposition.UID.V, currentCompositionID)
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to upgrade current Composition UID: %w", err)
	}
	nextComposition.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	compositionVersion := NewOriginalVersion(nextComposition.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_COMPOSITION(nextComposition), utils.Some(currentCompositionID))
	contribution := NewContribution("Composition updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_COMPOSITION_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(nextComposition.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	contributionData, err := sonic.Marshal(contribution)
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to marshal contribution data: %w", err)
	}
	compositionVersionData, err := sonic.Marshal(compositionVersion)
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to marshal composition version data: %w", err)
	}

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`, contribution.UID.Value, ehrID)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contributionData)

	// Insert COMPOSITION
	batch.Queue(`INSERT INTO openehr.tbl_composition (id, version_int, versioned_composition_id, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`, nextComposition.UID.V.OBJECT_VERSION_ID().Value, nextComposition.UID.V.OBJECT_VERSION_ID().VersionTreeID().Int(), nextComposition.UID.V.OBJECT_VERSION_ID().UID(), ehrID, contribution.UID.Value)
	nextComposition.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_composition_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, nextComposition.UID.V.OBJECT_VERSION_ID().Value, compositionVersionData)

	// Update EHR with contribution ref
	batch.Queue(`
		UPDATE openehr.tbl_ehr_data
		SET data = jsonb_insert(data, '{contributions, -1}', $1::jsonb, true)
		WHERE id = $2
	`, rm.OBJECT_REF{
		Type:      rm.CONTRIBUTION_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(contribution.UID),
	}, ehrID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to execute batch insert for Composition update: %w", err)
	}

	err = br.Close()
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to close batch result for Composition update: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextComposition, nil
}

func (s *Service) DeleteVersionedComposition(ctx context.Context, ehrID uuid.UUID, versionedCompositionID uuid.UUID) error {
	contribution := NewContribution("Composition deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      rm.COMPOSITION_TYPE,
				ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
					Value: versionedCompositionID.String(),
				}),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`, contribution.UID.Value, ehrID)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, contribution.UID.Value, contribution)

	// Delete COMPOSITION (todo return 1 and check if deleted?)
	batch.Queue(`DELETE FROM openehr.tbl_versioned_object WHERE ehr_id = $1 AND versioned_composition_id = $2`, ehrID, versionedCompositionID)

	// Update EHR, add contribution and remove composition ref from list
	batch.Queue(`
		UPDATE openehr.tbl_ehr
		SET data = jsonb_insert(data, '{contributions, -1}', $1::jsonb, true) #- (
			SELECT ARRAY['compositions', (idx-1)::text]
			FROM jsonb_array_elements(data->'compositions') WITH ORDINALITY arr(item, idx)
			WHERE item->'id'->>'value' = $2::text
			LIMIT 1
		)
		WHERE id = $3
	`, rm.OBJECT_REF{
		Type:      rm.CONTRIBUTION_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(contribution.UID),
	}, versionedCompositionID.String(),
		ehrID,
	)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("failed to execute batch insert for Versioned Composition deletion: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("failed to close batch result for Versioned Composition deletion: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) GetVersionedCompositionRawJSON(ctx context.Context, ehrID uuid.UUID, versionedCompositionID string) ([]byte, error) {
	query := `
		SELECT voc.data 
		FROM openehr.tbl_versioned_object vo
		JOIN openehr.tbl_versioned_composition voc ON voc.id = vo.id
		WHERE vo.ehr_id = $1
		  AND vo.id = $2
		LIMIT 1
	`
	args := []any{ehrID, versionedCompositionID}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrCompositionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Versioned Composition by ID from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetVersionedCompositionRevisionHistoryRawJSON(ctx context.Context, ehrID uuid.UUID, versionedCompositionID string) ([]byte, error) {
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
			WHERE c.ehr_id = $1
				AND version->>'type' = 'COMPOSITION'
				AND version->'id'->>'value' LIKE $2 || '%'
			GROUP BY version->'id'->>'value'
		) grouped
	`
	args := []any{ehrID, versionedCompositionID}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, ErrCompositionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Revision History from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetCompositionAtTimeRawJSON(ctx context.Context, ehrID uuid.UUID, versionedCompositionID uuid.UUID, filterAtTime time.Time) ([]byte, error) {
	query := `
		SELECT cd.data 
		FROM openehr.tbl_composition c
		JOIN openehr.tbl_composition_data cd ON cd.id = c.id
		WHERE c.ehr_id = $1 AND c.versioned_composition_id = $2
	`
	args := []any{ehrID, versionedCompositionID}

	if !filterAtTime.IsZero() {
		query += `AND c.created_at <= $3 `
		args = append(args, filterAtTime)
	}

	query += `ORDER BY c.created_at DESC LIMIT 1`
	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrCompositionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Composition version at time from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetCompositionByIDRawJSON(ctx context.Context, ehrID uuid.UUID, versionedCompositionID uuid.UUID, compositionID string) ([]byte, error) {
	query := `
		SELECT cd.data 
		FROM openehr.tbl_composition c
		JOIN openehr.tbl_composition_data cd ON cd.id = c.id
		WHERE c.ehr_id = $1 AND c.versioned_composition_id = $2 AND c.id = $3
		LIMIT 1
	`
	args := []any{ehrID, versionedCompositionID, compositionID}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrCompositionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Composition by ID from database: %w", err)
	}

	return data, nil
}

func (s *Service) ValidateDirectory(ctx context.Context, ehrID uuid.UUID, directory rm.FOLDER) error {
	validateErr := directory.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Directory validation can be added here
	folderQueue := make([]rm.FOLDER, 0)
	folderQueue = append(folderQueue, directory)
	pathQueue := make([]string, 0)
	pathQueue = append(pathQueue, "$")

	for len(folderQueue) > 0 {
		currentFolder := folderQueue[0]
		folderQueue = folderQueue[1:]
		currentPath := pathQueue[0]
		pathQueue = pathQueue[1:]

		// Add folders
		if currentFolder.Folders.E {
			folderQueue = append(folderQueue, currentFolder.Folders.V...)
			for i := range currentFolder.Folders.V {
				pathQueue = append(pathQueue, fmt.Sprintf("%s.folders[%d]", currentPath, i))
			}
		}

		// Process object references
		for i, currentRef := range currentFolder.Items.V {
			if currentRef.Namespace != config.NAMESPACE_LOCAL {
				// Skip external references
				continue
			}

			itemPath := fmt.Sprintf("%s.items[%d]", currentPath, i)

			// Handle different reference types
			switch currentRef.Type {
			case rm.COMPOSITION_TYPE:
				id, ok := currentRef.ID.Value.(*rm.OBJECT_VERSION_ID)
				if !ok {
					validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
						Model:          rm.COMPOSITION_TYPE,
						Path:           itemPath,
						Message:        "Mismatch between type and id provided",
						Recommendation: "Ensure the ID is of type OBJECT_VERSION_ID",
					})
				}

				exists, err := s.ExistsComposition(ctx, ehrID, id.Value)
				if err != nil {
					return fmt.Errorf("failed to validate existence of Composition: %w", err)
				}
				if !exists {
					validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
						Model:          rm.COMPOSITION_TYPE,
						Path:           itemPath,
						Message:        "COMPOSITION does not exist for this EHR in the system",
						Recommendation: "Ensure the composition exists for this EHR",
					})
				}
			case rm.VERSIONED_COMPOSITION_TYPE:
				id, ok := currentRef.ID.Value.(*rm.HIER_OBJECT_ID)
				if !ok {
					validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
						Model:          rm.VERSIONED_COMPOSITION_TYPE,
						Path:           itemPath,
						Message:        "Mismatch between type and id provided",
						Recommendation: "Ensure the ID is of type HIER_OBJECT_ID",
					})
				}

				exists, err := s.ExistsComposition(ctx, ehrID, id.Value)
				if err != nil {
					return fmt.Errorf("failed to validate existence of Composition: %w", err)
				}
				if !exists {
					validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
						Model:          rm.COMPOSITION_TYPE,
						Path:           itemPath,
						Message:        "COMPOSITION does not exist for this EHR in the system",
						Recommendation: "Ensure the composition exists for this EHR",
					})
				}
			}
		}
	}

	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	return nil
}

func (s *Service) CreateDirectory(ctx context.Context, ehrID uuid.UUID, directory rm.FOLDER) (rm.FOLDER, error) {
	err := s.ValidateDirectory(ctx, ehrID, directory)
	if err != nil {
		return rm.FOLDER{}, err
	}

	exists, err := s.ExistsDirectory(ctx, ehrID)
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to check if Directory exists: %w", err)
	}
	if exists {
		return rm.FOLDER{}, ErrDirectoryAlreadyExists
	}

	if directory.UID.E {
		var versionedFolderIDstr string
		switch directory.UID.V.Kind {
		case rm.UID_BASED_ID_kind_HIER_OBJECT_ID:
			versionedFolderIDstr = directory.UID.V.HIER_OBJECT_ID().Value
		case rm.UID_BASED_ID_kind_OBJECT_VERSION_ID:
			versionedFolderIDstr = directory.UID.V.OBJECT_VERSION_ID().UID()
		default:
			return rm.FOLDER{}, fmt.Errorf("unsupported UID kind: %d", directory.UID.V.Kind)
		}

		versionedFolderID, err := uuid.Parse(versionedFolderIDstr)
		if err != nil {
			return rm.FOLDER{}, fmt.Errorf("invalid directory UID: %w", err)
		}

		exists, err := s.ExistsVersionedObject(ctx, versionedFolderID)
		if err != nil {
			return rm.FOLDER{}, fmt.Errorf("failed to check if versioned directory exists: %w", err)
		}
		if exists {
			return rm.FOLDER{}, ErrDirectoryAlreadyExists
		}
	} else {
		// Assign a new UID if not provided
		newUID, err := uuid.NewRandom()
		if err != nil {
			return rm.FOLDER{}, fmt.Errorf("failed to generate new UUID for Directory: %w", err)
		}
		directory.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&rm.OBJECT_VERSION_ID{
			Value: fmt.Sprintf("%s::%s::%d", newUID.String(), config.NAMESPACE_LOCAL, 1),
		}))
	}

	versionedFolder := NewVersionedFolder(uuid.MustParse(directory.UID.V.OBJECT_VERSION_ID().UID()), ehrID)
	folderVersion := NewOriginalVersion(directory.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_FOLDER(directory), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Directory created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_FOLDER_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(directory.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`, contribution.UID.Value, ehrID)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert VERSIONED_FOLDER
	batch.Queue(`INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id) VALUES ($1, $2, $3)`, versionedFolder.UID.Value, rm.VERSIONED_FOLDER_TYPE, ehrID)
	versionedFolder.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_versioned_folder_data (id, data) VALUES ($1, $2)`, versionedFolder.UID.Value, versionedFolder)

	// Insert FOLDER
	batch.Queue(`INSERT INTO openehr.tbl_folder (id, version_int, version_object_id, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`, folderVersion.UID.Value, directory.UID.V.OBJECT_VERSION_ID().VersionTreeID().Int(), directory.UID.V.OBJECT_VERSION_ID().UID(), ehrID, contribution.UID.Value)
	folderVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_folder_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, folderVersion.UID.Value, folderVersion)

	// Update EHR, add contribution ref to list
	batch.Queue(`UPDATE openehr.tbl_ehr
		SET data = jsonb_insert(
			jsonb_insert(
				jsonb_set(data, '{directory}', $2::jsonb)
			, '{folders, 0}', $2::jsonb)
		, '{contributions, -1}', $1::jsonb, true)
		WHERE id = $3
	`, rm.OBJECT_REF{
		Type:      rm.CONTRIBUTION_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(contribution.UID),
	}, rm.OBJECT_REF{
		Type:      rm.VERSIONED_FOLDER_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(versionedFolder.UID),
	}, ehrID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to execute batch insert for Directory creation: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to close batch result for Directory creation: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return directory, nil
}

func (s *Service) ExistsDirectory(ctx context.Context, ehrID uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_versioned_object vo WHERE vo.ehr_id = $1 AND vo.type = $2 LIMIT 1`
	args := []any{ehrID, rm.VERSIONED_FOLDER_TYPE}

	var exists int
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		if err == database.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if Directory exists in database: %w", err)
	}

	return true, nil
}

func (s *Service) GetDirectoryID(ctx context.Context, ehrID uuid.UUID) (rm.OBJECT_VERSION_ID, error) {
	query := `SELECT id FROM openehr.tbl_folder WHERE ehr_id = $1 ORDER BY version_int DESC LIMIT 1`
	args := []any{ehrID}

	var id string
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		if err == database.ErrNoRows {
			return rm.OBJECT_VERSION_ID{}, ErrDirectoryNotFound
		}
		return rm.OBJECT_VERSION_ID{}, fmt.Errorf("failed to fetch Directory ID from database: %w", err)
	}

	return rm.OBJECT_VERSION_ID{Value: id}, nil
}

func (s *Service) GetDirectoryRawJSON(ctx context.Context, ehrID uuid.UUID) ([]byte, error) {
	query := `
		SELECT fd.data
        FROM openehr.tbl_folder f
        JOIN openehr.tbl_folder_data fd ON fd.id = f.id
        WHERE f.ehr_id = $2
        ORDER BY f.created_at DESC
        LIMIT 1
	`
	args := []any{rm.FOLDER_TYPE, ehrID}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrDirectoryNotFound
		}
		return nil, fmt.Errorf("failed to fetch Directory by EHR ID from database: %w", err)
	}

	return data, nil
}

func (s *Service) UpdateDirectory(ctx context.Context, ehrID uuid.UUID, currentDirectoryID rm.OBJECT_VERSION_ID, nextDirectory rm.FOLDER) (rm.FOLDER, error) {
	err := s.ValidateDirectory(ctx, ehrID, nextDirectory)
	if err != nil {
		return rm.FOLDER{}, err
	}

	if !nextDirectory.UID.E {
		nextDirectory.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: currentDirectoryID.UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(nextDirectory.UID.V, currentDirectoryID)
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to upgrade current Directory UID: %w", err)
	}
	nextDirectory.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	folderVersion := NewOriginalVersion(nextDirectory.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_FOLDER(nextDirectory), utils.Some(currentDirectoryID))
	contribution := NewContribution("Directory updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_FOLDER_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(nextDirectory.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`, contribution.UID.Value, ehrID)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert FOLDER
	batch.Queue(`INSERT INTO openehr.tbl_folder (id, version_int, version_object_id, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`, folderVersion.UID.Value, nextDirectory.UID.V.OBJECT_VERSION_ID().VersionTreeID().Int(), nextDirectory.UID.V.OBJECT_VERSION_ID().UID(), ehrID, contribution.UID.Value)
	folderVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_folder_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, folderVersion.UID.Value, folderVersion)

	// Update EHR with contribution ref
	batch.Queue(`UPDATE openehr.tbl_ehr
		SET data = jsonb_insert(data, '{contributions, -1}', $1::jsonb, true)
		WHERE id = $2
	`, rm.OBJECT_REF{
		Type:      rm.CONTRIBUTION_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(contribution.UID),
	}, ehrID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to execute batch insert for Directory update: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to close batch result for Directory update: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextDirectory, nil
}

func (s *Service) DeleteDirectory(ctx context.Context, ehrID uuid.UUID, versionedFolderID uuid.UUID) error {
	contribution := NewContribution("Directory deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      rm.FOLDER_TYPE,
				ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
					Value: versionedFolderID.String(),
				}),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`, contribution.UID.Value, ehrID)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Delete FOLDER (todo return 1 and check if deleted?)
	batch.Queue(`DELETE FROM openehr.tbl_versioned_object WHERE ehr_id = $1 AND id = $2`, ehrID, versionedFolderID)

	// Update EHR, add contribution ref to list and remove directory reference
	// Folder reference is deleted as just the first entry, like the openehr docs specify
	batch.Queue(`
		UPDATE openehr.tbl_ehr
		SET data = jsonb_insert(data, '{contributions, -1}', $1::jsonb, true) #- '{directory}' #- '{folders, 0}'
		WHERE id = $2
	`, rm.OBJECT_REF{
		Type:      rm.CONTRIBUTION_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(contribution.UID),
	}, ehrID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("failed to execute batch insert for Directory deletion: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("failed to close batch result for Directory deletion: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) GetFolderAtTimeRawJSON(ctx context.Context, ehrID uuid.UUID, filterAtTime time.Time, path string) ([]byte, error) {
	jsonPath := "$"
	for part := range strings.SplitSeq(path, "/") {
		jsonPath += fmt.Sprintf(`.folders ? (@.name.value == "%s")`, part)
	}

	query := fmt.Sprintf(`
		SELECT jsonb_path_query_first(fd.data, '%s') 
		FROM openehr.tbl_folder f
		JOIN openehr.tbl_folder_data fd ON fd.id = f.id
		WHERE f.ehr_id = $1
		  AND fd.data @? $2
	`, jsonPath)
	args := []any{ehrID, jsonPath}

	if !filterAtTime.IsZero() {
		query += `AND ov.created_at <= $3 `
		args = append(args, filterAtTime)
	}

	query += `ORDER BY ov.version_int DESC LIMIT 1`

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == database.ErrNoRows {
			if path != "" {
				return nil, ErrFolderNotFoundInDirectory
			}
			return nil, ErrDirectoryNotFound
		}
		return nil, fmt.Errorf("failed to fetch Folder at time from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetFolderInDirectoryByIDRawJSON(ctx context.Context, ehrID uuid.UUID, folderID uuid.UUID, path string) ([]byte, error) {
	jsonPath := "$"
	for part := range strings.SplitSeq(path, "/") {
		jsonPath += fmt.Sprintf(`.folders ? (@.name.value == "%s")`, part)
	}

	query := fmt.Sprintf(`
		SELECT jsonb_path_query_first(fd.data, '%s') 
		FROM openehr.tbl_folder f
		JOIN openehr.tbl_folder_data fd ON fd.id = f.id
		WHERE f.ehr_id = $1
		  AND fd.data @? $2
		  AND f.id = $3
		LIMIT 1
	`, jsonPath)
	args := []any{ehrID, jsonPath, folderID}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == database.ErrNoRows {
			if path != "" {
				return nil, ErrFolderNotFoundInDirectory
			}
			return nil, ErrDirectoryNotFound
		}
		return nil, fmt.Errorf("failed to fetch Folder at time from database: %w", err)
	}

	return data, nil
}

func (s *Service) ValidateTags(ctx context.Context, tags []rm.ITEM_TAG) error {
	for i := range tags {
		validateErr := tags[i].Validate("$")
		if len(validateErr.Errs) > 0 {
			return validateErr
		}
	}
	return nil
}

func (s *Service) GetEHRTagsRawJSON(ctx context.Context, ehrID uuid.UUID) ([]byte, error) {
	query := `
		SELECT COALESCE(jsonb_agg(data), '[]'::jsonb)
		FROM (
			SELECT data
			FROM openehr.tbl_ehr_status_tag
			WHERE ehr_id = $1
			UNION ALL
			SELECT data
			FROM openehr.tbl_composition_tag
			WHERE ehr_id = $1
		)
	`
	row := s.DB.QueryRow(ctx, query, ehrID)

	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, fmt.Errorf("failed to fetch EHR tags by EHR ID from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetCompositionTagsRawJSON(ctx context.Context, ehrID uuid.UUID, compositionID string) ([]byte, error) {
	query := `
		SELECT jsonb_agg(data)
		FROM openehr.tbl_composition_tag
		WHERE ehr_id = $1 AND composition_id = $2
	`
	row := s.DB.QueryRow(ctx, query, ehrID, compositionID)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == database.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch Composition tags by ID from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetVersionedCompositionTagsRawJSON(ctx context.Context, ehrID uuid.UUID, versionedCompositionID uuid.UUID) ([]byte, error) {
	query := `
		SELECT jsonb_agg(data)
		FROM openehr.tbl_versioned_composition_tag
		WHERE ehr_id = $1 AND versioned_composition_id = $2
	`
	row := s.DB.QueryRow(ctx, query, ehrID, versionedCompositionID)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == database.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch Versioned Composition tags by ID from database: %w", err)
	}

	return data, nil
}

func (s *Service) ReplaceVersionedCompositionTags(ctx context.Context, ehrID, versionedCompositionID uuid.UUID, tags []rm.ITEM_TAG) ([]rm.ITEM_TAG, error) {
	err := s.ValidateTags(ctx, tags)
	if err != nil {
		return nil, err
	}

	// Make sure target and object ref are correct
	for i := range tags {
		if tags[i].Target.Kind != rm.UID_BASED_ID_kind_HIER_OBJECT_ID {
			return nil, fmt.Errorf("invalid tag target UID kind: expected HIER_OBJECT_ID, got %d", tags[i].Target.Kind)
		}
		if tags[i].Target.HIER_OBJECT_ID().Value != versionedCompositionID.String() {
			return nil, fmt.Errorf("invalid tag target UID value: expected %s, got %s", versionedCompositionID.String(), tags[i].Target.HIER_OBJECT_ID().Value)
		}
		if tags[i].OwnerID.ID.Kind != rm.OBJECT_ID_kind_HIER_OBJECT_ID {
			return nil, fmt.Errorf("invalid tag owner ID kind: expected HIER_OBJECT_ID, got %d", tags[i].OwnerID.ID.Kind)
		}
		if tags[i].OwnerID.ID.HIER_OBJECT_ID().Value != ehrID.String() {
			return nil, fmt.Errorf("invalid tag owner ID value: expected %s, got %s", ehrID.String(), tags[i].OwnerID.ID.HIER_OBJECT_ID().Value)
		}
	}

	contribution := NewContribution("Versioned Composition tags replaced", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_COMPOSITION_TYPE,
				Namespace: rm.Namespace_local,
				ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
					Value: versionedCompositionID.String(),
				}),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`, contribution.UID.Value, ehrID)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Delete existing tags
	batch.Queue(`DELETE FROM openehr.tbl_versioned_composition_tag WHERE ehr_id = $1 AND versioned_composition_id = $2`, ehrID, versionedCompositionID)

	// Insert new tags
	for _, tag := range tags {
		batch.Queue(`INSERT INTO openehr.tbl_versioned_composition_tag (versioned_composition_id, key, data, ehr_id) VALUES ($1, $2, $3, $4)`, versionedCompositionID, tag.Key, tag, ehrID)
	}

	// Update EHR with contribution ref
	batch.Queue(`UPDATE openehr.tbl_ehr SET data = jsonb_insert(data, '{contributions, -1}', $1::jsonb, true) WHERE id = $2`, rm.OBJECT_REF{
		Type:      rm.CONTRIBUTION_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(contribution.UID),
	}, ehrID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return nil, fmt.Errorf("failed to execute batch insert for Versioned Composition tags replacement: %w", err)
	}
	err = br.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close batch result for Versioned Composition tags replacement: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tags, nil
}

func (s *Service) ReplaceCompositionTags(ctx context.Context, ehrID uuid.UUID, compositionID string, tags []rm.ITEM_TAG) ([]rm.ITEM_TAG, error) {
	err := s.ValidateTags(ctx, tags)
	if err != nil {
		return nil, err
	}

	// Make sure target and object ref are correct
	for i := range tags {
		if tags[i].Target.Kind != rm.UID_BASED_ID_kind_OBJECT_VERSION_ID {
			return nil, fmt.Errorf("invalid tag target UID kind: expected OBJECT_VERSION_ID, got %s", tags[i].Target.Kind.String())
		}
		if tags[i].Target.OBJECT_VERSION_ID().Value != compositionID {
			return nil, fmt.Errorf("invalid tag target UID value: expected %s, got %s", compositionID, tags[i].Target.OBJECT_VERSION_ID().Value)
		}
		if tags[i].OwnerID.ID.Kind != rm.OBJECT_ID_kind_HIER_OBJECT_ID {
			return nil, fmt.Errorf("invalid tag owner ID kind: expected HIER_OBJECT_ID, got %s", tags[i].OwnerID.ID.Kind.String())
		}
		if tags[i].OwnerID.ID.HIER_OBJECT_ID().Value != ehrID.String() {
			return nil, fmt.Errorf("invalid tag owner ID value: expected %s, got %s", ehrID.String(), tags[i].OwnerID.ID.HIER_OBJECT_ID().Value)
		}
	}

	contribution := NewContribution("Composition tags replaced", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.COMPOSITION_TYPE,
				Namespace: rm.Namespace_local,
				ID: rm.OBJECT_ID_from_OBJECT_VERSION_ID(rm.OBJECT_VERSION_ID{
					Value: compositionID,
				}),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`, contribution.UID.Value, ehrID)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Delete existing tags
	batch.Queue(`
		DELETE FROM openehr.tbl_composition_tag ct
		JOIN openehr.tbl_composition c ON c.id = ct.composition_id
		WHERE ehr_id = $1 AND data->'target'->'id'->>'value' = $2
	`, ehrID, compositionID)

	for _, tag := range tags {
		batch.Queue(`
			INSERT INTO openehr.tbl_tag (ehr_id, data)
			VALUES ($1, $2)
		`, ehrID, tag)
	}

	// Update EHR with contribution ref
	batch.Queue(`UPDATE openehr.tbl_ehr
		SET data = jsonb_insert(data, '{contributions, -1}', $1::jsonb, true)
		WHERE id = $2
	`, rm.OBJECT_REF{
		Type:      rm.CONTRIBUTION_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(contribution.UID),
	}, ehrID)

	// Update EHR with new tags
	batch.Queue(`UPDATE openehr.tbl_ehr
		SET data = jsonb_set(
			data,
			'{tags}',
			(
				SELECT jsonb_agg(item)
				FROM (
					SELECT item
					FROM jsonb_array_elements(data->'tags') AS item
					WHERE item->'target''id'->>'value' != $2::text
					UNION ALL
					SELECT item
					FROM unnest($3::jsonb[]) AS item
				)
			)
		)
		WHERE id = $1
	`, ehrID, compositionID, tags)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return nil, fmt.Errorf("failed to execute batch insert for Composition tags replacement: %w", err)
	}
	err = br.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close batch result for Composition tags replacement: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tags, nil
}

func (s *Service) DeleteVersionedCompositionTagByKey(ctx context.Context, ehrID uuid.UUID, versionedCompositionID uuid.UUID, key string) error {
	contribution := NewContribution("Versioned Composition tag deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_COMPOSITION_TYPE,
				Namespace: rm.Namespace_local,
				ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
					Value: versionedCompositionID.String(),
				}),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`, contribution.UID.Value, ehrID)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Delete the tag
	batch.Queue(`
		DELETE FROM openehr.tbl_tag
		WHERE ehr_id = $1
		  AND data->'target'->'id'->>'value' = $2
		  AND data->>'key' = $3
	`, ehrID, versionedCompositionID.String(), key)

	// Update EHR
	batch.Queue(`UPDATE openehr.tbl_ehr
		SET data = jsonb_insert(data, '{contributions, -1}', $1::jsonb, true)
		WHERE id = $2
	`, rm.OBJECT_REF{
		Type:      rm.CONTRIBUTION_TYPE,
		Namespace: rm.Namespace_local,
		ID:        rm.OBJECT_ID_from_HIER_OBJECT_ID(contribution.UID),
	}, ehrID)

	// Update EHR with new tags
	batch.Queue(`UPDATE openehr.tbl_ehr
		SET data = jsonb_set(
			data,
			'{tags}',
			(
				SELECT jsonb_agg(item)
				FROM jsonb_array_elements(data->'tags') AS item
				WHERE item->'target''id'->>'value' != $2::text
					AND item->>'key' != $3::text
			)
		)
		WHERE id = $1
	`, ehrID, versionedCompositionID, key)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("failed to execute batch insert for Versioned Composition tag deletion: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("failed to close batch result for Versioned Composition tag deletion: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) ValidateAgent(ctx context.Context, agent rm.AGENT) error {
	validateErr := agent.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Agent validation can be added here

	return nil
}

func (s *Service) CreateAgent(ctx context.Context, agent rm.AGENT) (rm.AGENT, error) {
	err := s.ValidateAgent(ctx, agent)
	if err != nil {
		return rm.AGENT{}, err
	}

	if agent.UID.E {
		// Check if agent with the same version ID already exists
		var versionedPartyIDStr string
		switch agent.UID.V.Kind {
		case rm.UID_BASED_ID_kind_HIER_OBJECT_ID:
			versionedPartyIDStr = agent.UID.V.HIER_OBJECT_ID().Value
		case rm.UID_BASED_ID_kind_OBJECT_VERSION_ID:
			versionedPartyIDStr = agent.UID.V.OBJECT_VERSION_ID().UID()
		default:
			return rm.AGENT{}, fmt.Errorf("unsupported UID kind: %d", agent.UID.V.Kind)
		}

		versionedPartyID, err := uuid.Parse(versionedPartyIDStr)
		if err != nil {
			return rm.AGENT{}, fmt.Errorf("invalid agent UID: %w", err)
		}

		exists, err := s.ExistsVersionedObject(ctx, versionedPartyID)
		if err != nil {
			return rm.AGENT{}, fmt.Errorf("failed to check existing agent: %w", err)
		}
		if exists {
			return rm.AGENT{}, ErrAgentAlreadyExists
		}
	} else {
		// Assign a new UID if not provided
		agent.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&rm.OBJECT_VERSION_ID{
			Value: fmt.Sprintf("%s::%s::%d", uuid.New().String(), config.NAMESPACE_LOCAL, 1),
		}))
	}

	versionedParty := NewVersionedParty(uuid.MustParse(agent.UID.V.OBJECT_VERSION_ID().UID()))
	agentVersion := NewOriginalVersion(agent.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_AGENT(agent), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Agent created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(agent.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert VERSIONED_PARTY
	batch.Queue(`INSERT INTO openehr.tbl_versioned_object (id, type) VALUES ($1, $2)`, versionedParty.UID.Value, rm.VERSIONED_PARTY_TYPE)
	versionedParty.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_versioned_party_data (id, data) VALUES ($1, $2)`, versionedParty.UID.Value, versionedParty)

	// Insert AGENT
	batch.Queue(`INSERT INTO openehr.tbl_agent (id, version_int, versioned_party_id, contribution_id) VALUES ($1, $2, $3, $4)`, agentVersion.UID.Value, agentVersion.UID.VersionTreeID().Int(), versionedParty.UID.Value, contribution.UID.Value)
	agentVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_agent_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, agentVersion.UID.Value, agentVersion)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to execute batch insert for Agent creation: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to close batch result for Agent creation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return agent, nil
}

func (s *Service) GetAgentID(ctx context.Context, versionedPartyID uuid.UUID) (rm.OBJECT_VERSION_ID, error) {
	query := `
		SELECT a.id 
		FROM openehr.tbl_agent a
		WHERE a.versioned_party_id = $1
		ORDER BY version_int DESC
		LIMIT 1
	`
	row := s.DB.QueryRow(ctx, query, versionedPartyID)

	var agentID string
	err := row.Scan(&agentID)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.OBJECT_VERSION_ID{}, ErrAgentNotFound
		}
		return rm.OBJECT_VERSION_ID{}, fmt.Errorf("failed to get latest agent by versioned party ID: %w", err)
	}

	return rm.OBJECT_VERSION_ID{Value: agentID}, nil
}

func (s *Service) GetAgentByVersionedPartyIDRawJSON(ctx context.Context, versionedPartyID uuid.UUID) ([]byte, error) {
	query := `
		SELECT ad.data 
		FROM openehr.tbl_agent a
		JOIN openehr.tbl_agent_data ad ON ad.id = a.id 
		WHERE a.versioned_party_id = $1
		ORDER BY version_int DESC
		LIMIT 1
	`

	row := s.DB.QueryRow(ctx, query, versionedPartyID)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrAgentNotFound
		}
		return nil, fmt.Errorf("failed to get latest agent by versioned party ID: %w", err)
	}

	return data, nil
}

func (s *Service) GetAgentAtVersionRawJSON(ctx context.Context, agentID string) ([]byte, error) {
	query := `
		SELECT data 
		FROM openehr.tbl_agent_data
		WHERE id = $1
		LIMIT 1
	`

	row := s.DB.QueryRow(ctx, query, agentID)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrAgentNotFound
		}
		return nil, fmt.Errorf("failed to get agent at version from database: %w", err)
	}

	return data, nil
}

func (s *Service) UpdateAgent(ctx context.Context, currentAgentID rm.OBJECT_VERSION_ID, nextAgent rm.AGENT) (rm.AGENT, error) {
	err := s.ValidateAgent(ctx, nextAgent)
	if err != nil {
		return rm.AGENT{}, err
	}

	if !nextAgent.UID.E {
		nextAgent.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: currentAgentID.UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(nextAgent.UID.V, currentAgentID)
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to upgrade current Agent UID: %w", err)
	}
	nextAgent.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	agentVersion := NewOriginalVersion(nextAgent.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_AGENT(nextAgent), utils.Some(currentAgentID))
	contribution := NewContribution("Agent updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(nextAgent.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert AGENT
	batch.Queue(`INSERT INTO openehr.tbl_agent (id, version_int, versioned_party_id, type, contribution_id) VALUES ($1, $2, $3, $4, $5)`, agentVersion.UID.Value, 1, uuid.MustParse(strings.Split(nextAgent.UID.V.OBJECT_VERSION_ID().UID(), "::")[0]), rm.AGENT_TYPE, contribution.UID.Value)
	agentVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_agent_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, agentVersion.UID.Value, agentVersion)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to execute batch insert for Agent update: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to close batch result for Agent update: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextAgent, nil
}

func (s *Service) DeleteAgent(ctx context.Context, versionedPartyID uuid.UUID) error {
	contribution := NewContribution("Agent deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      rm.AGENT_TYPE,
				ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
					Value: versionedPartyID.String(),
				}),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Delete VERSIONED_PARTY
	batch.Queue(`DELETE FROM openehr.tbl_versioned_object WHERE id = $1`, versionedPartyID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("failed to execute batch insert for Agent deletion: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("failed to close batch result for Agent deletion: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) ValidatePerson(ctx context.Context, person rm.PERSON) error {
	validateErr := person.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Person validation can be added here

	return nil
}

func (s *Service) CreatePerson(ctx context.Context, person rm.PERSON) (rm.PERSON, error) {
	err := s.ValidatePerson(ctx, person)
	if err != nil {
		return rm.PERSON{}, err
	}

	if person.UID.E {
		// Check if Version Object ID already exists
		var versionedPartyIDStr string
		switch person.UID.V.Kind {
		case rm.UID_BASED_ID_kind_HIER_OBJECT_ID:
			versionedPartyIDStr = person.UID.V.HIER_OBJECT_ID().Value
		case rm.UID_BASED_ID_kind_OBJECT_VERSION_ID:
			versionedPartyIDStr = person.UID.V.OBJECT_VERSION_ID().UID()
		default:
			return rm.PERSON{}, fmt.Errorf("unsupported UID kind: %d", person.UID.V.Kind)
		}

		versionedPartyID, err := uuid.Parse(versionedPartyIDStr)
		if err != nil {
			return rm.PERSON{}, fmt.Errorf("invalid person UID: %w", err)
		}

		exists, err := s.ExistsVersionedObject(ctx, versionedPartyID)
		if err != nil {
			return rm.PERSON{}, fmt.Errorf("failed to check existing person: %w", err)
		}
		if exists {
			return rm.PERSON{}, ErrPersonAlreadyExists
		}
	} else {
		// Assign a new UID if not provided
		person.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&rm.OBJECT_VERSION_ID{
			Value: fmt.Sprintf("%s::%s::%d", uuid.New().String(), config.NAMESPACE_LOCAL, 1),
		}))
	}

	versionedParty := NewVersionedParty(uuid.MustParse(person.UID.V.OBJECT_VERSION_ID().UID()))
	personVersion := NewOriginalVersion(person.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_PERSON(person), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Person created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(person.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert VERSIONED_PARTY
	batch.Queue(`INSERT INTO openehr.tbl_versioned_object (id, type) VALUES ($1, $2)`, versionedParty.UID.Value, rm.VERSIONED_PARTY_TYPE)
	versionedParty.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_versioned_party_data (id, data) VALUES ($1, $2)`, versionedParty.UID.Value, versionedParty)

	// Insert PERSON
	batch.Queue(`INSERT INTO openehr.tbl_person (id, version_int, versioned_party_id, contribution_id) VALUES ($1, $2, $3, $4)`, personVersion.UID.Value, personVersion.UID.VersionTreeID().Int(), versionedParty.UID.Value, contribution.UID.Value)
	personVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_person_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, personVersion.UID.Value, personVersion)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to execute batch insert for Person creation: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to close batch result for Person creation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return person, nil
}

func (s *Service) GetCurrentPersonID(ctx context.Context, versionedPartyID uuid.UUID) (rm.OBJECT_VERSION_ID, error) {
	query := `
		SELECT p.id 
		FROM openehr.tbl_person p
		WHERE p.versioned_party_id = $1
		ORDER BY version_int DESC
		LIMIT 1
	`
	row := s.DB.QueryRow(ctx, query, versionedPartyID)

	var personID string
	err := row.Scan(&personID)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.OBJECT_VERSION_ID{}, ErrPersonNotFound
		}
		return rm.OBJECT_VERSION_ID{}, fmt.Errorf("failed to get latest person by versioned party ID: %w", err)
	}

	return rm.OBJECT_VERSION_ID{Value: personID}, nil
}

func (s *Service) GetPersonByVersionedPartyIDRawJSON(ctx context.Context, versionedPartyID uuid.UUID) ([]byte, error) {
	query := `
		SELECT pd.data 
		FROM openehr.tbl_person p
		JOIN openehr.tbl_person_data pd ON pd.id = p.id 
		WHERE p.versioned_party_id = $1
		ORDER BY version_int DESC
		LIMIT 1
	`

	row := s.DB.QueryRow(ctx, query, versionedPartyID)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrPersonNotFound
		}
		return nil, fmt.Errorf("failed to get latest person by versioned party ID: %w", err)
	}

	return data, nil
}

func (s *Service) GetPersonByIDRawJSON(ctx context.Context, personID string) ([]byte, error) {
	query := `SELECT data FROM openehr.tbl_person_data WHERE id = $1 LIMIT 1`

	var data []byte
	err := s.DB.QueryRow(ctx, query, personID).Scan(&data)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrPersonNotFound
		}
		return nil, fmt.Errorf("failed to get person: %w", err)
	}

	return data, nil
}

func (s *Service) UpdatePerson(ctx context.Context, currentPersonID rm.OBJECT_VERSION_ID, nextPerson rm.PERSON) (rm.PERSON, error) {
	err := s.ValidatePerson(ctx, nextPerson)
	if err != nil {
		return rm.PERSON{}, err
	}

	if !nextPerson.UID.E {
		nextPerson.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: currentPersonID.Value,
		}))
	}

	updatedID, err := UpgradeObjectVersionID(nextPerson.UID.V, currentPersonID)
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to upgrade current Person UID: %w", err)
	}
	nextPerson.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	personVersion := NewOriginalVersion(nextPerson.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_PERSON(nextPerson), utils.Some(currentPersonID))
	contribution := NewContribution("Person updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(nextPerson.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert PERSON
	batch.Queue(`INSERT INTO openehr.tbl_person (id, version_int, versioned_party_id, type, contribution_id) VALUES ($1, $2, $3, $4, $5)`, personVersion.UID.Value, 1, uuid.MustParse(strings.Split(nextPerson.UID.V.OBJECT_VERSION_ID().UID(), "::")[0]), rm.PERSON_TYPE, contribution.UID.Value)
	personVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_person_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, personVersion.UID.Value, personVersion)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to execute batch insert for Person creation: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to close batch result for Person creation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextPerson, nil
}

func (s *Service) DeletePerson(ctx context.Context, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Person deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      rm.PERSON_TYPE,
				ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
					Value: versionedObjectID.String(),
				}),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Delete PERSON
	batch.Queue(`DELETE FROM openehr.tbl_versioned_object WHERE id = $1`, versionedObjectID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("failed to execute batch insert for Person deletion: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("failed to close batch result for Person deletion: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) ValidateGroup(ctx context.Context, group rm.GROUP) error {
	validateErr := group.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Group validation can be added here

	return nil
}

func (s *Service) CreateGroup(ctx context.Context, group rm.GROUP) (rm.GROUP, error) {
	err := s.ValidateGroup(ctx, group)
	if err != nil {
		return rm.GROUP{}, err
	}

	if group.UID.E {
		// Check if versioned object ID already exists
		var versionedPartyIDStr string
		switch group.UID.V.Kind {
		case rm.UID_BASED_ID_kind_HIER_OBJECT_ID:
			versionedPartyIDStr = group.UID.V.HIER_OBJECT_ID().Value
		case rm.UID_BASED_ID_kind_OBJECT_VERSION_ID:
			versionedPartyIDStr = group.UID.V.OBJECT_VERSION_ID().UID()
		default:
			return rm.GROUP{}, fmt.Errorf("unsupported UID kind: %d", group.UID.V.Kind)
		}

		versionedPartyID, err := uuid.Parse(versionedPartyIDStr)
		if err != nil {
			return rm.GROUP{}, fmt.Errorf("invalid group UID: %w", err)
		}

		exists, err := s.ExistsVersionedObject(ctx, versionedPartyID)
		if err != nil {
			return rm.GROUP{}, fmt.Errorf("failed to check existing group: %w", err)
		}
		if exists {
			return rm.GROUP{}, ErrGroupAlreadyExists
		}
	} else {
		// Assign a new UID if not provided
		group.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&rm.OBJECT_VERSION_ID{
			Value: fmt.Sprintf("%s::%s::%d", uuid.New().String(), config.NAMESPACE_LOCAL, 1),
		}))
	}

	versionedParty := NewVersionedParty(uuid.MustParse(group.UID.V.OBJECT_VERSION_ID().UID()))
	groupVersion := NewOriginalVersion(group.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_GROUP(group), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Group created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(group.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert VERSIONED_PARTY
	batch.Queue(`INSERT INTO openehr.tbl_versioned_object (id, type) VALUES ($1, $2)`, versionedParty.UID.Value, rm.VERSIONED_PARTY_TYPE)
	versionedParty.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_versioned_party_data (id, data) VALUES ($1, $2)`, versionedParty.UID.Value, versionedParty)

	// Insert GROUP
	batch.Queue(`INSERT INTO openehr.tbl_group (id, version_int, versioned_party_id, contribution_id) VALUES ($1, $2, $3, $4)`, groupVersion.UID.Value, groupVersion.UID.VersionTreeID().Int(), versionedParty.UID.Value, contribution.UID.Value)
	groupVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_group_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, groupVersion.UID.Value, groupVersion)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to execute batch insert for Group creation: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to close batch result for Group creation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return group, nil
}

func (s *Service) GetGroupID(ctx context.Context, versionedPartyID uuid.UUID) (rm.OBJECT_VERSION_ID, error) {
	query := `
		SELECT g.id 
		FROM openehr.tbl_group g
		WHERE g.versioned_party_id = $1
		ORDER BY version_int DESC
		LIMIT 1
	`
	row := s.DB.QueryRow(ctx, query, versionedPartyID)

	var groupID string
	err := row.Scan(&groupID)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.OBJECT_VERSION_ID{}, ErrGroupNotFound
		}
		return rm.OBJECT_VERSION_ID{}, fmt.Errorf("failed to get latest group by versioned party ID: %w", err)
	}

	return rm.OBJECT_VERSION_ID{Value: groupID}, nil
}

func (s *Service) GetGroupByVersionedPartyIDRawJSON(ctx context.Context, versionedPartyID uuid.UUID) ([]byte, error) {
	query := `
		SELECT gd.data
		FROM openehr.tbl_group g
		JOIN openehr.tbl_group_data gd ON gd.id = g.id 
		WHERE g.versioned_object_id = $1
		ORDER BY g.version_int DESC
		LIMIT 1
	`

	row := s.DB.QueryRow(ctx, query, versionedPartyID)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, fmt.Errorf("failed to get latest group by versioned party ID: %w", err)
	}

	return data, nil
}

func (s *Service) GetGroupRawJSON(ctx context.Context, groupID string) ([]byte, error) {
	query := `SELECT data FROM tbl_group_data WHERE id = $1 LIMIT 1`

	var data []byte
	err := s.DB.QueryRow(ctx, query, groupID).Scan(&data)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	return data, nil
}

func (s *Service) UpdateGroup(ctx context.Context, currentGroupID rm.OBJECT_VERSION_ID, nextGroup rm.GROUP) (rm.GROUP, error) {
	err := s.ValidateGroup(ctx, nextGroup)
	if err != nil {
		return rm.GROUP{}, err
	}

	if !nextGroup.UID.E {
		nextGroup.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: currentGroupID.UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(nextGroup.UID.V, currentGroupID)
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to upgrade current Group UID: %w", err)
	}
	nextGroup.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	groupVersion := NewOriginalVersion(nextGroup.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_GROUP(nextGroup), utils.Some(currentGroupID))
	contribution := NewContribution("Group updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(nextGroup.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert GROUP
	batch.Queue(`INSERT INTO openehr.tbl_group (id, version_int, versioned_party_id, type, contribution_id) VALUES ($1, $2, $3, $4, $5)`, groupVersion.UID.Value, 1, uuid.MustParse(strings.Split(nextGroup.UID.V.OBJECT_VERSION_ID().UID(), "::")[0]), rm.GROUP_TYPE, contribution.UID.Value)
	groupVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_group_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, groupVersion.UID.Value, groupVersion)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to execute batch insert for Group creation: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to close batch result for Group creation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextGroup, nil
}

func (s *Service) DeleteGroup(ctx context.Context, versionedPartyID uuid.UUID) error {
	contribution := NewContribution("Group deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      rm.GROUP_TYPE,
				ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
					Value: versionedPartyID.String(),
				}),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Delete VERSIONED_PARTY
	batch.Queue(`DELETE FROM openehr.tbl_versioned_object WHERE id = $1`, versionedPartyID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("failed to execute batch insert for Group deletion: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("failed to close batch result for Group deletion: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) ValidateOrganisation(ctx context.Context, organisation rm.ORGANISATION) error {
	validateErr := organisation.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Organisation validation can be added here

	return nil
}

func (s *Service) CreateOrganisation(ctx context.Context, organisation rm.ORGANISATION) (rm.ORGANISATION, error) {
	err := s.ValidateOrganisation(ctx, organisation)
	if err != nil {
		return rm.ORGANISATION{}, err
	}

	if organisation.UID.E {
		// Check if versioned object ID already exists
		var versionedPartyIDStr string
		switch organisation.UID.V.Kind {
		case rm.UID_BASED_ID_kind_HIER_OBJECT_ID:
			versionedPartyIDStr = organisation.UID.V.HIER_OBJECT_ID().Value
		case rm.UID_BASED_ID_kind_OBJECT_VERSION_ID:
			versionedPartyIDStr = organisation.UID.V.OBJECT_VERSION_ID().UID()
		default:
			return rm.ORGANISATION{}, fmt.Errorf("unsupported UID kind: %d", organisation.UID.V.Kind)
		}

		versionedPartyID, err := uuid.Parse(versionedPartyIDStr)
		if err != nil {
			return rm.ORGANISATION{}, fmt.Errorf("invalid organisation UID: %w", err)
		}

		exists, err := s.ExistsVersionedObject(ctx, versionedPartyID)
		if err != nil {
			return rm.ORGANISATION{}, fmt.Errorf("failed to check existing organisation: %w", err)
		}
		if exists {
			return rm.ORGANISATION{}, ErrOrganisationAlreadyExists
		}
	} else {
		// Assign a new UID if not provided
		organisation.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&rm.OBJECT_VERSION_ID{
			Value: fmt.Sprintf("%s::%s::%d", uuid.New().String(), config.NAMESPACE_LOCAL, 1),
		}))
	}

	versionedParty := NewVersionedParty(uuid.MustParse(organisation.UID.V.OBJECT_VERSION_ID().UID()))
	organisationVersion := NewOriginalVersion(organisation.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_ORGANISATION(organisation), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Organisation created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(organisation.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert VERSIONED_PARTY
	batch.Queue(`INSERT INTO openehr.tbl_versioned_object (id, type) VALUES ($1, $2)`, versionedParty.UID.Value, rm.VERSIONED_PARTY_TYPE)
	versionedParty.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_versioned_party_data (id, data) VALUES ($1, $2)`, versionedParty.UID.Value, versionedParty)

	// Insert ORGANISATION
	batch.Queue(`INSERT INTO openehr.tbl_organisation (id, version_int, versioned_party_id, contribution_id) VALUES ($1, $2, $3, $4)`, organisationVersion.UID.Value, organisationVersion.UID.VersionTreeID().Int(), versionedParty.UID.Value, contribution.UID.Value)
	organisationVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_organisation_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, organisationVersion.UID.Value, organisationVersion)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to execute batch insert for Organisation creation: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to close batch result for Organisation creation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return organisation, nil
}

func (s *Service) GetCurrentOrganisationID(ctx context.Context, versionedPartyID uuid.UUID) (rm.OBJECT_VERSION_ID, error) {
	query := `
		SELECT o.id 
		FROM openehr.tbl_organisation o
		WHERE o.versioned_party_id = $1
		ORDER BY version_int DESC
		LIMIT 1
	`
	row := s.DB.QueryRow(ctx, query, versionedPartyID)

	var organisationID string
	err := row.Scan(&organisationID)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.OBJECT_VERSION_ID{}, ErrOrganisationNotFound
		}
		return rm.OBJECT_VERSION_ID{}, fmt.Errorf("failed to get latest organisation by versioned party ID: %w", err)
	}

	return rm.OBJECT_VERSION_ID{Value: organisationID}, nil
}

func (s *Service) GetOrganisationByVersionedPartyIDRawJSON(ctx context.Context, versionedPartyID uuid.UUID) ([]byte, error) {
	query := `
		SELECT od.data 
		FROM openehr.tbl_organisation o
		JOIN openehr.tbl_organisation_data od ON od.id = o.id 
		WHERE o.versioned_party_id = $1
		ORDER BY o.version_int DESC
		LIMIT 1
	`

	row := s.DB.QueryRow(ctx, query, versionedPartyID)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrOrganisationNotFound
		}
		return nil, fmt.Errorf("failed to get latest organisation by versioned party ID: %w", err)
	}

	return data, nil
}

func (s *Service) GetOrganisationByIDRawJSON(ctx context.Context, organisationID string) ([]byte, error) {
	query := `SELECT data FROM openehr.tbl_organisation_data WHERE id = $1 LIMIT 1`

	row := s.DB.QueryRow(ctx, query, organisationID)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, ErrOrganisationNotFound
		}
		return nil, fmt.Errorf("failed to get organisation: %w", err)
	}

	return data, nil
}

func (s *Service) UpdateOrganisation(ctx context.Context, currentOrganisationID rm.OBJECT_VERSION_ID, nextOrganisation rm.ORGANISATION) (rm.ORGANISATION, error) {
	err := s.ValidateOrganisation(ctx, nextOrganisation)
	if err != nil {
		return rm.ORGANISATION{}, err
	}

	if !nextOrganisation.UID.E {
		nextOrganisation.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: currentOrganisationID.UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(nextOrganisation.UID.V, currentOrganisationID)
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to upgrade current Organisation UID: %w", err)
	}
	nextOrganisation.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	organisationVersion := NewOriginalVersion(nextOrganisation.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_ORGANISATION(nextOrganisation), utils.Some(currentOrganisationID))
	contribution := NewContribution("Organisation updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(nextOrganisation.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert ORGANISATION
	batch.Queue(`INSERT INTO openehr.tbl_organisation (id, version_int, versioned_party_id, type, contribution_id) VALUES ($1, $2, $3, $4, $5)`, organisationVersion.UID.Value, 1, uuid.MustParse(strings.Split(nextOrganisation.UID.V.OBJECT_VERSION_ID().UID(), "::")[0]), rm.ORGANISATION_TYPE, contribution.UID.Value)
	organisationVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_organisation_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, organisationVersion.UID.Value, organisationVersion)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to execute batch insert for Organisation creation: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to close batch result for Organisation creation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextOrganisation, nil
}

func (s *Service) DeleteOrganisation(ctx context.Context, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Organisation deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      rm.ORGANISATION_TYPE,
				ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
					Value: versionedObjectID.String(),
				}),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Delete ORGANISATION
	batch.Queue(`DELETE FROM openehr.tbl_versioned_object WHERE id = $1`, versionedObjectID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("failed to execute batch insert for Organisation deletion: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("failed to close batch result for Organisation deletion: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) ValidateRole(ctx context.Context, role rm.ROLE) error {
	validateErr := role.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Organisation validation can be added here

	return nil
}

func (s *Service) CreateRole(ctx context.Context, role rm.ROLE) (rm.ROLE, error) {
	err := s.ValidateRole(ctx, role)
	if err != nil {
		return rm.ROLE{}, err
	}

	if role.UID.E {
		// Check if versioned object ID already exists
		var versionedPartyIDStr string
		switch role.UID.V.Kind {
		case rm.UID_BASED_ID_kind_HIER_OBJECT_ID:
			versionedPartyIDStr = role.UID.V.HIER_OBJECT_ID().Value
		case rm.UID_BASED_ID_kind_OBJECT_VERSION_ID:
			versionedPartyIDStr = role.UID.V.OBJECT_VERSION_ID().UID()
		default:
			return rm.ROLE{}, fmt.Errorf("unsupported UID kind: %d", role.UID.V.Kind)
		}

		versionedPartyID, err := uuid.Parse(versionedPartyIDStr)
		if err != nil {
			return rm.ROLE{}, fmt.Errorf("invalid role UID: %w", err)
		}
		exists, err := s.ExistsVersionedObject(ctx, versionedPartyID)
		if err != nil {
			return rm.ROLE{}, fmt.Errorf("failed to check existing role: %w", err)
		}
		if exists {
			return rm.ROLE{}, ErrRoleAlreadyExists
		}
	} else {
		// Assign a new UID if not provided
		role.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&rm.OBJECT_VERSION_ID{
			Value: fmt.Sprintf("%s::%s::%d", uuid.New().String(), config.NAMESPACE_LOCAL, 1),
		}))
	}

	versionedParty := NewVersionedParty(uuid.MustParse(role.UID.V.OBJECT_VERSION_ID().UID()))
	roleVersion := NewOriginalVersion(role.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_ROLE(role), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Role created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(role.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert VERSIONED_PARTY
	batch.Queue(`INSERT INTO openehr.tbl_versioned_object (id, type) VALUES ($1, $2)`, versionedParty.UID.Value, rm.VERSIONED_PARTY_TYPE)
	versionedParty.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_versioned_party_data (id, data) VALUES ($1, $2)`, versionedParty.UID.Value, versionedParty)

	// Insert ROLE
	batch.Queue(`INSERT INTO openehr.tbl_role (id, version_int, versioned_party_id, contribution_id) VALUES ($1, $2, $3, $4)`, roleVersion.UID.Value, roleVersion.UID.VersionTreeID().Int(), versionedParty.UID.Value, contribution.UID.Value)
	roleVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_role_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, roleVersion.UID.Value, roleVersion)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to execute batch insert for Role creation: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to close batch result for Role creation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return role, nil
}

func (s *Service) GetCurrentRoleVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (rm.ROLE, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var role rm.ROLE
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), rm.ROLE_TYPE).Scan(&role)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.ROLE{}, ErrRoleNotFound
		}
		return rm.ROLE{}, fmt.Errorf("failed to get latest role by versioned party ID: %w", err)
	}

	return role, nil
}

func (s *Service) GetRoleRawJSON(ctx context.Context, versionID string) (rm.ROLE, error) {
	query := `SELECT data FROM openehr.tbl_role_data WHERE id = $1 LIMIT 1`

	var role rm.ROLE
	err := s.DB.QueryRow(ctx, query, versionID, rm.ROLE_TYPE).Scan(&role)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.ROLE{}, ErrRoleNotFound
		}
		return rm.ROLE{}, fmt.Errorf("failed to get role at version: %w", err)
	}

	return role, nil
}

func (s *Service) UpdateRole(ctx context.Context, versionedPartyID uuid.UUID, role rm.ROLE) (rm.ROLE, error) {
	err := s.ValidateRole(ctx, role)
	if err != nil {
		return rm.ROLE{}, err
	}

	currentRole, err := s.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			return rm.ROLE{}, ErrRoleNotFound
		}
		return rm.ROLE{}, fmt.Errorf("failed to get current Role: %w", err)
	}
	currentRoleID := currentRole.UID.V.Value.(*rm.OBJECT_VERSION_ID)

	if !role.UID.E {
		role.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: currentRoleID.UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(currentRole.UID.V, *currentRoleID)
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to upgrade current Role UID: %w", err)
	}
	role.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	roleVersion := NewOriginalVersion(role.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_ROLE(role), utils.Some(*currentRoleID))
	contribution := NewContribution("Role updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(role.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Insert ROLE
	batch.Queue(`INSERT INTO openehr.tbl_role (id, version_int, versioned_party_id, type, contribution_id) VALUES ($1, $2, $3, $4, $5)`, roleVersion.UID.Value, roleVersion.UID.VersionTreeID().Int(), versionedPartyID, rm.ROLE_TYPE, contribution.UID.Value)
	roleVersion.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_role_data (id, data, version_data) VALUES ($1, ($2::jsonb)->'data', jsonb_set($2::jsonb, '{data}', 'null', true))`, roleVersion.UID.Value, roleVersion)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to execute batch insert for Role update: %w", err)
	}
	err = br.Close()
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to close batch result for Role update: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return role, nil
}

func (s *Service) DeleteRole(ctx context.Context, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Role deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      rm.ROLE_TYPE,
				ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
					Value: versionedObjectID.String(),
				}),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	batch := &pgx.Batch{}

	// Insert CONTRIBUTION
	batch.Queue(`INSERT INTO openehr.tbl_contribution (id) VALUES ($1)`, contribution.UID.Value)
	contribution.SetModelName()
	batch.Queue(`INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`, contribution.UID.Value, contribution)

	// Delete ROLE
	batch.Queue(`DELETE FROM openehr.tbl_versioned_object WHERE id = $1`, versionedObjectID)

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return fmt.Errorf("failed to execute batch insert for Role deletion: %w", err)
	}
	err = br.Close()
	if err != nil {
		return fmt.Errorf("failed to close batch result for Role deletion: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) GetVersionedPartyRawJSON(ctx context.Context, versionedPartyID uuid.UUID) ([]byte, error) {
	query := `SELECT data FROM openehr.tbl_versioned_object WHERE id = $1 AND type = $2 LIMIT 1`

	row := s.DB.QueryRow(ctx, query, versionedPartyID, rm.VERSIONED_PARTY_TYPE)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to get versioned party by ID: %w", err)
	}
	return data, nil
}

func (s *Service) GetVersionedPartyRevisionHistoryRawJSON(ctx context.Context, versionedObjectID uuid.UUID) ([]byte, error) {
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
	args := []any{[]string{rm.AGENT_TYPE, rm.PERSON_TYPE, rm.GROUP_TYPE, rm.ORGANISATION_TYPE, rm.ROLE_TYPE}, versionedObjectID}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrRevisionHistoryNotFound
		}
		return nil, fmt.Errorf("failed to fetch Revision History from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetVersionedPartyVersionAtTimeRawJSON(ctx context.Context, versionedPartyID uuid.UUID, filterAtTime time.Time) ([]byte, error) {
	query := `
		SELECT data 
		FROM (
			SELECT data, created_at
			FROM openehr.tbl_agent a
			JOIN openehr.tbl_agent_data ad ON ad.id = a.id
			WHERE a.versioned_party_id = $1
			UNION ALL
			SELECT data, created_at
			FROM openehr.tbl_person p
			JOIN openehr.tbl_person_data pd ON pd.id = p.id
			WHERE p.versioned_party_id = $1
			UNION ALL
			SELECT data, created_at
			FROM openehr.tbl_group g
			JOIN openehr.tbl_group_data gd ON gd.id = g.id
			WHERE g.versioned_party_id = $1
			UNION ALL
			SELECT data, created_at
			FROM openehr.tbl_organisation o
			JOIN openehr.tbl_organisation_data od ON od.id = o.id
			WHERE o.versioned_party_id = $1
			UNION ALL
			SELECT data, created_at
			FROM openehr.tbl_role r
			JOIN openehr.tbl_role_data rd ON rd.id = r.id
			WHERE r.versioned_party_id = $1
		)
	`
	args := []any{versionedPartyID}

	if !filterAtTime.IsZero() {
		query += `AND ov.created_at <= $2 `
		args = append(args, filterAtTime)
	}
	query += `ORDER BY ov.created_at DESC LIMIT 1`

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrVersionedPartyVersionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Party version at time from database: %w", err)
	}

	return data, nil
}

func (s *Service) GetContributionRawJSON(ctx context.Context, contributionID string, ehrID utils.Optional[uuid.UUID]) ([]byte, error) {
	query := `
		SELECT cd.data
		FROM openehr.tbl_contribution c
		JOIN openehr.tbl_contribution_data cd ON c.id = cd.id
		WHERE c.ehr_id = $1 AND c.id = $2
		LIMIT 1
	`
	args := []any{ehrID, contributionID}

	row := s.DB.QueryRow(ctx, query, args...)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrContributionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Contribution by ID from database: %w", err)
	}

	return data, nil
}

func (s *Service) ExistsVersionedObject(ctx context.Context, versionedObjectID uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_versioned_object WHERE id = $1 LIMIT 1`
	args := []any{versionedObjectID}

	var exists int
	err := s.DB.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if versioned object exists: %w", err)
	}
	return true, nil
}

// func (s *Service) SaveEHRWithTx(ctx context.Context, tx pgx.Tx, ehr rm.EHR) error {
// 	query := `INSERT INTO openehr.tbl_ehr (id, data) VALUES ($1, $2)`
// 	args := []any{ehr.EHRID.Value, ehr}
// 	_, err := tx.Exec(ctx, query, args...)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert EHR into the database: %w", err)
// 	}

// 	return nil
// }

// func (s *Service) SaveContributionWithTx(ctx context.Context, tx pgx.Tx, contribution rm.CONTRIBUTION, ehrID utils.Optional[uuid.UUID]) error {
// 	// Insert Contribution
// 	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
// 	args := []any{contribution.UID.Value, ehrID}
// 	_, err := tx.Exec(ctx, query, args...)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert contribution into the database: %w", err)
// 	}

// 	contribution.SetModelName()
// 	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
// 	args = []any{contribution.UID.Value, contribution}
// 	_, err = tx.Exec(ctx, query, args...)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert contribution data into the database: %w", err)
// 	}

// 	return nil
// }

// func (s *Service) SaveVersionedObjectWithTx(ctx context.Context, tx pgx.Tx, versionedObject any, ehrID utils.Optional[uuid.UUID]) error {
// 	var (
// 		modelType string
// 		id        string
// 	)
// 	switch v := versionedObject.(type) {
// 	case rm.VERSIONED_EHR_STATUS:
// 		v.SetModelName()
// 		modelType = rm.VERSIONED_EHR_STATUS_TYPE
// 		id = v.UID.Value
// 	case rm.VERSIONED_EHR_ACCESS:
// 		v.SetModelName()
// 		modelType = rm.VERSIONED_EHR_ACCESS_TYPE
// 		id = v.UID.Value
// 	case rm.VERSIONED_COMPOSITION:
// 		v.SetModelName()
// 		modelType = rm.VERSIONED_COMPOSITION_TYPE
// 		id = v.UID.Value
// 	case rm.VERSIONED_FOLDER:
// 		v.SetModelName()
// 		modelType = rm.VERSIONED_FOLDER_TYPE
// 		id = v.UID.Value
// 	case rm.VERSIONED_PARTY:
// 		v.SetModelName()
// 		modelType = rm.VERSIONED_PARTY_TYPE
// 		id = v.UID.Value
// 	default:
// 		return fmt.Errorf("unsupported versioned object type for creation: %T", versionedObject)
// 	}

// 	query := `INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id) VALUES ($1, $2, $3)`
// 	args := []any{id, modelType, ehrID}
// 	_, err := tx.Exec(ctx, query, args...)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert versioned object into the database: %w", err)
// 	}

// 	query = `INSERT INTO openehr.tbl_versioned_object_data (id, data) VALUES ($1, $2)`
// 	args = []any{id, versionedObject}
// 	_, err = tx.Exec(ctx, query, args...)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert versioned object data into the database: %w", err)
// 	}

// 	return nil
// }

// func (s *Service) SaveObjectVersionWithTx(ctx context.Context, tx pgx.Tx, version any, contributionID string, ehrID utils.Optional[uuid.UUID]) error {
// 	var data rm.OriginalVersionDataUnion
// 	switch v := version.(type) {
// 	case rm.ORIGINAL_VERSION:
// 		v.SetModelName()
// 		data = v.Data
// 	// After enabling, make sure to change the data path below in the INSERT statement
// 	// case rm.IMPORTED_VERSION:
// 	// 	object = v.Data
// 	default:
// 		return fmt.Errorf("unsupported version type for object version creation: %T", version)
// 	}

// 	var (
// 		modelType string
// 		id        rm.OBJECT_VERSION_ID
// 	)
// 	switch data.Kind {
// 	case rm.ORIGINAL_VERSION_data_kind_EHR_STATUS:
// 		modelType = rm.EHR_STATUS_TYPE
// 		id = *data.EHR_STATUS().UID.V.OBJECT_VERSION_ID()
// 	case rm.ORIGINAL_VERSION_data_kind_EHR_ACCESS:
// 		modelType = rm.EHR_ACCESS_TYPE
// 		id = *data.EHR_ACCESS().UID.V.OBJECT_VERSION_ID()
// 	case rm.ORIGINAL_VERSION_data_kind_COMPOSITION:
// 		modelType = rm.COMPOSITION_TYPE
// 		id = *data.COMPOSITION().UID.V.OBJECT_VERSION_ID()
// 	case rm.ORIGINAL_VERSION_data_kind_FOLDER:
// 		modelType = rm.FOLDER_TYPE
// 		id = *data.FOLDER().UID.V.OBJECT_VERSION_ID()
// 	case rm.ORIGINAL_VERSION_data_kind_ROLE:
// 		modelType = rm.ROLE_TYPE
// 		id = *data.ROLE().UID.V.OBJECT_VERSION_ID()
// 	case rm.ORIGINAL_VERSION_data_kind_PERSON:
// 		modelType = rm.PERSON_TYPE
// 		id = *data.PERSON().UID.V.OBJECT_VERSION_ID()
// 	case rm.ORIGINAL_VERSION_data_kind_AGENT:
// 		modelType = rm.AGENT_TYPE
// 		id = *data.AGENT().UID.V.OBJECT_VERSION_ID()
// 	case rm.ORIGINAL_VERSION_data_kind_GROUP:
// 		modelType = rm.GROUP_TYPE
// 		id = *data.GROUP().UID.V.OBJECT_VERSION_ID()
// 	case rm.ORIGINAL_VERSION_data_kind_ORGANISATION:
// 		modelType = rm.ORGANISATION_TYPE
// 		id = *data.ORGANISATION().UID.V.OBJECT_VERSION_ID()
// 	default:
// 		return fmt.Errorf("unsupported object type for version creation: %d", data.Kind)
// 	}

// 	query := `INSERT INTO openehr.tbl_object_version (id, version_int, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5, $6)`
// 	args := []any{id.Value, id.VersionTreeID().Int(), id.UID(), modelType, ehrID, contributionID}
// 	if _, err := tx.Exec(ctx, query, args...); err != nil {
// 		return fmt.Errorf("failed to insert object version into the database: %w", err)
// 	}

// 	query = `INSERT INTO openehr.tbl_object_version_data (id, version_data, object_data) VALUES ($1, jsonb_set($2, '{data}', 'null', true), $2->'data')`
// 	args = []any{id.Value, version}
// 	if _, err := tx.Exec(ctx, query, args...); err != nil {
// 		return fmt.Errorf("failed to insert object version data into the database: %w", err)
// 	}

// 	return nil
// }

func (s *Service) QueryWithStream(ctx context.Context, w io.Writer, aqlQuery string, aqlParams map[string]any) error {
	if aqlParams == nil {
		aqlParams = make(map[string]any)
	}

	sqlQuery, _, err := aql.ToSQL(aqlQuery, aqlParams)
	if err != nil {
		s.Logger.Error("internal error", "error", err)
		return err
	}

	s.Logger.DebugContext(ctx, "query error", "error", err, "aql", aqlQuery, "sql", strings.ReplaceAll(strings.ReplaceAll(sqlQuery, "\n", " "), "\t", " "))

	rows, err := s.DB.Query(ctx, sqlQuery)
	if err != nil {
		s.Logger.ErrorContext(ctx, "query error", "error", err, "aql", aqlQuery, "sql", strings.ReplaceAll(strings.ReplaceAll(sqlQuery, "\n", " "), "\t", " "))
		return err
	}

	// Stream results as JSON array
	_, _ = w.Write([]byte(`{"rows":[`))

	first := true
	for rows.Next() {
		var jsonData []byte
		if err := rows.Scan(&jsonData); err != nil {
			s.Logger.Error("scan error", "error", err)
			continue
		}

		if !first {
			_, _ = w.Write([]byte(","))
		}
		_, _ = w.Write(jsonData)
		first = false

		// Flush each row so client receives data progressively
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	_, _ = w.Write([]byte("]}"))
	return nil
}

func (s *Service) ListStoredQueries(ctx context.Context, filterName string) ([]StoredQuery, error) {
	var query strings.Builder
	var args []any

	query.WriteString(`
		SELECT COALESCE(jsonb_agg(
			jsonb_build_object(
				'name', name,
				'type', 'AQL', 
				'version', version,
				'saved', to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MSTZH:TZM'),
				'q', query
			)), '[]'::jsonb) as queries
		FROM openehr.tbl_query
	`)

	if filterName != "" {
		namespaceSeperator := strings.LastIndex(filterName, "::")
		if namespaceSeperator != -1 {
			query.WriteString(` WHERE name = $1`)
		} else {
			query.WriteString(` WHERE name LIKE '%' || $1`)
		}
		args = append(args, filterName)
	}

	var queries []StoredQuery
	if err := s.DB.QueryRow(ctx, query.String(), args...).Scan(&queries); err != nil {
		return nil, fmt.Errorf("error querying stored AQL queries: %w", err)
	}

	return queries, nil
}

func (s *Service) GetQueryByName(ctx context.Context, name string, filterVersion string) (StoredQuery, error) {
	var query strings.Builder
	var args []any

	query.WriteString(`SELECT name, version, query, created_at FROM openehr.tbl_query WHERE name = $1 `)
	args = append(args, name)

	if filterVersion != "" {
		query.WriteString(`AND version = $2 `)
		args = append(args, filterVersion)
	}

	query.WriteString(`ORDER BY created_at DESC LIMIT 1`)

	var storedQuery StoredQuery
	if err := s.DB.QueryRow(ctx, query.String(), args...).Scan(&storedQuery.Name, &storedQuery.Version, &storedQuery.Query, &storedQuery.Saved); err != nil {
		if err == database.ErrNoRows {
			return StoredQuery{}, ErrQueryNotFound
		}
		return StoredQuery{}, fmt.Errorf("error retrieving AQL query by name: %w", err)
	}

	storedQuery.Type = "AQL"
	return storedQuery, nil
}

func (s *Service) StoreQuery(ctx context.Context, name, version, aqlQuery string) error {
	// Store the new query
	_, err := s.DB.Exec(ctx, `INSERT INTO openehr.tbl_query (id, name, version, query) VALUES ($1, $2, $3, $4) ON CONFLICT (name, version) DO UPDATE SET query = EXCLUDED.query`,
		uuid.New(),
		name,
		version,
		aqlQuery,
	)
	if err != nil {
		return fmt.Errorf("error storing AQL query: %w", err)
	}

	return nil
}

func NewVersionedEHRAccess(id, ehrID uuid.UUID) rm.VERSIONED_EHR_ACCESS {
	return rm.VERSIONED_EHR_ACCESS{
		UID: rm.HIER_OBJECT_ID{
			Value: id.String(),
		},
		OwnerID: rm.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      rm.EHR_TYPE,
			ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
				Value: ehrID.String(),
			}),
		},
		TimeCreated: rm.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func NewEHRAccess(id uuid.UUID) rm.EHR_ACCESS {
	return rm.EHR_ACCESS{
		UID: utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&rm.OBJECT_VERSION_ID{
			Type_: utils.Some(rm.OBJECT_VERSION_ID_TYPE),
			Value: fmt.Sprintf("%s::%s::1", id.String(), config.SYSTEM_ID_GOPENEHR),
		})),
		Name: rm.DV_TEXT_from_DV_TEXT(rm.DV_TEXT{
			Value: "EHR Access",
		}),
		ArchetypeNodeID: "openEHR-EHR-EHR_ACCESS.generic.v1",
	}
}

func NewVersionedEHRStatus(id, ehrID uuid.UUID) rm.VERSIONED_EHR_STATUS {
	return rm.VERSIONED_EHR_STATUS{
		UID: rm.HIER_OBJECT_ID{
			Value: id.String(),
		},
		OwnerID: rm.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      rm.EHR_TYPE,
			ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
				Value: ehrID.String(),
			}),
		},
		TimeCreated: rm.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func NewEHRStatus(id uuid.UUID) rm.EHR_STATUS {
	return rm.EHR_STATUS{
		// UID: utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&rm.OBJECT_VERSION_ID{
		// 	Type_: utils.Some(rm.OBJECT_VERSION_ID_TYPE),
		// 	Value: fmt.Sprintf("%s::%s::1", id.String(), config.SYSTEM_ID_GOPENEHR),
		// })),
		Name: rm.DV_TEXT_from_DV_TEXT(rm.DV_TEXT{
			Value: "EHR Status",
		}),
		ArchetypeNodeID: "openEHR-EHR-EHR_STATUS.generic.v1",
		Subject:         rm.PARTY_SELF{},
		IsQueryable:     true,
		IsModifiable:    true,
	}
}

func NewVersionedComposition(id, ehrID uuid.UUID) rm.VERSIONED_COMPOSITION {
	return rm.VERSIONED_COMPOSITION{
		UID: rm.HIER_OBJECT_ID{
			Value: id.String(),
		},
		OwnerID: rm.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      rm.EHR_TYPE,
			ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
				Value: ehrID.String(),
			}),
		},
		TimeCreated: rm.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func NewVersionedFolder(id, ehrID uuid.UUID) rm.VERSIONED_FOLDER {
	return rm.VERSIONED_FOLDER{
		UID: rm.HIER_OBJECT_ID{
			Value: id.String(),
		},
		OwnerID: rm.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      rm.EHR_TYPE,
			ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
				Value: ehrID.String(),
			}),
		},
		TimeCreated: rm.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func NewVersionedParty(uid uuid.UUID) rm.VERSIONED_PARTY {
	return rm.VERSIONED_PARTY{
		UID: rm.HIER_OBJECT_ID{
			Value: uid.String(),
		},
		OwnerID: rm.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      rm.ORGANISATION_TYPE,
			ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
				Value: uid.String(),
			}),
		},
		TimeCreated: rm.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func UpgradeObjectVersionID(current rm.UIDBasedIDUnion, previous rm.OBJECT_VERSION_ID) (rm.OBJECT_VERSION_ID, error) {
	switch current.Kind {
	case rm.UID_BASED_ID_kind_OBJECT_VERSION_ID:
		// Check version is incremented
		if current.OBJECT_VERSION_ID().VersionTreeID().CompareTo(previous.VersionTreeID()) <= 0 {
			return rm.OBJECT_VERSION_ID{}, ErrVersionLowerOrEqualToCurrent
		}
		return current.OBJECT_VERSION_ID(), nil
	case rm.UID_BASED_ID_kind_HIER_OBJECT_ID:
		currentUID := current.HIER_OBJECT_ID()

		// Add namespace and version to convert to OBJECT_VERSION_ID
		versionTreeID := previous.VersionTreeID()
		versionTreeID.Major++

		return rm.OBJECT_VERSION_ID{
			Type_: utils.Some(rm.OBJECT_VERSION_ID_TYPE),
			Value: fmt.Sprintf("%s::%s::%s", currentUID.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
		}, nil
	}

	return rm.OBJECT_VERSION_ID{}, fmt.Errorf("object UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %d", current.Kind)
}

func NewContribution(description string, auditChangeType terminology.AuditChangeType, versions []rm.OBJECT_REF) rm.CONTRIBUTION {
	contribution := rm.CONTRIBUTION{
		UID: rm.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: make([]rm.OBJECT_REF, 0),
		Audit: rm.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: rm.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			Description:   utils.Some(rm.DV_TEXT{Value: description}),
			Committer: rm.PARTY_PROXY_from_PARTY_SELF(rm.PARTY_SELF{
				Type_: utils.Some(rm.PARTY_SELF_TYPE),
			}),
		},
	}

	switch auditChangeType {
	case terminology.AUDIT_CHANGE_TYPE_CODE_CREATION:
		contribution.Audit.ChangeType = rm.DV_CODED_TEXT{
			Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_CREATION),
			DefiningCode: rm.CODE_PHRASE{
				CodeString: string(terminology.AUDIT_CHANGE_TYPE_CODE_CREATION),
				TerminologyID: rm.TERMINOLOGY_ID{
					Value: string(terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR),
				},
			},
		}
	case terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION:
		contribution.Audit.ChangeType = rm.DV_CODED_TEXT{
			Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION),
			DefiningCode: rm.CODE_PHRASE{
				CodeString: string(terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION),
				TerminologyID: rm.TERMINOLOGY_ID{
					Value: string(terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR),
				},
			},
		}
	case terminology.AUDIT_CHANGE_TYPE_CODE_DELETED:
		contribution.Audit.ChangeType = rm.DV_CODED_TEXT{
			Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_DELETED),
			DefiningCode: rm.CODE_PHRASE{
				CodeString: string(terminology.AUDIT_CHANGE_TYPE_CODE_DELETED),
				TerminologyID: rm.TERMINOLOGY_ID{
					Value: string(terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR),
				},
			},
		}
	default:
		panic("unsupported audit change type")
	}

	return contribution
}

func NewOriginalVersion(id rm.OBJECT_VERSION_ID, data rm.OriginalVersionDataUnion, precedingVersion utils.Optional[rm.OBJECT_VERSION_ID]) rm.ORIGINAL_VERSION {
	return rm.ORIGINAL_VERSION{
		UID:                 id,
		PrecedingVersionUID: precedingVersion,
		LifecycleState: rm.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: rm.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: rm.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: data,
	}
}
