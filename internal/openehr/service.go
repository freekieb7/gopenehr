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

func NewService(logger *telemetry.Logger, db *database.Database) Service {
	return Service{
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
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*ehrStatus.UID.V.OBJECT_VERSION_ID()),
			},
			{
				Type:      rm.VERSIONED_EHR_ACCESS_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*ehrAccess.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != database.ErrTxClosed {
			s.Logger.ErrorContext(ctx, "Failed to rollback transaction", "error", err)
		}
	}()

	err = s.SaveEHRWithTx(ctx, tx, ehrID)
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to save EHR: %w", err)
	}
	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedEHRStatus, utils.Some(ehrID))
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to save VERSIONED_EHR_STATUS: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, ehrStatusVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to save EHR_STATUS: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedEHRAccess, utils.Some(ehrID))
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to save VERSIONED_EHR_ACCESS: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, ehrAccessVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to save EHR_ACCESS: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return rm.EHR{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	ehr, err := s.GetEHR(ctx, ehrID)
	if err != nil {
		return rm.EHR{}, fmt.Errorf("failed to get EHR after creation: %w", err)
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

func (s *Service) GetEHR(ctx context.Context, id uuid.UUID) (rm.EHR, error) {
	query := `SELECT data FROM openehr.vw_ehr WHERE id = $1 LIMIT 1`
	args := []any{id}

	row := s.DB.QueryRow(ctx, query, args...)

	var ehr rm.EHR
	err := row.Scan(&ehr)
	if err != nil {
		if err == database.ErrNoRows {
			return rm.EHR{}, ErrEHRNotFound
		}
		return rm.EHR{}, fmt.Errorf("failed to fetch EHR from database: %w", err)
	}

	return ehr, nil
}

func (s *Service) GetEHRBySubject(ctx context.Context, subjectID, subjectNamespace string) (rm.EHR, error) {
	query := `
        SELECT ehr_id
        FROM openehr.tbl_object_version_data
        WHERE ov.type = $1
          AND object_data->'subject'->'external_ref'->>'namespace' = $2
		  AND object_data->'subject'->'external_ref'->'id'->>'value' = $3
        ORDER BY created_at DESC
        LIMIT 1
    `
	args := []any{rm.EHR_STATUS_TYPE, subjectNamespace, subjectID}

	row := s.DB.QueryRow(ctx, query, args...)

	var ehrID uuid.UUID
	err := row.Scan(&ehrID)
	if err != nil {
		if err == database.ErrNoRows {
			return rm.EHR{}, ErrEHRNotFound
		}
		return rm.EHR{}, fmt.Errorf("failed to fetch EHR by subject from database: %w", err)
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

func (s *Service) GetEHRStatus(ctx context.Context, ehrID uuid.UUID, filterOnTime time.Time, filterOnVersionID string) (rm.EHR_STATUS, error) {
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
	args = []any{rm.EHR_STATUS_TYPE, ehrID}
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

	var ehrStatus rm.EHR_STATUS
	err := row.Scan(&ehrStatus)
	if err != nil {
		if err == database.ErrNoRows {
			return rm.EHR_STATUS{}, ErrEHRStatusNotFound
		}
		return rm.EHR_STATUS{}, fmt.Errorf("failed to fetch EHR Status from database: %w", err)
	}

	return ehrStatus, nil
}

func (s *Service) UpdateEHRStatus(ctx context.Context, ehrID uuid.UUID, previousEHRStatus, nextEHRStatus rm.EHR_STATUS) (rm.EHR_STATUS, error) {
	if err := s.ValidateEHRStatus(ctx, nextEHRStatus); err != nil {
		return rm.EHR_STATUS{}, err
	}

	// Ensure EHR Status contains a UID to check/upgrade
	if !nextEHRStatus.UID.E {
		nextEHRStatus.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: previousEHRStatus.UID.V.OBJECT_VERSION_ID().UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(nextEHRStatus.UID.V, *previousEHRStatus.UID.V.OBJECT_VERSION_ID())
	if err != nil {
		return rm.EHR_STATUS{}, fmt.Errorf("failed to upgrade current EHR Status UID: %w", err)
	}
	nextEHRStatus.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	ehrStatusVersion := NewOriginalVersion(nextEHRStatus.ObjectVersionID(), rm.ORIGINAL_VERSION_DATA_from_EHR_STATUS(nextEHRStatus), utils.Some(*previousEHRStatus.UID.V.OBJECT_VERSION_ID()))
	contribution := NewContribution("EHR Status updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_EHR_STATUS_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*nextEHRStatus.UID.V.OBJECT_VERSION_ID()),
			},
		},
	)

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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return rm.EHR_STATUS{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, ehrStatusVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return rm.EHR_STATUS{}, fmt.Errorf("failed to save ehr status version: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return rm.EHR_STATUS{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextEHRStatus, nil
}

func (s *Service) GetVersionedEHRStatus(ctx context.Context, ehrID uuid.UUID) (rm.VERSIONED_EHR_STATUS, error) {
	query := `
		SELECT vod.data 
		FROM openehr.tbl_versioned_object vo
		JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id
		WHERE vo.type = $1 
		  AND vo.ehr_id = $2
		LIMIT 1
	`
	args := []any{rm.VERSIONED_EHR_STATUS_TYPE, ehrID}

	row := s.DB.QueryRow(ctx, query, args...)

	var versionedEHRStatus rm.VERSIONED_EHR_STATUS
	err := row.Scan(&versionedEHRStatus)
	if err != nil {
		if err == database.ErrNoRows {
			return rm.VERSIONED_EHR_STATUS{}, ErrEHRNotFound
		}
		return rm.VERSIONED_EHR_STATUS{}, fmt.Errorf("failed to fetch Versioned EHR Status from database: %w", err)
	}

	return versionedEHRStatus, nil
}

func (s *Service) GetVersionedEHRStatusRevisionHistory(ctx context.Context, ehrID uuid.UUID) (rm.REVISION_HISTORY, error) {
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

	var revisionHistory rm.REVISION_HISTORY
	err := row.Scan(&revisionHistory)
	if err != nil {
		if err == database.ErrNoRows {
			return rm.REVISION_HISTORY{}, ErrEHRNotFound
		}
		return rm.REVISION_HISTORY{}, fmt.Errorf("failed to fetch Revision History from database: %w", err)
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
	args = []any{rm.EHR_STATUS_TYPE, ehrID}
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
	compositionVersion := NewOriginalVersion(*composition.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_COMPOSITION(composition), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Composition created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_COMPOSITION_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*composition.UID.V.OBJECT_VERSION_ID()),
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

	// Insert Contribution
	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to save contribution: %w", err)
	}

	// Insert Versioned Composition
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedComposition, utils.Some(ehrID))
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to save versioned composition: %w", err)
	}

	// Insert Composition Version
	err = s.SaveObjectVersionWithTx(ctx, tx, compositionVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to save composition version: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return composition, nil
}

func (s *Service) GetComposition(ctx context.Context, ehrID uuid.UUID, objectVersionID string) (rm.COMPOSITION, error) {
	query := `
		SELECT ovd.object_data 
		FROM openehr.tbl_object_version ov 
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.type = $1 AND ov.ehr_id = $2 AND ov.id = $3 LIMIT 1
	`
	args := []any{rm.COMPOSITION_TYPE, ehrID, objectVersionID}

	var composition rm.COMPOSITION
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&composition); err != nil {
		if err == database.ErrNoRows {
			return rm.COMPOSITION{}, ErrCompositionNotFound
		}
		return rm.COMPOSITION{}, fmt.Errorf("failed to fetch Composition by ID from database: %w", err)
	}

	return composition, nil
}

func (s *Service) GetCurrentCompositionByVersionedCompositionID(ctx context.Context, ehrID uuid.UUID, versionedCompositionID uuid.UUID) (rm.COMPOSITION, error) {
	query := `
		SELECT ovd.object_data 
		FROM openehr.tbl_object_version ov 
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.type = $1 AND ov.ehr_id = $2 AND ov.versioned_object_id = $3 
		ORDER BY ov.created_at DESC
		LIMIT 1
	`
	args := []any{rm.COMPOSITION_TYPE, ehrID, versionedCompositionID}

	var composition rm.COMPOSITION
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&composition); err != nil {
		if err == database.ErrNoRows {
			return rm.COMPOSITION{}, ErrCompositionNotFound
		}
		return rm.COMPOSITION{}, fmt.Errorf("failed to fetch Composition by ID from database: %w", err)
	}

	return composition, nil
}

func (s *Service) UpdateComposition(ctx context.Context, ehrID uuid.UUID, previousComposition, nextComposition rm.COMPOSITION) (rm.COMPOSITION, error) {
	err := s.ValidateComposition(ctx, nextComposition)
	if err != nil {
		return rm.COMPOSITION{}, err
	}

	// Ensure Composition contains a UID to check/upgrade
	if !nextComposition.UID.E {
		nextComposition.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: previousComposition.UID.V.OBJECT_VERSION_ID().UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(nextComposition.UID.V, *previousComposition.UID.V.OBJECT_VERSION_ID())
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to upgrade current Composition UID: %w", err)
	}
	nextComposition.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	compositionVersion := NewOriginalVersion(*nextComposition.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_COMPOSITION(nextComposition), utils.Some(*previousComposition.UID.V.OBJECT_VERSION_ID()))
	contribution := NewContribution("Composition updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_COMPOSITION_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*nextComposition.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, compositionVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to save composition: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.COMPOSITION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextComposition, nil
}

func (s *Service) DeleteComposition(ctx context.Context, ehrID uuid.UUID, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Composition deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      rm.COMPOSITION_TYPE,
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

func (s *Service) GetVersionedComposition(ctx context.Context, ehrID uuid.UUID, versionedObjectID string) (rm.VERSIONED_COMPOSITION, error) {
	query := `
		SELECT vod.data 
		FROM openehr.tbl_versioned_object vo
		JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id
		WHERE vo.type = $1
		  AND vo.ehr_id = $2 
		  AND vo.id = $3
		LIMIT 1
	`
	args := []any{rm.VERSIONED_COMPOSITION_TYPE, ehrID, versionedObjectID}
	row := s.DB.QueryRow(ctx, query, args...)

	var versionedComposition rm.VERSIONED_COMPOSITION
	if err := row.Scan(&versionedComposition); err != nil {
		if err == database.ErrNoRows {
			return rm.VERSIONED_COMPOSITION{}, ErrCompositionNotFound
		}
		return rm.VERSIONED_COMPOSITION{}, fmt.Errorf("failed to fetch Versioned Composition by ID from database: %w", err)
	}

	return versionedComposition, nil
}

func (s *Service) GetVersionedCompositionRevisionHistory(ctx context.Context, ehrID uuid.UUID, versionedObjectID string) (rm.REVISION_HISTORY, error) {
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

	var revisionHistory rm.REVISION_HISTORY
	err := row.Scan(&revisionHistory)
	if err != nil {
		if err == database.ErrNoRows {
			return rm.REVISION_HISTORY{}, ErrCompositionNotFound
		}
		return rm.REVISION_HISTORY{}, fmt.Errorf("failed to fetch Revision History from database: %w", err)
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
	args = []any{rm.COMPOSITION_TYPE, ehrID, versionedObjectID}
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
	folderVersion := NewOriginalVersion(*directory.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_FOLDER(directory), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Directory created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_FOLDER_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*directory.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedFolder, utils.Some(ehrID))
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to save versioned folder: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, folderVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to save folder version: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return directory, nil
}

func (s *Service) ExistsDirectory(ctx context.Context, ehrID uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_object_version ov WHERE ov.ehr_id = $1 AND ov.type = $2 LIMIT 1`
	args := []any{ehrID, rm.FOLDER_TYPE}

	var exists int
	if err := s.DB.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		if err == database.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if Directory exists in database: %w", err)
	}

	return true, nil
}

func (s *Service) GetDirectory(ctx context.Context, ehrID uuid.UUID) (rm.FOLDER, error) {
	query := `
		SELECT ovd.object_data
        FROM openehr.tbl_object_version ov
        JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id
        WHERE ov.type = $1
          AND ov.ehr_id = $2
        ORDER BY ov.created_at DESC
        LIMIT 1
	`
	args := []any{rm.FOLDER_TYPE, ehrID}
	row := s.DB.QueryRow(ctx, query, args...)

	var directory rm.FOLDER
	if err := row.Scan(&directory); err != nil {
		if err == database.ErrNoRows {
			return rm.FOLDER{}, ErrDirectoryNotFound
		}
		return rm.FOLDER{}, fmt.Errorf("failed to fetch Directory by EHR ID from database: %w", err)
	}

	return directory, nil
}

func (s *Service) UpdateDirectory(ctx context.Context, ehrID uuid.UUID, previousDirectory, nextDirectory rm.FOLDER) (rm.FOLDER, error) {
	err := s.ValidateDirectory(ctx, ehrID, nextDirectory)
	if err != nil {
		return rm.FOLDER{}, err
	}

	if !nextDirectory.UID.E {
		nextDirectory.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: previousDirectory.UID.V.OBJECT_VERSION_ID().UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(previousDirectory.UID.V, *previousDirectory.UID.V.OBJECT_VERSION_ID())
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to upgrade current Directory UID: %w", err)
	}
	nextDirectory.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	folderVersion := NewOriginalVersion(*nextDirectory.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_FOLDER(nextDirectory), utils.Some(*previousDirectory.UID.V.OBJECT_VERSION_ID()))
	contribution := NewContribution("Directory updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_FOLDER_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*nextDirectory.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.Some(ehrID))
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, folderVersion, contribution.UID.Value, utils.Some(ehrID))
	if err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to save folder version: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.FOLDER{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextDirectory, nil
}

func (s *Service) DeleteDirectory(ctx context.Context, ehrID uuid.UUID, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Directory deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      rm.FOLDER_TYPE,
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

func (s *Service) GetFolderInDirectoryVersion(ctx context.Context, ehrID uuid.UUID, filterAtTime time.Time, filterVersionID string, filterPathParts []string) (rm.FOLDER, error) {
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
	args = []any{rm.FOLDER_TYPE, ehrID, jsonPath}
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

	var folder rm.FOLDER
	if err := row.Scan(&folder); err != nil {
		if err == database.ErrNoRows {
			if len(filterPathParts) > 0 {
				return rm.FOLDER{}, ErrFolderNotFoundInDirectory
			}
			return rm.FOLDER{}, ErrDirectoryNotFound
		}
		return rm.FOLDER{}, fmt.Errorf("failed to fetch Folder at time from database: %w", err)
	}

	return folder, nil
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
	agentVersion := NewOriginalVersion(*agent.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_AGENT(agent), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Agent created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*agent.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedParty, utils.None[uuid.UUID]())
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to save versioned party: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, agentVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to save agent: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return agent, nil
}

func (s *Service) GetCurrentAgentVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (rm.AGENT, error) {
	query := `
		SELECT ovd.object_data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var agent rm.AGENT
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), rm.AGENT_TYPE).Scan(&agent)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.AGENT{}, ErrAgentNotFound
		}
		return rm.AGENT{}, fmt.Errorf("failed to get latest agent by versioned party ID: %w", err)
	}

	return agent, nil
}

func (s *Service) GetAgentAtVersion(ctx context.Context, versionID string) (rm.AGENT, error) {
	query := `
		SELECT ovd.object_data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var agent rm.AGENT
	err := s.DB.QueryRow(ctx, query, versionID, rm.AGENT_TYPE).Scan(&agent)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.AGENT{}, ErrAgentNotFound
		}
		return rm.AGENT{}, fmt.Errorf("failed to get agent at version: %w", err)
	}

	return agent, nil
}

func (s *Service) UpdateAgent(ctx context.Context, currentAgent, nextAgent rm.AGENT) (rm.AGENT, error) {
	err := s.ValidateAgent(ctx, nextAgent)
	if err != nil {
		return rm.AGENT{}, err
	}

	if !nextAgent.UID.E {
		nextAgent.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: currentAgent.UID.V.OBJECT_VERSION_ID().UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(currentAgent.UID.V, *currentAgent.UID.V.OBJECT_VERSION_ID())
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to upgrade current Agent UID: %w", err)
	}
	nextAgent.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	agentVersion := NewOriginalVersion(*nextAgent.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_AGENT(nextAgent), utils.Some(*currentAgent.UID.V.OBJECT_VERSION_ID()))
	contribution := NewContribution("Agent updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*nextAgent.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, agentVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to save agent: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.AGENT{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextAgent, nil
}

func (s *Service) DeleteAgent(ctx context.Context, versionedObjectID string) error {
	// Check if agent exists
	_, err := s.GetCurrentAgentVersionByVersionedPartyID(ctx, uuid.MustParse(strings.Split(versionedObjectID, "::")[0]))
	if err != nil {
		if errors.Is(err, ErrAgentNotFound) {
			return ErrAgentNotFound
		}
		return fmt.Errorf("failed to get agent before deletion: %w", err)
	}

	contribution := NewContribution("Agent deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      rm.AGENT_TYPE,
				ID: rm.OBJECT_ID_from_HIER_OBJECT_ID(rm.HIER_OBJECT_ID{
					Value: versionedObjectID,
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
	personVersion := NewOriginalVersion(*person.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_PERSON(person), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Person created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*person.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedParty, utils.None[uuid.UUID]())
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to save versioned party: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, personVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to save person: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return person, nil
}

func (s *Service) GetCurrentPersonVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (rm.PERSON, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var person rm.PERSON
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), rm.PERSON_TYPE).Scan(&person)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.PERSON{}, ErrPersonNotFound
		}
		return rm.PERSON{}, fmt.Errorf("failed to get latest person by versioned party ID: %w", err)
	}

	return person, nil
}

func (s *Service) GetPersonAtVersion(ctx context.Context, versionID string) (rm.PERSON, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var person rm.PERSON
	err := s.DB.QueryRow(ctx, query, versionID, rm.GROUP_TYPE).Scan(&person)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.PERSON{}, ErrPersonNotFound
		}
		return rm.PERSON{}, fmt.Errorf("failed to get person at version: %w", err)
	}

	return person, nil
}

func (s *Service) UpdatePerson(ctx context.Context, previousPerson, nextPerson rm.PERSON) (rm.PERSON, error) {
	err := s.ValidatePerson(ctx, nextPerson)
	if err != nil {
		return rm.PERSON{}, err
	}

	if !nextPerson.UID.E {
		nextPerson.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: previousPerson.UID.V.HIER_OBJECT_ID().Value,
		}))
	}

	updatedID, err := UpgradeObjectVersionID(previousPerson.UID.V, *previousPerson.UID.V.OBJECT_VERSION_ID())
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to upgrade current Person UID: %w", err)
	}
	nextPerson.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	personVersion := NewOriginalVersion(*nextPerson.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_PERSON(nextPerson), utils.Some(*previousPerson.UID.V.OBJECT_VERSION_ID()))
	contribution := NewContribution("Person updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*nextPerson.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, personVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to save person: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.PERSON{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextPerson, nil
}

func (s *Service) DeletePerson(ctx context.Context, versionedObjectID uuid.UUID) error {
	// Check if person exists
	_, err := s.GetCurrentPersonVersionByVersionedPartyID(ctx, versionedObjectID)
	if err != nil {
		if errors.Is(err, ErrPersonNotFound) {
			return ErrPersonNotFound
		}
		return fmt.Errorf("failed to get person before deletion: %w", err)
	}

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
	groupVersion := NewOriginalVersion(*group.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_GROUP(group), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Group created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*group.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedParty, utils.None[uuid.UUID]())
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to save versioned party: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, groupVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to save group: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return group, nil
}

func (s *Service) GetCurrentGroupVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (rm.GROUP, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var group rm.GROUP
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), rm.GROUP_TYPE).Scan(&group)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.GROUP{}, ErrGroupNotFound
		}
		return rm.GROUP{}, fmt.Errorf("failed to get latest group by versioned party ID: %w", err)
	}

	return group, nil
}

func (s *Service) GetGroupAtVersion(ctx context.Context, versionID string) (rm.GROUP, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var group rm.GROUP
	err := s.DB.QueryRow(ctx, query, versionID, rm.GROUP_TYPE).Scan(&group)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.GROUP{}, ErrGroupNotFound
		}
		return rm.GROUP{}, fmt.Errorf("failed to get group at version: %w", err)
	}

	return group, nil
}

func (s *Service) UpdateGroup(ctx context.Context, currentGroup, nextGroup rm.GROUP) (rm.GROUP, error) {
	err := s.ValidateGroup(ctx, nextGroup)
	if err != nil {
		return rm.GROUP{}, err
	}

	if !currentGroup.UID.E {
		nextGroup.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: currentGroup.UID.V.OBJECT_VERSION_ID().UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(nextGroup.UID.V, *currentGroup.UID.V.OBJECT_VERSION_ID())
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to upgrade current Group UID: %w", err)
	}
	nextGroup.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	groupVersion := NewOriginalVersion(*nextGroup.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_GROUP(nextGroup), utils.Some(*currentGroup.UID.V.OBJECT_VERSION_ID()))
	contribution := NewContribution("Group updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*nextGroup.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, groupVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to save group: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.GROUP{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextGroup, nil
}

func (s *Service) DeleteGroup(ctx context.Context, versionedObjectID uuid.UUID) error {
	contribution := NewContribution("Group deleted", terminology.AUDIT_CHANGE_TYPE_CODE_DELETED,
		[]rm.OBJECT_REF{
			{
				Namespace: config.NAMESPACE_LOCAL,
				Type:      rm.GROUP_TYPE,
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
	organisationVersion := NewOriginalVersion(*organisation.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_ORGANISATION(organisation), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Organisation created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*organisation.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedParty, utils.None[uuid.UUID]())
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to save versioned party: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, organisationVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to save organisation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return organisation, nil
}

func (s *Service) GetCurrentOrganisationVersionByVersionedPartyID(ctx context.Context, versionedPartyID uuid.UUID) (rm.ORGANISATION, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.versioned_object_id = $1 AND ov.type = $2
		ORDER BY ov.created_at DESC
		LIMIT 1
	`

	var organisation rm.ORGANISATION
	err := s.DB.QueryRow(ctx, query, versionedPartyID.String(), rm.AGENT_TYPE).Scan(&organisation)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.ORGANISATION{}, ErrOrganisationNotFound
		}
		return rm.ORGANISATION{}, fmt.Errorf("failed to get latest organisation by versioned party ID: %w", err)
	}

	return organisation, nil
}

func (s *Service) GetOrganisationAtVersion(ctx context.Context, versionID string) (rm.ORGANISATION, error) {
	query := `
		SELECT ovd.object_data 
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

	var organisation rm.ORGANISATION
	err := s.DB.QueryRow(ctx, query, versionID, rm.ORGANISATION_TYPE).Scan(&organisation)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return rm.ORGANISATION{}, ErrOrganisationNotFound
		}
		return rm.ORGANISATION{}, fmt.Errorf("failed to get organisation at version: %w", err)
	}

	return organisation, nil
}

func (s *Service) UpdateOrganisation(ctx context.Context, currentOrganisation, nextOrganisation rm.ORGANISATION) (rm.ORGANISATION, error) {
	err := s.ValidateOrganisation(ctx, nextOrganisation)
	if err != nil {
		return rm.ORGANISATION{}, err
	}

	if !nextOrganisation.UID.E {
		nextOrganisation.UID = utils.Some(rm.UID_BASED_ID_from_HIER_OBJECT_ID(&rm.HIER_OBJECT_ID{
			Value: currentOrganisation.UID.V.Value.(*rm.OBJECT_VERSION_ID).UID(),
		}))
	}

	updatedID, err := UpgradeObjectVersionID(nextOrganisation.UID.V, *currentOrganisation.UID.V.OBJECT_VERSION_ID())
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to upgrade current Organisation UID: %w", err)
	}
	nextOrganisation.UID = utils.Some(rm.UID_BASED_ID_from_OBJECT_VERSION_ID(&updatedID))

	organisationVersion := NewOriginalVersion(*nextOrganisation.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_ORGANISATION(nextOrganisation), utils.Some(*currentOrganisation.UID.V.OBJECT_VERSION_ID()))
	contribution := NewContribution("Organisation updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*currentOrganisation.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, organisationVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to save organisation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rm.ORGANISATION{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextOrganisation, nil
}

func (s *Service) DeleteOrganisation(ctx context.Context, versionedObjectID uuid.UUID) error {
	// Check if organisation exists
	_, err := s.GetCurrentOrganisationVersionByVersionedPartyID(ctx, versionedObjectID)
	if err != nil {
		if errors.Is(err, ErrOrganisationNotFound) {
			return ErrOrganisationNotFound
		}
		return fmt.Errorf("failed to get organisation before deletion: %w", err)
	}

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
	roleVersion := NewOriginalVersion(*role.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_ROLE(role), utils.None[rm.OBJECT_VERSION_ID]())
	contribution := NewContribution("Role created", terminology.AUDIT_CHANGE_TYPE_CODE_CREATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*role.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveVersionedObjectWithTx(ctx, tx, versionedParty, utils.None[uuid.UUID]())
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to save versioned party: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, roleVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to save role: %w", err)
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

func (s *Service) GetRoleAtVersion(ctx context.Context, versionID string) (rm.ROLE, error) {
	query := `
		SELECT ovd.object_data
		FROM openehr.tbl_object_version ov
		JOIN openehr.tbl_object_version_data ovd ON ov.id = ovd.id 
		WHERE ov.id = $1 AND ov.type = $2
		LIMIT 1
	`

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

	roleVersion := NewOriginalVersion(*role.UID.V.OBJECT_VERSION_ID(), rm.ORIGINAL_VERSION_DATA_from_ROLE(role), utils.Some(*currentRoleID))
	contribution := NewContribution("Role updated", terminology.AUDIT_CHANGE_TYPE_CODE_MODIFICATION,
		[]rm.OBJECT_REF{
			{
				Type:      rm.VERSIONED_PARTY_TYPE,
				Namespace: rm.Namespace_local,
				ID:        rm.OBJECT_ID_from_OBJECT_VERSION_ID(*role.UID.V.OBJECT_VERSION_ID()),
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

	err = s.SaveContributionWithTx(ctx, tx, contribution, utils.None[uuid.UUID]())
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to save contribution: %w", err)
	}
	err = s.SaveObjectVersionWithTx(ctx, tx, roleVersion, contribution.UID.Value, utils.None[uuid.UUID]())
	if err != nil {
		return rm.ROLE{}, fmt.Errorf("failed to save role: %w", err)
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

func (s *Service) GetVersionedParty(ctx context.Context, versionedObjectID uuid.UUID) (rm.VERSIONED_PARTY, error) {
	query := `
		SELECT vod.data 
		FROM openehr.tbl_versioned_object vo
		JOIN openehr.tbl_versioned_object_data vod ON vo.id = vod.id
		WHERE vo.ehr_id IS NULL AND vo.type = $1 AND vo.id = $2 
		LIMIT 1
	`

	var versionedParty rm.VERSIONED_PARTY
	err := s.DB.QueryRow(ctx, query, rm.VERSIONED_PARTY_TYPE, versionedObjectID).Scan(&versionedParty)
	if err != nil {
		return rm.VERSIONED_PARTY{}, fmt.Errorf("failed to get versioned party by ID: %w", err)
	}
	return versionedParty, nil
}

func (s *Service) GetVersionedPartyRevisionHistory(ctx context.Context, versionedObjectID uuid.UUID) (rm.REVISION_HISTORY, error) {
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
	args := []any{[]string{rm.AGENT_TYPE, rm.PERSON_TYPE, rm.GROUP_TYPE, rm.ORGANISATION_TYPE, rm.ROLE_TYPE}, versionedObjectID}
	row := s.DB.QueryRow(ctx, query, args...)

	var revisionHistory rm.REVISION_HISTORY
	if err := row.Scan(&revisionHistory); err != nil {
		if err == database.ErrNoRows {
			return rm.REVISION_HISTORY{}, ErrRevisionHistoryNotFound
		}
		return rm.REVISION_HISTORY{}, fmt.Errorf("failed to fetch Revision History from database: %w", err)
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

func (s *Service) GetContribution(ctx context.Context, contributionID string, ehrID utils.Optional[uuid.UUID]) (rm.CONTRIBUTION, error) {
	query := `
		SELECT cd.data
		FROM openehr.tbl_contribution c
		JOIN openehr.tbl_contribution_data cd ON c.id = cd.id
		WHERE c.ehr_id = $1 AND c.id = $2
		LIMIT 1
	`
	args := []any{ehrID, contributionID}
	row := s.DB.QueryRow(ctx, query, args...)

	var contribution rm.CONTRIBUTION
	if err := row.Scan(&contribution); err != nil {
		if err == database.ErrNoRows {
			return rm.CONTRIBUTION{}, ErrContributionNotFound
		}
		return rm.CONTRIBUTION{}, fmt.Errorf("failed to fetch Contribution by ID from database: %w", err)
	}

	return contribution, nil
}

func (s *Service) ExistsVersionedObject(ctx context.Context, versionedObjectID uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM openehr.tbl_versioned_object vo WHERE vo.id = $1 LIMIT 1`
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

func (s *Service) SaveEHRWithTx(ctx context.Context, tx pgx.Tx, ehrID uuid.UUID) error {
	query := `INSERT INTO openehr.tbl_ehr (id) VALUES ($1)`
	args := []any{ehrID}
	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert EHR into the database: %w", err)
	}

	return nil
}

func (s *Service) SaveContributionWithTx(ctx context.Context, tx pgx.Tx, contribution rm.CONTRIBUTION, ehrID utils.Optional[uuid.UUID]) error {
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
	case rm.VERSIONED_EHR_STATUS:
		v.SetModelName()
		modelType = rm.VERSIONED_EHR_STATUS_TYPE
		id = v.UID.Value
	case rm.VERSIONED_EHR_ACCESS:
		v.SetModelName()
		modelType = rm.VERSIONED_EHR_ACCESS_TYPE
		id = v.UID.Value
	case rm.VERSIONED_COMPOSITION:
		v.SetModelName()
		modelType = rm.VERSIONED_COMPOSITION_TYPE
		id = v.UID.Value
	case rm.VERSIONED_FOLDER:
		v.SetModelName()
		modelType = rm.VERSIONED_FOLDER_TYPE
		id = v.UID.Value
	case rm.VERSIONED_PARTY:
		v.SetModelName()
		modelType = rm.VERSIONED_PARTY_TYPE
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
	var data rm.OriginalVersionDataUnion
	switch v := version.(type) {
	case rm.ORIGINAL_VERSION:
		v.SetModelName()
		data = v.Data
	// After enabling, make sure to change the data path below in the INSERT statement
	// case rm.IMPORTED_VERSION:
	// 	object = v.Data
	default:
		return fmt.Errorf("unsupported version type for object version creation: %T", version)
	}

	var (
		modelType string
		id        rm.OBJECT_VERSION_ID
	)
	switch data.Kind {
	case rm.ORIGINAL_VERSION_data_kind_EHR_STATUS:
		modelType = rm.EHR_STATUS_TYPE
		id = *data.EHR_STATUS().UID.V.OBJECT_VERSION_ID()
	case rm.ORIGINAL_VERSION_data_kind_EHR_ACCESS:
		modelType = rm.EHR_ACCESS_TYPE
		id = *data.EHR_ACCESS().UID.V.OBJECT_VERSION_ID()
	case rm.ORIGINAL_VERSION_data_kind_COMPOSITION:
		modelType = rm.COMPOSITION_TYPE
		id = *data.COMPOSITION().UID.V.OBJECT_VERSION_ID()
	case rm.ORIGINAL_VERSION_data_kind_FOLDER:
		modelType = rm.FOLDER_TYPE
		id = *data.FOLDER().UID.V.OBJECT_VERSION_ID()
	case rm.ORIGINAL_VERSION_data_kind_ROLE:
		modelType = rm.ROLE_TYPE
		id = *data.ROLE().UID.V.OBJECT_VERSION_ID()
	case rm.ORIGINAL_VERSION_data_kind_PERSON:
		modelType = rm.PERSON_TYPE
		id = *data.PERSON().UID.V.OBJECT_VERSION_ID()
	case rm.ORIGINAL_VERSION_data_kind_AGENT:
		modelType = rm.AGENT_TYPE
		id = *data.AGENT().UID.V.OBJECT_VERSION_ID()
	case rm.ORIGINAL_VERSION_data_kind_GROUP:
		modelType = rm.GROUP_TYPE
		id = *data.GROUP().UID.V.OBJECT_VERSION_ID()
	case rm.ORIGINAL_VERSION_data_kind_ORGANISATION:
		modelType = rm.ORGANISATION_TYPE
		id = *data.ORGANISATION().UID.V.OBJECT_VERSION_ID()
	default:
		return fmt.Errorf("unsupported object type for version creation: %d", data.Kind)
	}

	query := `INSERT INTO openehr.tbl_object_version (id, versioned_object_id, type, ehr_id, contribution_id) VALUES ($1, $2, $3, $4, $5)`
	args := []any{id.Value, id.UID(), modelType, ehrID, contributionID}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert object version into the database: %w", err)
	}

	query = `INSERT INTO openehr.tbl_object_version_data (id, version_data, object_data) VALUES ($1, jsonb_set($2, '{data}', 'null', true), $2->'data')`
	args = []any{id.Value, version}
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
		return *current.OBJECT_VERSION_ID(), nil
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
	default:
		return rm.CONTRIBUTION{}
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
