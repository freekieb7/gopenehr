package terminology

// Term Mapping Purpose vocabulary codes
// This vocabulary codifies purposes for term mappings as used in openEHR coded text data.
// Used in: TERM_MAPPING.purpose

const (
	TERM_MAPPING_PURPOSE_CODE_PUBLIC_HEALTH  string = "669" // public health
	TERM_MAPPING_PURPOSE_CODE_REIMBURSEMENT  string = "670" // reimbursement
	TERM_MAPPING_PURPOSE_CODE_RESEARCH_STUDY string = "671" // research study
)

// TermMappingPurposeNames maps term mapping purpose codes to their display names
var TermMappingPurposeNames = map[string]string{
	TERM_MAPPING_PURPOSE_CODE_PUBLIC_HEALTH:  "public health",
	TERM_MAPPING_PURPOSE_CODE_REIMBURSEMENT:  "reimbursement",
	TERM_MAPPING_PURPOSE_CODE_RESEARCH_STUDY: "research study",
}

// IsValidTermMappingPurposeCode checks if the provided code is a valid term mapping purpose
func IsValidTermMappingPurposeCode(code string) bool {
	_, exists := TermMappingPurposeNames[code]
	return exists
}

// GetTermMappingPurposeName returns the display name for a term mapping purpose code
func GetTermMappingPurposeName(code string) string {
	if name, exists := TermMappingPurposeNames[code]; exists {
		return name
	}
	return ""
}
