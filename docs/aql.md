> AQL is still under development. It contains some experimental features and may change in the future.

# Introduction
AQL is a SQL-like language tailored for openEHR data.

This document describes the AQL dialect supported by this system and how to write effective queries against EHR content and demographics data.
After reading this document, you should be able to write queries that:
- Select data from EHR content and demographics
- Filter and sort data
- Perform aggregations
- Use functions
- Use parameters
- Use wildcards and set membership
- Use LIMIT and OFFSET

# Basics
## Selecting data
Selecting data is done using the SELECT clause and a FROM tree. The SELECT clause specifies the columns to be returned and the FROM tree specifies the data sources to be queried.

## Query syntax
- ```SELECT [DISTINCT] <column list>```

## Examples
- ```SELECT * FROM EHR```
- ```SELECT e FROM EHR e```
- ```SELECT e/ehr_id FROM EHR e```
- ```SELECT c/uid/value, c/name/value FROM COMPOSITION c```
- ```SELECT MAX(c/time_created) FROM COMPOSITION c```
- ```SELECT DISTINCT c/archetype_node_id FROM COMPOSITION```
- ```SELECT COUNT(*) FROM EHR```

# FROM tree
The FROM tree specifies the data sources to be queried. The FROM tree consists of a list of sources, each of which can be a composition, a version, a person, an organisation, a group, or an EHR.

## Query syntax
- ```FROM <source> [AS <alias>]? [CONTAINS <source> [AS <alias>]? [AND | OR]? ...```

## Examples
- ```SELECT * FROM EHR e```
- ```SELECT * FROM EHR e CONTAINS COMPOSITION c```
- ```SELECT * FROM EHR e CONTAINS COMPOSITION c CONTAINS OBSERVATION o```
- ```SELECT * FROM EHR e CONTAINS (COMPOSITION c CONTAINS OBSERVATION o OR EVALUATION v)```
- ```SELECT * FROM EHR e CONTAINS (COMPOSITION c CONTAINS OBSERVATION o AND EVALUATION v)```
- ```SELECT * FROM EHR e CONTAINS (COMPOSITION c CONTAINS OBSERVATION o AND (EVALUATION v CONTAINS CLUSTER cl)```
- ```SELECT * FROM EHR e CONTAINS (EHR_STATUS es AND CONTAINS COMPOSITION c CONTAINS OBSERVATION o)```

# WHERE clause
The WHERE clause specifies the conditions to be applied to the data sources.

## Query syntax
- ```WHERE <predicate>```

## Examples
- ```SELECT * FROM EHR e WHERE e/ehr_id/value = '...'```
- ```SELECT * FROM EHR e WHERE e/ehr_status/state = 'active'```
- ```SELECT * FROM EHR e WHERE e/ehr_id/value = '...' AND e/ehr_status/state = 'active' OR e/ehr_status/state = 'suspended'```
- ```SELECT * FROM EHR e WHERE e/ehr_id/value = '...' AND e/ehr_status/state = 'active' OR e/ehr_status/state = 'suspended' AND e/ehr_status/last_change_time > '2021-12-21T15:19:31.649+01:00'```

# ORDER BY clause
The ORDER BY clause specifies the order in which the results should be returned.

## Query syntax
- ```ORDER BY <column> [ASC | DESC]```

## Examples
- ```SELECT * FROM EHR e ORDER BY e/ehr_id/value ASC```
- ```SELECT * FROM EHR e ORDER BY e/ehr_id/value DESC```
- ```SELECT * FROM EHR e ORDER BY e/ehr_id/value ASC, e/ehr_status/state DESC```

# LIMIT clause
The LIMIT clause specifies the maximum number of results to be returned.

## Query syntax
- ```LIMIT <number> OFFSET <number>```

## Examples
- ```SELECT * FROM EHR e LIMIT 10```
- ```SELECT * FROM EHR e LIMIT 10 OFFSET 20```

# Parameters
It is possible to have parameterized values in your query, making the query dynamic and reusable.

## Query Syntax
- ```SELECT <parameter>```
- ```WHERE o/name/value = <parameter>```
- ```WHERE o/data[<parameter>]/events[<parameter> and name/value = <parameter>] = <parameter>```

# More examples
### Scenario: Get all blood glucose values and their corresponding subject ids, where blood glucose > 11 mmol/L or blood glucose >= 200 mg/dL
```
SELECT
    e/ehr_status/subject/external_ref/id/value as subjectId,
    a/items[at0001]/value as analyteName,
    a/items[at0001]/value as analyteResult
FROM EHR e
    CONTAINS COMPOSITION c
        CONTAINS OBSERVATION o[openEHR-EHR-OBSERVATION.laboratory_test_result.v1]
            CONTAINS CLUSTER a[openEHR-EHR-CLUSTER.laboratory_test_analyte.v1]
WHERE
    (a/items[at0001]/value/defining_code/code_string matches {'14743-9','2345-7'} AND a/items[at0001]/value/defining_code/terminology_id = 'LOINC')
    AND
    ((a/items[at0024]/value/magnitude > 11 AND a/items[at0024]/value/units matches {'mmol/L'})
        OR (a/items[at0024]/value/magnitude >= 200 AND a/items[at0024]/value/units matches {'mg/dL'}))
```

# References
For more information on what AQL has to offer, read [this page](https://specifications.openehr.org/releases/QUERY/latest/AQL.html) explaining AQL in more details.