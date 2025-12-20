package migration

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var _ Migration = (*SetupOpenEHR)(nil)

type SetupOpenEHR struct{}

func (m *SetupOpenEHR) Version() uint64 {
	return 20251113195000
}

func (m *SetupOpenEHR) Name() string {
	return "Setup OpenEHR Schema and Tables"
}

func (m *SetupOpenEHR) Up(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `CREATE SCHEMA openehr;`)
	if err != nil {
		return fmt.Errorf("failed to create openehr schema: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_ehr (
			id UUID PRIMARY KEY,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_ehr table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_ehr_data (
			id UUID PRIMARY KEY REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_ehr_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_ehr_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_ehr_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_contribution (
			id UUID PRIMARY KEY,
			ehr_id UUID REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_contribution table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_contribution_data (
			id UUID PRIMARY KEY REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
			data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_contribution_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_contribution_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_contribution_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_versioned_object (
			id UUID PRIMARY KEY,
			type TEXT NOT NULL,
			ehr_id UUID REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_versioned_object table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_versioned_object_data (
			id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_versioned_object_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_versioned_object_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_versioned_object_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_versioned_ehr_status_tag (
			versioned_ehr_status_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			key TEXT NOT NULL,
			data JSONB NOT NULL,
			ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			PRIMARY KEY (versioned_ehr_status_id, key)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_versioned_ehr_status_tag table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_versioned_ehr_access_data (
			id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_versioned_ehr_access_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_versioned_ehr_access_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_versioned_ehr_access_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_versioned_composition_data (
			id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_versioned_composition_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_versioned_composition_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_versioned_composition_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_versioned_composition_tag (
			versioned_composition_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			key TEXT NOT NULL,
			data JSONB NOT NULL,
			ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			PRIMARY KEY (versioned_composition_id, key)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_versioned_composition_tag table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_versioned_folder_data (
			id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_versioned_folder_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_versioned_folder_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_versioned_folder_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_versioned_party_data (
			id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_versioned_party_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_versioned_party_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_versioned_party_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_versioned_party_tag (
			versioned_party_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			key TEXT NOT NULL,
			data JSONB NOT NULL,
			ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			PRIMARY KEY (versioned_party_id, key)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_versioned_party_tag table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_ehr_status (
			id TEXT PRIMARY KEY,
			version_int int NOT NULL,
			versioned_ehr_status_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
			local_ref_versioned_party_id UUID REFERENCES openehr.tbl_versioned_object(id) ON DELETE SET NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_ehr_status table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_ehr_status_data (
			id TEXT PRIMARY KEY,
			data JSONB NOT NULL,
			version_data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_ehr_status_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_ehr_status_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_ehr_status_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_ehr_status_data ALTER COLUMN version_data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_ehr_status_data.version_data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_ehr_status_tag (
			ehr_status_id TEXT NOT NULL REFERENCES openehr.tbl_ehr_status(id) ON DELETE CASCADE,
			key TEXT NOT NULL,
			data JSONB NOT NULL,
			ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			PRIMARY KEY (ehr_status_id, key)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_ehr_status_tag table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_ehr_access (
			id TEXT PRIMARY KEY,
			version_int int NOT NULL,
			versioned_ehr_access_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_ehr_access table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_ehr_access_data (
			id TEXT PRIMARY KEY,
			data JSONB NOT NULL,
			version_data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_ehr_access_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_ehr_access_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_ehr_access_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_ehr_access_data ALTER COLUMN version_data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_ehr_access_data.version_data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_composition (
			id TEXT PRIMARY KEY,
			version_int INT NOT NULL,
			versioned_composition_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_composition table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_composition_data (
			id TEXT PRIMARY KEY,
			data JSONB NOT NULL,
			version_data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_composition_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_composition_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_composition_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_composition_data ALTER COLUMN version_data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_composition_data.version_data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_composition_tag (
			composition_id TEXT NOT NULL REFERENCES openehr.tbl_composition(id) ON DELETE CASCADE,
			key TEXT NOT NULL,
			data JSONB NOT NULL,
			ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			PRIMARY KEY (composition_id, key)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_composition_tag table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_folder (
			id TEXT PRIMARY KEY,
			version_int INT NOT NULL,
			versioned_folder_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
			contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_folder table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_folder_data (
			id TEXT PRIMARY KEY,
			data JSONB NOT NULL,
			version_data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_folder_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_folder_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_folder_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_folder_data ALTER COLUMN version_data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_folder_data.version_data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_person (
			id TEXT PRIMARY KEY,
			version_int INT NOT NULL,
			versioned_party_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_person table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_person_data (
			id TEXT PRIMARY KEY,
			data JSONB NOT NULL,
			version_data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_person_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_person_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_person_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_person_data ALTER COLUMN version_data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_person_data.version_data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_person_tag (
			person_id TEXT NOT NULL REFERENCES openehr.tbl_person(id) ON DELETE CASCADE,
			key TEXT NOT NULL,
			data JSONB NOT NULL,
			PRIMARY KEY (person_id, key)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_person_tag table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_agent (
			id TEXT PRIMARY KEY,
			version_int INT NOT NULL,
			versioned_party_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_agent table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_agent_data (
			id TEXT PRIMARY KEY,
			data JSONB NOT NULL,
			version_data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_agent_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_agent_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_agent_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_agent_data ALTER COLUMN version_data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_agent_data.version_data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_agent_tag (
			agent_id TEXT NOT NULL REFERENCES openehr.tbl_agent(id) ON DELETE CASCADE,
			key TEXT NOT NULL,
			data JSONB NOT NULL,
			PRIMARY KEY (agent_id, key)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_agent_tag table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_group (
			id TEXT PRIMARY KEY,
			version_int INT NOT NULL,
			versioned_party_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_group table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_group_data (
			id TEXT PRIMARY KEY,
			data JSONB NOT NULL,
			version_data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_group_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_group_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_group_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_group_data ALTER COLUMN version_data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_group_data.version_data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_group_tag (
			group_id TEXT NOT NULL REFERENCES openehr.tbl_group(id) ON DELETE CASCADE,
			key TEXT NOT NULL,
			data JSONB NOT NULL,
			PRIMARY KEY (group_id, key)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_group_tag table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_organisation (
			id TEXT PRIMARY KEY,
			version_int INT NOT NULL,
			versioned_party_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
			contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_organisation table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_organisation_data (
			id TEXT PRIMARY KEY,
			data JSONB NOT NULL,
			version_data JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_organisation_data table: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_organisation_data ALTER COLUMN data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_organisation_data.data column: %w", err)
	}

	_, err = tx.Exec(ctx, `ALTER TABLE openehr.tbl_organisation_data ALTER COLUMN version_data SET COMPRESSION lz4;`)
	if err != nil {
		return fmt.Errorf("failed to set compression on tbl_organisation_data.version_data column: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_organisation_tag (
			organisation_id TEXT NOT NULL REFERENCES openehr.tbl_organisation(id) ON DELETE CASCADE,
			key TEXT NOT NULL,
			data JSONB NOT NULL,
			PRIMARY KEY (organisation_id, key)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_organisation_tag table: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE openehr.tbl_query (
			name TEXT NOT NULL,
			version INT NOT NULL,
			query TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (name, version)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tbl_query table: %w", err)
	}

	return nil
}

func (m *SetupOpenEHR) Down(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `DROP SCHEMA openehr CASCADE;`)
	if err != nil {
		return fmt.Errorf("failed to drop openehr schema: %w", err)
	}

	return nil
}
