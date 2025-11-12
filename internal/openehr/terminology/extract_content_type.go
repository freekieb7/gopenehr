package terminology

// Extract Content Type vocabulary codes
// This vocabulary codifies the type of the content required for an Extract specification.
// Used in: EXTRACT_SPEC.extract_type

const (
	EXTRACT_CONTENT_TYPE_CODE_OPENEHR_EHR             string = "803" // openEHR EHR
	EXTRACT_CONTENT_TYPE_CODE_OPENEHR_DEMOGRAPHIC     string = "804" // openEHR Demographic
	EXTRACT_CONTENT_TYPE_CODE_OPENEHR_SYNCHRONISATION string = "805" // openEHR synchronisation
	EXTRACT_CONTENT_TYPE_CODE_OPENEHR_GENERIC         string = "806" // openEHR generic
	EXTRACT_CONTENT_TYPE_CODE_GENERIC_EMR             string = "807" // generic EMR
	EXTRACT_CONTENT_TYPE_CODE_OTHER                   string = "808" // other
)

// ExtractContentTypeNames maps extract content type codes to their display names
var ExtractContentTypeNames = map[string]string{
	EXTRACT_CONTENT_TYPE_CODE_OPENEHR_EHR:             "openEHR EHR",
	EXTRACT_CONTENT_TYPE_CODE_OPENEHR_DEMOGRAPHIC:     "openEHR Demographic",
	EXTRACT_CONTENT_TYPE_CODE_OPENEHR_SYNCHRONISATION: "openEHR synchronisation",
	EXTRACT_CONTENT_TYPE_CODE_OPENEHR_GENERIC:         "openEHR generic",
	EXTRACT_CONTENT_TYPE_CODE_GENERIC_EMR:             "generic EMR",
	EXTRACT_CONTENT_TYPE_CODE_OTHER:                   "other",
}

// IsValidExtractContentTypeCode checks if the provided code is a valid extract content type
func IsValidExtractContentTypeCode(code string) bool {
	_, exists := ExtractContentTypeNames[code]
	return exists
}

// GetExtractContentTypeName returns the display name for an extract content type code
func GetExtractContentTypeName(code string) string {
	if name, exists := ExtractContentTypeNames[code]; exists {
		return name
	}
	return ""
}
