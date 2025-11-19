CREATE SCHEMA openehr;

CREATE TABLE openehr.tbl_template (
    id TEXT PRIMARY KEY,
    concept TEXT NOT NULL,
    archetype_id TEXT NOT NULL,
    data JSONB NOT NULL,
    raw BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_template_archetype_id ON openehr.tbl_template(archetype_id);

CREATE TABLE openehr.tbl_ehr (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ehr_data ON openehr.tbl_ehr USING GIN (data);
CREATE INDEX idx_ehr_created_at ON openehr.tbl_ehr(created_at);
CREATE INDEX idx_ehr_data_type ON openehr.tbl_ehr USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_contribution (
    id TEXT PRIMARY KEY,
    ehr_id TEXT REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contribution_ehr_id ON openehr.tbl_contribution(ehr_id);
CREATE INDEX idx_contribution_data ON openehr.tbl_contribution USING GIN (data);
CREATE INDEX idx_contribution_created_at ON openehr.tbl_contribution(created_at);
CREATE INDEX idx_contribution_data_type ON openehr.tbl_contribution USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_versioned_object (
    id TEXT PRIMARY KEY,
    ehr_id TEXT REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    data JSONB NOT NULL,
    is_latest BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_versioned_object_ehr_id ON openehr.tbl_versioned_object(ehr_id);
CREATE INDEX idx_versioned_object_type ON openehr.tbl_versioned_object(type);
CREATE INDEX idx_versioned_object_data ON openehr.tbl_versioned_object USING GIN (data);
CREATE INDEX idx_versioned_object_created_at ON openehr.tbl_versioned_object(created_at);
CREATE INDEX idx_versioned_object_data_type ON openehr.tbl_versioned_object USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_ehr_status (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    ehr_id TEXT NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    is_latest BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ehr_status_versioned_object_id ON openehr.tbl_ehr_status(versioned_object_id);
CREATE INDEX idx_ehr_status_data ON openehr.tbl_ehr_status USING GIN (data);
CREATE INDEX idx_ehr_status_ehr_id ON openehr.tbl_ehr_status(ehr_id);
CREATE INDEX idx_ehr_status_created_at ON openehr.tbl_ehr_status(created_at);
CREATE INDEX idx_ehr_status_data_type ON openehr.tbl_ehr_status USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_ehr_access (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    ehr_id TEXT NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    is_latest BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ehr_access_versioned_object_id ON openehr.tbl_ehr_access(versioned_object_id);
CREATE INDEX idx_ehr_access_data ON openehr.tbl_ehr_access USING GIN (data);
CREATE INDEX idx_ehr_access_ehr_id ON openehr.tbl_ehr_access(ehr_id);
CREATE INDEX idx_ehr_access_created_at ON openehr.tbl_ehr_access(created_at);
CREATE INDEX idx_ehr_access_data_type ON openehr.tbl_ehr_access USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_composition (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    ehr_id TEXT NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    is_latest BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_composition_versioned_object_id ON openehr.tbl_composition(versioned_object_id);
CREATE INDEX idx_composition_data ON openehr.tbl_composition USING GIN (data);
CREATE INDEX idx_composition_ehr_id ON openehr.tbl_composition(ehr_id);
CREATE INDEX idx_composition_created_at ON openehr.tbl_composition(created_at);
CREATE INDEX idx_composition_data_type ON openehr.tbl_composition USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_folder (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    ehr_id TEXT NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    is_latest BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_folder_versioned_object_id ON openehr.tbl_folder(versioned_object_id);
CREATE INDEX idx_folder_data ON openehr.tbl_folder USING GIN (data);
CREATE INDEX idx_folder_ehr_id ON openehr.tbl_folder(ehr_id);
CREATE INDEX idx_folder_created_at ON openehr.tbl_folder(created_at);
CREATE INDEX idx_folder_data_type ON openehr.tbl_folder USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_role (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    is_latest BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_role_versioned_object_id ON openehr.tbl_role(versioned_object_id);
CREATE INDEX idx_role_data ON openehr.tbl_role USING GIN (data);
CREATE INDEX idx_role_created_at ON openehr.tbl_role(created_at);
CREATE INDEX idx_role_data_type ON openehr.tbl_role USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_person (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    is_latest BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_person_versioned_object_id ON openehr.tbl_person(versioned_object_id);
CREATE INDEX idx_person_data ON openehr.tbl_person USING GIN (data);
CREATE INDEX idx_person_created_at ON openehr.tbl_person(created_at);
CREATE INDEX idx_person_data_type ON openehr.tbl_person USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_agent (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    is_latest BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agent_versioned_object_id ON openehr.tbl_agent(versioned_object_id);
CREATE INDEX idx_agent_data ON openehr.tbl_agent USING GIN (data);
CREATE INDEX idx_agent_created_at ON openehr.tbl_agent(created_at);
CREATE INDEX idx_agent_data_type ON openehr.tbl_agent USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_group (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    is_latest BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_group_versioned_object_id ON openehr.tbl_group(versioned_object_id);
CREATE INDEX idx_group_data ON openehr.tbl_group USING GIN (data);
CREATE INDEX idx_group_created_at ON openehr.tbl_group(created_at);
CREATE INDEX idx_group_data_type ON openehr.tbl_group USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_organisation (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    is_latest BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_organisation_versioned_object_id ON openehr.tbl_organisation(versioned_object_id);
CREATE INDEX idx_organisation_data ON openehr.tbl_organisation USING GIN (data);
CREATE INDEX idx_organisation_created_at ON openehr.tbl_organisation(created_at);
CREATE INDEX idx_organisation_data_type ON openehr.tbl_organisation USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE openehr.tbl_query (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    query TEXT NOT NULL,
    version TEXT NOT NULL,
    is_latest BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (name, version)
);

CREATE INDEX idx_query_name ON openehr.tbl_query(name);
CREATE INDEX idx_query_created_at ON openehr.tbl_query(created_at);