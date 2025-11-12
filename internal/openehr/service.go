package openehr

import (
	"context"
	"fmt"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/google/uuid"
)

type Service struct {
	DB *database.Database
}

func (s *Service) CreateEHR(ctx context.Context) (EHR, error) {
	return s.CreateEHRWithID(ctx, uuid.NewString())
}

func (s *Service) CreateEHRWithID(ctx context.Context, id string) (EHR, error) {
	newEhr := EHR{
		EHRID: HIER_OBJECT_ID{
			Value: id,
		},
		EHRStatus: OBJECT_REF{
			Namespace: "local",
			Type:      VERSIONED_EHR_STATUS_MODEL_NAME,
			ID: X_OBJECT_ID{
				Value: &HIER_OBJECT_ID{
					Type_: Some(HIER_OBJECT_ID_MODEL_NAME),
					Value: uuid.New().String(),
				},
			},
		},
		EHRAccess: OBJECT_REF{
			Namespace: "local",
			Type:      VERSIONED_EHR_ACCESS_MODEL_NAME,
			ID: X_OBJECT_ID{
				Value: &HIER_OBJECT_ID{
					Type_: Some(HIER_OBJECT_ID_MODEL_NAME),
					Value: uuid.New().String(),
				},
			},
		},
		TimeCreated: DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	if errs := newEhr.Validate("$"); len(errs) > 0 {
		return EHR{}, fmt.Errorf("ehr is invalid: %v", errs)
	}

	query := `INSERT INTO tbl_openehr_ehr (id, data) VALUES ($1, $2)`
	_, err := s.DB.Exec(ctx, query, newEhr.EHRID.Value, newEhr)
	if err != nil {
		return EHR{}, fmt.Errorf("failed to insert ehr into the database: %w", err)
	}

	return newEhr, err
}

func (s *Service) GetEHRByID(ctx context.Context, id string) (EHR, error) {
	query := `SELECT data FROM tbl_openehr_ehr WHERE id = $1`
	row := s.DB.QueryRow(ctx, query, id)

	var ehr EHR
	if err := row.Scan(&ehr); err != nil {
		return EHR{}, fmt.Errorf("failed to fetch EHR from database: %w", err)
	}

	return ehr, nil
}
