package service

import (
	"context"
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
)

type EHRService struct {
	Logger *slog.Logger
	DB     *database.Database
}

func NewEHRService(logger *slog.Logger, db *database.Database) EHRService {
	return EHRService{
		Logger: logger,
		DB:     db,
	}
}

func (s *EHRService) CreateEHR(ctx context.Context, ehrID uuid.UUID, ehrStatus openehr.EHR_STATUS) (openehr.EHR, error) {
	// Validate EHR Status
	if err := s.ValidateEHRStatus(ctx, ehrStatus); err != nil {
		return openehr.EHR{}, err
	}

	// Provide ID when EHR Status does not have one
	if !ehrStatus.UID.E {
		ehrStatus.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", uuid.NewString(), config.SYSTEM_ID_GOPENEHR),
			},
		})
	}

	switch ehrStatus.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Add namespace and version to convert to OBJECT_VERSION_ID
		hierID := ehrStatus.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		ehrStatus.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", hierID.Value, config.SYSTEM_ID_GOPENEHR),
			},
		})
	default:
		return openehr.EHR{}, fmt.Errorf("EHR Status UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", ehrStatus.UID.V.Value)
	}

	// Prepare Versioned EHR Status
	versionedEhrStatus := openehr.VERSIONED_EHR_STATUS{
		UID: openehr.HIER_OBJECT_ID{
			Value: ehrStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.EHR_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: ehrID.String(),
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Prepare EHR Status Version
	ehrStatusVersion := openehr.ORIGINAL_VERSION{
		UID: *ehrStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &ehrStatus,
	}

	// Build Versioned EHR Access
	versionedEhrAccess := openehr.VERSIONED_EHR_ACCESS{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.EHR_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: ehrID.String(),
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Build EHR Access
	ehrAccess := openehr.EHR_ACCESS{
		UID: util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", versionedEhrAccess.UID.Value, config.SYSTEM_ID_GOPENEHR),
			},
		}),
		Name: openehr.X_DV_TEXT{
			Value: &openehr.DV_TEXT{
				Value: "EHR Access",
			},
		},
		ArchetypeNodeID: "openEHR-EHR-EHR_ACCESS.generic.v1",
	}

	// Build EHR Access Version
	ehrAccessVersion := openehr.ORIGINAL_VERSION{
		UID: *ehrAccess.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &ehrAccess,
	}

	// Prepare Contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.EHR_STATUS_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &ehrStatusVersion.UID,
				},
			},
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.EHR_ACCESS_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &ehrAccessVersion.UID,
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
				Value: "EHR created",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				},
			},
		},
	}

	// Create EHR
	ehr := openehr.EHR{
		SystemID: openehr.HIER_OBJECT_ID{
			Value: config.SYSTEM_ID_GOPENEHR,
		},
		EHRID: openehr.HIER_OBJECT_ID{
			Value: ehrID.String(),
		},
		EHRStatus: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.EHR_STATUS_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &ehrStatusVersion.UID,
			},
		},
		EHRAccess: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.EHR_ACCESS_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &ehrAccessVersion.UID,
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Insert EHR
	query := `INSERT INTO openehr.tbl_ehr (id) VALUES ($1)`
	args := []any{ehrID}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr into the database: %w", err)
	}

	ehr.SetModelName()
	query = `INSERT INTO openehr.tbl_ehr_data (id, data) VALUES ($1, $2)`
	args = []any{ehrID, ehr}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr data into the database: %w", err)
	}

	// Insert Contribution
	query = `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, ehrID}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	contribution.SetModelName()
	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert contribution data into the database: %w", err)
	}

	// Insert Versioned EHR Status
	query = `INSERT INTO openehr.tbl_versioned_object (id, type, object_type,  ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args = []any{versionedEhrStatus.UID.Value, openehr.VERSIONED_EHR_STATUS_MODEL_NAME, openehr.EHR_STATUS_MODEL_NAME, ehrID, contribution.UID.Value}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr status into the database: %w", err)
	}

	versionedEhrStatus.SetModelName()
	query = `INSERT INTO openehr.tbl_versioned_object_data (id, data) VALUES ($1, $2)`
	args = []any{versionedEhrStatus.UID.Value, versionedEhrStatus}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr status data into the database: %w", err)
	}

	// Insert EHR Status Version
	query = `INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args = []any{ehrStatusVersion.UID.Value, versionedEhrStatus.UID.Value, openehr.EHR_STATUS_MODEL_NAME, ehrID, contribution.UID.Value}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr status into the database: %w", err)
	}

	ehrStatusVersion.SetModelName()
	query = `INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)`
	args = []any{ehrStatusVersion.UID.Value, ehrStatusVersion}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr status version data into the database: %w", err)
	}

	// Insert Versioned EHR Access
	query = `INSERT INTO openehr.tbl_versioned_object (id, type, object_type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args = []any{versionedEhrAccess.UID.Value, openehr.VERSIONED_EHR_ACCESS_MODEL_NAME, openehr.EHR_ACCESS_MODEL_NAME, ehrID, contribution.UID.Value}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr access into the database: %w", err)
	}

	versionedEhrAccess.SetModelName()
	query = `INSERT INTO openehr.tbl_versioned_object_data (id, data) VALUES ($1, $2)`
	args = []any{versionedEhrAccess.UID.Value, versionedEhrAccess}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr access data into the database: %w", err)
	}

	// Insert EHR Access Version
	query = `INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args = []any{ehrAccessVersion.UID.Value, versionedEhrAccess.UID.Value, openehr.EHR_ACCESS_MODEL_NAME, ehrID, contribution.UID.Value}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr access into the database: %w", err)
	}

	ehrAccessVersion.SetModelName()
	query = `INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)`
	args = []any{ehrAccessVersion.UID.Value, ehrAccessVersion}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr access version data into the database: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ehr, nil
}

func (s *EHRService) GetEHR(ctx context.Context, id string) (openehr.EHR, error) {
	query := `
		SELECT ed.data
		FROM openehr.tbl_ehr e
		JOIN openehr.tbl_ehr_data ed ON e.id = ed.id
		WHERE e.id = $1
	`
	args := []any{id}

	row := s.DB.QueryRow(ctx, query, args...)

	var ehr openehr.EHR
	err := row.Scan(&ehr)
	if err != nil {
		if err == database.ErrNoRows {
			return openehr.EHR{}, ErrEHRNotFound
		}
		return openehr.EHR{}, fmt.Errorf("failed to fetch EHR from database: %w", err)
	}

	return ehr, nil
}

func (s *EHRService) GetEHRBySubject(ctx context.Context, subjectID, subjectNamespace string) (openehr.EHR, error) {
	query := `
        SELECT ed.data 
        FROM openehr.tbl_ehr e 
        JOIN openehr.tbl_ehr_data ed ON e.id = ed.id
        JOIN openehr.tbl_object_version ov ON ov.ehr_id = e.id
        JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id
        WHERE ov.type = $1
          AND jsonb_path_exists(
              ovd.object_data, 
              '$.subject.external_ref ? (@.namespace == $namespace && @.id.value == $subject)',
              jsonb_build_object('namespace', $2::text, 'subject', $3::text)
          )
        ORDER BY ov.created_at DESC
        LIMIT 1
    `
	args := []any{openehr.EHR_STATUS_MODEL_NAME, subjectNamespace, subjectID}

	row := s.DB.QueryRow(ctx, query, args...)

	var ehr openehr.EHR
	err := row.Scan(&ehr)
	if err != nil {
		if err == database.ErrNoRows {
			return openehr.EHR{}, ErrEHRNotFound
		}
		return openehr.EHR{}, fmt.Errorf("failed to fetch EHR by subject from database: %w", err)
	}

	return ehr, nil
}

func (s *EHRService) DeleteEHR(ctx context.Context, id string) error {
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

func (s *EHRService) DeleteEHRBulk(ctx context.Context, ids []string) error {
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

func (s *EHRService) NewEHRStatus() openehr.EHR_STATUS {
	return openehr.EHR_STATUS{
		Name: openehr.X_DV_TEXT{
			Value: &openehr.DV_TEXT{
				Value: "EHR Status",
			},
		},
		ArchetypeNodeID: "openEHR-EHR-EHR_STATUS.generic.v1",
		Subject:         openehr.PARTY_SELF{},
		IsQueryable:     true,
		IsModifiable:    true,
	}
}

func (s *EHRService) ValidateEHRStatus(ctx context.Context, ehrStatus openehr.EHR_STATUS) error {
	validateErr := ehrStatus.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional EHR Status validation
	if ehrStatus.Subject.ExternalRef.E {
		externalRef := ehrStatus.Subject.ExternalRef.V
		if externalRef.Namespace == "local" {

			attrPath := "$.subject.external_ref"

			switch v := externalRef.ID.Value.(type) {
			case *openehr.HIER_OBJECT_ID:
				// Must be a valid type
				if externalRef.Type != openehr.VERSIONED_PARTY_MODEL_NAME {
					validateErr.Errs = append(validateErr.Errs, util.ValidationError{
						Model:          openehr.EHR_STATUS_MODEL_NAME,
						Path:           attrPath + ".type",
						Message:        fmt.Sprintf("invalid subject external_ref type: %s", externalRef.Type),
						Recommendation: "Ensure external ref type is VERSIONED_PARTY",
					})
				}

				// Must be a valid UUID
				if err := uuid.Validate(v.Value); err != nil {
					validateErr.Errs = append(validateErr.Errs, util.ValidationError{
						Model:          openehr.EHR_STATUS_MODEL_NAME,
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
						validateErr.Errs = append(validateErr.Errs, util.ValidationError{
							Model:          openehr.EHR_STATUS_MODEL_NAME,
							Path:           attrPath + ".id.value",
							Message:        "Subject external ref id " + v.Value + " with type " + externalRef.Type + " does not exist in tbl_versioned_party",
							Recommendation: "Ensure external ref id exists in tbl_versioned_party",
						})
					} else {
						return err
					}
				}
			default:
				validateErr.Errs = append(validateErr.Errs, util.ValidationError{
					Model:          openehr.EHR_STATUS_MODEL_NAME,
					Path:           attrPath + ".id",
					Message:        fmt.Sprintf("Unsupported subject external_ref id type: %s", v.GetModelName()),
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

func (s *EHRService) GetEHRStatus(ctx context.Context, ehrID string, filterOnTime time.Time, filterOnVersionID string) (openehr.EHR_STATUS, error) {
	// Build query
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id
		WHERE ov.type = $1
		  AND ov.ehr_id = $2 
	`)
	args = []any{openehr.EHR_STATUS_MODEL_NAME, ehrID}
	argNum += 2

	if !filterOnTime.IsZero() {
		query.WriteString(fmt.Sprintf(`AND ov.created_at <= $%d `, argNum))
		args = append(args, filterOnTime)
		argNum++
	}

	if filterOnVersionID != "" {
		query.WriteString(fmt.Sprintf(`AND ov.id = $%d `, argNum))
		args = append(args, filterOnVersionID)
	}

	query.WriteString(`ORDER BY ov.created_at DESC LIMIT 1`)

	row := s.DB.QueryRow(ctx, query.String(), args...)

	var ehrStatus openehr.EHR_STATUS
	err := row.Scan(&ehrStatus)
	if err != nil {
		if err == database.ErrNoRows {
			return openehr.EHR_STATUS{}, ErrEHRStatusNotFound
		}
		return openehr.EHR_STATUS{}, fmt.Errorf("failed to fetch EHR Status from database: %w", err)
	}

	return ehrStatus, nil
}

func (s *EHRService) UpdateEHRStatus(ctx context.Context, ehrID string, ehrStatus openehr.EHR_STATUS) (openehr.EHR_STATUS, error) {
	// Validate request EHR Status
	if err := s.ValidateEHRStatus(ctx, ehrStatus); err != nil {
		return openehr.EHR_STATUS{}, err
	}

	// Get current EHR Status
	currentEHRStatus, err := s.GetEHRStatus(ctx, ehrID, time.Time{}, "")
	if err != nil {
		return openehr.EHR_STATUS{}, fmt.Errorf("failed to get current EHR Status: %w", err)
	}
	currentEHRStatusID := currentEHRStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID)

	// Provide ID when EHR Status does not have one
	if !ehrStatus.UID.E {
		ehrStatus.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: currentEHRStatusID.UID(),
			},
		})
	}

	// Handle EHR Status UID types
	switch ehrStatus.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check version is incremented
		currentVersionID := currentEHRStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		newVersionID := ehrStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		if newVersionID.VersionTreeID().CompareTo(currentVersionID.VersionTreeID()) <= 0 {
			return openehr.EHR_STATUS{}, ErrEHRStatusVersionLowerOrEqualToCurrent
		}
	case *openehr.HIER_OBJECT_ID:
		// Check HIER_OBJECT_ID matches current EHR Status UID
		hierID := ehrStatus.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		if currentEHRStatusID.UID() != hierID.Value {
			return openehr.EHR_STATUS{}, ErrInvalidEHRStatusUIDMismatch
		}

		// Increment version
		versionTreeID := currentEHRStatusID.VersionTreeID()
		versionTreeID.Major++

		// Convert to OBJECT_VERSION_ID
		ehrStatus.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", hierID.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.EHR_STATUS{}, fmt.Errorf("EHR Status UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", ehrStatus.UID.V.Value)
	}

	// Prepare Original Version of EHR Status
	ehrStatusVersion := openehr.ORIGINAL_VERSION{
		UID: *ehrStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		PrecedingVersionUID: util.Some(*currentEHRStatusID),
		Data:                &ehrStatus,
	}

	// Prepare contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.EHR_STATUS_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &ehrStatusVersion.UID,
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
				Value: "EHR Status updated",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				},
			},
		},
	}

	// Start transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.EHR_STATUS{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Insert Contribution
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
	args := []any{contribution.UID.Value, ehrID}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR_STATUS{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	contribution.SetModelName()
	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR_STATUS{}, fmt.Errorf("failed to insert contribution data into the database: %w", err)
	}

	// Insert EHR Status Version
	query = `INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args = []any{ehrStatusVersion.UID.Value, ehrStatusVersion.UID.UID(), openehr.EHR_STATUS_MODEL_NAME, ehrID, contribution.UID.Value}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR_STATUS{}, fmt.Errorf("failed to insert ehr status into the database: %w", err)
	}

	ehrStatusVersion.SetModelName()
	query = `INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)`
	args = []any{ehrStatusVersion.UID.Value, ehrStatusVersion}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR_STATUS{}, fmt.Errorf("failed to insert ehr status version data into the database: %w", err)
	}

	// Update EHR Status reference in EHR
	query = `UPDATE openehr.tbl_ehr_data SET data = jsonb_set(data, '{ehr_status,id,value}', to_jsonb($1::text)) WHERE id = $2`
	args = []any{ehrStatusVersion.UID.Value, ehrID}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR_STATUS{}, fmt.Errorf("failed to update ehr status reference in ehr: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return openehr.EHR_STATUS{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ehrStatus, nil
}

func (s *EHRService) GetVersionedEHRStatus(ctx context.Context, ehrID string) (openehr.VERSIONED_EHR_STATUS, error) {
	query := `
		SELECT vod.data 
		FROM openehr.tbl_versioned_object vo
		JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id
		WHERE vo.type = $1 
		  AND vo.ehr_id = $2
		LIMIT 1
	`
	args := []any{openehr.VERSIONED_EHR_STATUS_MODEL_NAME, ehrID}

	row := s.DB.QueryRow(ctx, query, args...)

	var versionedEHRStatus openehr.VERSIONED_EHR_STATUS
	err := row.Scan(&versionedEHRStatus)
	if err != nil {
		if err == database.ErrNoRows {
			return openehr.VERSIONED_EHR_STATUS{}, ErrEHRNotFound
		}
		return openehr.VERSIONED_EHR_STATUS{}, fmt.Errorf("failed to fetch Versioned EHR Status from database: %w", err)
	}

	return versionedEHRStatus, nil
}

func (s *EHRService) GetVersionedEHRStatusRevisionHistory(ctx context.Context, ehrID string) (openehr.REVISION_HISTORY, error) {
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
			WHERE c.ehr_id = $1
				AND version->>'type' = 'EHR_STATUS'
			GROUP BY version->'id'->>'value'
		) grouped
    `
	args := []any{ehrID}

	row := s.DB.QueryRow(ctx, query, args...)

	var revisionHistory openehr.REVISION_HISTORY
	err := row.Scan(&revisionHistory)
	if err != nil {
		if err == database.ErrNoRows {
			return openehr.REVISION_HISTORY{}, ErrEHRNotFound
		}
		return openehr.REVISION_HISTORY{}, fmt.Errorf("failed to fetch Revision History from database: %w", err)
	}

	return revisionHistory, nil
}

func (s *EHRService) GetVersionedEHRStatusVersionAsJSON(ctx context.Context, ehrID uuid.UUID, filterAtTime time.Time, filterVersionID string) ([]byte, error) {
	// Fetch EHR Status at given time
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id
		WHERE ov.type = $1
		  AND ov.ehr_id = $2 
	`)
	args = []any{openehr.EHR_STATUS_MODEL_NAME, ehrID}
	argNum += 2

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

	var rawEhrStatusJSON []byte
	err := row.Scan(&rawEhrStatusJSON)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR Status at time from database: %w", err)
	}

	return rawEhrStatusJSON, nil
}

func (s *EHRService) ValidateComposition(ctx context.Context, composition openehr.COMPOSITION) error {
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

func (s *EHRService) ExistsComposition(ctx context.Context, versionID string, filterOnEHRID uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.type = $1 AND ov.id = $2 `
	args := []any{openehr.COMPOSITION_MODEL_NAME, versionID}

	if filterOnEHRID != uuid.Nil {
		query += `AND ov.ehr_id = $3 `
		args = append(args, filterOnEHRID)
	}

	query += `LIMIT 1`

	var exists int
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		if err == database.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if Composition exists in database: %w", err)
	}

	return true, nil
}

func (s *EHRService) CreateComposition(ctx context.Context, ehrID uuid.UUID, composition openehr.COMPOSITION) (openehr.COMPOSITION, error) {
	// Validate composition
	if err := s.ValidateComposition(ctx, composition); err != nil {
		return openehr.COMPOSITION{}, err
	}

	// Provide UID if not set
	if !composition.UID.E {
		composition.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::gopenehr::1", uuid.NewString()),
			},
		})
	}

	// Handle UID types
	switch composition.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// OK
	case *openehr.HIER_OBJECT_ID:
		// Convert to OBJECT_VERSION_ID
		hierID := composition.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		composition.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", hierID.Value, config.SYSTEM_ID_GOPENEHR),
			},
		})
	default:
		return openehr.COMPOSITION{}, fmt.Errorf("composition UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", composition.UID.V.Value)
	}

	// Extract UID type
	compositionID, ok := composition.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return openehr.COMPOSITION{}, fmt.Errorf("composition UID must be of type OBJECT_VERSION_ID, got %T", composition.UID.V.Value)
	}

	// Check if Composition already exists within the EHR
	exists, err := s.ExistsComposition(ctx, compositionID.Value, uuid.Nil)
	if err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to check if composition exists: %w", err)
	}
	if exists {
		return openehr.COMPOSITION{}, ErrCompositionAlreadyExists
	}

	// Prepare Versioned Composition
	versionedComposition := openehr.VERSIONED_COMPOSITION{
		UID: openehr.HIER_OBJECT_ID{
			Value: strings.Split(composition.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, "::")[0],
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.EHR_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: ehrID.String(),
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Prepare Original Version of Composition
	compositionVersion := openehr.ORIGINAL_VERSION{
		UID: openehr.OBJECT_VERSION_ID{
			Value: composition.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value,
		},
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &composition,
	}

	// Prepare Contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.COMPOSITION_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &compositionVersion.UID,
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
				Value: "Composition created",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Insert Contribution
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
	args := []any{contribution.UID.Value, ehrID}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	contribution.SetModelName()
	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert contribution data into the database: %w", err)
	}

	// Insert Versioned Composition
	query = `INSERT INTO openehr.tbl_versioned_object (id, type, object_type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args = []any{versionedComposition.UID.Value, openehr.VERSIONED_COMPOSITION_MODEL_NAME, openehr.COMPOSITION_MODEL_NAME, ehrID, contribution.UID.Value}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert versioned composition into the database: %w", err)
	}

	versionedComposition.SetModelName()
	query = `INSERT INTO openehr.tbl_versioned_object_data (id, data) VALUES ($1, $2)`
	args = []any{versionedComposition.UID.Value, versionedComposition}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert versioned composition data into the database: %w", err)
	}

	// Insert Composition Version
	query = `INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args = []any{compositionVersion.UID.Value, versionedComposition.UID.Value, openehr.COMPOSITION_MODEL_NAME, ehrID, contribution.UID.Value}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert composition into the database: %w", err)
	}

	compositionVersion.SetModelName()
	query = `INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)`
	args = []any{compositionVersion.UID.Value, compositionVersion}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert composition data into the database: %w", err)
	}

	// Update EHR
	query = `UPDATE openehr.tbl_ehr_data
			SET data = jsonb_set(data, '{compositions}', 
				COALESCE(data->'compositions', '[]'::jsonb) || to_jsonb($1::jsonb)
			) 
			WHERE id = $2`
	compositionRef := openehr.OBJECT_REF{
		Namespace: config.NAMESPACE_LOCAL,
		Type:      openehr.COMPOSITION_MODEL_NAME,
		ID: openehr.X_OBJECT_ID{
			Value: &compositionVersion.UID,
		},
	}
	args = []any{compositionRef, ehrID}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to update EHR with new composition reference: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return composition, nil
}

func (s *EHRService) GetComposition(ctx context.Context, ehrID uuid.UUID, uidBasedID string) (openehr.COMPOSITION, error) {
	query := `
		SELECT ovd.object_data 
		FROM openehr.tbl_object_version ov 
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.type = $1 AND ov.ehr_id = $2`
	args := []any{openehr.COMPOSITION_MODEL_NAME, ehrID}

	if strings.Count(uidBasedID, "::") == 2 {
		// UID is of type OBJECT_VERSION_ID
		query += `AND ov.id = $3 `
	} else {
		// UID is of type HIER_OBJECT_ID
		query += `AND ov.versioned_object_id = $3 `
	}
	args = append(args, uidBasedID)

	query += `ORDER BY ov.created_at DESC LIMIT 1`

	var composition openehr.COMPOSITION
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&composition); err != nil {
		if err == database.ErrNoRows {
			return openehr.COMPOSITION{}, ErrCompositionNotFound
		}
		return openehr.COMPOSITION{}, fmt.Errorf("failed to fetch Composition by ID from database: %w", err)
	}

	return composition, nil
}

func (s *EHRService) UpdateComposition(ctx context.Context, ehrID uuid.UUID, composition openehr.COMPOSITION) (openehr.COMPOSITION, error) {
	// Provide UID when Composition does not have one
	if !composition.UID.E {
		return openehr.COMPOSITION{}, ErrCompositionUIDNotProvided
	}

	// Get current Composition
	var currentCompositionID *openehr.OBJECT_VERSION_ID

	// Handle Composition UID types
	switch v := composition.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		currentComposition, err := s.GetComposition(ctx, ehrID, v.Value)
		if err != nil {
			if err == ErrCompositionNotFound {
				return openehr.COMPOSITION{}, ErrCompositionNotFound
			}
			return openehr.COMPOSITION{}, fmt.Errorf("failed to get current Composition: %w", err)
		}
		currentCompositionID = currentComposition.UID.V.Value.(*openehr.OBJECT_VERSION_ID)

		// Check version is incremented
		newVersionID := composition.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		if newVersionID.VersionTreeID().CompareTo(currentCompositionID.VersionTreeID()) <= 0 {
			return openehr.COMPOSITION{}, ErrCompositionVersionLowerOrEqualToCurrent
		}
	case *openehr.HIER_OBJECT_ID:
		currentComposition, err := s.GetComposition(ctx, ehrID, v.Value)
		if err != nil {
			if err == ErrCompositionNotFound {
				return openehr.COMPOSITION{}, ErrCompositionNotFound
			}
			return openehr.COMPOSITION{}, fmt.Errorf("failed to get current Composition: %w", err)
		}
		currentCompositionID = currentComposition.UID.V.Value.(*openehr.OBJECT_VERSION_ID)

		// Check HIER_OBJECT_ID matches current EHR Status UID
		if currentCompositionID.UID() != v.Value {
			return openehr.COMPOSITION{}, ErrInvalidCompositionUIDMismatch
		}

		// Increment version
		versionTreeID := currentCompositionID.VersionTreeID()
		versionTreeID.Major++

		// Convert to OBJECT_VERSION_ID
		composition.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", v.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.COMPOSITION{}, fmt.Errorf("composition UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", composition.UID.V.Value)
	}

	// Prepare Original Version of Composition
	compositionVersion := openehr.ORIGINAL_VERSION{
		UID: *composition.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		PrecedingVersionUID: util.Some(*currentCompositionID),
		Data:                &composition,
	}

	// Prepare contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.COMPOSITION_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &compositionVersion.UID,
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
				Value: "Composition updated",
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
		return openehr.COMPOSITION{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	// Insert Contribution
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
	args := []any{contribution.UID.Value, ehrID}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	contribution.SetModelName()
	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert contribution data into the database: %w", err)
	}

	// Insert Composition
	query = `INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args = []any{compositionVersion.UID.Value, compositionVersion.UID.UID(), openehr.COMPOSITION_MODEL_NAME, ehrID, contribution.UID.Value}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert updated composition into the database: %w", err)
	}

	compositionVersion.SetModelName()
	query = `INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)`
	args = []any{compositionVersion.UID.Value, compositionVersion}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert updated composition data into the database: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return composition, nil
}

func (s *EHRService) DeleteComposition(ctx context.Context, ehrID uuid.UUID, objectVersionID string) error {
	versionedObjectID := strings.Split(objectVersionID, "::")[0]

	// Create contribution for deletion
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.COMPOSITION_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &openehr.OBJECT_VERSION_ID{
						Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
						Value: objectVersionID,
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
				Value: "Composition deleted",
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

	// Insert Contribution
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
	args := []any{contribution.UID.Value, ehrID}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	contribution.SetModelName()
	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert contribution data into the database: %w", err)
	}

	// Delete Composition
	var deleted uint8
	query = `DELETE FROM openehr.tbl_versioned_object WHERE id = $1 AND type = $2 RETURNING 1`
	args = []any{versionedObjectID, openehr.VERSIONED_COMPOSITION_MODEL_NAME}

	err = tx.QueryRow(ctx, query, args...).Scan(&deleted)
	if err != nil {
		if err == database.ErrNoRows {
			return ErrCompositionNotFound
		}
		return fmt.Errorf("failed to delete Composition from database: %w", err)
	}

	// Remove reference from EHR
	query = `UPDATE openehr.tbl_ehr_data
		SET data = jsonb_set(data, '{compositions}', 
			COALESCE(
				(
					SELECT jsonb_agg(comp) 
					FROM jsonb_array_elements(data->'compositions') AS comp
					WHERE comp->'id'->>'value' != $1
				), '[]'::jsonb
			)
		) 
		WHERE id = $2`
	args = []any{versionedObjectID, ehrID}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to update EHR to remove composition reference: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *EHRService) GetVersionedComposition(ctx context.Context, ehrID, versionedObjectID string) (openehr.VERSIONED_COMPOSITION, error) {
	query := `
		SELECT vod.data 
		FROM openehr.tbl_versioned_object vo
		JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id
		WHERE vo.type = $1
		  AND vo.ehr_id = $2 
		  AND vo.id = $3
		LIMIT 1
	`
	args := []any{openehr.VERSIONED_COMPOSITION_MODEL_NAME, ehrID, versionedObjectID}
	row := s.DB.QueryRow(ctx, query, args...)

	var versionedComposition openehr.VERSIONED_COMPOSITION
	if err := row.Scan(&versionedComposition); err != nil {
		if err == database.ErrNoRows {
			return openehr.VERSIONED_COMPOSITION{}, ErrCompositionNotFound
		}
		return openehr.VERSIONED_COMPOSITION{}, fmt.Errorf("failed to fetch Versioned Composition by ID from database: %w", err)
	}

	return versionedComposition, nil
}

func (s *EHRService) GetVersionedCompositionRevisionHistory(ctx context.Context, ehrID, versionedObjectID string) (openehr.REVISION_HISTORY, error) {
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
			WHERE c.ehr_id = $1
				AND version->>'type' = 'COMPOSITION'
				AND version->'id'->>'value' LIKE $2 || '%'
			GROUP BY version->'id'->>'value'
		) grouped
	`
	args := []any{ehrID, versionedObjectID}
	row := s.DB.QueryRow(ctx, query, args...)

	var revisionHistory openehr.REVISION_HISTORY
	if err := row.Scan(&revisionHistory); err != nil {
		if err == database.ErrNoRows {
			return openehr.REVISION_HISTORY{}, ErrCompositionNotFound
		}
		return openehr.REVISION_HISTORY{}, fmt.Errorf("failed to fetch Revision History from database: %w", err)
	}

	return revisionHistory, nil
}

func (s *EHRService) GetVersionedCompositionVersionJSON(ctx context.Context, ehrID, versionedObjectID string, filterAtTime time.Time, filterVersionID string) ([]byte, error) {
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`
		SELECT ovd.data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id
		WHERE ov.type = $1
		  AND ov.ehr_id = $2
		  AND ov.versioned_object_id = $3
	`)
	args = []any{openehr.COMPOSITION_MODEL_NAME, ehrID, versionedObjectID}
	argNum += 3

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

	var compositionVersionJSON []byte
	if err := row.Scan(&compositionVersionJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrCompositionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Composition version at time from database: %w", err)
	}

	return compositionVersionJSON, nil
}

func (s *EHRService) ValidateDirectory(ctx context.Context, ehrID uuid.UUID, directory openehr.FOLDER) error {
	validateErr := directory.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Directory validation can be added here
	folderQueue := make([]openehr.FOLDER, 0)
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
			case openehr.COMPOSITION_MODEL_NAME:
				id, ok := currentRef.ID.Value.(*openehr.OBJECT_VERSION_ID)
				if !ok {
					validateErr.Errs = append(validateErr.Errs, util.ValidationError{
						Model:          openehr.COMPOSITION_MODEL_NAME,
						Path:           itemPath,
						Message:        "Mismatch between type and id provided",
						Recommendation: "Ensure the ID is of type OBJECT_VERSION_ID",
					})
				}

				exists, err := s.ExistsComposition(ctx, id.Value, ehrID)
				if err != nil {
					return fmt.Errorf("failed to validate existence of Composition: %w", err)
				}
				if !exists {
					validateErr.Errs = append(validateErr.Errs, util.ValidationError{
						Model:          openehr.COMPOSITION_MODEL_NAME,
						Path:           itemPath,
						Message:        "COMPOSITION does not exist for this EHR in the system",
						Recommendation: "Ensure the composition exists for this EHR",
					})
				}
			case openehr.VERSIONED_COMPOSITION_MODEL_NAME:
				id, ok := currentRef.ID.Value.(*openehr.HIER_OBJECT_ID)
				if !ok {
					validateErr.Errs = append(validateErr.Errs, util.ValidationError{
						Model:          openehr.VERSIONED_COMPOSITION_MODEL_NAME,
						Path:           itemPath,
						Message:        "Mismatch between type and id provided",
						Recommendation: "Ensure the ID is of type HIER_OBJECT_ID",
					})
				}

				exists, err := s.ExistsComposition(ctx, id.Value, ehrID)
				if err != nil {
					return fmt.Errorf("failed to validate existence of Composition: %w", err)
				}
				if !exists {
					validateErr.Errs = append(validateErr.Errs, util.ValidationError{
						Model:          openehr.COMPOSITION_MODEL_NAME,
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

func (s *EHRService) CreateDirectory(ctx context.Context, ehrID uuid.UUID, directory openehr.FOLDER) (openehr.FOLDER, error) {
	// Validate directory
	if err := s.ValidateDirectory(ctx, ehrID, directory); err != nil {
		return openehr.FOLDER{}, err
	}

	// Provide ID when Directory does not have one
	if !directory.UID.E {
		directory.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: uuid.NewString(),
			},
		})
	}

	switch directory.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Add namespace and version to convert to OBJECT_VERSION_ID
		hierID := directory.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		directory.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", hierID.Value, config.SYSTEM_ID_GOPENEHR),
			},
		})
	default:
		return openehr.FOLDER{}, fmt.Errorf("directory UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", directory.UID.V.Value)
	}

	// Prepare Versioned Directory
	versionedDirectory := openehr.VERSIONED_FOLDER{
		UID: openehr.HIER_OBJECT_ID{
			Value: directory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(),
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.EHR_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: ehrID.String(),
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Prepare Original Version of Directory
	directoryVersion := openehr.ORIGINAL_VERSION{
		UID: *directory.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &directory,
	}

	// Prepare Contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.FOLDER_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &directoryVersion.UID,
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
				Value: "Directory created",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Insert Contribution
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
	args := []any{contribution.UID.Value, ehrID}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	contribution.SetModelName()
	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert contribution data into the database: %w", err)
	}

	// Insert Versioned Directory
	query = `INSERT INTO openehr.tbl_versioned_object (id, type, object_type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args = []any{versionedDirectory.UID.Value, openehr.VERSIONED_FOLDER_MODEL_NAME, openehr.FOLDER_MODEL_NAME, ehrID, contribution.UID.Value}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert versioned directory into the database: %w", err)
	}

	versionedDirectory.SetModelName()
	query = `INSERT INTO openehr.tbl_versioned_object_data (id, data) VALUES ($1, $2)`
	args = []any{versionedDirectory.UID.Value, versionedDirectory}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert versioned directory data into the database: %w", err)
	}

	// Insert Directory Version
	query = `INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args = []any{directoryVersion.UID.Value, versionedDirectory.UID.Value, openehr.FOLDER_MODEL_NAME, ehrID, contribution.UID.Value}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert directory into the database: %w", err)
	}

	directoryVersion.SetModelName()
	query = `INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)`
	args = []any{directoryVersion.UID.Value, directoryVersion}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert directory data into the database: %w", err)
	}

	// Update EHR
	directoryRef := openehr.OBJECT_REF{
		Namespace: config.NAMESPACE_LOCAL,
		Type:      openehr.FOLDER_MODEL_NAME,
		ID: openehr.X_OBJECT_ID{
			Value: &directoryVersion.UID,
		},
	}

	directoryRef.SetModelName()
	query = `UPDATE openehr.tbl_ehr_data
            SET data = jsonb_set(data, '{directory}', to_jsonb($1::jsonb)) 
            WHERE id = $2`
	args = []any{directoryRef, ehrID}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to update EHR with new directory reference: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return directory, nil
}

func (s *EHRService) GetDirectory(ctx context.Context, ehrID uuid.UUID) (openehr.FOLDER, error) {
	// Fetch Directory by EHR ID
	query := `
		SELECT ovd.object_data
        FROM openehr.tbl_object_version ov
        JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id
        WHERE ov.type = $1
          AND ov.ehr_id = $2
        ORDER BY ov.created_at DESC
        LIMIT 1
	`
	args := []any{openehr.FOLDER_MODEL_NAME, ehrID}
	row := s.DB.QueryRow(ctx, query, args...)

	var directory openehr.FOLDER
	if err := row.Scan(&directory); err != nil {
		if err == database.ErrNoRows {
			return openehr.FOLDER{}, ErrDirectoryNotFound
		}
		return openehr.FOLDER{}, fmt.Errorf("failed to fetch Directory by EHR ID from database: %w", err)
	}

	return directory, nil
}

func (s *EHRService) UpdateDirectory(ctx context.Context, ehrID uuid.UUID, directory openehr.FOLDER) (openehr.FOLDER, error) {
	// Validate Directory
	if err := s.ValidateDirectory(ctx, ehrID, directory); err != nil {
		return openehr.FOLDER{}, err
	}

	// Get current Directory
	currentDirectory, err := s.GetDirectory(ctx, ehrID)
	if err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to get current Directory: %w", err)
	}
	currentDirectoryID := currentDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID)

	// Provide ID when Directory does not have one
	if !directory.UID.E {
		directory.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: currentDirectoryID.UID(),
			},
		})
	}

	// Handle Directory UID types
	switch v := directory.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check version is incremented
		if v.VersionTreeID().CompareTo(currentDirectoryID.VersionTreeID()) <= 0 {
			return openehr.FOLDER{}, ErrDirectoryVersionLowerOrEqualToCurrent
		}
	case *openehr.HIER_OBJECT_ID:
		// Check HIER_OBJECT_ID matches current Directory UID
		if currentDirectoryID.UID() != v.Value {
			return openehr.FOLDER{}, ErrInvalidDirectoryUIDMismatch
		}

		// Increment version
		versionTreeID := currentDirectoryID.VersionTreeID()
		versionTreeID.Major++

		// Convert to OBJECT_VERSION_ID
		directory.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", v.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.FOLDER{}, fmt.Errorf("directory UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", directory.UID.V.Value)
	}

	// Prepare Original Version of Directory
	directoryVersion := openehr.ORIGINAL_VERSION{
		UID: *directory.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		PrecedingVersionUID: util.Some(*currentDirectoryID),
		Data:                &directory,
	}

	// Prepare contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.FOLDER_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &directoryVersion.UID,
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
				Value: "Directory updated",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	// Insert Contribution
	contribution.SetModelName()
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id, data) VALUES ($1, $2, $3)`
	args := []any{contribution.UID.Value, ehrID, contribution}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	// Insert Directory Version
	directoryVersion.SetModelName()
	query = `INSERT INTO openehr.tbl_folder (id, versioned_object_id, ehr_id, contribution_id, data) VALUES ($1, $2, $3, $4, $5)`
	args = []any{directoryVersion.UID.Value, directoryVersion.UID.UID(), ehrID, contribution.UID.Value, directoryVersion}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert updated directory into the database: %w", err)
	}

	return directory, nil
}

func (s *EHRService) DeleteDirectory(ctx context.Context, ehrID uuid.UUID) error {
	// Check if directory exists
	directory, err := s.GetDirectory(ctx, ehrID)
	if err != nil {
		return fmt.Errorf("failed to check if Directory exists: %w", err)
	}

	// Create contribution for deletion
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.FOLDER_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: directory.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
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
				Value: "Directory deleted",
			}),
			Committer: openehr.X_PARTY_PROXY{
				Value: &openehr.PARTY_SELF{
					Type_: util.Some(openehr.PARTY_SELF_MODEL_NAME),
				},
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

	// Insert Contribution
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
	args := []any{contribution.UID.Value, ehrID}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	contribution.SetModelName()
	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert contribution data into the database: %w", err)
	}

	// Delete Directory
	query = `DELETE FROM openehr.tbl_versioned_object WHERE id = $1 AND ehr_id = $2 AND type = $3`
	args = []any{directory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID(), ehrID, openehr.VERSIONED_FOLDER_MODEL_NAME}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to delete Directory from database: %w", err)
	}

	// Remove reference from EHR
	query = `UPDATE openehr.tbl_ehr_data
            SET data = data - 'directory'
            WHERE id = $1`
	args = []any{ehrID}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to update EHR to remove directory reference: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *EHRService) GetFolderInDirectoryVersion(ctx context.Context, ehrID string, filterAtTime time.Time, filterVersionID string, filterPathParts []string) (openehr.FOLDER, error) {
	var queryBuilder strings.Builder
	var args []any
	argNum := 1

	jsonPath := "$.**.data"

	for _, part := range filterPathParts {
		jsonPath += fmt.Sprintf(`.folders ? (@.name.value == "%s")`, part)
	}

	queryBuilder.WriteString(fmt.Sprintf(`
		SELECT jsonb_path_query_first(ovd.data, '%s') 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id
		WHERE ov.type = $1 
		  AND ov.ehr_id = $2
		  AND ovd.data @? $3
	`, jsonPath))
	args = []any{openehr.FOLDER_MODEL_NAME, ehrID, jsonPath}
	argNum += 3

	if !filterAtTime.IsZero() {
		queryBuilder.WriteString(fmt.Sprintf(`AND ov.created_at <= $%d `, argNum))
		args = append(args, filterAtTime)
		argNum++
	}

	if filterVersionID != "" {
		queryBuilder.WriteString(fmt.Sprintf(`AND ov.id = $%d `, argNum))
		args = append(args, filterVersionID)
	}

	queryBuilder.WriteString(`ORDER BY ov.created_at DESC LIMIT 1`)
	row := s.DB.QueryRow(ctx, queryBuilder.String(), args...)

	var folder openehr.FOLDER
	if err := row.Scan(&folder); err != nil {
		if err == database.ErrNoRows {
			if len(filterPathParts) > 0 {
				return openehr.FOLDER{}, ErrFolderNotFoundInDirectory
			}
			return openehr.FOLDER{}, ErrDirectoryNotFound
		}
		return openehr.FOLDER{}, fmt.Errorf("failed to fetch Folder at time from database: %w", err)
	}

	return folder, nil
}

func (s *EHRService) CreateContribution(ctx context.Context, ehrID uuid.UUID, contribution openehr.CONTRIBUTION) (openehr.CONTRIBUTION, error) {
	// Insert Contribution
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
	args := []any{contribution.UID.Value, ehrID}
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

func (s *EHRService) GetContribution(ctx context.Context, ehrID uuid.UUID, contributionID string) (openehr.CONTRIBUTION, error) {
	query := `
		SELECT cd.data
		FROM openehr.tbl_contribution c
		JOIN openehr.tbl_contribution_data cd ON c.id = cd.id
		WHERE c.ehr_id = $1 AND c.id = $2
		LIMIT 1
	`
	args := []any{ehrID, contributionID}
	row := s.DB.QueryRow(ctx, query, args...)

	var contribution openehr.CONTRIBUTION
	if err := row.Scan(&contribution); err != nil {
		if err == database.ErrNoRows {
			return openehr.CONTRIBUTION{}, ErrContributionNotFound
		}
		return openehr.CONTRIBUTION{}, fmt.Errorf("failed to fetch Contribution by ID from database: %w", err)
	}

	return contribution, nil
}
