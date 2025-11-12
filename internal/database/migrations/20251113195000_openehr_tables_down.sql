-- Migration: 20251030204017_openehr_tables
-- Created: 2025-10-30T20:40:17+01:00
-- Description: openehr_tables (rollback)

-- Add your down migration SQL here

DROP TABLE tbl_patient_ehr_link;
DROP TABLE tbl_openehr_organisation;
DROP TABLE tbl_openehr_group;
DROP TABLE tbl_openehr_agent;
DROP TABLE tbl_openehr_person;
DROP TABLE tbl_openehr_role;
DROP TABLE tbl_openehr_folder_item;
DROP TABLE tbl_openehr_folder;
DROP TABLE tbl_openehr_composition;
DROP TABLE tbl_openehr_ehr_access;
DROP TABLE tbl_openehr_ehr_status;
DROP TABLE tbl_openehr_versioned_object;
DROP TABLE tbl_openehr_contribution_version;
DROP TABLE tbl_openehr_contribution;
DROP TABLE tbl_openehr_ehr;
DROP TABLE tbl_openehr_template;