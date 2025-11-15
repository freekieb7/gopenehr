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
	ErrEHRNotFound            = fmt.Errorf("EHR not found")
	ErrEHRAlreadyExists       = fmt.Errorf("EHR already exists")
	ErrEHRStatusNotFound      = fmt.Errorf("EHR Status not found")
	ErrEHRStatusAlreadyExists = fmt.Errorf("EHR Status already exists")
)

type EHRService struct {
	Logger *slog.Logger
	DB     *database.Database
}

func (s *EHRService) CreateEHR(ctx context.Context, ehrStatus util.Optional[openehr.EHR_STATUS]) (openehr.EHR, error) {
	return s.CreateEHRWithID(ctx, ehrStatus, uuid.NewString())
}

func (s *EHRService) CreateEHRWithID(ctx context.Context, ehrStatus util.Optional[openehr.EHR_STATUS], ehrID string) (openehr.EHR, error) {
	// Check if EHR with the given ID already exists
	exists, err := s.CheckEHRExists(ctx, ehrID)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to check if EHR exists: %w", err)
	}
	if exists {
		return openehr.EHR{}, ErrEHRAlreadyExists
	}

	// Check if provided EHR Status UID already exists
	if ehrStatus.E {
		var ehrStatusID string
		switch v := ehrStatus.V.UID.V.Value.(type) {
		case *openehr.OBJECT_VERSION_ID:
			// OK
			ehrStatusID = v.Value
		default:
			return openehr.EHR{}, fmt.Errorf("EHR Status UID must be of type OBJECT_VERSION_ID, got %T", ehrStatus.V.UID.V)
		}

		// Check if EHR Status UID is already exists
		var exists bool
		if err := s.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM tbl_openehr_ehr_status WHERE id = $1)", ehrStatusID).Scan(&exists); err != nil {
			return openehr.EHR{}, fmt.Errorf("failed to check if EHR Status UID exists: %w", err)
		}
		if exists {
			return openehr.EHR{}, ErrEHRStatusAlreadyExists
		}
	}

	// Create EHR

	// 1. Prepare EHR Status
	// 1.1 Build Versioned EHR Status
	newVersionedEhrStatus := openehr.VERSIONED_EHR_STATUS{
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
	// If EHR Status is provided, use it; otherwise, create a default one
	if ehrStatus.E {
		newVersionedEhrStatus.UID.Value = strings.Split(ehrStatus.V.UID.V.Value.(*openehr.OBJECT_VERSION_ID).Value, "::")[0]
	}

	if errs := newVersionedEhrStatus.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new Versioned EHR Status: %v", errs)
	}

	// 1.2. Build EHR Status
	newEhrStatus := openehr.EHR_STATUS{
		UID: util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::gopenehr::1", newVersionedEhrStatus.UID.Value),
			},
		}),
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
	// If EHR Status is provided, use it
	if ehrStatus.E {
		newEhrStatus = ehrStatus.V
	}

	// 1.3. Build Original Version of EHR Status
	newEhrStatusVersion := openehr.ORIGINAL_VERSION{
		UID: *newEhrStatus.UID.V.Value.(*openehr.OBJECT_VERSION_ID),
		LifecycleState: openehr.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: openehr.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: openehr.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: &newEhrStatus,
	}

	if errs := newEhrStatusVersion.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new EHR Status Version: %v", errs)
	}

	// 2. Prepare EHR Access

	// 2.1 Build Versioned EHR Access
	newVersionedEhrAccess := openehr.VERSIONED_EHR_ACCESS{
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

	if errs := newVersionedEhrAccess.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new Versioned EHR Access: %v", errs)
	}

	// 2.2 Build EHR Access
	newEhrAccess := openehr.EHR_ACCESS{
		UID: util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::gopenehr::1", newVersionedEhrAccess.UID.Value),
			},
		}),
		Name: openehr.X_DV_TEXT{
			Value: &openehr.DV_TEXT{
				Value: "EHR Access",
			},
		},
		ArchetypeNodeID: "openEHR-EHR-EHR_ACCESS.generic.v1",
	}

	// 2.3 Build Original Version of EHR Access
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

	if errs := newEhrAccessVersion.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new EHR Access Version: %v", errs)
	}

	// 3. Prepare EHR
	newEhr := openehr.EHR{
		EHRID: openehr.HIER_OBJECT_ID{
			Value: ehrID,
		},
		EHRStatus: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.VERSIONED_EHR_STATUS_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: newVersionedEhrStatus.UID.Value,
				},
			},
		},
		EHRAccess: openehr.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      openehr.VERSIONED_EHR_ACCESS_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: newVersionedEhrAccess.UID.Value,
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	if errs := newEhr.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new EHR: %v", errs)
	}

	// 4. Prepare Contribution
	contribution := openehr.CONTRIBUTION{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: []openehr.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.EHR_STATUS_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &openehr.OBJECT_VERSION_ID{
						Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
						Value: newEhrStatusVersion.UID.Value,
					},
				},
			},
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      openehr.EHR_ACCESS_MODEL_NAME,
				ID: openehr.X_OBJECT_ID{
					Value: &openehr.OBJECT_VERSION_ID{
						Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
						Value: newEhrAccessVersion.UID.Value,
					},
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

	if errs := contribution.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new Contribution: %v", errs)
	}

	// Forces setting the model name in _type field
	newVersionedEhrAccess.SetModelName()
	newVersionedEhrStatus.SetModelName()
	newEhrAccessVersion.SetModelName()
	newEhrStatusVersion.SetModelName()
	newEhr.SetModelName()
	contribution.SetModelName()

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
	query = `INSERT INTO tbl_openehr_versioned_object (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{newVersionedEhrStatus.UID.Value, ehrID, newVersionedEhrStatus}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr status into the database: %w", err)
	}

	// Insert EHR Status Version
	query = `INSERT INTO tbl_openehr_ehr_status (id, versioned_object_id, ehr_id, data) 
         VALUES ($1, $2, $3, $4)`
	args = []any{newEhrStatusVersion.UID.Value, newVersionedEhrStatus.UID.Value, ehrID, newEhrStatusVersion}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr status into the database: %w", err)
	}

	// Insert Versioned EHR Access
	query = `INSERT INTO tbl_openehr_versioned_object (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{newVersionedEhrAccess.UID.Value, ehrID, newVersionedEhrAccess}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr access into the database: %w", err)
	}

	// Insert EHR Access Version
	// Store version metadata (without .data) and data separately
	query = `INSERT INTO tbl_openehr_ehr_access (id, versioned_object_id, ehr_id, data) 
         VALUES ($1, $2, $3, $4)`
	args = []any{newEhrAccessVersion.UID.Value, newVersionedEhrAccess.UID.Value, ehrID, newEhrAccessVersion}
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

func (s *EHRService) CheckEHRExists(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM tbl_openehr_ehr WHERE id = $1)`
	args := []any{id}
	var exists bool
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check if EHR exists in database: %w", err)
	}
	return exists, nil
}

func (s *EHRService) GetRawEHRByID(ctx context.Context, id string) ([]byte, error) {
	query := `SELECT data FROM tbl_openehr_ehr WHERE id = $1`
	args := []any{id}
	row := s.DB.QueryRow(ctx, query, args...)

	var rawEhrJSON []byte
	if err := row.Scan(&rawEhrJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR from database: %w", err)
	}

	return rawEhrJSON, nil
}

func (s *EHRService) GetRawEHRBySubject(ctx context.Context, subjectID, subjectNamespace string) ([]byte, error) {
	query := `
        SELECT e.data 
        FROM tbl_openehr_ehr e 
        JOIN (
			SELECT id, ehr_id, jsonb_path_query_first(data, '$.** ? (@._type == "EHR_STATUS")') data, created_at 
			FROM tbl_openehr_ehr_status
		) es ON e.id = es.ehr_id
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
	// Check if EHR exists
	exists, err := s.CheckEHRExists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if EHR exists: %w", err)
	}
	if !exists {
		return ErrEHRNotFound
	}

	// Delete EHR
	query := `DELETE FROM tbl_openehr_ehr WHERE id = $1`
	args := []any{id}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to delete EHR from database: %w", err)
	}
	return nil
}

func (s *EHRService) DeleteMultipleEHRs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	// Build query with variable number of placeholders
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

func (s *EHRService) GetRawEHRStatus(ctx context.Context, ehrID string, filterOnTime time.Time, filterOnVersionID string) ([]byte, error) {
	// Check if EHR exists
	exists, err := s.CheckEHRExists(ctx, ehrID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing EHR: %w", err)
	}
	if !exists {
		return nil, ErrEHRNotFound
	}

	// Build query
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`
		SELECT es.data 
		FROM (
			SELECT id, ehr_id, jsonb_path_query_first(data, '$.** ? (@._type == "EHR_STATUS")') data, created_at 
			FROM tbl_openehr_ehr_status
		) es 
		WHERE ehr_id = $1 
	`)
	args = []any{ehrID}
	argNum++

	if filterOnVersionID != "" {
		query.WriteString(fmt.Sprintf(`AND es.id = $%d `, argNum))
		args = append(args, filterOnVersionID)
		argNum++
	}

	if !filterOnTime.IsZero() {
		query.WriteString(fmt.Sprintf(`AND es.created_at <= $%d `, argNum))
		args = append(args, filterOnTime)
	}

	query.WriteString(`ORDER BY es.created_at DESC LIMIT 1`)
	row := s.DB.QueryRow(ctx, query.String(), args...)

	s.Logger.InfoContext(ctx, "get ehr status", "query", query.String())

	var rawEhrStatusJSON []byte
	if err := row.Scan(&rawEhrStatusJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR Status from database: %w", err)
	}

	return rawEhrStatusJSON, nil
}

func (s *EHRService) UpdateEHRStatus(ctx context.Context, ehrID string, newStatus openehr.EHR_STATUS) error {
	if !newStatus.UID.E {
		return fmt.Errorf("EHR Status UID must be set for update")
	}

	var ehrStatusId string
	switch v := newStatus.UID.V.Value.(type) {
	case *openehr.OBJECT_VERSION_ID:
		// OK
		ehrStatusId = v.Value
	default:
		return fmt.Errorf("EHR Status UID must be of type OBJECT_VERSION_ID, got %T", v)
	}

	if errs := newStatus.Validate("$"); len(errs) > 0 {
		return fmt.Errorf("validation errors for updated EHR Status: %v", errs)
	}

	// Extract EHR Status UID (without version)
	versionedEhrStatusId := strings.Split(ehrStatusId, "::")[0]

	// Check if EHR exists
	rawEHR, err := s.GetRawEHRByID(ctx, ehrID)
	if err != nil {
		if err == ErrEHRNotFound {
			return ErrEHRNotFound
		}

		return fmt.Errorf("failed to check existing EHR: %w", err)
	}

	var ehr openehr.EHR
	if err := json.Unmarshal(rawEHR, &ehr); err != nil {
		return fmt.Errorf("failed to unmarshal existing EHR JSON: %w", err)
	}

	// Check if EHR Status belongs to the EHR
	if versionedEhrStatusId != ehr.EHRStatus.ID.Value.(*openehr.HIER_OBJECT_ID).Value {
		return fmt.Errorf("EHR Status UID %s does not belong to EHR ID %s", versionedEhrStatusId, ehrID)
	}

	// Prepare Original Version of EHR Status
	newEhrStatusVersion := openehr.ORIGINAL_VERSION{
		UID: openehr.OBJECT_VERSION_ID{
			Value: ehrStatusId,
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
		Data: &newStatus,
	}
	newEhrStatusVersion.SetModelName()

	if errs := newEhrStatusVersion.Validate("$"); len(errs) > 0 {
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
					Value: &newEhrStatusVersion.UID,
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

	// Insert new EHR Status
	query := `INSERT INTO tbl_openehr_ehr_status (id, versioned_object_id, ehr_id, data) VALUES ($1, $2, $3, $4)`
	args := []any{newEhrStatusVersion.UID.Value, versionedEhrStatusId, ehrID, newEhrStatusVersion}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert ehr status into the database: %w", err)
	}

	// Insert Contribution
	query = `INSERT INTO tbl_openehr_contribution (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{contribution.UID.Value, ehrID, contribution}
	if _, err := s.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert contribution into the database: %w", err)
	}

	return nil
}

func (s *EHRService) GetRawVersionedEHRStatus(ctx context.Context, ehrID string) ([]byte, error) {
	// Check if EHR exists
	exists, err := s.CheckEHRExists(ctx, ehrID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing EHR: %w", err)
	}
	if !exists {
		return nil, ErrEHRNotFound
	}

	// Fetch Versioned EHR Status
	query := `SELECT vo.data FROM tbl_openehr_versioned_object vo WHERE vo.ehr_id = $1 AND vo.data->>'_type' = $2`
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

func (s *EHRService) GetRawVersionedEHRStatusRevisionHistory(ctx context.Context, ehrID string) ([]byte, error) {
	// Check if EHR exists
	exists, err := s.CheckEHRExists(ctx, ehrID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing EHR: %w", err)
	}
	if !exists {
		return nil, ErrEHRNotFound
	}

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

func (s *EHRService) GetRawVersionedEHRStatusVersionAtTime(ctx context.Context, ehrID string, atTime time.Time) ([]byte, error) {
	// Check if EHR exists
	exists, err := s.CheckEHRExists(ctx, ehrID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing EHR: %w", err)
	}
	if !exists {
		return nil, ErrEHRNotFound
	}

	// Fetch EHR Status at given time
	var query strings.Builder
	var args []any
	argNum := 1

	query.WriteString(`SELECT data FROM tbl_openehr_ehr_status WHERE ehr_id = $1 `)
	args = []any{ehrID}
	argNum++

	if !atTime.IsZero() {
		query.WriteString(fmt.Sprintf(` AND created_at <= $%d`, argNum))
		args = append(args, atTime)
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

func (s *EHRService) GetRawVersionedEHRStatusByVersionID(ctx context.Context, ehrID, versionID string) ([]byte, error) {
	// Check if EHR exists
	exists, err := s.CheckEHRExists(ctx, ehrID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing EHR: %w", err)
	}
	if !exists {
		return nil, ErrEHRNotFound
	}

	// Fetch EHR Status by version ID
	query := `SELECT data FROM tbl_openehr_ehr_status WHERE ehr_id = $1 AND id = $2`
	args := []any{ehrID, versionID}
	row := s.DB.QueryRow(ctx, query, args...)

	var rawEhrStatusJSON []byte
	if err := row.Scan(&rawEhrStatusJSON); err != nil {
		if err == database.ErrNoRows {
			return nil, ErrEHRNotFound
		}
		return nil, fmt.Errorf("failed to fetch EHR Status by version ID from database: %w", err)
	}

	return rawEhrStatusJSON, nil
}
