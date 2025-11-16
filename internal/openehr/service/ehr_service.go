package service

import (
	"context"
	"encoding/json"
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
	ErrEHRNotFound               = fmt.Errorf("EHR not found")
	ErrEHRAlreadyExists          = fmt.Errorf("EHR already exists")
	ErrEHRStatusNotFound         = fmt.Errorf("EHR Status not found")
	ErrEHRStatusAlreadyExists    = fmt.Errorf("EHR Status already exists")
	ErrCompositionAlreadyExists  = fmt.Errorf("composition already exists")
	ErrCompositionNotFound       = fmt.Errorf("composition not found")
	ErrDirectoryAlreadyExists    = fmt.Errorf("directory already exists")
	ErrDirectoryNotFound         = fmt.Errorf("directory not found")
	ErrContributionAlreadyExists = fmt.Errorf("contribution already exists")
	ErrContributionNotFound      = fmt.Errorf("contribution not found")
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

func (s *EHRService) CreateEHR(ctx context.Context, ehrStatus util.Optional[openehr.EHR_STATUS]) (openehr.EHR, error) {
	return s.CreateEHRWithID(ctx, ehrStatus, uuid.NewString())
}

func (s *EHRService) CreateEHRWithID(ctx context.Context, providedEHRStatus util.Optional[openehr.EHR_STATUS], ehrID string) (openehr.EHR, error) {
	// Validate EHR Status
	if providedEHRStatus.E {
		if providedEHRStatus.V.UID.E {
			ehrStatusVersionID, ok := providedEHRStatus.V.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
			if !ok {
				return openehr.EHR{}, fmt.Errorf("EHR Status UID must be of type OBJECT_VERSION_ID, got %T", providedEHRStatus.V.UID.V.Value)
			}

			// Check if EHR Status with given version ID already exists
			var exists bool
			if err := s.DB.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM tbl_openehr_ehr_status WHERE id = $1)`, ehrStatusVersionID.UID()).Scan(&exists); err != nil {
				return openehr.EHR{}, fmt.Errorf("failed to check if EHR Status exists in database: %w", err)
			}
			if exists {
				return openehr.EHR{}, ErrEHRStatusAlreadyExists
			}
		}

		if errs := providedEHRStatus.V.Validate("$"); len(errs) > 0 {
			return openehr.EHR{}, fmt.Errorf("validation errors for provided EHR Status: %v", errs)
		}
	}

	// Prepare EHR Status
	var ehrStatus openehr.EHR_STATUS
	if providedEHRStatus.E {
		ehrStatus = providedEHRStatus.V
	} else {
		ehrStatus = openehr.EHR_STATUS{
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

	if !ehrStatus.UID.E {
		ehrStatus.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", uuid.NewString(), config.SYSTEM_ID_GOPENEHR),
			},
		})
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
					Value: ehrID,
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
	newEhrStatusVersion := openehr.ORIGINAL_VERSION{
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

	newEhrStatusVersion.SetModelName()
	if errs := newEhrStatusVersion.Validate("$"); len(errs) > 0 {
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
					Value: ehrID,
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
	newEhrAccess := openehr.EHR_ACCESS{
		UID: util.Some(openehr.X_UID_BASED_ID{
			Value: &versionedEhrAccess.UID,
		}),
		Name: openehr.X_DV_TEXT{
			Value: &openehr.DV_TEXT{
				Value: "EHR Access",
			},
		},
		ArchetypeNodeID: "openEHR-EHR-EHR_ACCESS.generic.v1",
	}

	// Build EHR Access Version
	newEhrAccessVersion := openehr.ORIGINAL_VERSION{
		UID: *newEhrAccess.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &newEhrAccess,
	}

	newEhrAccessVersion.SetModelName()
	if errs := newEhrAccessVersion.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new EHR Access Version: %v", errs)
	}

	// Prepare EHR
	newEhr := openehr.EHR{
		EHRID: openehr.HIER_OBJECT_ID{
			Value: ehrID,
		},
		EHRStatus: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.VERSIONED_EHR_STATUS_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &versionedEhrStatus.UID,
			},
		},
		EHRAccess: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.VERSIONED_EHR_ACCESS_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &versionedEhrAccess.UID,
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	newEhr.SetModelName()
	if errs := newEhr.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new EHR: %v", errs)
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
					Value: &newEhrStatusVersion.UID,
				},
			},
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.EHR_ACCESS_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &newEhrAccessVersion.UID,
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
				}, // TODO make this configurable, could also be set to internal person
			},
		},
	}

	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
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
	query := `INSERT INTO tbl_openehr_ehr (id, data) VALUES ($1, $2)`
	args := []any{ehrID, newEhr}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr into the database: %w", err)
	}

	// Insert Versioned EHR Status
	query = `INSERT INTO tbl_openehr_versioned_object (id, ehr_id, type, data) VALUES ($1, $2, $3, $4)`
	args = []any{versionedEhrStatus.UID.Value, ehrID, openehr.EHR_STATUS_MODEL_NAME, versionedEhrStatus}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr status into the database: %w", err)
	}

	// Insert EHR Status Version
	query = `INSERT INTO tbl_openehr_ehr_status (id, versioned_object_id, ehr_id, data) 
         VALUES ($1, $2, $3, $4)`
	args = []any{newEhrStatusVersion.UID.Value, versionedEhrStatus.UID.Value, ehrID, newEhrStatusVersion}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr status into the database: %w", err)
	}

	// Insert Versioned EHR Access
	query = `INSERT INTO tbl_openehr_versioned_object (id, ehr_id, type, data) VALUES ($1, $2, $3, $4)`
	args = []any{versionedEhrAccess.UID.Value, ehrID, openehr.EHR_ACCESS_MODEL_NAME, versionedEhrAccess}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr access into the database: %w", err)
	}

	// Insert EHR Access Version
	query = `INSERT INTO tbl_openehr_ehr_access (id, versioned_object_id, ehr_id, data) 
         VALUES ($1, $2, $3, $4)`
	args = []any{newEhrAccessVersion.UID.Value, versionedEhrAccess.UID.Value, ehrID, newEhrAccessVersion}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr access into the database: %w", err)
	}

	// Insert Contribution
	query = `INSERT INTO tbl_openehr_contribution (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{contribution.UID.Value, ehrID, contribution}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newEhr, nil
}

// func (s *EHRService) CheckEHRExists(ctx context.Context, id string) (bool, error) {
// 	query := `SELECT EXISTS (SELECT 1 FROM tbl_openehr_ehr WHERE id = $1)`
// 	args := []any{id}
// 	var exists bool
// 	if err := s.DB.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
// 		return false, fmt.Errorf("failed to check if EHR exists in database: %w", err)
// 	}
// 	return exists, nil
// }

func (s *EHRService) GetEHRAsJSON(ctx context.Context, id string) ([]byte, error) {
	query := `SELECT data FROM tbl_openehr_ehr WHERE id = $1`
	args := []any{id}
	row := s.DB.QueryRow(ctx, query, args...)

	var ehrJSON []byte
	if err := row.Scan(&ehrJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR from database: %w", err)
	}

	return ehrJSON, nil
}

func (s *EHRService) GetEHRBySubjectAsJSON(ctx context.Context, subjectID, subjectNamespace string) ([]byte, error) {
	query := `
        SELECT e.data 
        FROM tbl_openehr_ehr e 
        JOIN tbl_openehr_ehr_status es ON e.id = es.ehr_id
        WHERE es.data->'subject'->'external_ref'->>'namespace' = $1
          AND es.data->'subject'->'external_ref'->'id'->>'value' = $2
		ORDER BY es.created_at DESC
        LIMIT 1
    `
	args := []any{subjectNamespace, subjectID}
	row := s.DB.QueryRow(ctx, query, args...)

	var rawEhrJSON []byte
	if err := row.Scan(&rawEhrJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR by subject from database: %w", err)
	}

	return rawEhrJSON, nil
}

func (s *EHRService) DeleteEHRByID(ctx context.Context, id string) error {
	var deleted uint8
	query := `DELETE FROM tbl_openehr_ehr WHERE id = $1 RETURNING 1`
	args := []any{id}
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&deleted); err != nil {
		if err == database.ErrNoRows {
			return ErrEHRNotFound
		}

		return fmt.Errorf("failed to delete EHR from database: %w", err)
	}
	return nil
}

func (s *EHRService) DeleteMultipleEHRs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(`DELETE FROM tbl_openehr_ehr WHERE id IN (%s)`, strings.Join(placeholders, ", "))

	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to delete multiple EHRs from database: %w", err)
	}
	return nil
}

// GetEHRStatus retrieves the latest EHR Status for the given EHR ID.
func (s *EHRService) GetEHRStatus(ctx context.Context, ehrID string) (openehr.EHR_STATUS, error) {
	query := `SELECT data FROM tbl_openehr_ehr_status WHERE ehr_id = $1 ORDER BY created_at DESC LIMIT 1`
	args := []any{ehrID}

	row := s.DB.QueryRow(ctx, query, args...)

	var ehrStatus openehr.EHR_STATUS
	if err := row.Scan(&ehrStatus); err != nil {
		if err == database.ErrNoRows {
			return openehr.EHR_STATUS{}, ErrEHRStatusNotFound
		}
		return openehr.EHR_STATUS{}, fmt.Errorf("failed to fetch EHR Status from database: %w", err)
	}

	return ehrStatus, nil
}

func (s *EHRService) GetEHRStatusAsJSON(ctx context.Context, ehrID string, filterOnTime time.Time, filterOnVersionID string) ([]byte, error) {
	// Build query
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`
		SELECT data 
		FROM tbl_openehr_ehr_status
		WHERE ehr_id = $1 
	`)
	args = []any{ehrID}
	argNum++

	if !filterOnTime.IsZero() {
		query.WriteString(fmt.Sprintf(`AND created_at <= $%d `, argNum))
		args = append(args, filterOnTime)
		argNum++
	}

	if filterOnVersionID != "" {
		query.WriteString(fmt.Sprintf(`AND versioned_object_id = $%d `, argNum))
		args = append(args, filterOnVersionID)
	}

	query.WriteString(`ORDER BY created_at DESC LIMIT 1`)
	row := s.DB.QueryRow(ctx, query.String(), args...)

	var rawEhrStatusJSON []byte
	if err := row.Scan(&rawEhrStatusJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR Status from database: %w", err)
	}

	return rawEhrStatusJSON, nil
}

func (s *EHRService) UpdateEHRStatus(ctx context.Context, ehrID string, ehrStatus openehr.EHR_STATUS) error {
	// Validate EHR Status
	if !ehrStatus.UID.E {
		return fmt.Errorf("EHR Status UID must be set for update")
	}

	ehrStatusVersionID, ok := ehrStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID)
	if !ok {
		return fmt.Errorf("EHR Status UID must be of type OBJECT_VERSION_ID, got %T", ehrStatus.UID.V.Value)
	}

	if errs := ehrStatus.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for updated EHR Status: %v", errs)
	}

	// Prepare Original Version of EHR Status
	ehrStatusVersion := openehr.ORIGINAL_VERSION{
		UID: *ehrStatusVersionID,
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
		return fmt.Errorf("validation errors for new EHR Status Version: %v", errs)
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
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Insert new EHR Status
	query := `INSERT INTO tbl_openehr_ehr_status (id, versioned_object_id, ehr_id, data) VALUES ($1, $2, $3, $4)`
	args := []any{ehrStatusVersion.UID.Value, ehrStatusVersion.UID.UID(), ehrID, ehrStatusVersion}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert ehr status into the database: %w", err)
	}

	// Insert Contribution
	query = `INSERT INTO tbl_openehr_contribution (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{contribution.UID.Value, ehrID, contribution}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *EHRService) GetVersionedEHRStatusAsJSON(ctx context.Context, ehrID string) ([]byte, error) {
	query := `
		SELECT vo.data 
		FROM tbl_openehr_versioned_object vo 
		WHERE vo.ehr_id = $1 AND vo.data->>'_type' = $2
		LIMIT 1
	`
	args := []any{ehrID, openehr.VERSIONED_EHR_STATUS_MODEL_NAME}
	row := s.DB.QueryRow(ctx, query, args...)

	var rawVersionedEhrStatusJSON []byte
	if err := row.Scan(&rawVersionedEhrStatusJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch Versioned EHR Status from database: %w", err)
	}

	return rawVersionedEhrStatusJSON, nil
}

func (s *EHRService) GetVersionedEHRStatusRevisionHistoryAsJSON(ctx context.Context, ehrID string) ([]byte, error) {
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
			FROM tbl_openehr_contribution c,
				jsonb_array_elements(c.data->'versions') as version
			WHERE c.ehr_id = $1
				AND version->>'type' = 'EHR_STATUS'
			GROUP BY version->'id'->>'value'
		) grouped
    `
	args := []any{ehrID}
	row := s.DB.QueryRow(ctx, query, args...)

	var rawRevisionHistoryDataJSON []byte
	if err := row.Scan(&rawRevisionHistoryDataJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch Revision History from database: %w", err)
	}

	return rawRevisionHistoryDataJSON, nil
}

func (s *EHRService) GetVersionedEHRStatusVersionAsJSON(ctx context.Context, ehrID string, filterAtTime time.Time, filterVersionID string) ([]byte, error) {
	// Fetch EHR Status at given time
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`SELECT data FROM tbl_openehr_ehr_status WHERE ehr_id = $1 `)
	args = []any{ehrID}
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

func (s *EHRService) ExistsEHRComposition(ctx context.Context, compositionID string) (bool, error) {
	// Check if Composition Exists
	var existsComp bool
	if err := s.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM tbl_openehr_composition WHERE id = $1)", compositionID).Scan(&existsComp); err != nil {
		return false, fmt.Errorf("failed to check if Composition UID exists: %w", err)
	}

	return existsComp, nil
}

func (s *EHRService) CreateComposition(ctx context.Context, ehrID string, composition openehr.COMPOSITION) (openehr.COMPOSITION, error) {
	if !composition.UID.E {
		composition.UID = util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::gopenehr::1", uuid.NewString()),
			},
		})
	} else {
		switch v := composition.UID.V.Value.(type) {
		case *openehr.OBJECT_VERSION_ID:
			// OK
		default:
			return openehr.COMPOSITION{}, fmt.Errorf("composition UID must be of type OBJECT_VERSION_ID, got %T", v)
		}

		// Check if Composition Exists
		existsComp, err := s.ExistsEHRComposition(ctx, composition.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value)
		if err != nil {
			return openehr.COMPOSITION{}, fmt.Errorf("failed to check if Composition exists: %w", err)
		}
		if existsComp {
			return openehr.COMPOSITION{}, ErrCompositionAlreadyExists
		}
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
					Value: ehrID,
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

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.COMPOSITION{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	contribution.SetModelName()

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

	// Insert Versioned Composition
	query := `INSERT INTO tbl_openehr_versioned_object (id, ehr_id, type, data) VALUES ($1, $2, $3, $4)`
	args := []any{versionedComposition.UID.Value, ehrID, openehr.COMPOSITION_MODEL_NAME, versionedComposition}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert versioned composition into the database: %w", err)
	}

	// Insert Composition Version
	query = `INSERT INTO tbl_openehr_composition (id, versioned_object_id, ehr_id, data) 
		 VALUES ($1, $2, $3, $4)`
	args = []any{compositionVersion.UID.Value, versionedComposition.UID.Value, ehrID, compositionVersion}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert composition into the database: %w", err)
	}

	// Insert Contribution
	query = `INSERT INTO tbl_openehr_contribution (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{contribution.UID.Value, ehrID, contribution}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	// Update EHR
	query = `UPDATE tbl_openehr_ehr 
			SET data = jsonb_set(data, '{compositions}', 
				COALESCE(data->'compositions', '[]'::jsonb) || to_jsonb($1::jsonb)
			) 
			WHERE id = $2`
	compositionRef := openehr.OBJECT_REF{
		Namespace: config.NAMESPACE_LOCAL,
		Type:      openehr.VERSIONED_COMPOSITION_MODEL_NAME,
		ID: openehr.X_OBJECT_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: versionedComposition.UID.Value,
			},
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

func (s *EHRService) GetComposition(ctx context.Context, compositionID string, filterOnEHRID string) (openehr.COMPOSITION, error) {
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`
		SELECT c.data 
		FROM tbl_openehr_composition c
		WHERE c.id = $1
	`)
	args = []any{compositionID}
	argNum++

	if filterOnEHRID != "" {
		query.WriteString(fmt.Sprintf(` AND c.ehr_id = $%d`, argNum))
		args = append(args, filterOnEHRID)
	}

	query.WriteString(` LIMIT 1`)

	var composition openehr.COMPOSITION
	if err := s.DB.QueryRow(ctx, query.String(), args...).Scan(&composition); err != nil {
		if err == database.ErrNoRows {
			return openehr.COMPOSITION{}, ErrCompositionNotFound
		}
		return openehr.COMPOSITION{}, fmt.Errorf("failed to fetch Composition by ID from database: %w", err)
	}

	return composition, nil
}

func (s *EHRService) GetCurrentCompositionByVersionedObjectID(ctx context.Context, versionedObjectID string) (openehr.COMPOSITION, error) {
	var composition openehr.COMPOSITION
	if err := s.DB.QueryRow(ctx, `
		SELECT c.data 
		FROM tbl_openehr_composition c
		WHERE c.versioned_object_id = $1
		ORDER BY c.created_at DESC
		LIMIT 1
	`, versionedObjectID).Scan(&composition); err != nil {
		if err == database.ErrNoRows {
			return openehr.COMPOSITION{}, ErrCompositionNotFound
		}
		return openehr.COMPOSITION{}, fmt.Errorf("failed to fetch Composition by Versioned Object ID from database: %w", err)
	}

	return composition, nil
}

// GetComposition retrieves the composition as the latest version when providing versioned object id, or the specified ID version.
func (s *DemographicService) GetComposition(ctx context.Context, ehrID, uidBasedID string) (openehr.COMPOSITION, error) {
	var composition openehr.COMPOSITION

	query := "SELECT data FROM tbl_openehr_composition WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM tbl_openehr_composition WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&composition)
	if err != nil {
		return openehr.COMPOSITION{}, fmt.Errorf("failed to get composition by ID: %w", err)
	}
	return composition, nil
}

// GetCompositionAsJSON retrieves the composition as the latest version when providing versioned object id, or the specified ID version, and returns it as the raw JSON representation.
// This is useful for scenarios where the raw JSON is needed without unmarshalling into Go structs.
func (s *EHRService) GetCompositionAsJSON(ctx context.Context, ehrID, uidBasedID string) ([]byte, error) {
	var compositionJSON []byte

	query := "SELECT data FROM tbl_openehr_composition WHERE versioned_object_id = $1 ORDER BY created_at DESC LIMIT 1"
	if strings.Contains(uidBasedID, "::") {
		query = "SELECT data FROM tbl_openehr_composition WHERE id = $1 LIMIT 1"
	}

	err := s.DB.QueryRow(ctx, query, uidBasedID).Scan(&compositionJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to get composition by ID: %w", err)
	}
	return compositionJSON, nil
}

func (s *EHRService) GetRawCurrentCompositionByVersionedObjectID(ctx context.Context, ehrID, versionedObjectID string) ([]byte, error) {
	// Fetch Composition by Versioned Object ID
	query := `
		SELECT c.data 
		FROM tbl_openehr_composition c
		WHERE c.ehr_id = $1 AND c.versioned_object_id = $2
		ORDER BY c.created_at DESC
		LIMIT 1
	`
	args := []any{ehrID, versionedObjectID}
	row := s.DB.QueryRow(ctx, query, args...)

	var rawCompositionJSON []byte
	if err := row.Scan(&rawCompositionJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrCompositionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Composition by Versioned Object ID from database: %w", err)
	}

	return rawCompositionJSON, nil
}

func (s *EHRService) UpdateCompositionByID(ctx context.Context, ehrID string, updatedComposition openehr.COMPOSITION) error {
	// Validate updated Composition UID
	if errs := updatedComposition.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for updated Composition UID: %v", errs)
	}

	if !updatedComposition.UID.E {
		return fmt.Errorf("composition UID must be set for update")
	}

	// Check if Composition exists
	existsComp, err := s.ExistsEHRComposition(ctx, updatedComposition.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value)
	if err != nil {
		return fmt.Errorf("failed to check if Composition exists: %w", err)
	}
	if existsComp {
		return ErrCompositionAlreadyExists
	}

	if errs := updatedComposition.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for updated Composition: %v", errs)
	}

	// Prepare Original Version of Composition
	compositionVersion := openehr.ORIGINAL_VERSION{
		UID: *updatedComposition.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &updatedComposition,
	}
	compositionVersion.SetModelName()

	// Just for ensurance that models are valid
	if errs := compositionVersion.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for new Composition Version: %v", errs)
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
		return fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Insert Composition Version
	query := `INSERT INTO tbl_openehr_composition (id, versioned_object_id, ehr_id, data) 
		 VALUES ($1, $2, $3, $4)`
	args := []any{compositionVersion.UID.Value, strings.Split(compositionVersion.UID.Value, "::")[0], ehrID, compositionVersion}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert updated composition into the database: %w", err)
	}

	// Insert Contribution
	query = `INSERT INTO tbl_openehr_contribution (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{contribution.UID.Value, ehrID, contribution}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	return nil
}

func (s *EHRService) DeleteComposition(ctx context.Context, versionedObjectID string) error {
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

	// Delete Composition
	query := `DELETE FROM tbl_openehr_versioned_object WHERE id = $1 AND type = $2 RETURNING ehr_id`
	args := []any{versionedObjectID, openehr.COMPOSITION_MODEL_NAME}

	var ehrID string
	if err := tx.QueryRow(ctx, query, args...).Scan(&ehrID); err != nil {
		if err == database.ErrNoRows {
			return ErrCompositionNotFound
		}
		return fmt.Errorf("failed to delete Composition from database: %w", err)
	}

	// Insert Contribution
	query = `INSERT INTO tbl_openehr_contribution (id, data) VALUES ($1, $2)`
	args = []any{contribution.UID.Value, contribution}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	// Cleanup if there are no more versions of the Versioned Composition
	var existsVersions bool
	query = `
		SELECT EXISTS (
			SELECT 1 FROM tbl_openehr_composition c
			WHERE c.versioned_object_id = $1
		)
	`
	args = []any{versionedObjectID}
	if err := tx.QueryRow(ctx, query, args...).Scan(&existsVersions); err != nil {
		return fmt.Errorf("failed to check for existing Composition versions: %w", err)
	}

	if !existsVersions {
		// Delete Versioned Composition
		query = `DELETE FROM tbl_openehr_versioned_object WHERE ehr_id = $1 AND id = $2`
		args = []any{ehrID, versionedObjectID}
		if _, err := tx.Exec(ctx, query, args...); err != nil {
			return fmt.Errorf("failed to delete Versioned Composition from database: %w", err)
		}

		// Also remove reference from EHR
		query = `UPDATE tbl_openehr_ehr 
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
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *EHRService) GetVersionedCompositionByIDAsJSON(ctx context.Context, ehrID, versionedObjectID string) ([]byte, error) {
	query := `
		SELECT vo.data 
		FROM tbl_openehr_versioned_object vo
		WHERE vo.ehr_id = $1 AND vo.id = $2 AND vo.data->>'_type' = $3
		LIMIT 1
	`
	args := []any{ehrID, versionedObjectID, openehr.VERSIONED_COMPOSITION_MODEL_NAME}
	row := s.DB.QueryRow(ctx, query, args...)

	var versionedCompositionJSON []byte
	if err := row.Scan(&versionedCompositionJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrCompositionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Versioned Composition by ID from database: %w", err)
	}

	return versionedCompositionJSON, nil
}

func (s *EHRService) GetVersionedCompositionRevisionHistoryAsJSON(ctx context.Context, ehrID, versionedObjectID string) ([]byte, error) {
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
			FROM tbl_openehr_contribution c,
				jsonb_array_elements(c.data->'versions') as version
			WHERE c.ehr_id = $1
				AND version->>'type' = 'COMPOSITION'
				AND version->'id'->>'value' LIKE $2 || '%'
			GROUP BY version->'id'->>'value'
		) grouped
	`
	args := []any{ehrID, versionedObjectID}
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

func (s *EHRService) GetVersionedCompositionVersionAsJSON(ctx context.Context, ehrID, versionedObjectID string, filterAtTime time.Time, filterVersionID string) ([]byte, error) {
	// Fetch Composition at given time
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`SELECT data FROM tbl_openehr_composition WHERE ehr_id = $1 AND versioned_object_id = $2 `)
	args = []any{ehrID, versionedObjectID}
	argNum += 2

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

	var compositionVersionJSON []byte
	if err := row.Scan(&compositionVersionJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrCompositionNotFound
		}
		return nil, fmt.Errorf("failed to fetch Composition version at time from database: %w", err)
	}

	return compositionVersionJSON, nil
}

func (s *EHRService) ExistsEHRDirectoryForEHR(ctx context.Context, ehrID string) (bool, error) {
	// Check if Directory Exists
	// Any folder for the given EHR ID indicates the presence of a directory
	var existsDir bool
	if err := s.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM tbl_openehr_folder WHERE ehr_id = $1)", ehrID).Scan(&existsDir); err != nil {
		return false, fmt.Errorf("failed to check if Directory UID exists: %w", err)
	}

	return existsDir, nil
}

func (s *EHRService) CreateDirectory(ctx context.Context, ehrID string) (openehr.FOLDER, error) {
	// Check if Directory Exists
	existsDir, err := s.ExistsEHRDirectoryForEHR(ctx, ehrID)
	if err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to check if Directory exists: %w", err)
	}
	if existsDir {
		return openehr.FOLDER{}, ErrDirectoryAlreadyExists
	}

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
					Value: ehrID,
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

	// Insert Versioned Directory
	query := `INSERT INTO tbl_openehr_versioned_object (id, ehr_id, type, data) VALUES ($1, $2, $3, $4)`
	args := []any{versionedDirectory.UID.Value, ehrID, openehr.FOLDER_MODEL_NAME, versionedDirectory}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert versioned directory into the database: %w", err)
	}

	// Insert Directory Version
	query = `INSERT INTO tbl_openehr_folder (id, versioned_object_id, ehr_id, data) 
		 VALUES ($1, $2, $3, $4)`
	args = []any{directoryVersion.UID.Value, versionedDirectory.UID.Value, ehrID, directoryVersion}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert directory into the database: %w", err)
	}

	// Insert Contribution
	query = `INSERT INTO tbl_openehr_contribution (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{contribution.UID.Value, ehrID, contribution}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return openehr.FOLDER{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	// Update EHR
	query = `UPDATE tbl_openehr_ehr 
            SET data = jsonb_set(data, '{directory}', to_jsonb($1::jsonb)) 
            WHERE id = $2`
	directoryRef := openehr.OBJECT_REF{
		Namespace: config.NAMESPACE_LOCAL,
		Type:      openehr.VERSIONED_FOLDER_MODEL_NAME,
		ID: openehr.X_OBJECT_ID{
			Value: &openehr.HIER_OBJECT_ID{
				Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
				Value: versionedDirectory.UID.Value,
			},
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

func (s *EHRService) GetRawDirectory(ctx context.Context, ehrID string) ([]byte, error) {
	// Fetch Directory by EHR ID
	query := `
		SELECT f.data 
		FROM tbl_openehr_folder f
		WHERE f.ehr_id = $1
		ORDER BY f.created_at DESC
		LIMIT 1
	`
	args := []any{ehrID}
	row := s.DB.QueryRow(ctx, query, args...)

	var rawDirectoryJSON []byte
	if err := row.Scan(&rawDirectoryJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrDirectoryNotFound
		}
		return nil, fmt.Errorf("failed to fetch Directory by EHR ID from database: %w", err)
	}

	return rawDirectoryJSON, nil
}

func (s *EHRService) UpdateDirectory(ctx context.Context, ehrID string, updatedDirectory openehr.FOLDER) error {
	// Validate updated Directory UID
	if errs := updatedDirectory.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for updated Directory UID: %v", errs)
	}

	if !updatedDirectory.UID.E {
		return fmt.Errorf("directory UID must be set for update")
	}

	// Check if Directory exists
	var currectDirectory openehr.FOLDER
	rawDirectoryJSON, err := s.GetRawDirectory(ctx, ehrID)
	if err != nil {
		return fmt.Errorf("failed to check if Directory exists: %w", err)
	}
	if err := json.Unmarshal(rawDirectoryJSON, &currectDirectory); err != nil {
		return fmt.Errorf("failed to unmarshal current Directory JSON: %w", err)
	}

	// Ensure that the versioned object ID matches
	if currectDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID() != updatedDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).UID() {
		return ErrDirectoryNotFound
	}

	// Ensure that the version is incremented
	currectVersion := currectDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).VersionTreeID()
	updatedVersion := updatedDirectory.UID.V.Value.(*openehr.OBJECT_VERSION_ID).VersionTreeID()
	if updatedVersion <= currectVersion {
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
	query := `INSERT INTO tbl_openehr_folder (id, versioned_object_id, ehr_id, data) 
		 VALUES ($1, $2, $3, $4)`
	args := []any{directoryVersion.UID.Value, directoryVersion.UID.UID(), ehrID, directoryVersion}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert updated directory into the database: %w", err)
	}

	// Insert Contribution
	query = `INSERT INTO tbl_openehr_contribution (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{contribution.UID.Value, ehrID, contribution}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	return nil
}

func (s *EHRService) DeleteDirectory(ctx context.Context, ehrID string) error {
	// Check if directory exists
	var versionedObjectID string
	err := s.DB.QueryRow(ctx, "SELECT id FROM tbl_openehr_versioned_object WHERE ehr_id = $1 AND type = $2", ehrID, openehr.FOLDER_MODEL_NAME).Scan(&versionedObjectID)
	if err != nil {
		if err == database.ErrNoRows {
			return ErrDirectoryNotFound
		}
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
					Value: &openehr.HIER_OBJECT_ID{
						Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
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
				Value: "Directory deleted",
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

	// Delete Directory
	query := `DELETE FROM tbl_openehr_versioned_object WHERE id = $1 AND ehr_id = $2 AND type = $3`
	args := []any{versionedObjectID, ehrID, openehr.FOLDER_MODEL_NAME}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to delete Directory from database: %w", err)
	}

	// Insert Contribution
	query = `INSERT INTO tbl_openehr_contribution (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{contribution.UID.Value, ehrID, contribution}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	// Remove reference from EHR
	query = `UPDATE tbl_openehr_ehr 
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

func (s *EHRService) GetFolderInDirectoryVersionAsJSON(ctx context.Context, ehrID string, filterAtTime time.Time, filterVersionID string) ([]byte, error) {
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`SELECT data FROM tbl_openehr_folder WHERE ehr_id = $1 `)
	args = []any{ehrID}
	argNum++

	if !filterAtTime.IsZero() {
		query.WriteString(fmt.Sprintf(` AND created_at <= $%d`, argNum))
		args = append(args, filterAtTime)
		argNum++
	}

	if filterVersionID != "" {
		query.WriteString(fmt.Sprintf(` AND id = $%d`, argNum))
		args = append(args, filterVersionID)
		argNum++
	}

	query.WriteString(` ORDER BY created_at DESC LIMIT 1`)
	row := s.DB.QueryRow(ctx, query.String(), args...)

	var rawFolderJSON []byte
	if err := row.Scan(&rawFolderJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrDirectoryNotFound
		}
		return nil, fmt.Errorf("failed to fetch Folder at time from database: %w", err)
	}

	return rawFolderJSON, nil
}

func (s *EHRService) CreateContribution(ctx context.Context, ehrID string, contribution openehr.CONTRIBUTION) (openehr.CONTRIBUTION, error) {
	if contribution.UID.Value == "" {
		contribution.UID.Value = uuid.NewString()
	}

	// Validate Contribution
	contribution.SetModelName()
	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.CONTRIBUTION{}, fmt.Errorf("validation errors for Contribution: %v", errs)
	}

	// Insert Contribution
	query := `INSERT INTO tbl_openehr_contribution (id, ehr_id, data) VALUES ($1, $2, $3)`
	args := []any{contribution.UID.Value, ehrID, contribution}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return openehr.CONTRIBUTION{}, fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	return contribution, nil
}

func (s *EHRService) GetContributionAsJSON(ctx context.Context, ehrID, contributionID string) ([]byte, error) {
	query := `
		SELECT c.data 
		FROM tbl_openehr_contribution c
		WHERE c.ehr_id = $1 AND c.id = $2
		LIMIT 1
	`
	args := []any{ehrID, contributionID}
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
