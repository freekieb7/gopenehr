roadmap:
- custom encode / decode (xml,json,flat)
 x json
 - xml
 - flat
- value validation
- endpoints

AQL ideas
===================
TABLES 
- EHR
- EHR_STATUS
- EHR_ACCESS
- COMPOSITION
- FOLDER
- CONTRIBUTION
- VERSIONED_OBJECT ?

SELECT * FROM EHR
SELECT * FROM COMPOSITION
SELECT * FROM FOLDER
SELECT * FROM EHR_STATUS

SELECT * FROM EHR LEFT JOIN COMPOSITION
SELECT c.content.name::DV_TEXT FROM COMPOSITION c
SELECT c.**.name FROM COMPOSITION c
SELECT c.content[*] ? (@.name = "test") .value FROM COMPOSITION c
SELECT c FROM COMPOSITION c WHERE c.content[*] ANY SECTION
SELECT c FROM COM

