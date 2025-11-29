package openehr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/openehr/aql"
	"github.com/freekieb7/gopenehr/internal/openehr/model"
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

	ErrVersionLowerOrEqualToCurrent = fmt.Errorf("EHR Status version must be incremented")

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

func NewService(logger *telemetry.Logger, db *database.Database) Service {
	return Service{
		Logger: logger,
		DB:     db,
	}
}

func (s *Service) CreateEHR(ctx context.Context, ehrID uuid.UUID, ehrStatus model.EHR_STATUS) (model.EHR, error) {
	err := s.ValidateEHRStatus(ctx, ehrStatus)
	if err != nil {
		return model.EHR{}, err
	}

	err = UpgradeObjectVersionID(&ehrStatus.UID, utils.None[model.OBJECT_VERSION_ID]())
	if err != nil {
		return model.EHR{}, fmt.Errorf("failed to upgrade EHR Status UID: %w", err)
	}

	versionedEHRStatus := NewVersionedEHRStatus(ehrStatus.UID.V.Value.(*model.OBJECT_VERSION_ID).UID(), ehrID)
	versionedEHRAccess := NewVersionedEHRAccess(uuid.New(), ehrID)
	ehrAccess := NewEHRAccess(uuid.MustParse(versionedEHRAccess.UID.Value))
	ehrStatusVersion := NewOriginalVersion(&ehrStatus, utils.None[model.OBJECT_VERSION_ID]())
	ehrAccessVersion := NewOriginalVersion(&ehrAccess, utils.None[model.OBJECT_VERSION_ID]())
	contribution := NewContribution("EHR created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]model.OBJECT_REF{
			ehrStatusVersion.ObjectRef(),
			ehrAccessVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.EHR{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveEHRWithTx(ctx, tx, ehrID)
	if err != nil {
		return model.EHR{}, fmt.Errorf("failed to save EHR: %w", err)
	}
	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return model.EHR{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedEHRStatus, utils.Some(ehrID))
	if err != nil {
		return model.EHR{}, fmt.Errorf("failed to save VERSIONED_EHR_STATUS: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, ehrStatusVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return model.EHR{}, fmt.Errorf("failed to save EHR_STATUS: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedEHRAccess, utils.Some(ehrID))
	if err != nil {
		return model.EHR{}, fmt.Errorf("failed to save VERSIONED_EHR_ACCESS: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, ehrAccessVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return model.EHR{}, fmt.Errorf("failed to save EHR_ACCESS: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return model.EHR{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	ehr, err := s.GetEHR(ctx, ehrID)
	if err != nil {
		return model.EHR{}, fmt.Errorf("failed to get EHR after creation: %w", err)
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

func (s *Service) GetEHR(ctx context.Context, id uuid.UUID) (model.EHR, error) {
	query := `SELECT data FROM openehr.vw_ehr WHERE id = $1 LIMIT 1`
	args := []any{id}

	row := s.DB.QueryRow(ctx, query, args...)

	var ehr model.EHR
	err := row.Scan(&ehr)
	if err != nil {
		if err == database.ErrNoRows {
			return model.EHR{}, ErrEHRNotFound
		}
		return model.EHR{}, fmt.Errorf("failed to fetch EHR from database: %w", err)
	}

	return ehr, nil
}

func (s *Service) GetEHRBySubject(ctx context.Context, subjectID, subjectNamespace string) (model.EHR, error) {
	query := `
        SELECT ehr_id
        FROM openehr.tbl_object_version_data
        WHERE ov.type = $1
          AND object_data->'subject'->'external_ref'->>'namespace' = $2
		  AND object_data->'subject'->'external_ref'->'id'->>'value' = $3
        ORDER BY created_at DESC
        LIMIT 1
    `
	args := []any{model.EHR_STATUS_MODEL_NAME, subjectNamespace, subjectID}

	row := s.DB.QueryRow(ctx, query, args...)

	var ehrID uuid.UUID
	err := row.Scan(&ehrID)
	if err != nil {
		if err == database.ErrNoRows {
			return model.EHR{}, ErrEHRNotFound
		}
		return model.EHR{}, fmt.Errorf("failed to fetch EHR by subject from database: %w", err)
	}

	// We could get the EHR directory from the database, but this is cache friendly
	return s.GetEHR(ctx, ehrID)
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

func (s *Service) ValidateEHRStatus(ctx context.Context, ehrStatus model.EHR_STATUS) error {
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
			case *model.HIER_OBJECT_ID:
				// Must be a valid type
				if externalRef.Type != model.VERSIONED_PARTY_MODEL_NAME {
					validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
						Model:          model.EHR_STATUS_MODEL_NAME,
						Path:           attrPath + ".type",
						Message:        fmt.Sprintf("invalid subject external_ref type: %s", externalRef.Type),
						Recommendation: "Ensure external ref type is VERSIONED_PARTY",
					})
				}

				// Must be a valid UUID
				if err := uuid.Validate(v.Value); err != nil {
					validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
						Model:          model.EHR_STATUS_MODEL_NAME,
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
							Model:          model.EHR_STATUS_MODEL_NAME,
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
					Model:          model.EHR_STATUS_MODEL_NAME,
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

func (s *Service) GetEHRStatus(ctx context.Context, ehrID uuid.UUID, filterOnTime time.Time, filterOnVersionID string) (model.EHR_STATUS, error) {
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
	args = []any{model.EHR_STATUS_MODEL_NAME, ehrID}
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

	var ehrStatus model.EHR_STATUS
	err := row.Scan(&ehrStatus)
	if err != nil {
		if err == database.ErrNoRows {
			return model.EHR_STATUS{}, ErrEHRStatusNotFound
		}
		return model.EHR_STATUS{}, fmt.Errorf("failed to fetch EHR Status from database: %w", err)
	}

	return ehrStatus, nil
}

func (s *Service) UpdateEHRStatus(ctx context.Context, ehrID uuid.UUID, ehrStatus model.EHR_STATUS) (model.EHR_STATUS, error) {
	if err := s.ValidateEHRStatus(ctx, ehrStatus); err != nil {
		return model.EHR_STATUS{}, err
	}

	currentEHRStatus, err := s.GetEHRStatus(ctx, ehrID, time.Time{}, "")
	if err != nil {
		return model.EHR_STATUS{}, fmt.Errorf("failed to get current EHR Status: %w", err)
	}
	currentEHRStatusID := currentEHRStatus.UID.V.Value.(*model.OBJECT_VERSION_ID)

	err = UpgradeObjectVersionID(&currentEHRStatus.UID, utils.Some(*currentEHRStatusID))
	if err != nil {
		return model.EHR_STATUS{}, fmt.Errorf("failed to upgrade current EHR Status UID: %w", err)
	}

	ehrStatusVersion := NewOriginalVersion(&ehrStatus, utils.Some(*currentEHRStatusID))
	contribution := NewContribution("EHR Status updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]model.OBJECT_REF{
			ehrStatusVersion.ObjectRef(),
		},
	)

	// Start transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.EHR_STATUS{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return model.EHR_STATUS{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, ehrStatusVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return model.EHR_STATUS{}, fmt.Errorf("failed to save ehr status version: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return model.EHR_STATUS{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ehrStatus, nil
}

func (s *Service) GetVersionedEHRStatus(ctx context.Context, ehrID uuid.UUID) (model.VERSIONED_EHR_STATUS, error) {
	query := `
		SELECT vod.data 
		FROM openehr.tbl_versioned_object vo
		JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id
		WHERE vo.type = $1 
		  AND vo.ehr_id = $2
		LIMIT 1
	`
	args := []any{model.VERSIONED_EHR_STATUS_MODEL_NAME, ehrID}

	row := s.DB.QueryRow(ctx, query, args...)

	var versionedEHRStatus model.VERSIONED_EHR_STATUS
	err := row.Scan(&versionedEHRStatus)
	if err != nil {
		if err == database.ErrNoRows {
			return model.VERSIONED_EHR_STATUS{}, ErrEHRNotFound
		}
		return model.VERSIONED_EHR_STATUS{}, fmt.Errorf("failed to fetch Versioned EHR Status from database: %w", err)
	}

	return versionedEHRStatus, nil
}

func (s *Service) GetVersionedEHRStatusRevisionHistory(ctx context.Context, ehrID uuid.UUID) (model.REVISION_HISTORY, error) {
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

	var revisionHistory model.REVISION_HISTORY
	err := row.Scan(&revisionHistory)
	if err != nil {
		if err == database.ErrNoRows {
			return model.REVISION_HISTORY{}, ErrEHRNotFound
		}
		return model.REVISION_HISTORY{}, fmt.Errorf("failed to fetch Revision History from database: %w", err)
	}

	return revisionHistory, nil
}

func (s *Service) GetVersionedEHRStatusVersionAsJSON(ctx context.Context, ehrID uuid.UUID, filterAtTime time.Time, filterVersionID string) ([]byte, error) {
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
	args = []any{model.EHR_STATUS_MODEL_NAME, ehrID}
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

func (s *Service) ValidateComposition(ctx context.Context, composition model.COMPOSITION) error {
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
	args := []any{ehrID, model.COMPOSITION_MODEL_NAME, versionID}

	var exists int
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		if err == database.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if Composition exists in database: %w", err)
	}

	return true, nil
}

func (s *Service) CreateComposition(ctx context.Context, ehrID uuid.UUID, composition model.COMPOSITION) (model.COMPOSITION, error) {
	err := s.ValidateComposition(ctx, composition)
	if err != nil {
		return model.COMPOSITION{}, err
	}

	err = UpgradeObjectVersionID(&composition.UID, utils.None[model.OBJECT_VERSION_ID]())
	if err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to upgrade composition UID: %w", err)
	}

	exists, err := s.ExistsComposition(ctx, ehrID, composition.UID.V.ValueAsString())
	if err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to check if composition exists: %w", err)
	}
	if exists {
		return model.COMPOSITION{}, ErrCompositionAlreadyExists
	}

	versionedComposition := NewVersionedComposition(composition.UID.V.Value.(*model.OBJECT_VERSION_ID).UID(), ehrID)
	compositionVersion := NewOriginalVersion(&composition, utils.None[model.OBJECT_VERSION_ID]())
	contribution := NewContribution("Composition created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]model.OBJECT_REF{
			compositionVersion.ObjectRef(),
		},
	)

	// Begin transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	// Insert Contribution
	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to save contribution: %w", err)
	}

	// Insert Versioned Composition
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedComposition, utils.Some(ehrID))
	if err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to save versioned composition: %w", err)
	}

	// Insert Composition Version
	err = s.SaveObjectVersionWithTx(ctx, tx, compositionVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to save composition version: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return composition, nil
}

func (s *Service) GetComposition(ctx context.Context, ehrID uuid.UUID, uidBasedID string) (model.COMPOSITION, error) {
	query := `
		SELECT ovd.object_data 
		FROM openehr.tbl_object_version ov 
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.type = $1 AND ov.ehr_id = $2
	`
	args := []any{model.COMPOSITION_MODEL_NAME, ehrID}

	if strings.Count(uidBasedID, "::") == 2 {
		// UID is of type OBJECT_VERSION_ID
		query += `AND ov.id = $3 `
	} else {
		// UID is of type HIER_OBJECT_ID
		query += `AND ov.versioned_object_id = $3 `
	}
	args = append(args, uidBasedID)

	query += `ORDER BY ov.created_at DESC LIMIT 1`

	var composition model.COMPOSITION
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&composition); err != nil {
		if err == database.ErrNoRows {
			return model.COMPOSITION{}, ErrCompositionNotFound
		}
		return model.COMPOSITION{}, fmt.Errorf("failed to fetch Composition by ID from database: %w", err)
	}

	return composition, nil
}

func (s *Service) UpdateComposition(ctx context.Context, ehrID uuid.UUID, composition model.COMPOSITION) (model.COMPOSITION, error) {
	if !composition.UID.E {
		return model.COMPOSITION{}, ErrCompositionUIDNotProvided
	}

	currentComposition, err := s.GetComposition(ctx, ehrID, composition.UID.V.ValueAsString())
	if err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to get current Composition: %w", err)
	}
	currentCompositionID := currentComposition.UID.V.Value.(*model.OBJECT_VERSION_ID)

	compositionVersion := NewOriginalVersion(&composition, utils.Some(*currentCompositionID))
	contribution := NewContribution("Composition updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]model.OBJECT_REF{
			compositionVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, compositionVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to save composition version: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.COMPOSITION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return composition, nil
}

func (s *Service) DeleteComposition(ctx context.Context, ehrID uuid.UUID, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Composition deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]model.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      model.COMPOSITION_MODEL_NAME,
				ID: model.X_OBJECT_ID{
					Value: &model.HIER_OBJECT_ID{
						Value: versionedObjectID.String(),
					},
				},
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.DeleteVersionedObjectWithTx(ctx, tx, versionedObjectID)
	if err != nil {
		return fmt.Errorf("failed to delete composition version: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) GetVersionedComposition(ctx context.Context, ehrID uuid.UUID, versionedObjectID string) (model.VERSIONED_COMPOSITION, error) {
	query := `
		SELECT vod.data 
		FROM openehr.tbl_versioned_object vo
		JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id
		WHERE vo.type = $1
		  AND vo.ehr_id = $2 
		  AND vo.id = $3
		LIMIT 1
	`
	args := []any{model.VERSIONED_COMPOSITION_MODEL_NAME, ehrID, versionedObjectID}
	row := s.DB.QueryRow(ctx, query, args...)

	var versionedComposition model.VERSIONED_COMPOSITION
	if err := row.Scan(&versionedComposition); err != nil {
		if err == database.ErrNoRows {
			return model.VERSIONED_COMPOSITION{}, ErrCompositionNotFound
		}
		return model.VERSIONED_COMPOSITION{}, fmt.Errorf("failed to fetch Versioned Composition by ID from database: %w", err)
	}

	return versionedComposition, nil
}

func (s *Service) GetVersionedCompositionRevisionHistory(ctx context.Context, ehrID uuid.UUID, versionedObjectID string) (model.REVISION_HISTORY, error) {
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

	var revisionHistory model.REVISION_HISTORY
	err := row.Scan(&revisionHistory)
	if err != nil {
		if err == database.ErrNoRows {
			return model.REVISION_HISTORY{}, ErrCompositionNotFound
		}
		return model.REVISION_HISTORY{}, fmt.Errorf("failed to fetch Revision History from database: %w", err)
	}

	return revisionHistory, nil
}

func (s *Service) GetVersionedCompositionVersionJSON(ctx context.Context, ehrID uuid.UUID, versionedObjectID string, filterAtTime time.Time, filterVersionID string) ([]byte, error) {
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
	args = []any{model.COMPOSITION_MODEL_NAME, ehrID, versionedObjectID}
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

func (s *Service) ValidateDirectory(ctx context.Context, ehrID uuid.UUID, directory model.FOLDER) error {
	validateErr := directory.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Directory validation can be added here
	folderQueue := make([]model.FOLDER, 0)
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
			case model.COMPOSITION_MODEL_NAME:
				id, ok := currentRef.ID.Value.(*model.OBJECT_VERSION_ID)
				if !ok {
					validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
						Model:          model.COMPOSITION_MODEL_NAME,
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
						Model:          model.COMPOSITION_MODEL_NAME,
						Path:           itemPath,
						Message:        "COMPOSITION does not exist for this EHR in the system",
						Recommendation: "Ensure the composition exists for this EHR",
					})
				}
			case model.VERSIONED_COMPOSITION_MODEL_NAME:
				id, ok := currentRef.ID.Value.(*model.HIER_OBJECT_ID)
				if !ok {
					validateErr.Errs = append(validateErr.Errs, outil.ValidationError{
						Model:          model.VERSIONED_COMPOSITION_MODEL_NAME,
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
						Model:          model.COMPOSITION_MODEL_NAME,
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

func (s *Service) CreateDirectory(ctx context.Context, ehrID uuid.UUID, directory model.FOLDER) (model.FOLDER, error) {
	err := s.ValidateDirectory(ctx, ehrID, directory)
	if err != nil {
		return model.FOLDER{}, err
	}

	exists, err := s.ExistsDirectory(ctx, ehrID)
	if err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to check if Directory exists: %w", err)
	}
	if exists {
		return model.FOLDER{}, ErrDirectoryAlreadyExists
	}

	// Upgrade Directory UID to OBJECT_VERSION_ID
	err = UpgradeObjectVersionID(&directory.UID, utils.None[model.OBJECT_VERSION_ID]())
	if err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to upgrade directory UID: %w", err)
	}
	versionedFolder := NewVersionedFolder(directory.UID.V.Value.(*model.OBJECT_VERSION_ID).UID(), ehrID)
	folderVersion := NewOriginalVersion(&directory, utils.None[model.OBJECT_VERSION_ID]())
	contribution := NewContribution("Directory created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]model.OBJECT_REF{
			folderVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedFolder, utils.Some(ehrID))
	if err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to save versioned folder: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, folderVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to save folder version: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return directory, nil
}

func (s *Service) ExistsDirectory(ctx context.Context, ehrID uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.ehr_id = $1 AND ov.type = $2 LIMIT 1`
	args := []any{ehrID, model.FOLDER_MODEL_NAME}

	var exists int
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		if err == database.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if Directory exists in database: %w", err)
	}

	return true, nil
}

func (s *Service) GetDirectory(ctx context.Context, ehrID uuid.UUID) (model.FOLDER, error) {
	query := `
		SELECT ovd.object_data
        FROM openehr.tbl_object_version ov
        JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id
        WHERE ov.type = $1
          AND ov.ehr_id = $2
        ORDER BY ov.created_at DESC
        LIMIT 1
	`
	args := []any{model.FOLDER_MODEL_NAME, ehrID}
	row := s.DB.QueryRow(ctx, query, args...)

	var directory model.FOLDER
	if err := row.Scan(&directory); err != nil {
		if err == database.ErrNoRows {
			return model.FOLDER{}, ErrDirectoryNotFound
		}
		return model.FOLDER{}, fmt.Errorf("failed to fetch Directory by EHR ID from database: %w", err)
	}

	return directory, nil
}

func (s *Service) UpdateDirectory(ctx context.Context, ehrID uuid.UUID, directory model.FOLDER) (model.FOLDER, error) {
	err := s.ValidateDirectory(ctx, ehrID, directory)
	if err != nil {
		return model.FOLDER{}, err
	}

	currentDirectory, err := s.GetDirectory(ctx, ehrID)
	if err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to get current Directory: %w", err)
	}
	currentDirectoryID := currentDirectory.UID.V.Value.(*model.OBJECT_VERSION_ID)

	err = UpgradeObjectVersionID(&currentDirectory.UID, utils.Some(*currentDirectoryID))
	if err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to upgrade current Directory UID: %w", err)
	}

	folderVersion := NewOriginalVersion(&directory, utils.Some(*currentDirectoryID))
	contribution := NewContribution("Directory updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]model.OBJECT_REF{
			folderVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", rbErr)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, folderVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to save folder version: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.FOLDER{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return directory, nil
}

func (s *Service) DeleteDirectory(ctx context.Context, ehrID uuid.UUID, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Directory deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]model.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      model.FOLDER_MODEL_NAME,
				ID: model.X_OBJECT_ID{
					Value: &model.HIER_OBJECT_ID{
						Value: versionedObjectID.String(),
					},
				},
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.DeleteVersionedObjectWithTx(ctx, tx, versionedObjectID)
	if err != nil {
		return fmt.Errorf("failed to delete directory version: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) GetFolderInDirectoryVersion(ctx context.Context, ehrID uuid.UUID, filterAtTime time.Time, filterVersionID string, filterPathParts []string) (model.FOLDER, error) {
	var queryBuilder strings.Builder
	var args []any
	argNum := 1

	jsonPath := "$"
	for _, part := range filterPathParts {
		jsonPath += fmt.Sprintf(`.folders ? (@.name.value == "%s")`, part)
	}

	queryBuilder.WriteString(fmt.Sprintf(`
		SELECT jsonb_path_query_first(ovd.data, '%s') 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id
		WHERE ov.type = $1 
		  AND ov.ehr_id = $2
		  AND ovd.object_data @? $3
	`, jsonPath))
	args = []any{model.FOLDER_MODEL_NAME, ehrID, jsonPath}
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

	var folder model.FOLDER
	if err := row.Scan(&folder); err != nil {
		if err == database.ErrNoRows {
			if len(filterPathParts) > 0 {
				return model.FOLDER{}, ErrFolderNotFoundInDirectory
			}
			return model.FOLDER{}, ErrDirectoryNotFound
		}
		return model.FOLDER{}, fmt.Errorf("failed to fetch Folder at time from database: %w", err)
	}

	return folder, nil
}

func (s *Service) ValidateAgent(ctx context.Context, agent model.AGENT) error {
	validateErr := agent.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Agent validation can be added here

	return nil
}

func (s *Service) ExistsAgent(ctx context.Context, versionID string) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.type = $1 AND ov.id = $2 LIMIT 1`

	var exists int
	err := s.DB.QueryRow(ctx, query, model.AGENT_MODEL_NAME, versionID).Scan(&exists)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if agent exists: %w", err)
	}
	return true, nil
}

func (s *Service) CreateAgent(ctx context.Context, agent model.AGENT) (model.AGENT, error) {
	err := s.ValidateAgent(ctx, agent)
	if err != nil {
		return model.AGENT{}, err
	}

	err = UpgradeObjectVersionID(&agent.UID, utils.None[model.OBJECT_VERSION_ID]())
	if err != nil {
		return model.AGENT{}, fmt.Errorf("failed to upgrade agent UID: %w", err)
	}

	exists, err := s.ExistsAgent(ctx, agent.UID.V.ValueAsString())
	if err != nil {
		return model.AGENT{}, fmt.Errorf("failed to check existing agent: %w", err)
	}
	if exists {
		return model.AGENT{}, ErrAgentAlreadyExists
	}

	versionedParty := NewVersionedParty(agent.UID.V.Value.(*model.OBJECT_VERSION_ID).UID())
	agentVersion := NewOriginalVersion(&agent, utils.None[model.OBJECT_VERSION_ID]())
	contribution := NewContribution("Agent created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]model.OBJECT_REF{
			agentVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.AGENT{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return model.AGENT{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedParty, utils.None[uuid.UUID]())
	if err != nil {
		return model.AGENT{}, fmt.Errorf("failed to save versioned party: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, agentVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return model.AGENT{}, fmt.Errorf("failed to save agent: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.AGENT{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return agent, nil
}

func (s *Service) GetCurrentAgentVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (model.AGENT, error) {
	query := `
		SELECT ovd.object_data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var agent model.AGENT
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), model.AGENT_MODEL_NAME).Scan(&agent)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return model.AGENT{}, ErrAgentNotFound
		}
		return model.AGENT{}, fmt.Errorf("failed to get latest agent by versioned party ID: %w", err)
	}

	return agent, nil
}

func (s *Service) GetAgentAtVersion(ctx context.Context, versionID string) (model.AGENT, error) {
	query := `
		SELECT ovd.object_data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var agent model.AGENT
	err := s.DB.QueryRow(ctx, query, versionID, model.AGENT_MODEL_NAME).Scan(&agent)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return model.AGENT{}, ErrAgentNotFound
		}
		return model.AGENT{}, fmt.Errorf("failed to get agent at version: %w", err)
	}

	return agent, nil
}

func (s *Service) UpdateAgent(ctx context.Context, versionedPartyID uuid.UUID, agent model.AGENT) (model.AGENT, error) {
	err := s.ValidateAgent(ctx, agent)
	if err != nil {
		return model.AGENT{}, err
	}

	currentAgent, err := s.GetCurrentAgentVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if errors.Is(err, ErrAgentNotFound) {
			return model.AGENT{}, ErrAgentNotFound
		}
		return model.AGENT{}, fmt.Errorf("failed to get current Agent: %w", err)
	}
	currentAgentID := currentAgent.UID.V.Value.(*model.OBJECT_VERSION_ID)

	err = UpgradeObjectVersionID(&currentAgent.UID, utils.Some(*currentAgentID))
	if err != nil {
		return model.AGENT{}, fmt.Errorf("failed to upgrade current Agent UID: %w", err)
	}

	agentVersion := NewOriginalVersion(&agent, utils.Some(*currentAgentID))
	contribution := NewContribution("Agent updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]model.OBJECT_REF{
			agentVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.AGENT{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return model.AGENT{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, agentVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return model.AGENT{}, fmt.Errorf("failed to save agent: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.AGENT{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return agent, nil
}

func (s *Service) DeleteAgent(ctx context.Context, versionedObjectID string) error {
	contribution := NewContribution("Agent deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]model.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      model.AGENT_MODEL_NAME,
				ID: model.X_OBJECT_ID{
					Value: &model.OBJECT_VERSION_ID{
						Type_: utils.Some(model.OBJECT_VERSION_ID_MODEL_NAME),
						Value: versionedObjectID,
					},
				},
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.DeleteVersionedObjectWithTx(ctx, tx, uuid.MustParse(strings.Split(versionedObjectID, "::")[0]))
	if err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) ValidatePerson(ctx context.Context, person model.PERSON) error {
	validateErr := person.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Person validation can be added here

	return nil
}

func (s *Service) ExistsPerson(ctx context.Context, versionID string) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.type = $1 AND ov.id = $2 LIMIT 1`

	var exists int
	err := s.DB.QueryRow(ctx, query, model.PERSON_MODEL_NAME, versionID).Scan(&exists)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if person exists: %w", err)
	}
	return true, nil
}

func (s *Service) CreatePerson(ctx context.Context, person model.PERSON) (model.PERSON, error) {
	err := s.ValidatePerson(ctx, person)
	if err != nil {
		return model.PERSON{}, err
	}

	err = UpgradeObjectVersionID(&person.UID, utils.None[model.OBJECT_VERSION_ID]())
	if err != nil {
		return model.PERSON{}, fmt.Errorf("failed to upgrade person UID: %w", err)
	}

	exists, err := s.ExistsPerson(ctx, person.UID.V.ValueAsString())
	if err != nil {
		return model.PERSON{}, fmt.Errorf("failed to check existing person: %w", err)
	}
	if exists {
		return model.PERSON{}, ErrPersonAlreadyExists
	}

	versionedParty := NewVersionedParty(person.UID.V.Value.(*model.OBJECT_VERSION_ID).UID())
	personVersion := NewOriginalVersion(&person, utils.None[model.OBJECT_VERSION_ID]())
	contribution := NewContribution("Person created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]model.OBJECT_REF{
			personVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.PERSON{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return model.PERSON{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedParty, utils.None[uuid.UUID]())
	if err != nil {
		return model.PERSON{}, fmt.Errorf("failed to save versioned party: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, personVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return model.PERSON{}, fmt.Errorf("failed to save person: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.PERSON{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return person, nil
}

func (s *Service) GetCurrentPersonVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (model.PERSON, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var person model.PERSON
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), model.PERSON_MODEL_NAME).Scan(&person)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return model.PERSON{}, ErrPersonNotFound
		}
		return model.PERSON{}, fmt.Errorf("failed to get latest person by versioned party ID: %w", err)
	}

	return person, nil
}

func (s *Service) GetPersonAtVersion(ctx context.Context, versionID string) (model.PERSON, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var person model.PERSON
	err := s.DB.QueryRow(ctx, query, versionID, model.GROUP_MODEL_NAME).Scan(&person)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return model.PERSON{}, ErrPersonNotFound
		}
		return model.PERSON{}, fmt.Errorf("failed to get person at version: %w", err)
	}

	return person, nil
}

func (s *Service) UpdatePerson(ctx context.Context, versionedPartyID uuid.UUID, person model.PERSON) (model.PERSON, error) {
	err := s.ValidatePerson(ctx, person)
	if err != nil {
		return model.PERSON{}, err
	}

	currentPerson, err := s.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if errors.Is(err, ErrPersonNotFound) {
			return model.PERSON{}, ErrPersonNotFound
		}
		return model.PERSON{}, fmt.Errorf("failed to get current Person: %w", err)
	}
	currentPersonID := currentPerson.UID.V.Value.(*model.OBJECT_VERSION_ID)

	err = UpgradeObjectVersionID(&currentPerson.UID, utils.Some(*currentPersonID))
	if err != nil {
		return model.PERSON{}, fmt.Errorf("failed to upgrade current Person UID: %w", err)
	}

	personVersion := NewOriginalVersion(&person, utils.Some(*currentPersonID))
	contribution := NewContribution("Person updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]model.OBJECT_REF{
			personVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.PERSON{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return model.PERSON{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, personVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return model.PERSON{}, fmt.Errorf("failed to save person: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.PERSON{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return person, nil
}

func (s *Service) DeletePerson(ctx context.Context, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Person deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]model.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      model.PERSON_MODEL_NAME,
				ID: model.X_OBJECT_ID{
					Value: &model.HIER_OBJECT_ID{
						Type_: utils.Some(model.HIER_OBJECT_ID_MODEL_NAME),
						Value: versionedObjectID.String(),
					},
				},
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.DeleteVersionedObjectWithTx(ctx, tx, versionedObjectID)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) ValidateGroup(ctx context.Context, group model.GROUP) error {
	validateErr := group.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Group validation can be added here

	return nil
}

func (s *Service) ExistsGroup(ctx context.Context, versionID string) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.type = $1 AND ov.id = $2 LIMIT 1`

	var exists int
	err := s.DB.QueryRow(ctx, query, model.GROUP_MODEL_NAME, versionID).Scan(&exists)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if group exists: %w", err)
	}
	return true, nil
}

func (s *Service) CreateGroup(ctx context.Context, group model.GROUP) (model.GROUP, error) {
	err := s.ValidateGroup(ctx, group)
	if err != nil {
		return model.GROUP{}, err
	}

	err = UpgradeObjectVersionID(&group.UID, utils.None[model.OBJECT_VERSION_ID]())
	if err != nil {
		return model.GROUP{}, fmt.Errorf("failed to upgrade group UID: %w", err)
	}

	exists, err := s.ExistsGroup(ctx, group.UID.V.ValueAsString())
	if err != nil {
		return model.GROUP{}, fmt.Errorf("failed to check existing group: %w", err)
	}
	if exists {
		return model.GROUP{}, ErrGroupAlreadyExists
	}

	versionedParty := NewVersionedParty(group.UID.V.Value.(*model.OBJECT_VERSION_ID).UID())
	groupVersion := NewOriginalVersion(&group, utils.None[model.OBJECT_VERSION_ID]())
	contribution := NewContribution("Group created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]model.OBJECT_REF{
			groupVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.GROUP{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return model.GROUP{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedParty, utils.None[uuid.UUID]())
	if err != nil {
		return model.GROUP{}, fmt.Errorf("failed to save versioned party: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, groupVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return model.GROUP{}, fmt.Errorf("failed to save group: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.GROUP{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return group, nil
}

func (s *Service) GetCurrentGroupVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (model.GROUP, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var group model.GROUP
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), model.GROUP_MODEL_NAME).Scan(&group)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return model.GROUP{}, ErrGroupNotFound
		}
		return model.GROUP{}, fmt.Errorf("failed to get latest group by versioned party ID: %w", err)
	}

	return group, nil
}

func (s *Service) GetGroupAtVersion(ctx context.Context, versionID string) (model.GROUP, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var group model.GROUP
	err := s.DB.QueryRow(ctx, query, versionID, model.GROUP_MODEL_NAME).Scan(&group)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return model.GROUP{}, ErrGroupNotFound
		}
		return model.GROUP{}, fmt.Errorf("failed to get group at version: %w", err)
	}

	return group, nil
}

func (s *Service) UpdateGroup(ctx context.Context, versionedPartyID uuid.UUID, group model.GROUP) (model.GROUP, error) {
	err := s.ValidateGroup(ctx, group)
	if err != nil {
		return model.GROUP{}, err
	}

	currentGroup, err := s.GetCurrentGroupVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if errors.Is(err, ErrGroupNotFound) {
			return model.GROUP{}, ErrGroupNotFound
		}
		return model.GROUP{}, fmt.Errorf("failed to get current Group: %w", err)
	}
	currentGroupID := currentGroup.UID.V.Value.(*model.OBJECT_VERSION_ID)

	err = UpgradeObjectVersionID(&currentGroup.UID, utils.Some(*currentGroupID))
	if err != nil {
		return model.GROUP{}, fmt.Errorf("failed to upgrade current Group UID: %w", err)
	}

	groupVersion := NewOriginalVersion(&group, utils.Some(*currentGroupID))
	contribution := NewContribution("Group updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]model.OBJECT_REF{
			groupVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.GROUP{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return model.GROUP{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, groupVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return model.GROUP{}, fmt.Errorf("failed to save group: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.GROUP{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return group, nil
}

func (s *Service) DeleteGroup(ctx context.Context, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Group deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]model.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      model.GROUP_MODEL_NAME,
				ID: model.X_OBJECT_ID{
					Value: &model.HIER_OBJECT_ID{
						Type_: utils.Some(model.OBJECT_VERSION_ID_MODEL_NAME),
						Value: versionedObjectID.String(),
					},
				},
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.DeleteVersionedObjectWithTx(ctx, tx, versionedObjectID)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) ValidateOrganisation(ctx context.Context, organisation model.ORGANISATION) error {
	validateErr := organisation.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Organisation validation can be added here

	return nil
}

func (s *Service) ExistsOrganisation(ctx context.Context, versionID string) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.type = $1 AND ov.id = $2 LIMIT 1`

	var exists int
	err := s.DB.QueryRow(ctx, query, model.ORGANISATION_MODEL_NAME, versionID).Scan(&exists)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if organisation exists: %w", err)
	}
	return true, nil
}

func (s *Service) CreateOrganisation(ctx context.Context, organisation model.ORGANISATION) (model.ORGANISATION, error) {
	err := s.ValidateOrganisation(ctx, organisation)
	if err != nil {
		return model.ORGANISATION{}, err
	}

	err = UpgradeObjectVersionID(&organisation.UID, utils.None[model.OBJECT_VERSION_ID]())
	if err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to upgrade organisation UID: %w", err)
	}

	exists, err := s.ExistsOrganisation(ctx, organisation.UID.V.ValueAsString())
	if err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to check existing organisation: %w", err)
	}
	if exists {
		return model.ORGANISATION{}, ErrOrganisationAlreadyExists
	}

	versionedParty := NewVersionedParty(organisation.UID.V.Value.(*model.OBJECT_VERSION_ID).UID())
	organisationVersion := NewOriginalVersion(&organisation, utils.None[model.OBJECT_VERSION_ID]())
	contribution := NewContribution("Organisation created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]model.OBJECT_REF{
			organisationVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedParty, utils.None[uuid.UUID]())
	if err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to save versioned party: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, organisationVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to save organisation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return organisation, nil
}

func (s *Service) GetCurrentOrganisationVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (model.ORGANISATION, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var organisation model.ORGANISATION
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), model.AGENT_MODEL_NAME).Scan(&organisation)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return model.ORGANISATION{}, ErrOrganisationNotFound
		}
		return model.ORGANISATION{}, fmt.Errorf("failed to get latest organisation by versioned party ID: %w", err)
	}

	return organisation, nil
}

func (s *Service) GetOrganisationAtVersion(ctx context.Context, versionID string) (model.ORGANISATION, error) {
	query := `
		SELECT ovd.object_data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var organisation model.ORGANISATION
	err := s.DB.QueryRow(ctx, query, versionID, model.ORGANISATION_MODEL_NAME).Scan(&organisation)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return model.ORGANISATION{}, ErrOrganisationNotFound
		}
		return model.ORGANISATION{}, fmt.Errorf("failed to get organisation at version: %w", err)
	}

	return organisation, nil
}

func (s *Service) UpdateOrganisation(ctx context.Context, versionedPartyID uuid.UUID, organisation model.ORGANISATION) (model.ORGANISATION, error) {
	err := s.ValidateOrganisation(ctx, organisation)
	if err != nil {
		return model.ORGANISATION{}, err
	}

	currentOrganisation, err := s.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if errors.Is(err, ErrOrganisationNotFound) {
			return model.ORGANISATION{}, ErrOrganisationNotFound
		}
		return model.ORGANISATION{}, fmt.Errorf("failed to get current Organisation: %w", err)
	}
	currentOrganisationID := currentOrganisation.UID.V.Value.(*model.OBJECT_VERSION_ID)

	err = UpgradeObjectVersionID(&currentOrganisation.UID, utils.Some(*currentOrganisationID))
	if err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to upgrade current Organisation UID: %w", err)
	}

	organisationVersion := NewOriginalVersion(&organisation, utils.Some(*currentOrganisationID))
	contribution := NewContribution("Organisation updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]model.OBJECT_REF{
			organisationVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, organisationVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to save organisation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.ORGANISATION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return organisation, nil
}

func (s *Service) DeleteOrganisation(ctx context.Context, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Organisation deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]model.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      model.ORGANISATION_MODEL_NAME,
				ID: model.X_OBJECT_ID{
					Value: &model.HIER_OBJECT_ID{
						Type_: utils.Some(model.HIER_OBJECT_ID_MODEL_NAME),
						Value: versionedObjectID.String(),
					},
				},
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.DeleteVersionedObjectWithTx(ctx, tx, versionedObjectID)
	if err != nil {
		return fmt.Errorf("failed to delete organisation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) ValidateRole(ctx context.Context, role model.ROLE) error {
	validateErr := role.Validate("$")
	if len(validateErr.Errs) > 0 {
		return validateErr
	}

	// Additional Organisation validation can be added here

	return nil
}

func (s *Service) ExistsRole(ctx context.Context, versionID string) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.type = $1 AND ov.id = $2 LIMIT 1`

	var exists int
	err := s.DB.QueryRow(ctx, query, model.ROLE_MODEL_NAME, versionID).Scan(&exists)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if role exists: %w", err)
	}
	return true, nil
}

func (s *Service) CreateRole(ctx context.Context, role model.ROLE) (model.ROLE, error) {
	err := s.ValidateRole(ctx, role)
	if err != nil {
		return model.ROLE{}, err
	}

	err = UpgradeObjectVersionID(&role.UID, utils.None[model.OBJECT_VERSION_ID]())
	if err != nil {
		return model.ROLE{}, fmt.Errorf("failed to upgrade role UID: %w", err)
	}

	exists, err := s.ExistsRole(ctx, role.UID.V.ValueAsString())
	if err != nil {
		return model.ROLE{}, fmt.Errorf("failed to check existing role: %w", err)
	}
	if exists {
		return model.ROLE{}, ErrRoleAlreadyExists
	}

	versionedParty := NewVersionedParty(role.UID.V.Value.(*model.OBJECT_VERSION_ID).UID())
	roleVersion := NewOriginalVersion(&role, utils.None[model.OBJECT_VERSION_ID]())
	contribution := NewContribution("Role created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]model.OBJECT_REF{
			roleVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.ROLE{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return model.ROLE{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedParty, utils.None[uuid.UUID]())
	if err != nil {
		return model.ROLE{}, fmt.Errorf("failed to save versioned party: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, roleVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return model.ROLE{}, fmt.Errorf("failed to save role: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.ROLE{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return role, nil
}

func (s *Service) GetCurrentRoleVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (model.ROLE, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var role model.ROLE
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), model.ROLE_MODEL_NAME).Scan(&role)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return model.ROLE{}, ErrRoleNotFound
		}
		return model.ROLE{}, fmt.Errorf("failed to get latest role by versioned party ID: %w", err)
	}

	return role, nil
}

func (s *Service) GetRoleAtVersion(ctx context.Context, versionID string) (model.ROLE, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var role model.ROLE
	err := s.DB.QueryRow(ctx, query, versionID, model.ROLE_MODEL_NAME).Scan(&role)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return model.ROLE{}, ErrRoleNotFound
		}
		return model.ROLE{}, fmt.Errorf("failed to get role at version: %w", err)
	}

	return role, nil
}

func (s *Service) UpdateRole(ctx context.Context, versionedPartyID uuid.UUID, role model.ROLE) (model.ROLE, error) {
	err := s.ValidateRole(ctx, role)
	if err != nil {
		return model.ROLE{}, err
	}

	currentRole, err := s.GetCurrentRoleVersionByVersionedPartyID(ctx, versionedPartyID)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			return model.ROLE{}, ErrRoleNotFound
		}
		return model.ROLE{}, fmt.Errorf("failed to get current Role: %w", err)
	}
	currentRoleID := currentRole.UID.V.Value.(*model.OBJECT_VERSION_ID)

	err = UpgradeObjectVersionID(&currentRole.UID, utils.Some(*currentRoleID))
	if err != nil {
		return model.ROLE{}, fmt.Errorf("failed to upgrade current Role UID: %w", err)
	}

	roleVersion := NewOriginalVersion(&role, utils.Some(*currentRoleID))
	contribution := NewContribution("Role updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]model.OBJECT_REF{
			roleVersion.ObjectRef(),
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return model.ROLE{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return model.ROLE{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, roleVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return model.ROLE{}, fmt.Errorf("failed to save role: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.ROLE{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return role, nil
}

func (s *Service) DeleteRole(ctx context.Context, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Role deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]model.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      model.ROLE_MODEL_NAME,
				ID: model.X_OBJECT_ID{
					Value: &model.HIER_OBJECT_ID{
						Type_: utils.Some(model.HIER_OBJECT_ID_MODEL_NAME),
						Value: versionedObjectID.String(),
					},
				},
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.DeleteVersionedObjectWithTx(ctx, tx, versionedObjectID)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) GetVersionedParty(ctx context.Context, versionedObjectID uuid.UUID) (model.VERSIONED_PARTY, error) {
	query := `
		SELECT vod.data 
		FROM openehr.tbl_versioned_object vo
		JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id
		WHERE vo.ehr_id IS NULL AND vo.type = $1 AND vo.id = $2 
		LIMIT 1
	`

	var versionedParty model.VERSIONED_PARTY
	err := s.DB.QueryRow(ctx, query, model.VERSIONED_PARTY_MODEL_NAME, versionedObjectID).Scan(&versionedParty)
	if err != nil {
		return model.VERSIONED_PARTY{}, fmt.Errorf("failed to get versioned party by ID: %w", err)
	}
	return versionedParty, nil
}

func (s *Service) GetVersionedPartyRevisionHistory(ctx context.Context, versionedObjectID uuid.UUID) (model.REVISION_HISTORY, error) {
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
	args := []any{[]string{model.AGENT_MODEL_NAME, model.PERSON_MODEL_NAME, model.GROUP_MODEL_NAME, model.ORGANISATION_MODEL_NAME, model.ROLE_MODEL_NAME}, versionedObjectID}
	row := s.DB.QueryRow(ctx, query, args...)

	var revisionHistory model.REVISION_HISTORY
	if err := row.Scan(&revisionHistory); err != nil {
		if err == database.ErrNoRows {
			return model.REVISION_HISTORY{}, ErrRevisionHistoryNotFound
		}
		return model.REVISION_HISTORY{}, fmt.Errorf("failed to fetch Revision History from database: %w", err)
	}

	return revisionHistory, nil
}

func (s *Service) GetVersionedPartyVersionJSON(ctx context.Context, versionedObjectID uuid.UUID, filterAtTime time.Time, filterVersionID string) ([]byte, error) {
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

func (s *Service) GetContribution(ctx context.Context, contributionID string, ehrID utils.Optional[uuid.UUID]) (model.CONTRIBUTION, error) {
	query := `
		SELECT cd.data
		FROM openehr.tbl_contribution c
		JOIN openehr.tbl_contribution_data cd ON c.id = cd.id
		WHERE c.ehr_id = $1 AND c.id = $2
		LIMIT 1
	`
	args := []any{ehrID, contributionID}
	row := s.DB.QueryRow(ctx, query, args...)

	var contribution model.CONTRIBUTION
	if err := row.Scan(&contribution); err != nil {
		if err == database.ErrNoRows {
			return model.CONTRIBUTION{}, ErrContributionNotFound
		}
		return model.CONTRIBUTION{}, fmt.Errorf("failed to fetch Contribution by ID from database: %w", err)
	}

	return contribution, nil
}

func (s *Service) SaveEHRWithTx(ctx context.Context, tx pgx.Tx, ehrID uuid.UUID) error {
	query := `INSERT INTO openehr.tbl_ehr (id) VALUES ($1)`
	args := []any{ehrID}
	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert EHR into the database: %w", err)
	}

	return nil
}

func (s *Service) SaveContributionWithTx(ctx context.Context, tx pgx.Tx, contribution model.CONTRIBUTION, ehrID utils.Optional[uuid.UUID]) error {
	// Insert Contribution
	query := `INSERT INTO openehr.tbl_contribution (id, ehr_id) VALUES ($1, $2)`
	args := []any{contribution.UID.Value, ehrID}
	_, err := tx.Exec(ctx, query, args...)
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

	return nil
}

func (s *Service) SaveVersionedObjectWithTx(ctx context.Context, tx pgx.Tx, versionedObject any, ehrID utils.Optional[uuid.UUID]) error {
	var (
		modelType string
		id        string
	)
	switch v := versionedObject.(type) {
	case model.VERSIONED_EHR_STATUS:
		v.SetModelName()
		modelType = model.VERSIONED_EHR_STATUS_MODEL_NAME
		id = v.UID.Value
	case model.VERSIONED_EHR_ACCESS:
		v.SetModelName()
		modelType = model.VERSIONED_EHR_ACCESS_MODEL_NAME
		id = v.UID.Value
	case model.VERSIONED_COMPOSITION:
		v.SetModelName()
		modelType = model.VERSIONED_COMPOSITION_MODEL_NAME
		id = v.UID.Value
	case model.VERSIONED_FOLDER:
		v.SetModelName()
		modelType = model.VERSIONED_FOLDER_MODEL_NAME
		id = v.UID.Value
	case model.VERSIONED_PARTY:
		v.SetModelName()
		modelType = model.VERSIONED_PARTY_MODEL_NAME
		id = v.UID.Value
	default:
		return fmt.Errorf("unsupported versioned object type for creation: %T", versionedObject)
	}

	query := `INSERT INTO openehr.tbl_versioned_object (id, type, ehr_id) VALUES ($1, $2, $3)`
	args := []any{id, modelType, ehrID}
	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert versioned object into the database: %w", err)
	}

	query = `INSERT INTO openehr.tbl_versioned_object_data (id, data) VALUES ($1, $2)`
	args = []any{id, versionedObject}
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert versioned object data into the database: %w", err)
	}

	return nil
}

func (s *Service) SaveObjectVersionWithTx(ctx context.Context, tx pgx.Tx, version any, contributionID string, ehrID utils.Optional[uuid.UUID]) error {
	var object any
	switch v := version.(type) {
	case model.ORIGINAL_VERSION:
		v.SetModelName()
		object = v.Data
	// After enabling, make sure to change the data path below in the INSERT statement
	// case model.IMPORTED_VERSION:
	// 	object = v.Data
	default:
		return fmt.Errorf("unsupported version type for object version creation: %T", version)
	}

	var (
		modelType         string
		id                string
		versionedObjectID uuid.UUID
	)
	switch v := object.(type) {
	case *model.EHR_STATUS:
		modelType = model.EHR_STATUS_MODEL_NAME
		id = v.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
		versionedObjectID = v.UID.V.Value.(*model.OBJECT_VERSION_ID).UID()
	case *model.EHR_ACCESS:
		modelType = model.EHR_ACCESS_MODEL_NAME
		id = v.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
		versionedObjectID = v.UID.V.Value.(*model.OBJECT_VERSION_ID).UID()
	case *model.COMPOSITION:
		modelType = model.COMPOSITION_MODEL_NAME
		id = v.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
		versionedObjectID = v.UID.V.Value.(*model.OBJECT_VERSION_ID).UID()
	case *model.FOLDER:
		modelType = model.FOLDER_MODEL_NAME
		id = v.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
		versionedObjectID = v.UID.V.Value.(*model.OBJECT_VERSION_ID).UID()
	case *model.ROLE:
		modelType = model.ROLE_MODEL_NAME
		id = v.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
		versionedObjectID = v.UID.V.Value.(*model.OBJECT_VERSION_ID).UID()
	case *model.PERSON:
		modelType = model.PERSON_MODEL_NAME
		id = v.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
		versionedObjectID = v.UID.V.Value.(*model.OBJECT_VERSION_ID).UID()
	case *model.AGENT:
		modelType = model.AGENT_MODEL_NAME
		id = v.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
		versionedObjectID = v.UID.V.Value.(*model.OBJECT_VERSION_ID).UID()
	case *model.GROUP:
		modelType = model.GROUP_MODEL_NAME
		id = v.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
		versionedObjectID = v.UID.V.Value.(*model.OBJECT_VERSION_ID).UID()
	case *model.ORGANISATION:
		modelType = model.ORGANISATION_MODEL_NAME
		id = v.UID.V.Value.(*model.OBJECT_VERSION_ID).Value
		versionedObjectID = v.UID.V.Value.(*model.OBJECT_VERSION_ID).UID()
	default:
		return fmt.Errorf("unsupported object type for version creation: %T", object)
	}

	query := `INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args := []any{id, versionedObjectID, modelType, ehrID, contributionID}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert object version into the database: %w", err)
	}

	query = `INSERT INTO openehr.tbl_object_version_data (id, version_data, object_data) VALUES ($1, jsonb_set($2, '{data}', 'null', true), $2->'data')`
	args = []any{id, version}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert object version data into the database: %w", err)
	}

	return nil
}

func (s *Service) DeleteVersionedObjectWithTx(ctx context.Context, tx pgx.Tx, versionedObjectID uuid.UUID) error {
	var deleted uint8
	row := tx.QueryRow(ctx, `DELETE FROM openehr.tbl_versioned_object WHERE id = $1 RETURNING 1`, versionedObjectID)
	err := row.Scan(&deleted)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return ErrVersionedObjectNotFound
		}
		return fmt.Errorf("failed to delete versioned object by object version ID from the database: %w", err)
	}

	return nil
}

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

func NewVersionedEHRAccess(id, ehrID uuid.UUID) model.VERSIONED_EHR_ACCESS {
	return model.VERSIONED_EHR_ACCESS{
		UID: model.HIER_OBJECT_ID{
			Value: id.String(),
		},
		OwnerID: model.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      model.EHR_MODEL_NAME,
			ID: model.X_OBJECT_ID{
				Value: &model.HIER_OBJECT_ID{
					Type_: utils.Some(model.HIER_OBJECT_ID_MODEL_NAME),
					Value: ehrID.String(),
				},
			},
		},
		TimeCreated: model.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func NewEHRAccess(id uuid.UUID) model.EHR_ACCESS {
	return model.EHR_ACCESS{
		UID: utils.Some(model.X_UID_BASED_ID{
			Value: &model.OBJECT_VERSION_ID{
				Type_: utils.Some(model.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", id.String(), config.SYSTEM_ID_GOPENEHR),
			},
		}),
		Name: model.X_DV_TEXT{
			Value: &model.DV_TEXT{
				Value: "EHR Access",
			},
		},
		ArchetypeNodeID: "openEHR-EHR-EHR_ACCESS.generic.v1",
	}
}

func NewVersionedEHRStatus(id, ehrID uuid.UUID) model.VERSIONED_EHR_STATUS {
	return model.VERSIONED_EHR_STATUS{
		UID: model.HIER_OBJECT_ID{
			Value: id.String(),
		},
		OwnerID: model.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      model.EHR_MODEL_NAME,
			ID: model.X_OBJECT_ID{
				Value: &model.HIER_OBJECT_ID{
					Type_: utils.Some(model.HIER_OBJECT_ID_MODEL_NAME),
					Value: ehrID.String(),
				},
			},
		},
		TimeCreated: model.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func NewEHRStatus(id uuid.UUID) model.EHR_STATUS {
	return model.EHR_STATUS{
		UID: utils.Some(model.X_UID_BASED_ID{
			Value: &model.OBJECT_VERSION_ID{
				Type_: utils.Some(model.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", id.String(), config.SYSTEM_ID_GOPENEHR),
			},
		}),
		Name: model.X_DV_TEXT{
			Value: &model.DV_TEXT{
				Value: "EHR Status",
			},
		},
		ArchetypeNodeID: "openEHR-EHR-EHR_STATUS.generic.v1",
		Subject:         model.PARTY_SELF{},
		IsQueryable:     true,
		IsModifiable:    true,
	}
}

func NewVersionedComposition(id, ehrID uuid.UUID) model.VERSIONED_COMPOSITION {
	return model.VERSIONED_COMPOSITION{
		UID: model.HIER_OBJECT_ID{
			Value: id.String(),
		},
		OwnerID: model.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      model.EHR_MODEL_NAME,
			ID: model.X_OBJECT_ID{
				Value: &model.HIER_OBJECT_ID{
					Type_: utils.Some(model.HIER_OBJECT_ID_MODEL_NAME),
					Value: ehrID.String(),
				},
			},
		},
		TimeCreated: model.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func NewVersionedFolder(id, ehrID uuid.UUID) model.VERSIONED_FOLDER {
	return model.VERSIONED_FOLDER{
		UID: model.HIER_OBJECT_ID{
			Value: id.String(),
		},
		OwnerID: model.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      model.EHR_MODEL_NAME,
			ID: model.X_OBJECT_ID{
				Value: &model.HIER_OBJECT_ID{
					Type_: utils.Some(model.HIER_OBJECT_ID_MODEL_NAME),
					Value: ehrID.String(),
				},
			},
		},
		TimeCreated: model.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func NewVersionedParty(uid uuid.UUID) model.VERSIONED_PARTY {
	return model.VERSIONED_PARTY{
		UID: model.HIER_OBJECT_ID{
			Value: uid.String(),
		},
		OwnerID: model.OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      model.ORGANISATION_MODEL_NAME,
			ID: model.X_OBJECT_ID{
				Value: &model.HIER_OBJECT_ID{
					Type_: utils.Some(model.HIER_OBJECT_ID_MODEL_NAME),
					Value: uid.String(),
				},
			},
		},
		TimeCreated: model.DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

func UpgradeObjectVersionID(currentUID *utils.Optional[model.X_UID_BASED_ID], previousUID utils.Optional[model.OBJECT_VERSION_ID]) error {
	// Provide ID when EHR Status does not have one
	if !currentUID.E {
		*currentUID = utils.Some(model.X_UID_BASED_ID{
			Value: &model.OBJECT_VERSION_ID{
				Type_: utils.Some(model.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", uuid.NewString(), config.SYSTEM_ID_GOPENEHR),
			},
		})
	}

	switch v := currentUID.V.Value.(type) {
	case *model.OBJECT_VERSION_ID:
		// valid type
		if previousUID.E {
			// Check version is incremented
			if v.VersionTreeID().CompareTo(previousUID.V.VersionTreeID()) <= 0 {
				return ErrVersionLowerOrEqualToCurrent
			}
		}
	case *model.HIER_OBJECT_ID:
		// Add namespace and version to convert to OBJECT_VERSION_ID
		hierID := currentUID.V.Value.(*model.HIER_OBJECT_ID)
		*currentUID = utils.Some(model.X_UID_BASED_ID{
			Value: &model.OBJECT_VERSION_ID{
				Type_: utils.Some(model.OBJECT_VERSION_ID_MODEL_NAME),
				Value: fmt.Sprintf("%s::%s::1", hierID.Value, config.SYSTEM_ID_GOPENEHR),
			},
		})
	default:
		return fmt.Errorf("EHR Status UID must be of type OBJECT_VERSION_ID or HIER_OBJECT_ID, got %T", currentUID.V.Value)
	}
	return nil
}

func NewContribution(description string, auditChangeType terminology.AuditChangeType, versions []model.OBJECT_REF) model.CONTRIBUTION {
	contribution := model.CONTRIBUTION{
		UID: model.HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		Versions: make([]model.OBJECT_REF, 0),
		Audit: model.AUDIT_DETAILS{
			SystemID:      config.SYSTEM_ID_GOPENEHR,
			TimeCommitted: model.DV_DATE_TIME{Value: time.Now().UTC().Format(time.RFC3339)},
			Description:   utils.Some(model.DV_TEXT{Value: description}),
			Committer: model.X_PARTY_PROXY{
				Value: &model.PARTY_SELF{
					Type_: utils.Some(model.PARTY_SELF_MODEL_NAME),
				},
			},
		},
	}

	switch auditChangeType {
	case terminology.AUDIT_CHANGE_TYPE_CODE_CREATION:
		contribution.Audit.ChangeType = model.DV_CODED_TEXT{
			Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_CREATION),
			DefiningCode: model.CODE_PHRASE{
				CodeString: string(terminology.AUDIT_CHANGE_TYPE_CODE_CREATION),
				TerminologyID: model.TERMINOLOGY_ID{
					Value: string(terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR),
				},
			},
		}
	case terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION:
		contribution.Audit.ChangeType = model.DV_CODED_TEXT{
			Value: terminology.GetAuditChangeTypeName(terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION),
			DefiningCode: model.CODE_PHRASE{
				CodeString: string(terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION),
				TerminologyID: model.TERMINOLOGY_ID{
					Value: string(terminology.AUDIT_CHANGE_TYPE_TERMINOLOGY_ID_OPENEHR),
				},
			},
		}
	default:
		return model.CONTRIBUTION{}
	}

	return contribution
}

func NewOriginalVersion(data model.VersionModel, precedingVersion utils.Optional[model.OBJECT_VERSION_ID]) model.ORIGINAL_VERSION {
	return model.ORIGINAL_VERSION{
		UID:                 data.ObjectVersionID(),
		PrecedingVersionUID: precedingVersion,
		LifecycleState: model.DV_CODED_TEXT{
			Value: terminology.VersionLifecycleStateNames[terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE],
			DefiningCode: model.CODE_PHRASE{
				CodeString: terminology.VERSION_LIFECYCLE_STATE_CODE_COMPLETE,
				TerminologyID: model.TERMINOLOGY_ID{
					Value: terminology.VERSION_LIFECYCLE_STATE_TERMINOLOGY_ID_OPENEHR,
				},
			},
		},
		Data: data,
	}
}
