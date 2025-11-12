package terminology

// Version Lifecycle State vocabulary constants
// This vocabulary codifies lifecycle states of Compositions or other elements of the health record.
// Used in: ORIGINAL_VERSION.lifecycle_state, IMPORTED_VERSION.lifecycle_state, VERSIONED_OBJECT.lifecycle_state
// External reference: openEHR terminology version_lifecycle_state
const (
	VERSION_LIFECYCLE_STATE_CODE_COMPLETE   string = "532"
	VERSION_LIFECYCLE_STATE_CODE_INCOMPLETE string = "553"
	VERSION_LIFECYCLE_STATE_CODE_DELETED    string = "523"
	VERSION_LIFECYCLE_STATE_CODE_INACTIVE   string = "800"
	VERSION_LIFECYCLE_STATE_CODE_ABANDONED  string = "801"
)

// VersionLifecycleStateNames maps version lifecycle state codes to their human-readable names
var VersionLifecycleStateNames = map[string]string{
	VERSION_LIFECYCLE_STATE_CODE_COMPLETE:   "complete",
	VERSION_LIFECYCLE_STATE_CODE_INCOMPLETE: "incomplete",
	VERSION_LIFECYCLE_STATE_CODE_DELETED:    "deleted",
	VERSION_LIFECYCLE_STATE_CODE_INACTIVE:   "inactive",
	VERSION_LIFECYCLE_STATE_CODE_ABANDONED:  "abandoned",
}

// IsValidVersionLifecycleStateCode checks if the given code is a valid version lifecycle state
func IsValidVersionLifecycleStateCode(code string) bool {
	_, exists := VersionLifecycleStateNames[code]
	return exists
}

// GetVersionLifecycleStateName returns the human-readable name for the given version lifecycle state code
func GetVersionLifecycleStateName(code string) string {
	if name, exists := VersionLifecycleStateNames[code]; exists {
		return name
	}
	return ""
}
