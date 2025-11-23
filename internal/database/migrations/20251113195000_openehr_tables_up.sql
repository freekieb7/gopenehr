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
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_ehr_data (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);
    
-- ========= Contribution Table ==========

CREATE TABLE openehr.tbl_contribution (
    id UUID PRIMARY KEY,
    ehr_id UUID REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_contribution_data (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

-- ========= Versioned Object Table ==========

CREATE TABLE openehr.tbl_versioned_object (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    ehr_id UUID REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_versioned_object_data (
    id UUID PRIMARY KEY REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

-- ========= Object Version Table ==========

CREATE TABLE openehr.tbl_object_version (
    id TEXT PRIMARY KEY,
    versioned_object_id UUID NOT NULL REFERENCES openehr.tbl_versioned_object(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    ehr_id UUID REFERENCES openehr.tbl_ehr(id) ON DELETE CASCADE,
    contribution_id UUID NOT NULL REFERENCES openehr.tbl_contribution(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE openehr.tbl_object_version_data (
    id TEXT PRIMARY KEY REFERENCES openehr.tbl_object_version(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

-- -- ========= Tag Table ==========
-- CREATE TABLE openehr.tbl_object_version_tag (
--     object_version_id TEXT NOT NULL,
--     tag TEXT NOT NULL,
--     major INT NOT NULL,
--     minor INT NOT NULL,
--     patch INT NOT NULL
-- );

-- CREATE UNIQUE INDEX idx_tbl_object_version_tag ON openehr.tbl_object_version_tag USING btree (object_version_id, tag);

-- ALTER TABLE openehr.tbl_object_version_tag ADD CONSTRAINT pk_tbl_object_version_tag PRIMARY KEY USING INDEX idx_tbl_object_version_tag;
-- ALTER TABLE openehr.tbl_object_version_tag ADD CONSTRAINT fk_tbl_object_version_tag_object_version FOREIGN KEY (object_version_id) REFERENCES openehr.tbl_object_version(id) ON DELETE CASCADE;