package terminology

// Extract Action Type vocabulary codes
// This vocabulary codifies the action types of a Request for an Extract specification.
// Used in: EXTRACT_ACTION_REQUEST.action

const (
	EXTRACT_ACTION_TYPE_CODE_CANCEL   string = "809" // cancel
	EXTRACT_ACTION_TYPE_CODE_RESEND   string = "810" // resend
	EXTRACT_ACTION_TYPE_CODE_SEND_NEW string = "811" // send new
)

// ExtractActionTypeNames maps extract action type codes to their display names
var ExtractActionTypeNames = map[string]string{
	EXTRACT_ACTION_TYPE_CODE_CANCEL:   "cancel",
	EXTRACT_ACTION_TYPE_CODE_RESEND:   "resend",
	EXTRACT_ACTION_TYPE_CODE_SEND_NEW: "send new",
}

// IsValidExtractActionTypeCode checks if the provided code is a valid extract action type
func IsValidExtractActionTypeCode(code string) bool {
	_, exists := ExtractActionTypeNames[code]
	return exists
}

// GetExtractActionTypeName returns the display name for an extract action type code
func GetExtractActionTypeName(code string) string {
	if name, exists := ExtractActionTypeNames[code]; exists {
		return name
	}
	return ""
}
