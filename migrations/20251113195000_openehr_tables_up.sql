CREATE SCHEMA openehr;

-- ========== Template Table (TODO is not done) ==========

-- CREATE TABLE openehr.tbl_template (
--     id TEXT NOT NULL,
--     concept TEXT NOT NULL,
--     archetype_id TEXT NOT NULL,
--     data JSONB NOT NULL,
--     raw BYTEA NOT NULL,
--     created_at TIMESTAMP NOT NULL DEFAULT NOW()
-- );

-- CREATE UNIQUE INDEX idx_template_id ON openehr.tbl_template USING btree (id);
-- CREATE INDEX idx_template_archetype_id ON openehr.tbl_template(archetype_id);

-- ALTER TABLE openehr.tbl_template ADD CONSTRAINT pk_tbl_template PRIMARY KEY USING INDEX idx_template_id;

-- ========= EHR Table ==========

CREATE TABLE openehr.tbl_ehr (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for time-based EHR queries
CREATE INDEX idx_ehr_created_at ON openehr.tbl_ehr USING btree (created_at DESC);

-- ========= Contribution Table ==========

CREATE TABLE openehr.tbl_contribution (
    id UUID PRIMARY KEY,
    ehr_id UUID REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for filtering contributions by EHR
CREATE INDEX idx_contribution_ehr_id ON openehr.tbl_contribution USING btree (ehr_id) WHERE ehr_id IS NOT NULL;

-- Index for time-based contribution queries
CREATE INDEX idx_contribution_created_at ON openehr.tbl_contribution USING btree (created_at DESC);

-- Composite index for EHR + time queries (most common pattern)
CREATE INDEX idx_contribution_ehr_created ON openehr.tbl_contribution 
    USING btree (ehr_id, created_at DESC) WHERE ehr_id IS NOT NULL;

-- Index for demographic contributions (no EHR association)
CREATE INDEX idx_contribution_null_ehr ON openehr.tbl_contribution 
    USING btree (created_at DESC) WHERE ehr_id IS NULL;

CREATE TABLE openehr.tbl_contribution_data (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_contribution_data ALTER COLUMN data SET COMPRESSION lz4;

-- GIN index for JSONB queries on contribution data
CREATE INDEX idx_contribution_data_gin ON openehr.tbl_contribution_data USING gin (data jsonb_path_ops);

-- ========= Versioned Object Table ==========

CREATE TABLE openehr.tbl_versioned_object (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    ehr_id UUID REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for filtering by type (very common in AQL)
CREATE INDEX idx_versioned_object_type ON openehr.tbl_versioned_object 
    USING btree (type);

-- Index for EHR-related versioned objects
CREATE INDEX idx_versioned_object_ehr_id ON openehr.tbl_versioned_object 
    USING btree (ehr_id) WHERE ehr_id IS NOT NULL;

-- Composite index for EHR + type queries (common AQL pattern)
CREATE INDEX idx_versioned_object_ehr_type ON openehr.tbl_versioned_object 
    USING btree (ehr_id, type) WHERE ehr_id IS NOT NULL;

-- Index for time-based queries
CREATE INDEX idx_versioned_object_created_at ON openehr.tbl_versioned_object 
    USING btree (created_at DESC);

-- Index for demographic versioned objects (VERSIONED_PARTY)
CREATE INDEX idx_versioned_object_demographic ON openehr.tbl_versioned_object 
    USING btree (type, created_at DESC) WHERE ehr_id IS NULL;

CREATE TABLE openehr.tbl_versioned_object_data (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_versioned_object_data ALTER COLUMN data SET COMPRESSION lz4;

-- GIN index for general JSONB queries
CREATE INDEX idx_versioned_object_data_gin ON openehr.tbl_versioned_object_data 
    USING gin (data jsonb_path_ops);

-- Index for UID lookups (common in openEHR)
CREATE INDEX idx_versioned_object_data_uid ON openehr.tbl_versioned_object_data 
    USING btree ((data->'uid'->>'value'));

-- ========= Object Version Table ==========

CREATE TABLE openehr.tbl_object_version (
    id TEXT PRIMARY KEY,
    versioned_object_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    ehr_id UUID REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- CRITICAL: Index for getting latest version of each versioned object (DISTINCT ON optimization)
CREATE INDEX idx_object_version_latest ON openehr.tbl_object_version 
    USING btree (versioned_object_id, created_at DESC);

-- Index for version history queries (get all versions of a versioned object)
CREATE INDEX idx_object_version_versioned_object ON openehr.tbl_object_version 
    USING btree (versioned_object_id, created_at DESC);

-- Index for filtering by type (common in AQL)
CREATE INDEX idx_object_version_type ON openehr.tbl_object_version 
    USING btree (type);

-- Partial indexes for specific types (faster for common queries)
CREATE INDEX idx_object_version_composition ON openehr.tbl_object_version 
    USING btree (versioned_object_id, created_at DESC) WHERE type = 'COMPOSITION';

CREATE INDEX idx_object_version_ehr_status ON openehr.tbl_object_version 
    USING btree (versioned_object_id, created_at DESC) WHERE type = 'EHR_STATUS';

CREATE INDEX idx_object_version_person ON openehr.tbl_object_version 
    USING btree (versioned_object_id, created_at DESC) WHERE type = 'PERSON';

-- Index for EHR-related object versions
CREATE INDEX idx_object_version_ehr_id ON openehr.tbl_object_version 
    USING btree (ehr_id) WHERE ehr_id IS NOT NULL;

-- Composite index for EHR + type queries
CREATE INDEX idx_object_version_ehr_type ON openehr.tbl_object_version 
    USING btree (ehr_id, type, created_at DESC) WHERE ehr_id IS NOT NULL;

-- Index for contribution lookups
CREATE INDEX idx_object_version_contribution_id ON openehr.tbl_object_version 
    USING btree (contribution_id);

-- Index for time-based queries
CREATE INDEX idx_object_version_created_at ON openehr.tbl_object_version 
    USING btree (created_at DESC);

-- Index for demographic object versions (no EHR association)
CREATE INDEX idx_object_version_demographic ON openehr.tbl_object_version 
    USING btree (type, created_at DESC) WHERE ehr_id IS NULL;

CREATE TABLE openehr.tbl_object_version_data (
    id TEXT PRIMARY KEY REFERENCES openehr.tbl_object_version(id) ON DELETE CASCADE,
    version_data JSONB NOT NULL,
    object_data JSONB NOT NULL
);

ALTER TABLE openehr.tbl_object_version_data ALTER COLUMN version_data SET COMPRESSION lz4;
ALTER TABLE openehr.tbl_object_version_data ALTER COLUMN object_data SET COMPRESSION lz4;

-- GIN index for general JSONB queries on version_data
CREATE INDEX idx_object_version_version_data_gin ON openehr.tbl_object_version_data 
    USING gin (version_data jsonb_path_ops);

-- GIN index for queries on extracted object_data (most common)
CREATE INDEX idx_object_version_object_data_gin ON openehr.tbl_object_version_data 
    USING gin (object_data jsonb_path_ops);

-- Index for archetype_node_id queries (very common in AQL)
CREATE INDEX idx_object_version_data_archetype ON openehr.tbl_object_version_data 
    USING btree ((object_data->>'archetype_node_id'));

-- Index for name queries
CREATE INDEX idx_object_version_data_name ON openehr.tbl_object_version_data 
    USING btree ((object_data->'name'->>'value'));

-- Index for subject reference lookups (EHR_STATUS → PERSON join)
CREATE INDEX idx_object_version_data_subject ON openehr.tbl_object_version_data 
    USING btree (
        (object_data->'subject'->'external_ref'->'id'->>'value'),
        (object_data->'subject'->'external_ref'->>'namespace'),
        (object_data->'subject'->'external_ref'->>'type')
    );

-- Partial index for COMPOSITION category (encounter, event, persistent)
CREATE INDEX idx_object_version_data_composition_category ON openehr.tbl_object_version_data 
    USING btree ((object_data->'category'->>'value'))
    WHERE object_data->>'_type' IN ('COMPOSITION');

-- Index for time-based queries on compositions (context/start_time)
CREATE INDEX idx_object_version_data_composition_time ON openehr.tbl_object_version_data 
    USING btree ((object_data->'context'->'start_time'->>'value'))
    WHERE object_data->>'_type' = 'COMPOSITION';

-- Index for folder items (for FOLDER → COMPOSITION joins)
CREATE INDEX idx_object_version_data_folder_items ON openehr.tbl_object_version_data 
    USING gin ((object_data->'items'))
    WHERE object_data->>'_type' = 'FOLDER';

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

-- ========= EHR View ==========
CREATE VIEW openehr.vw_ehr AS
SELECT
    e.id,
    jsonb_build_object(
        'system_id', jsonb_build_object(
            'value', 'gopenehr'
        ),
        'ehr_id', jsonb_build_object(
            'value', e.id
        ),
        'ehr_status', jsonb_build_object(
            'namespace', 'local',
            'type', 'EHR_STATUS',
            'id', jsonb_build_object(
                '_type', 'OBJECT_VERSION_ID',
                'value', es.id
            )
        ),
        'ehr_access', jsonb_build_object(
            'namespace', 'local',
            'type', 'EHR_ACCESS',
            'id', jsonb_build_object(
                '_type', 'OBJECT_VERSION_ID',
                'value', ea.id
            )
        ),
        'contributions', coalesce(contrib.contributions, '[]'::jsonb),
        'compositions', coalesce(comp.compositions, '[]'::jsonb),
        'directory', jsonb_build_object(
            'namespace', 'local',
            'type', 'FOLDER',
            'id', jsonb_build_object(
                '_type', 'OBJECT_VERSION_ID',
                'value', fo.id
            )
        ),
        'time_created', jsonb_build_object(
            'value', to_char(e.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.USZ')
        )
    ) AS data,
    e.created_at
FROM openehr.tbl_ehr e

JOIN openehr.tbl_object_version es 
  ON es.ehr_id = e.id AND es.type = 'EHR_STATUS'
JOIN openehr.tbl_object_version ea 
  ON ea.ehr_id = e.id AND ea.type = 'EHR_ACCESS'
LEFT JOIN openehr.tbl_object_version fo 
  ON fo.ehr_id = e.id AND fo.type = 'FOLDER'
LEFT JOIN LATERAL (
    SELECT jsonb_agg(jsonb_build_object(
        'namespace', 'local',
        'type', 'CONTRIBUTION',
        'id', jsonb_build_object(
            '_type', 'HIER_OBJECT_ID',
            'value', id
        )
    )) AS contributions
    FROM openehr.tbl_contribution
    WHERE ehr_id = e.id
) contrib ON true
LEFT JOIN LATERAL (
    SELECT jsonb_agg(jsonb_build_object(
        'namespace', 'local',
        'type', 'VERSIONED_COMPOSITION',
        'id', jsonb_build_object(
            '_type', 'HIER_OBJECT_ID',
            'value', id
        )
    )) AS compositions
    FROM openehr.tbl_versioned_object
    WHERE ehr_id = e.id
      AND type = 'VERSIONED_COMPOSITION'
) comp ON true;

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

-- Increase statistics target for better query planning
ALTER TABLE openehr.tbl_object_version ALTER COLUMN versioned_object_id SET STATISTICS 1000;
ALTER TABLE openehr.tbl_object_version ALTER COLUMN type SET STATISTICS 1000;
ALTER TABLE openehr.tbl_object_version ALTER COLUMN ehr_id SET STATISTICS 1000;
ALTER TABLE openehr.tbl_object_version ALTER COLUMN created_at SET STATISTICS 1000;

ALTER TABLE openehr.tbl_versioned_object ALTER COLUMN type SET STATISTICS 1000;
ALTER TABLE openehr.tbl_versioned_object ALTER COLUMN ehr_id SET STATISTICS 1000;

ALTER TABLE openehr.tbl_contribution ALTER COLUMN ehr_id SET STATISTICS 1000;
ALTER TABLE openehr.tbl_contribution ALTER COLUMN created_at SET STATISTICS 1000;

-- Enable parallel query execution for large scans
ALTER TABLE openehr.tbl_object_version SET (parallel_workers = 4);
ALTER TABLE openehr.tbl_object_version_data SET (parallel_workers = 4);
ALTER TABLE openehr.tbl_versioned_object SET (parallel_workers = 4);
ALTER TABLE openehr.tbl_versioned_object_data SET (parallel_workers = 4);