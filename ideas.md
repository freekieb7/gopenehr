gopenEHR
===================

A custom / idealistic version of what an openEHR service should be.

### The goal ...
- Fast
- Simple
- Lightweight
- Open Source

### At the cost of ...
- Being out of compliance with the official openEHR specs
    - Removing, improving and adding features to create a better product
- Using the minimal amount of libraries we can get away with
    - Just write the damn code, not all problems have to be solved by relying on packages.

### What are the deviations?
- Deviation: Adding more 'List XXX' endpoints
- Reason: 'Get XXX by ID' works if you always know what's in the system, but what do you do if you don't know, or want to get multiple ?

## Roadmap
- Implement EHR endpoints
- Implement Query endpoints
    - Only allow execution, no storage 
- Full support for HTTP application/json
- JAQL (JSON path oriented AQL)

## Graveyard
- 

## JAQL
Tables:
- SELECT * FROM EHR
- SELECT * FROM EHR_STATUS
- SELECT * FROM EHR_ACCESS
- SELECT * FROM COMPOSITION
- SELECT * FROM FOLDER
- SELECT * FROM CONTRIBUTION
- SELECT * FROM EHR LEFT JOIN COMPOSITION

SELECT c.content.name::DV_TEXT FROM COMPOSITION c

SELECT c.**.name FROM COMPOSITION c

SELECT c.content[*] ? (@.name = "test") .value FROM COMPOSITION c

SELECT c FROM COMPOSITION c WHERE c.content[*] ANY SECTION

SELECT c FROM COM

