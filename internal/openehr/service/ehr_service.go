package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/openehr"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/google/uuid"
)

var (
	ErrEHRNotFound      = fmt.Errorf("EHR not found")
	ErrEHRAlreadyExists = fmt.Errorf("EHR already exists")
)

type EHR struct {
	Logger *slog.Logger
	DB     *database.Database
}

func (s *EHR) CreateEHR(ctx context.Context) (openehr.EHR, error) {
	return s.CreateEHRWithID(ctx, uuid.NewString())
}

func (s *EHR) CreateEHRWithID(ctx context.Context, ehrID string) (openehr.EHR, error) {
	// Check if EHR with the given ID already exists
	existingEHR, err := s.GetEHRByID(ctx, ehrID)
	if err == nil && existingEHR.EHRID.Value != "" {
		return openehr.EHR{}, ErrEHRAlreadyExists
	}

	ehrStatusUid := uuid.New().String()
	ehrStatusVersionUid := fmt.Sprintf("%s::gopenehr::1", ehrStatusUid)
	ehrAccessUid := uuid.New().String()
	ehrAccessVersionUid := fmt.Sprintf("%s::gopenehr::1", ehrAccessUid)

	// Create Versioned EHR Status
	newVersionedEhrStatus := openehr.VERSIONED_EHR_STATUS{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.New().String(),
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: "local",
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

	if errs := newVersionedEhrStatus.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new Versioned EHR Status: %v", errs)
	}

	// Create Versioned EHR Access
	newVersionedEhrAccess := openehr.VERSIONED_EHR_ACCESS{
		UID: openehr.HIER_OBJECT_ID{
			Value: uuid.New().String(),
		},
		OwnerID: openehr.OBJECT_REF{
			Namespace: "local",
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

	// Create EHR Status
	newEhrStatus := openehr.EHR_STATUS{
		UID: util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: ehrStatusVersionUid,
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

	if errs := newEhrStatus.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new EHR Status: %v", errs)
	}

	// Create EHR Access
	newEhrAccess := openehr.EHR_ACCESS{
		UID: util.Some(openehr.X_UID_BASED_ID{
			Value: &openehr.OBJECT_VERSION_ID{
				Type_: util.Some(openehr.OBJECT_VERSION_ID_MODEL_NAME),
				Value: ehrAccessVersionUid,
			},
		}),
		Name: openehr.X_DV_TEXT{
			Value: &openehr.DV_TEXT{
				Value: "EHR Access",
			},
		},
		ArchetypeNodeID: "openEHR-EHR-EHR_ACCESS.generic.v1",
	}

	if errs := newEhrAccess.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("validation errors for new EHR Access: %v", errs)
	}

	// Create new EHR
	newEhr := openehr.EHR{
		EHRID: openehr.HIER_OBJECT_ID{
			Value: ehrID,
		},
		EHRStatus: openehr.OBJECT_REF{
			Namespace: "local",
			Type:      openehr.VERSIONED_EHR_STATUS_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: ehrStatusUid,
				},
			},
		},
		EHRAccess: openehr.OBJECT_REF{
			Namespace: "local",
			Type:      openehr.VERSIONED_EHR_ACCESS_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: ehrAccessUid,
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

	// Forces setting the model name in _type field
	newEhr.SetModelName()
	newEhrStatus.SetModelName()
	newEhrAccess.SetModelName()

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
	args = []any{ehrStatusUid, ehrID, newVersionedEhrStatus}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr status into the database: %w", err)
	}

	// Insert Versioned EHR Access
	query = `INSERT INTO tbl_openehr_versioned_object (id, ehr_id, data) VALUES ($1, $2, $3)`
	args = []any{ehrAccessUid, ehrID, newVersionedEhrAccess}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert versioned ehr access into the database: %w", err)
	}

	// Insert EHR Status
	query = `INSERT INTO tbl_openehr_ehr_status (id, versioned_object_id, ehr_id, data) VALUES ($1, $2, $3, $4)`
	args = []any{ehrStatusVersionUid, ehrStatusUid, ehrID, newEhrStatus}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr status into the database: %w", err)
	}

	// Insert EHR Access
	query = `INSERT INTO tbl_openehr_ehr_access (id, versioned_object_id, ehr_id, data) VALUES ($1, $2, $3, $4)`
	args = []any{ehrAccessVersionUid, ehrAccessUid, ehrID, newEhrAccess}
	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr access into the database: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newEhr, nil
}

func (s *EHR) GetEHRByID(ctx context.Context, id string) (openehr.EHR, error) {
	query := `SELECT data FROM tbl_openehr_ehr WHERE id = $1`
	args := []any{id}
	row := s.DB.QueryRow(ctx, query, args...)

	var ehr openehr.EHR
	if err := row.Scan(&ehr); err != nil {
		if err == database.ErrNoRows {
			return openehr.EHR{}, ErrEHRNotFound
		}
		return openehr.EHR{}, fmt.Errorf("failed to fetch EHR from database: %w", err)
	}

	return ehr, nil
}

func (s *EHR) GetEHRBySubject(ctx context.Context, subjectID, subjectNamespace string) (openehr.EHR, error) {
	query := `
		SELECT e.data FROM tbl_openehr_ehr e 
		JOIN tbl_openehr_ehr_status es ON e.id = es.ehr_id
		WHERE es.data->'subject'->'external_ref'->'id'->>'value' = $1 AND es.data->'subject'->'external_ref'->'id'->>'namespace' = $2
	`
	args := []any{subjectID, subjectNamespace}
	row := s.DB.QueryRow(ctx, query, args...)

	var ehr openehr.EHR
	if err := row.Scan(&ehr); err != nil {
		if err == database.ErrNoRows {
			return openehr.EHR{}, ErrEHRNotFound
		}
		return openehr.EHR{}, fmt.Errorf("failed to fetch EHR by subject from database: %w", err)
	}

	return ehr, nil
}

func (s *EHR) DeleteEHRByID(ctx context.Context, id string) error {
	// Check if EHR exists
	_, err := s.GetEHRByID(ctx, id)
	if err != nil {
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

func (s *EHR) DeleteMultipleEHRs(ctx context.Context, ids []string) error {
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
