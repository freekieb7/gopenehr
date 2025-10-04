openEHR REST API - Getting Started Guide
==================
Practical primer for developers integrating with the openEHR REST APIs (EHR, Definition, Query) plus SMART on openEHR authentication.

# Overview
The openEHR REST APIs expose a clinical data platform built on a versioned Reference Model (RM). Core domains:
- <b>EHR API</b>: create and manage EHR containers, compositions, directory (folder), contributions.
- <b>Definition API</b>: upload and retrieve operational templates (ADL 1.4 / ADL2) and store/retrieve AQL stored queries.
- <b>Query API</b>: execute AQL (Archetype Query Language) ad-hoc or stored queries and receive RESULT_SET payloads.
- <b>SMART on openEHR</b>: OAuth2/OIDC based authentication, launch context, scopes mapping to openEHR resources.

Spec references (Release 1.0.3 unless stated):
- Overview: https://specifications.openehr.org/releases/ITS-REST/Release-1.0.3/overview.html
- EHR API: https://specifications.openehr.org/releases/ITS-REST/Release-1.0.3/ehr.html
- Query API: https://specifications.openehr.org/releases/ITS-REST/Release-1.0.3/query.html
- Definition API: https://specifications.openehr.org/releases/ITS-REST/Release-1.0.3/definition.html
- SMART on openEHR (auth): https://specifications.openehr.org/releases/ITS-REST/development/smart_app_launch.html

# Getting started
Before being able to access any of our API's, the user should first be authenticated. For that you need to have your application/client registered in our Authentication Service (**Smart-Auth**).

 [!tip] 
 Since our Authentication Service is a closed system, please contact us so we can register your application.

After you have obtained your Client ID, you can start using the Authentication Service to perform an [OAuth PKCE flow](https://specifications.openehr.org/releases/ITS-REST/development/smart_app_launch.html#_smart_authorization_flow). 

Choose a Patient or Practitioner as the type of user you want to authenticate as. When all the credentials are valid, and access is granted by the user, your are ready to query the services.

For a more in depth guide, follow the authentication guide.
<scalar-page-link filepath="docs/authentication.md">
</scalar-page-link>


<scalar-callout type="info"> Use the tokens inside the header 'Authentication: Bearer \<token\>' before attempting to query the service </scalar-callout>

# Typical Workflow
## 1. Create or Locate an EHR
Auto-generated EHR id:
```
POST /ehr
Prefer: return=representation
Authorization: Bearer <token>
```
Optional JSON EHR_STATUS body; else defaults applied.
Find by subject:
```
GET /ehr?subject_id={value}&subject_namespace={ns}
```

## 2. Upload Templates (Definition API)
ADL 1.4 XML:
```
POST /definition/template/adl1.4
Content-Type: application/xml
Prefer: return=minimal

<template xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.openehrorg/v1">
    <language>
        <terminology_id>
            <value>ISO_639-1</value>
        </terminology_id>
        <code_string>en</code_string>
    </language>
	....
</template>
```
List templates:
```
GET /definition/template/adl14
```
Retrieve template (web template JSON):
```
GET /definition/template/adl1.4/{template_id}
Accept: application/openehr.wt+json
```

## 3. Commit a COMPOSITION
```
POST /ehr/{ehr_id}/composition
Content-Type: application/json
Prefer: return=representation
Authorization: Bearer <token>

{
	"archetype_node_id": "openEHR-EHR-COMPOSITION.encounter.v1",
	"name": {
		"value": "Vital Signs"
	},
	"uid": {
		"_type": "OBJECT_VERSION_ID",
		"value": "8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1"
	},
	...
}
```
Body includes COMPOSITION with language, territory, category, context, composer, content[]. For simplified formats supply `openEHR-TEMPLATE_ID` if template metadata absent.
Response: 201 with Location and possibly full representation.

## 4. Retrieve or Update a COMPOSITION
Latest version (implicit):
```
GET /ehr/{ehr_id}/composition/{versioned_object_uid}
```
Specific version:
```
GET /ehr/{ehr_id}/composition/{version_uid}
```
Update (optimistic locking):
```
PUT /ehr/{ehr_id}/composition/{versioned_object_uid}
If-Match: "<latest_version_uid>"
```
Delete (logical):
```
DELETE /ehr/{ehr_id}/composition/{version_uid}
```

## 5. Directory (FOLDER)
Create:
```
POST /ehr/{ehr_id}/directory
```
Update with concurrency:
```
PUT /ehr/{ehr_id}/directory
If-Match: "<latest_version_uid>"
```
Fetch path inside directory at time:
```
GET /ehr/{ehr_id}/directory?path=episodes/a/b&version_at_time=YYYY-MM-DDThh:mm:ss.sss±hh:mm
```

## 6. Contribution (Batch commit)
```
POST /ehr/{ehr_id}/contribution
Content-Type: application/json
```
Supply versions array referencing operations and audit metadata.

## 7. Query Data (AQL)
Ad-hoc POST (recommended for length + params):
```
POST /query/aql
{
	"q": "SELECT ... WHERE o/data[...] > $temperature",
	"query_parameters": { "temperature": 38.5 }
}
```

# Quick Reference (Cheat Sheet)
| Action | Endpoint | Method | Notes |
| --- | --- | --- | --- |
| Create EHR | /ehr | POST | Optional EHR_STATUS body |
| Get EHR | /ehr/{ehr_id} | GET | 404 if missing |
| Create COMPOSITION | /ehr/{ehr_id}/composition | POST | 201 + Location |
| Get COMPOSITION (latest) | /ehr/{ehr_id}/composition/{versioned_object_uid} | GET | Uses container id |
| Get COMPOSITION (by version) | /ehr/{ehr_id}/composition/{version_uid} | GET | Specific version |
| Update COMPOSITION | /ehr/{ehr_id}/composition/{versioned_object_uid} | PUT | If-Match required |
| Delete COMPOSITION | /ehr/{ehr_id}/composition/{version_uid} | DELETE | Logical delete |
| Upload template ADL2 | /definition/template/adl2 | POST | text/plain body |
| List templates ADL2 | /definition/template/adl2 | GET | Returns metadata list |
| Get template (web template) | /definition/template/adl1.4/{id} | GET | Accept web template mime |
| Store query | /definition/query/{name} | PUT | text/plain AQL |
| Execute ad-hoc AQL | /query/aql | POST | JSON body |
| Execute stored query | /query/{qualified_query_name} | GET | Query params for variables |

# General info
## Base URL and Versioning
Implementations usually mount APIs under a versioned prefix, e.g.:
```
https://platform.example.com/openehr/rest/v1
```
All examples assume base = `https://openEHRSys.example.com/v1` (spec examples) or similar.

## Authentication (SMART on openEHR)
1. Discover configuration:
```
GET {platform-base}/.well-known/smart-configuration
```
2. Use Authorization Code + PKCE (public) or Authorization Code / JWT Bearer / Client Credentials (confidential) as advertised in `grant_types_supported`.
3. Request scopes combining:
	 - Launch/context: `launch`, `launch/patient` (patient+EHR context), optionally experimental `launch/episode`.
	 - Resource scopes (see Scopes section): e.g. `patient/composition-MyTemplate.v1.crud`.
	 - Identity (OIDC) scopes if required: `openid profile`.
4. After user auth + consent, exchange `code` at `token_endpoint` for `access_token` (+ optional `refresh_token`, `id_token`).
5. Extract openEHR context claims (if present) from token response: `ehrId`, optional `episodeId`.
6. Call openEHR REST endpoints with `Authorization: Bearer <token>`.

Embedded (iframe) launch: App receives `iss` and `launch` URL params, fetches smart configuration at `{iss}/.well-known/smart-configuration`, then initiates auth including the received `launch` value.

## Scopes (Resource Access Model)
Pattern: `<compartment>/<resource>.<permissions>`
- Compartments: `patient` (scoped to current EHR), `user` (user-wide), `system` (backend).
- Resources (examples):
	- `composition-<templateId>` (COMPOSITION instances of a template)
	- `template-<templateId>` (operational templates)
	- `aql-<queryName>` (stored queries or `*` for ad-hoc when allowed)
- Permissions letters:
	- `c` create, `r` read, `u` update, `d` delete, `s` search/execute
- Combine letters: `crud`, `cruds` etc.
Examples:
```
patient/composition-Vital_Signs.v1.crud
user/aql-org.openehr::compositions.rs
system/template-*.crud
```
Use wildcards judiciously; broad access should be minimized.

## Content Negotiation & Formats
Headers:
- Request body format: `Content-Type: application/json` or `application/xml` (must match server support). For simplified formats use media types: e.g. `application/openehr.wt.flat+json`.
- Response preference: `Prefer: return=representation` for full payload, otherwise minimal (default `return=minimal`).
- Accept header to request representation: `Accept: application/json` or `application/xml`.
Simplified (SDT / web template) JSON requires `openEHR-TEMPLATE_ID` header if COMPOSITION JSON does not carry template details.

## Core HTTP Headers (selected)
- `Authorization: Bearer <token>` (SMART)
- `ETag` (entity tag of versioned resource)
- `If-Match` (optimistic concurrency when updating/deleting versions)
- `Location` (URI of newly created or updated resource)
- `openEHR-VERSION.*` / `openEHR-AUDIT_DETAILS.*` (optional commit metadata)
- `openEHR-TEMPLATE_ID` (template id when simplified body omits it)
- `Prefer` (representation negotiation; can combine `resolve_refs`)

## Fundamental Resource Identification
- `ehr_id` (UUID, EHR.ehr_id.value)
- `versioned_object_uid` (container UID)
- `version_uid` (OBJECT_VERSION_ID: object::system::version)
- `uid_based_id` (either version_uid or versioned_object_uid depending context)

## Optimistic Concurrency & Versioning
- Always capture `ETag` (mirrors version_uid) from read responses of versioned resources.
- Supply that in `If-Match` for updates/deletes to prevent overwriting new versions.
- On 412 Precondition Failed, re-fetch to get latest `ETag`.

<!-- ## Audit & Commit Metadata
Optional custom headers can enrich audit trail when committing:
```
openEHR-AUDIT_DETAILS.change_type: code_string="249"  # creation
openEHR-AUDIT_DETAILS.description: value="Initial vitals capture"
openEHR-VERSION.lifecycle_state: code_string="532"    # complete
```
Server merges provided values with defaults. -->

<!-- ## DateTime Handling
- Use extended ISO 8601 (e.g. `2015-01-20T19:30:22.765+01:00`)
- Query params like `version_at_time` must be extended ISO 8601. -->

<!-- ## Error Handling Strategy
Check status codes; on 4xx parse optional JSON error body (when `Prefer: return=representation` supplied). Typical structure:
```
{
	"message": "Error message",
	"code": 90000,
	"errors": [ { "_type": "DV_CODED_TEXT", "value": "Error message", ... } ]
}
``` -->

## Security & Best Practices
- Use least privilege scopes; prefer template-specific over wildcard.
- Enforce PKCE for public clients.
- Rotate client secrets / keys (confidential clients) and support asymmetric JWT auth where possible.
- Validate `issuer` and `aud` claims in tokens (if `id_token` present) and match against configuration.
- Cache `.well-known/smart-configuration` with reasonable TTL; refresh on auth failures.

## Minimal Client Flow Outline
1. Discover SMART: fetch `/.well-known/smart-configuration`.
2. Build auth request (authorize endpoint) with scopes, PKCE, redirect_uri, state, code_challenge, launch (if iframe), aud (issuer/platform) and optional `launch/patient`.
3. Exchange code at token endpoint.
4. Use `access_token` with openEHR endpoints. Add `Prefer: return=representation` when needing full payload.
5. Handle 401 by refreshing token or re-authorizing.
6. Use AQL queries for analytics/reporting; restrict fetch + offset for pagination.

# Additional Notes
- Always store server-generated identifiers; never attempt to rewrite historical version_uids.
- COMPOSITION updates always create new VERSION internally; treat version_uid as immutable lineage.
- For performance, narrow AQL queries with path constraints and parameters instead of wide wildcard queries.

# Glossary (selected)
- AQL: Archetype Query Language used to query structured clinical data.
- COMPOSITION: Top-level clinical document instance (versioned).
- CONTRIBUTION: Logical commit grouping of one or more VERSIONs.
- EHR: Root container for a patient’s clinical record.
- VERSIONED_OBJECT: Container tracking successive VERSION states of a resource.
- RESULT_SET: Structured result of a query execution.

# Disclaimer
This getting started guide is a distilled practical aid; always refer to the normative specifications linked above for authoritative definitions, edge cases and full semantics.

