package terminology

// Extract Update Trigger Event Type vocabulary codes
// This vocabulary codifies the event names of an update Extract specification.
// Used in: EXTRACT_UPDATE_SPEC.trigger_events

const (
	EXTRACT_UPDATE_TRIGGER_EVENT_TYPE_CODE_ANY_CHANGE string = "812" // any change
	EXTRACT_UPDATE_TRIGGER_EVENT_TYPE_CODE_CORRECTION string = "813" // correction
	EXTRACT_UPDATE_TRIGGER_EVENT_TYPE_CODE_UPDATE     string = "814" // update
)

// ExtractUpdateTriggerEventTypeNames maps extract update trigger event type codes to their display names
var ExtractUpdateTriggerEventTypeNames = map[string]string{
	EXTRACT_UPDATE_TRIGGER_EVENT_TYPE_CODE_ANY_CHANGE: "any change",
	EXTRACT_UPDATE_TRIGGER_EVENT_TYPE_CODE_CORRECTION: "correction",
	EXTRACT_UPDATE_TRIGGER_EVENT_TYPE_CODE_UPDATE:     "update",
}

// IsValidExtractUpdateTriggerEventTypeCode checks if the provided code is a valid extract update trigger event type
func IsValidExtractUpdateTriggerEventTypeCode(code string) bool {
	_, exists := ExtractUpdateTriggerEventTypeNames[code]
	return exists
}

// GetExtractUpdateTriggerEventTypeName returns the display name for an extract update trigger event type code
func GetExtractUpdateTriggerEventTypeName(code string) string {
	if name, exists := ExtractUpdateTriggerEventTypeNames[code]; exists {
		return name
	}
	return ""
}
