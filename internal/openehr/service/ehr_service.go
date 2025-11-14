package service

import (
	"context"
	"fmt"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/openehr"
	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/google/uuid"
)

type EHR struct {
	DB *database.Database
}

func (s *EHR) CreateEHR(ctx context.Context) (openehr.EHR, error) {
	return s.CreateEHRWithID(ctx, uuid.NewString())
}

func (s *EHR) CreateEHRWithID(ctx context.Context, id string) (openehr.EHR, error) {
	newEhr := openehr.EHR{
		EHRID: openehr.HIER_OBJECT_ID{
			Value: id,
		},
		EHRStatus: openehr.OBJECT_REF{
			Namespace: "local",
			Type:      openehr.VERSIONED_EHR_STATUS_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: uuid.New().String(),
				},
			},
		},
		EHRAccess: openehr.OBJECT_REF{
			Namespace: "local",
			Type:      openehr.VERSIONED_EHR_ACCESS_MODEL_NAME,
			ID: openehr.X_OBJECT_ID{
				Value: &openehr.HIER_OBJECT_ID{
					Type_: util.Some(openehr.HIER_OBJECT_ID_MODEL_NAME),
					Value: uuid.New().String(),
				},
			},
		},
		TimeCreated: openehr.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}

	if errs := newEhr.Validate("$"); len(errs) > 0 {
		return openehr.EHR{}, fmt.Errorf("ehr is invalid: %v", errs)
	}

	query := `INSERT INTO tbl_openehr_ehr (id, data) VALUES ($1, $2)`
	_, err := s.DB.Exec(ctx, query, newEhr.EHRID.Value, newEhr)
	if err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to insert ehr into the database: %w", err)
	}

	return newEhr, err
}

func (s *EHR) GetEHRByID(ctx context.Context, id string) (openehr.EHR, error) {
	query := `SELECT data FROM tbl_openehr_ehr WHERE id = $1`
	row := s.DB.QueryRow(ctx, query, id)

	var ehr openehr.EHR
	if err := row.Scan(&ehr); err != nil {
		return openehr.EHR{}, fmt.Errorf("failed to fetch EHR from database: %w", err)
	}

	return ehr, nil
}
