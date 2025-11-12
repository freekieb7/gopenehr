package terminology

// Audit Change Type vocabulary constants
// This vocabulary codifies the kinds of changes to data which are recorded in audit trails.
// Used in: AUDIT_DETAILS.change_type
// External reference: openEHR terminology audit_change_type
const (
	AUDIT_CHANGE_TYPE_CODE_CREATION          string = "249"
	AUDIT_CHANGE_TYPE_CODE_AMENDMENT         string = "250"
	AUDIT_CHANGE_TYPE_CODE_MODIFICATION      string = "251"
	AUDIT_CHANGE_TYPE_CODE_SYNTHESIS         string = "252"
	AUDIT_CHANGE_TYPE_CODE_DELETED           string = "523"
	AUDIT_CHANGE_TYPE_CODE_ATTESTATION       string = "666"
	AUDIT_CHANGE_TYPE_CODE_RESTORATION       string = "816"
	AUDIT_CHANGE_TYPE_CODE_FORMAT_CONVERSION string = "817"
	AUDIT_CHANGE_TYPE_CODE_UNKNOWN           string = "253"
)

// AuditChangeTypeNames maps audit change type codes to their human-readable names
var AuditChangeTypeNames = map[string]string{
	AUDIT_CHANGE_TYPE_CODE_CREATION:          "creation",
	AUDIT_CHANGE_TYPE_CODE_AMENDMENT:         "amendment",
	AUDIT_CHANGE_TYPE_CODE_MODIFICATION:      "modification",
	AUDIT_CHANGE_TYPE_CODE_SYNTHESIS:         "synthesis",
	AUDIT_CHANGE_TYPE_CODE_DELETED:           "deleted",
	AUDIT_CHANGE_TYPE_CODE_ATTESTATION:       "attestation",
	AUDIT_CHANGE_TYPE_CODE_RESTORATION:       "restoration",
	AUDIT_CHANGE_TYPE_CODE_FORMAT_CONVERSION: "format conversion",
	AUDIT_CHANGE_TYPE_CODE_UNKNOWN:           "unknown",
}

// IsValidAuditChangeTypeCode checks if the given code is a valid audit change type
func IsValidAuditChangeTypeCode(code string) bool {
	_, exists := AuditChangeTypeNames[code]
	return exists
}

// GetAuditChangeTypeName returns the human-readable name for the given audit change type code
func GetAuditChangeTypeName(code string) string {
	if name, exists := AuditChangeTypeNames[code]; exists {
		return name
	}
	return ""
}
