CREATE SCHEMA openehr;

-- ========= EHR Table ==========

CREATE TABLE openehr.tbl_ehr (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_ehr_data (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_ehr_data ALTER COLUMN data SET COMPRESSION lz4;

-- ========= Contribution Table ==========

CREATE TABLE openehr.tbl_contribution (
    id UUID PRIMARY KEY,
    ehr_id UUID REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_contribution_data (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_contribution_data ALTER COLUMN data SET COMPRESSION lz4;

-- ========= Versioned Object Table ==========

CREATE TABLE openehr.tbl_versioned_object (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    ehr_id UUID REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_versioned_ehr_status (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_versioned_ehr_status ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_versioned_ehr_status_tag (
    versioned_ehr_status_id UUID NOT NULL REFERENCES openehr.tbl_versioned_ehr_status(id) ON DELETE CASCADE,
    key TEXT NOT NULL,
    data JSONB NOT NULL,
    ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    PRIMARY KEY (versioned_ehr_status_id, key)
);

CREATE TABLE openehr.tbl_versioned_ehr_access (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_versioned_ehr_access ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_versioned_composition (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_versioned_composition ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_versioned_composition_tag (
    versioned_composition_id UUID NOT NULL REFERENCES openehr.tbl_versioned_composition(id) ON DELETE CASCADE,
    key TEXT NOT NULL,
    data JSONB NOT NULL,
    ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    PRIMARY KEY (versioned_composition_id, key)
);

CREATE TABLE openehr.tbl_versioned_folder (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_versioned_folder ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_versioned_party (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_versioned_party ALTER COLUMN data SET COMPRESSION lz4;

-- ========= Object Tables ==========

CREATE TABLE openehr.tbl_ehr_status (
    id TEXT PRIMARY KEY,
    version_int int NOT NULL,
    versioned_ehr_status_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_ehr_status_data (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    version_data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_ehr_status_data ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_ehr_status_tag (
    ehr_status_id TEXT NOT NULL REFERENCES openehr.tbl_ehr_status(id) ON DELETE CASCADE,
    key TEXT NOT NULL,
    data JSONB NOT NULL,
    ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    PRIMARY KEY (ehr_status_id, key)
);

CREATE TABLE openehr.tbl_ehr_access (
    id TEXT PRIMARY KEY,
    version_int int NOT NULL,
    versioned_ehr_access_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_ehr_access_data (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    version_data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_ehr_access_data ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_composition (
    id TEXT PRIMARY KEY,
    version_int INT NOT NULL,
    versioned_composition_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_composition_data (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    version_data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_composition_data ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_composition_tag (
    composition_id TEXT NOT NULL REFERENCES openehr.tbl_composition(id) ON DELETE CASCADE,
    key TEXT NOT NULL,
    data JSONB NOT NULL,
    ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    PRIMARY KEY (composition_id, key)
);

CREATE TABLE openehr.tbl_folder (
    id TEXT PRIMARY KEY,
    version_int INT NOT NULL,
    versioned_folder_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    ehr_id UUID NOT NULL REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_folder_data (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    version_data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_folder_data ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_person (
    id TEXT PRIMARY KEY,
    version_int INT NOT NULL,
    versioned_party_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_person_data (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    version_data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_person_data ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_person_tag (
    person_id TEXT NOT NULL REFERENCES openehr.tbl_person(id) ON DELETE CASCADE,
    key TEXT NOT NULL,
    data JSONB NOT NULL,
    PRIMARY KEY (person_id, key)
);

CREATE TABLE openehr.tbl_agent (
    id TEXT PRIMARY KEY,
    version_int INT NOT NULL,
    versioned_party_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_agent_data (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    version_data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_agent_data ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_agent_tag (
    agent_id TEXT NOT NULL REFERENCES openehr.tbl_agent(id) ON DELETE CASCADE,
    key TEXT NOT NULL,
    data JSONB NOT NULL,
    PRIMARY KEY (agent_id, key)
);

CREATE TABLE openehr.tbl_group (
    id TEXT PRIMARY KEY,
    version_int INT NOT NULL,
    versioned_party_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_group_data (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    version_data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_group_data ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_group_tag (
    group_id TEXT NOT NULL REFERENCES openehr.tbl_group(id) ON DELETE CASCADE,
    key TEXT NOT NULL,
    data JSONB NOT NULL,
    PRIMARY KEY (group_id, key)
);

CREATE TABLE openehr.tbl_organisation (
    id TEXT PRIMARY KEY,
    version_int INT NOT NULL,
    versioned_party_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_organisation_data (
    id TEXT PRIMARY KEY,
    data JSONB NOT NULL,
    version_data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_organisation_data ALTER COLUMN data SET COMPRESSION lz4;

CREATE TABLE openehr.tbl_organisation_tag (
    organisation_id TEXT NOT NULL REFERENCES openehr.tbl_organisation(id) ON DELETE CASCADE,
    key TEXT NOT NULL,
    data JSONB NOT NULL,
    PRIMARY KEY (organisation_id, key)
);

-- ========= Query Table ==========

CREATE TABLE openehr.tbl_query (
    name TEXT NOT NULL,
    version INT NOT NULL,
    query TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (name, version)
);

-- Index for getting latest version of a named query
CREATE INDEX idx_query_name_version ON openehr.tbl_query 
    USING btree (name, version DESC);

-- Index for time-based query lookups
CREATE INDEX idx_query_created_at ON openehr.tbl_query 
    USING btree (created_at DESC);

-- ========== Audit Schemas ==========

CREATE SCHEMA audit;

CREATE TABLE audit.tbl_audit_log (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    actor_id UUID NOT NULL,
    actor_type TEXT NOT NULL,
    resource TEXT NOT NULL,
    action TEXT NOT NULL,
    success BOOLEAN NOT NULL,
    ip_address INET,
    user_agent TEXT,
    details JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Index for time-based audit log queries
ALTER TABLE audit.tbl_audit_log ALTER COLUMN details SET COMPRESSION lz4;

-- Index for time-based audit log queries
CREATE INDEX idx_audit_log_created_at ON audit.tbl_audit_log USING btree (created_at DESC);

-- ========== User Schema ==========

CREATE SCHEMA account;

CREATE TABLE account.tbl_account (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL, -- e.g., 'USER', 'SYSTEM'
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Insert initial SYSTEM account
INSERT INTO account.tbl_account (id, type, created_at, updated_at) VALUES
    ('00000000-0000-0000-0000-000000000001', 'SYSTEM', NOW(), NOW());

-- Unique index to ensure only one system account exists
CREATE UNIQUE INDEX one_system_account ON account.tbl_account (type) WHERE type = 'SYSTEM';


-- ========= Webhook Schema ==========
CREATE SCHEMA webhook;

CREATE TABLE webhook.tbl_event (
    id UUID PRIMARY KEY DEFAULT uuidv4(),
    type TEXT NOT NULL,  -- e.g., 'ehr_created', 'ehr_deleted'
    payload JSONB NOT NULL,    -- The payload sent to the webhook
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE webhook.tbl_subscription (
    id UUID PRIMARY KEY DEFAULT uuidv4(),
    url TEXT NOT NULL,
    secret TEXT NOT NULL,      
    event_types TEXT[] NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (url)
);

CREATE TABLE webhook.tbl_delivery (
    id UUID PRIMARY KEY DEFAULT uuidv4(),
    event_id UUID NOT NULL REFERENCES webhook.tbl_event(id) ON DELETE CASCADE,
    subscription_id UUID NOT NULL REFERENCES webhook.tbl_subscription(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    attempt_count INT NOT NULL,
    next_attempt_at TIMESTAMPTZ,
    last_attempt_at TIMESTAMPTZ,
    last_response_code INT,
    last_response_body TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (event_id, subscription_id)
);

CREATE INDEX idx_delivery_pending ON webhook.tbl_delivery (status, next_attempt_at);

-- ========= Additional Performance Tuning ==========

-- -- Increase statistics target for better query planning
-- ALTER TABLE openehr.tbl_object_version ALTER COLUMN versioned_object_id SET STATISTICS 1000;
-- ALTER TABLE openehr.tbl_object_version ALTER COLUMN type SET STATISTICS 1000;
-- ALTER TABLE openehr.tbl_object_version ALTER COLUMN ehr_id SET STATISTICS 1000;
-- ALTER TABLE openehr.tbl_object_version ALTER COLUMN created_at SET STATISTICS 1000;

-- ALTER TABLE openehr.tbl_versioned_object ALTER COLUMN type SET STATISTICS 1000;
-- ALTER TABLE openehr.tbl_versioned_object ALTER COLUMN ehr_id SET STATISTICS 1000;

-- ALTER TABLE openehr.tbl_contribution ALTER COLUMN ehr_id SET STATISTICS 1000;
-- ALTER TABLE openehr.tbl_contribution ALTER COLUMN created_at SET STATISTICS 1000;

-- -- Enable parallel query execution for large scans
-- ALTER TABLE openehr.tbl_object_version SET (parallel_workers = 4);
-- ALTER TABLE openehr.tbl_object_version_data SET (parallel_workers = 4);
-- ALTER TABLE openehr.tbl_versioned_object SET (parallel_workers = 4);
-- ALTER TABLE openehr.tbl_versioned_object_data SET (parallel_workers = 4);