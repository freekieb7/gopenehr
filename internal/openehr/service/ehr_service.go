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
	ErrDirectoryAlreadyExists                  = fmt.Errorf("directory already exists")
	ErrDirectoryNotFound                       = fmt.Errorf("directory not found")
	ErrContributionAlreadyExists               = fmt.Errorf("contribution already exists")
	ErrContributionNotFound                    = fmt.Errorf("contribution not found")
	ErrEHRStatusVersionLowerOrEqualToCurrent   = fmt.Errorf("EHR Status version must be incremented")
	ErrFolderNotFoundInDirectory               = fmt.Errorf("folder not found in directory")
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

	ehrStatus.SetModelName()
	if errs := ehrStatus.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new EHR Status: %v", errs)
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

	versionedEhrStatus.SetModelName()
	if errs := versionedEhrStatus.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new Versioned EHR Status: %v", errs)
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

	ehrStatusVersion.SetModelName()
	if errs := ehrStatusVersion.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new EHR Status Version: %v", errs)
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

	versionedEhrAccess.SetModelName()
	if errs := versionedEhrAccess.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new Versioned EHR Access: %v", errs)
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

	ehrAccessVersion.SetModelName()
	if errs := ehrAccessVersion.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new EHR Access Version: %v", errs)
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

	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
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

	ehr.SetModelName()
	if errs := ehr.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new EHR: %v", errs)
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

	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert contribution data into the database: %w", err)
	}

	// Insert Versioned EHR Status
	query = `INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4)`
	args = []any{versionedEhrStatus.UID.Value, openehr.VERSIONED_EHR_STATUS_MODEL_NAME, ehrID, contribution.UID.Value}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr status into the database: %w", err)
	}

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

	query = `INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)`
	args = []any{ehrStatusVersion.UID.Value, ehrStatusVersion}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr status version data into the database: %w", err)
	}

	// Insert Versioned EHR Access
	query = `INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4)`
	args = []any{versionedEhrAccess.UID.Value, openehr.VERSIONED_EHR_ACCESS_MODEL_NAME, ehrID, contribution.UID.Value}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr access into the database: %w", err)
	}

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
	currentEHRStatus, err := s.GetEHRStatus(ctx, ehrID, time.Time{}, "")
	if err != nil {
		return openehr.EHR_STATUS{}, fmt.Errorf("failed to get current EHR Status: %w", err)
	}
	currentEHRStatusID := currentEHRStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID)

	// Validate EHR Status
	if !ehrStatus.UID.E {
		ehrStatus.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: currentEHRStatusID.UID(),
			},
		})
	}

	switch ehrStatus.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check if version is being updated
		currentVersionID := currentEHRStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		newVersionID := ehrStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		if newVersionID.VersionTreeID().CompareTo(currentVersionID.VersionTreeID()) <= 0 {
			return openehr.EHR_STATUS{}, ErrEHRStatusVersionLowerOrEqualToCurrent
		}
		// valid type
	case *openehr.HIER_OBJECT_ID:
		// Add namespace and version to convert to OBJECT_VERSION_ID
		hierID := ehrStatus.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		versionTreeID := currentEHRStatusID.VersionTreeID()
		versionTreeID.Major++

		ehrStatus.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", hierID.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.EHR_STATUS{}, fmt.Errorf("EHR Status UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", ehrStatus.UID.V.Value)
	}

	if errs := ehrStatus.Validate("$"); len(errs) > 0 {
		return openehr.EHR_STATUS{}, fmt.Errorf("validation errors for updated EHR Status: %v", errs)
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
		Data: &ehrStatus,
	}

	ehrStatusVersion.SetModelName()
	if errs := ehrStatusVersion.Validate("$"); len(errs) > 0 {
		return openehr.EHR_STATUS{}, fmt.Errorf("validation errors for new EHR Status Version: %v", errs)
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

	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.EHR_STATUS{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
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

func (s *EHRService) CreateComposition(ctx context.Context, ehrID uuid.UUID, composition openehr.COMPOSITION) (openehr.COMPOSITION, error) {
	// Validate Composition
	if errs := composition.Validate("$"); len(errs) > 0 {
		return openehr.COMPOSITION{}, fmt.Errorf("validation errors for provided Composition: %v", errs)
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

	// Extract UID type
	compositionID, ok := composition.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return openehr.COMPOSITION{}, fmt.Errorf("composition UID must be of type OBJECT_VERSION_ID, got %T", composition.UID.V.Value)
	}

	// Check if Composition Exists
	_, err := s.GetComposition(ctx, ehrID, compositionID.UID())
	if err == nil {
		return openehr.COMPOSITION{}, ErrCompositionAlreadyExists
	}
	if err != ErrCompositionNotFound {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to check if Composition exists: %w", err)
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

	if errs := versionedComposition.Validate("$"); len(errs) > 0 {
		return openehr.COMPOSITION{}, fmt.Errorf("validation errors for new Versioned Composition: %v", errs)
	}

	versionedComposition.SetModelName()

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

	if errs := compositionVersion.Validate("$"); len(errs) > 0 {
		return openehr.COMPOSITION{}, fmt.Errorf("validation errors for new Composition Version: %v", errs)
	}

	compositionVersion.SetModelName()

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

	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.COMPOSITION{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
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
	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert contribution data into the database: %w", err)
	}

	// Insert Versioned Composition
	query = `INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4)`
	args = []any{versionedComposition.UID.Value, openehr.VERSIONED_COMPOSITION_MODEL_NAME, ehrID, contribution.UID.Value}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert versioned composition into the database: %w", err)
	}
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

func (s *EHRService) GetComposition(ctx context.Context, ehrID uuid.UUID, filterUIDBasedID string) (openehr.COMPOSITION, error) {
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
	args = []any{openehr.COMPOSITION_MODEL_NAME, ehrID}
	argNum += 2

	if filterUIDBasedID != "" {
		if strings.Count(filterUIDBasedID, "::") == 2 {
			// UID is of type OBJECT_VERSION_ID
			query.WriteString(fmt.Sprintf(`AND ov.id = $%d `, argNum))
		} else {
			// UID is of type HIER_OBJECT_ID
			query.WriteString(fmt.Sprintf(`AND ov.versioned_object_id = $%d `, argNum))
		}
		args = append(args, filterUIDBasedID)
	}

	query.WriteString(`ORDER BY ov.created_at DESC LIMIT 1`)

	var composition openehr.COMPOSITION
	if err := s.DB.QueryRow(ctx, query.String(), args...).Scan(&composition); err != nil {
		if err == database.ErrNoRows {
			return openehr.COMPOSITION{}, ErrCompositionNotFound
		}
		return openehr.COMPOSITION{}, fmt.Errorf("failed to fetch Composition by ID from database: %w", err)
	}

	return composition, nil
}

func (s *EHRService) UpdateComposition(ctx context.Context, ehrID uuid.UUID, composition openehr.COMPOSITION) (openehr.COMPOSITION, error) {
	// Validate Composition
	if !composition.UID.E {
		return openehr.COMPOSITION{}, fmt.Errorf("composition UID must be provided for update")
	}

	currentComposition, err := s.GetComposition(ctx, ehrID, "")
	if err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to get current Composition: %w", err)
	}

	switch composition.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// Check if version is being updated
		currentVersionID := currentComposition.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		newVersionID := composition.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
		if currentVersionID.VersionTreeID().CompareTo(newVersionID.VersionTreeID()) <= 0 {
			return openehr.COMPOSITION{}, ErrCompositionVersionLowerOrEqualToCurrent
		}
		// valid type

	case *openehr.HIER_OBJECT_ID:
		// Add namespace and version to convert to OBJECT_VERSION_ID
		hierID := composition.UID.V.Value.(*openehr.HIER_OBJECT_ID)
		versionTreeID := currentComposition.UID.V.Value.(*openehr.OBJECT_VERSION_ID).VersionTreeID()
		versionTreeID.Major++

		composition.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::%s", hierID.Value, config.SYSTEM_ID_GOPENEHR, versionTreeID.String()),
			},
		})
	default:
		return openehr.COMPOSITION{}, fmt.Errorf("composition UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", composition.UID.V.Value)
	}

	// Validate updated Composition UID
	if errs := composition.Validate("$"); len(errs) > 0 {
		return openehr.COMPOSITION{}, fmt.Errorf("validation errors for updated Composition UID: %v", errs)
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
		Data: &composition,
	}
	compositionVersion.SetModelName()

	// Just for ensurance that models are valid
	if errs := compositionVersion.Validate("$"); len(errs) > 0 {
		return openehr.COMPOSITION{}, fmt.Errorf("validation errors for new Composition Version: %v", errs)
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
	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.COMPOSITION{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
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

	// Insert Contribution
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
	args := []any{contribution.UID.Value, ehrID}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert contribution into the database: %w", err)
	}
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

func (s *EHRService) CreateDirectory(ctx context.Context, ehrID uuid.UUID, providedDirectory util.Optional[openehr.FOLDER]) (openehr.FOLDER, error) {
	// Check if Directory Exists
	_, err := s.GetDirectory(ctx, ehrID)
	if err == nil {
		return openehr.FOLDER{}, ErrDirectoryAlreadyExists
	}
	if err != ErrDirectoryNotFound {
		return openehr.FOLDER{}, fmt.Errorf("failed to check if Directory exists: %w", err)
	}

	// Prepare Directory
	directory := openehr.FOLDER{
		UID: util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::gopenehr::1", uuid.NewString()),
			},
		}),
		Name: openehr.X_DV_TEXT{
			Value: &openehr.DV_TEXT{
				Value: "Directory",
			},
		},
		ArchetypeNodeID: "openEHR-EHR-FOLDER.directory.v1",
	}

	if providedDirectory.E {
		if !providedDirectory.V.UID.E {
			providedDirectory.V.UID = directory.UID
		}

		directory = providedDirectory.V
	}

	// Validate Directory
	if errs := directory.Validate("$"); len(errs) > 0 {
		return openehr.FOLDER{}, fmt.Errorf("validation errors for provided Directory: %v", errs)
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

	if errs := versionedDirectory.Validate("$"); len(errs) > 0 {
		return openehr.FOLDER{}, fmt.Errorf("validation errors for new Versioned Directory: %v", errs)
	}

	versionedDirectory.SetModelName()

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

	if errs := directoryVersion.Validate("$"); len(errs) > 0 {
		return openehr.FOLDER{}, fmt.Errorf("validation errors for new Directory Version: %v", errs)
	}

	directoryVersion.SetModelName()

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

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.FOLDER{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	contribution.SetModelName()

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
	query = `INSERT INTO openehr.tbl_contribution_data (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert contribution data into the database: %w", err)
	}

	// Insert Versioned Directory
	query = `INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4)`
	args = []any{versionedDirectory.UID.Value, openehr.VERSIONED_FOLDER_MODEL_NAME, ehrID, contribution.UID.Value}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert versioned directory into the database: %w", err)
	}
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
	query = `INSERT INTO openehr.tbl_object_version_data (id, data) VALUES ($1, $2)`
	args = []any{directoryVersion.UID.Value, directoryVersion}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert directory data into the database: %w", err)
	}

	// Update EHR
	query = `UPDATE openehr.tbl_ehr_data
            SET data = jsonb_set(data, '{directory}', to_jsonb($1::jsonb)) 
            WHERE id = $2`
	directoryRef := openehr.OBJECT_REF{
		Namespace: config.NAMESPACE_LOCAL,
		Type:      openehr.FOLDER_MODEL_NAME,
		ID: openehr.X_OBJECT_ID{
			Value: &directoryVersion.UID,
		},
	}
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

func (s *EHRService) UpdateDirectory(ctx context.Context, ehrID uuid.UUID, updatedDirectory openehr.FOLDER) error {
	// Validate updated Directory UID
	if errs := updatedDirectory.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for updated Directory UID: %v", errs)
	}

	if !updatedDirectory.UID.E {
		return fmt.Errorf("directory UID must be set for update")
	}

	// Check if Directory exists
	currectDirectory, err := s.GetDirectory(ctx, ehrID)
	if err != nil {
		return fmt.Errorf("failed to check if Directory exists: %w", err)
	}

	// Ensure that the versioned object ID matches
	if currectDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID() != updatedDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID() {
		return ErrDirectoryNotFound
	}

	// Ensure that the version is incremented
	currectVersion := currectDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).VersionTreeID()
	updatedVersion := updatedDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).VersionTreeID()
	if updatedVersion.CompareTo(currectVersion) <= 0 {
		return fmt.Errorf("updated Directory version must be greater than current version")
	}

	// Validate updated Directory
	if errs := updatedDirectory.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for updated Directory: %v", errs)
	}

	// Prepare Original Version of Directory
	directoryVersion := openehr.ORIGINAL_VERSION{
		UID: *updatedDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &updatedDirectory,
	}
	directoryVersion.SetModelName()

	// Just for ensurance that models are valid
	if errs := directoryVersion.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for new Directory Version: %v", errs)
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
	contribution.SetModelName()

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Insert Directory Version
	query := `INSERT INTO openehr.tbl_folder (id, versioned_object_id, ehr_id, data) 
		 VALUES ($1, $2, $3, $4)`
	args := []any{directoryVersion.UID.Value, directoryVersion.UID.UID(), ehrID, directoryVersion}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert updated directory into the database: %w", err)
	}

	// Insert Contribution
	query = `INSERT INTO openehr.tbl_contribution (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{contribution.UID.Value, ehrID, contribution}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	return nil
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

	// Insert Contribution
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
	args := []any{contribution.UID.Value, ehrID}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert contribution into the database: %w", err)
	}
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

	queryBuilder.WriteString(fmt.Sprintf("SELECT jsonb_path_query_first(ovd.data, '%s') ", jsonPath))
	queryBuilder.WriteString(`
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id
		WHERE ov.type = $1 
		  AND ov.ehr_id = $2
		  AND ovd.data @? $3
	`)
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
	// Validate Contribution
	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.CONTRIBUTION{}, fmt.Errorf("validation errors for Contribution: %v", errs)
	}

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
