package terminology

// OpenEHR Support Terminology - Normal Statuses
// Based on: https://specifications.openehr.org/releases/TERM/latest/SupportTerminology.html#_normal_statuses
// External reference: HL7 Observation-Interpretation vocabulary
// Maps to some codes in HL7v2 User-defined table 0078 - Abnormal flags
// Used in: DV_ORDERED.normal_status

// Normal status constants
const (
	NORMAL_STATUS_HHH string = "HHH" // Critically high
	NORMAL_STATUS_HH  string = "HH"  // Very high
	NORMAL_STATUS_H   string = "H"   // High
	NORMAL_STATUS_N   string = "N"   // Normal
	NORMAL_STATUS_L   string = "L"   // Low
	NORMAL_STATUS_LL  string = "LL"  // Very low
	NORMAL_STATUS_LLL string = "LLL" // Critically low
)

// NormalStatusNames provides human-readable names for normal status codes
var NormalStatusNames = map[string]string{
	NORMAL_STATUS_HHH: "Critically high",
	NORMAL_STATUS_HH:  "Very high",
	NORMAL_STATUS_H:   "High",
	NORMAL_STATUS_N:   "Normal",
	NORMAL_STATUS_L:   "Low",
	NORMAL_STATUS_LL:  "Very low",
	NORMAL_STATUS_LLL: "Critically low",
}

// IsValidNormalStatus checks if the given normal status code is valid
func IsValidNormalStatus(status string) bool {
	_, exists := NormalStatusNames[status]
	return exists
}

// GetNormalStatusName returns the human-readable name for a normal status code
func GetNormalStatusName(status string) string {
	if name, exists := NormalStatusNames[status]; exists {
		return name
	}
	return "Unknown normal status"
}
