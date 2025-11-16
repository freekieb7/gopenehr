CREATE TABLE tbl_openehr_template (
    id TEXT PRIMARY KEY,
    concept TEXT NOT NULL,
    archetype_id TEXT NOT NULL,
    data JSONB NOT NULL,
    raw BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_template_archetype_id ON tbl_openehr_template(archetype_id);

CREATE TABLE tbl_openehr_ehr (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_ehr_data ON tbl_openehr_ehr USING GIN (data);
CREATE INDEX idx_openehr_ehr_created_at ON tbl_openehr_ehr(created_at);
CREATE INDEX idx_openehr_ehr_data_type ON tbl_openehr_ehr USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE tbl_openehr_contribution (
    id TEXT PRIMARY KEY,
    ehr_id TEXT REFERENCES tbl_openehr_ehr(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_contribution_ehr_id ON tbl_openehr_contribution(ehr_id);
CREATE INDEX idx_openehr_contribution_data ON tbl_openehr_contribution USING GIN (data);
CREATE INDEX idx_openehr_contribution_created_at ON tbl_openehr_contribution(created_at);
CREATE INDEX idx_openehr_contribution_data_type ON tbl_openehr_contribution USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

-- CREATE TABLE tbl_openehr_contribution_version (
--     contribution_id TEXT REFERENCES tbl_openehr_contribution(id) ON DELETE CASCADE,
--     object_id TEXT NOT NULL
-- );

-- CREATE INDEX idx_openehr_contribution_version_contribution_id ON tbl_openehr_contribution_version(contribution_id);
-- CREATE INDEX idx_openehr_contribution_version_object_id ON tbl_openehr_contribution_version(object_id);

CREATE TABLE tbl_openehr_versioned_object (
    id TEXT PRIMARY KEY,
    ehr_id TEXT REFERENCES tbl_openehr_ehr(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_versioned_object_ehr_id ON tbl_openehr_versioned_object(ehr_id);
CREATE INDEX idx_openehr_versioned_object_type ON tbl_openehr_versioned_object(type);
CREATE INDEX idx_openehr_versioned_object_data ON tbl_openehr_versioned_object USING GIN (data);
CREATE INDEX idx_openehr_versioned_object_created_at ON tbl_openehr_versioned_object(created_at);
CREATE INDEX idx_openehr_versioned_object_data_type ON tbl_openehr_versioned_object USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE tbl_openehr_ehr_status (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES tbl_openehr_versioned_object(id) ON DELETE CASCADE,
    ehr_id TEXT NOT NULL REFERENCES tbl_openehr_ehr(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_ehr_status_versioned_object_id ON tbl_openehr_ehr_status(versioned_object_id);
CREATE INDEX idx_openehr_ehr_status_ehr_id ON tbl_openehr_ehr_status(ehr_id);
CREATE INDEX idx_openehr_ehr_status_data ON tbl_openehr_ehr_status USING GIN (data);
CREATE INDEX idx_openehr_ehr_status_created_at ON tbl_openehr_ehr_status(created_at);
CREATE INDEX idx_openehr_ehr_status_data_type ON tbl_openehr_ehr_status USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE tbl_openehr_ehr_access (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES tbl_openehr_versioned_object(id) ON DELETE CASCADE,
    ehr_id TEXT NOT NULL REFERENCES tbl_openehr_ehr(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_ehr_access_versioned_object_id ON tbl_openehr_ehr_access(versioned_object_id);
CREATE INDEX idx_openehr_ehr_access_ehr_id ON tbl_openehr_ehr_access(ehr_id);
CREATE INDEX idx_openehr_ehr_access_data ON tbl_openehr_ehr_access USING GIN (data);
CREATE INDEX idx_openehr_ehr_access_created_at ON tbl_openehr_ehr_access(created_at);
CREATE INDEX idx_openehr_ehr_access_data_type ON tbl_openehr_ehr_access USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE tbl_openehr_composition (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES tbl_openehr_versioned_object(id) ON DELETE CASCADE,
    ehr_id TEXT NOT NULL REFERENCES tbl_openehr_ehr(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_composition_versioned_object_id ON tbl_openehr_composition(versioned_object_id);
CREATE INDEX idx_openehr_composition_ehr_id ON tbl_openehr_composition(ehr_id);
CREATE INDEX idx_openehr_composition_data ON tbl_openehr_composition USING GIN (data);
CREATE INDEX idx_openehr_composition_created_at ON tbl_openehr_composition(created_at);
CREATE INDEX idx_openehr_composition_data_type ON tbl_openehr_composition USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE tbl_openehr_folder (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES tbl_openehr_versioned_object(id) ON DELETE CASCADE,
    ehr_id TEXT NOT NULL REFERENCES tbl_openehr_ehr(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_folder_versioned_object_id ON tbl_openehr_folder(versioned_object_id);
CREATE INDEX idx_openehr_folder_ehr_id ON tbl_openehr_folder(ehr_id);
CREATE INDEX idx_openehr_folder_data ON tbl_openehr_folder USING GIN (data);
CREATE INDEX idx_openehr_folder_created_at ON tbl_openehr_folder(created_at);
CREATE INDEX idx_openehr_folder_data_type ON tbl_openehr_folder USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

-- CREATE TABLE tbl_openehr_folder_item (
--     folder_id TEXT NOT NULL REFERENCES tbl_openehr_folder(id) ON DELETE CASCADE,
--     object_id TEXT NOT NULL
-- );

-- CREATE INDEX idx_openehr_folder_item_folder_id ON tbl_openehr_folder_item(folder_id);
-- CREATE INDEX idx_openehr_folder_item_object_id ON tbl_openehr_folder_item(object_id);

CREATE TABLE tbl_openehr_role (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES tbl_openehr_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_role_versioned_object_id ON tbl_openehr_role(versioned_object_id);
CREATE INDEX idx_openehr_role_data ON tbl_openehr_role USING GIN (data);
CREATE INDEX idx_openehr_role_created_at ON tbl_openehr_role(created_at);
CREATE INDEX idx_openehr_role_data_type ON tbl_openehr_role USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE tbl_openehr_person (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES tbl_openehr_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_person_versioned_object_id ON tbl_openehr_person(versioned_object_id);
CREATE INDEX idx_openehr_person_data ON tbl_openehr_person USING GIN (data);
CREATE INDEX idx_openehr_person_created_at ON tbl_openehr_person(created_at);
CREATE INDEX idx_openehr_person_data_type ON tbl_openehr_person USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE tbl_openehr_agent (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES tbl_openehr_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_agent_versioned_object_id ON tbl_openehr_agent(versioned_object_id);
CREATE INDEX idx_openehr_agent_data ON tbl_openehr_agent USING GIN (data);
CREATE INDEX idx_openehr_agent_created_at ON tbl_openehr_agent(created_at);
CREATE INDEX idx_openehr_agent_data_type ON tbl_openehr_agent USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE tbl_openehr_group (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES tbl_openehr_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_group_versioned_object_id ON tbl_openehr_group(versioned_object_id);
CREATE INDEX idx_openehr_group_data ON tbl_openehr_group USING GIN (data);
CREATE INDEX idx_openehr_group_created_at ON tbl_openehr_group(created_at);
CREATE INDEX idx_openehr_group_data_type ON tbl_openehr_group USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

CREATE TABLE tbl_openehr_organisation (
    id TEXT PRIMARY KEY,
    versioned_object_id TEXT NOT NULL REFERENCES tbl_openehr_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_openehr_organisation_versioned_object_id ON tbl_openehr_organisation(versioned_object_id);
CREATE INDEX idx_openehr_organisation_data ON tbl_openehr_organisation USING GIN (data);
CREATE INDEX idx_openehr_organisation_created_at ON tbl_openehr_organisation(created_at);
CREATE INDEX idx_openehr_organisation_data_type ON tbl_openehr_organisation USING GIN (
    jsonb_path_query_array(data, '$.**._type')
);

-- CREATE TABLE tbl_patient_ehr_link (
--     patient_id TEXT NOT NULL REFERENCES tbl_openehr_person(id) ON DELETE CASCADE,
--     ehr_id TEXT NOT NULL REFERENCES tbl_openehr_ehr(id) ON DELETE CASCADE,
--     PRIMARY KEY (patient_id, ehr_id)
-- );