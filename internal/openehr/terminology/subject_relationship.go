package terminology

// Subject Relationship vocabulary codes
// This vocabulary codifies the relationship between the subject of care and some other party mentioned in the health record.
// Used in: PARTY_RELATED.relationship

const (
	SUBJECT_RELATIONSHIP_CODE_SELF                 string = "0"   // self
	SUBJECT_RELATIONSHIP_CODE_FOETUS               string = "3"   // foetus
	SUBJECT_RELATIONSHIP_CODE_MOTHER               string = "10"  // mother
	SUBJECT_RELATIONSHIP_CODE_FATHER               string = "9"   // father
	SUBJECT_RELATIONSHIP_CODE_DONOR                string = "6"   // donor
	SUBJECT_RELATIONSHIP_CODE_UNKNOWN              string = "253" // unknown
	SUBJECT_RELATIONSHIP_CODE_ADOPTED_DAUGHTER     string = "261" // adopted daughter
	SUBJECT_RELATIONSHIP_CODE_ADOPTED_SON          string = "260" // adopted son
	SUBJECT_RELATIONSHIP_CODE_ADOPTIVE_FATHER      string = "259" // adoptive father
	SUBJECT_RELATIONSHIP_CODE_ADOPTIVE_MOTHER      string = "258" // adoptive mother
	SUBJECT_RELATIONSHIP_CODE_BIOLOGICAL_FATHER    string = "256" // biological father
	SUBJECT_RELATIONSHIP_CODE_BIOLOGICAL_MOTHER    string = "255" // biological mother
	SUBJECT_RELATIONSHIP_CODE_BROTHER              string = "23"  // brother
	SUBJECT_RELATIONSHIP_CODE_CHILD                string = "28"  // child
	SUBJECT_RELATIONSHIP_CODE_COHABITEE            string = "265" // cohabitee
	SUBJECT_RELATIONSHIP_CODE_COUSIN               string = "257" // cousin
	SUBJECT_RELATIONSHIP_CODE_DAUGHTER             string = "29"  // daughter
	SUBJECT_RELATIONSHIP_CODE_GUARDIAN             string = "264" // guardian
	SUBJECT_RELATIONSHIP_CODE_MATERNAL_AUNT        string = "39"  // maternal aunt
	SUBJECT_RELATIONSHIP_CODE_MATERNAL_GRANDFATHER string = "8"   // maternal grandfather
	SUBJECT_RELATIONSHIP_CODE_MATERNAL_GRANDMOTHER string = "7"   // maternal grandmother
	SUBJECT_RELATIONSHIP_CODE_MATERNAL_UNCLE       string = "38"  // maternal uncle
	SUBJECT_RELATIONSHIP_CODE_NEONATE              string = "189" // neonate
	SUBJECT_RELATIONSHIP_CODE_PARENT               string = "254" // parent
	SUBJECT_RELATIONSHIP_CODE_PARTNER_SPOUSE       string = "22"  // partner/spouse
	SUBJECT_RELATIONSHIP_CODE_PATERNAL_AUNT        string = "41"  // paternal aunt
	SUBJECT_RELATIONSHIP_CODE_PATERNAL_GRANDFATHER string = "36"  // paternal grandfather
	SUBJECT_RELATIONSHIP_CODE_PATERNAL_GRANDMOTHER string = "37"  // paternal grandmother
	SUBJECT_RELATIONSHIP_CODE_PATERNAL_UNCLE       string = "40"  // paternal uncle
	SUBJECT_RELATIONSHIP_CODE_SIBLING              string = "27"  // sibling
	SUBJECT_RELATIONSHIP_CODE_SISTER               string = "24"  // sister
	SUBJECT_RELATIONSHIP_CODE_SON                  string = "31"  // son
	SUBJECT_RELATIONSHIP_CODE_STEP_FATHER          string = "263" // step father
	SUBJECT_RELATIONSHIP_CODE_STEP_MOTHER          string = "262" // step mother
	SUBJECT_RELATIONSHIP_CODE_STEP_OR_HALF_BROTHER string = "25"  // step or half brother
	SUBJECT_RELATIONSHIP_CODE_STEP_OR_HALF_SISTER  string = "26"  // step or half sister
)

// SubjectRelationshipNames maps subject relationship codes to their display names
var SubjectRelationshipNames = map[string]string{
	SUBJECT_RELATIONSHIP_CODE_SELF:                 "self",
	SUBJECT_RELATIONSHIP_CODE_FOETUS:               "foetus",
	SUBJECT_RELATIONSHIP_CODE_MOTHER:               "mother",
	SUBJECT_RELATIONSHIP_CODE_FATHER:               "father",
	SUBJECT_RELATIONSHIP_CODE_DONOR:                "donor",
	SUBJECT_RELATIONSHIP_CODE_UNKNOWN:              "unknown",
	SUBJECT_RELATIONSHIP_CODE_ADOPTED_DAUGHTER:     "adopted daughter",
	SUBJECT_RELATIONSHIP_CODE_ADOPTED_SON:          "adopted son",
	SUBJECT_RELATIONSHIP_CODE_ADOPTIVE_FATHER:      "adoptive father",
	SUBJECT_RELATIONSHIP_CODE_ADOPTIVE_MOTHER:      "adoptive mother",
	SUBJECT_RELATIONSHIP_CODE_BIOLOGICAL_FATHER:    "biological father",
	SUBJECT_RELATIONSHIP_CODE_BIOLOGICAL_MOTHER:    "biological mother",
	SUBJECT_RELATIONSHIP_CODE_BROTHER:              "brother",
	SUBJECT_RELATIONSHIP_CODE_CHILD:                "child",
	SUBJECT_RELATIONSHIP_CODE_COHABITEE:            "cohabitee",
	SUBJECT_RELATIONSHIP_CODE_COUSIN:               "cousin",
	SUBJECT_RELATIONSHIP_CODE_DAUGHTER:             "daughter",
	SUBJECT_RELATIONSHIP_CODE_GUARDIAN:             "guardian",
	SUBJECT_RELATIONSHIP_CODE_MATERNAL_AUNT:        "maternal aunt",
	SUBJECT_RELATIONSHIP_CODE_MATERNAL_GRANDFATHER: "maternal grandfather",
	SUBJECT_RELATIONSHIP_CODE_MATERNAL_GRANDMOTHER: "maternal grandmother",
	SUBJECT_RELATIONSHIP_CODE_MATERNAL_UNCLE:       "maternal uncle",
	SUBJECT_RELATIONSHIP_CODE_NEONATE:              "neonate",
	SUBJECT_RELATIONSHIP_CODE_PARENT:               "parent",
	SUBJECT_RELATIONSHIP_CODE_PARTNER_SPOUSE:       "partner/spouse",
	SUBJECT_RELATIONSHIP_CODE_PATERNAL_AUNT:        "paternal aunt",
	SUBJECT_RELATIONSHIP_CODE_PATERNAL_GRANDFATHER: "paternal grandfather",
	SUBJECT_RELATIONSHIP_CODE_PATERNAL_GRANDMOTHER: "paternal grandmother",
	SUBJECT_RELATIONSHIP_CODE_PATERNAL_UNCLE:       "paternal uncle",
	SUBJECT_RELATIONSHIP_CODE_SIBLING:              "sibling",
	SUBJECT_RELATIONSHIP_CODE_SISTER:               "sister",
	SUBJECT_RELATIONSHIP_CODE_SON:                  "son",
	SUBJECT_RELATIONSHIP_CODE_STEP_FATHER:          "step father",
	SUBJECT_RELATIONSHIP_CODE_STEP_MOTHER:          "step mother",
	SUBJECT_RELATIONSHIP_CODE_STEP_OR_HALF_BROTHER: "step or half brother",
	SUBJECT_RELATIONSHIP_CODE_STEP_OR_HALF_SISTER:  "step or half sister",
}

// IsValidSubjectRelationshipCode checks if the provided code is a valid subject relationship
func IsValidSubjectRelationshipCode(code string) bool {
	_, exists := SubjectRelationshipNames[code]
	return exists
}

// GetSubjectRelationshipName returns the display name for a subject relationship code
func GetSubjectRelationshipName(code string) string {
	if name, exists := SubjectRelationshipNames[code]; exists {
		return name
	}
	return ""
}
