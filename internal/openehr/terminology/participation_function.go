package terminology

// Participation Function vocabulary codes
// This vocabulary codifies functions of participation of parties in an interaction.
// Used in: PARTICIPATION.function

const (
	PARTICIPATION_FUNCTION_CODE_UNKNOWN string = "253" // unknown
)

// ParticipationFunctionNames maps participation function codes to their display names
var ParticipationFunctionNames = map[string]string{
	PARTICIPATION_FUNCTION_CODE_UNKNOWN: "unknown",
}

// IsValidParticipationFunctionCode checks if the provided code is a valid participation function
func IsValidParticipationFunctionCode(code string) bool {
	_, exists := ParticipationFunctionNames[code]
	return exists
}

// GetParticipationFunctionName returns the display name for a participation function code
func GetParticipationFunctionName(code string) string {
	if name, exists := ParticipationFunctionNames[code]; exists {
		return name
	}
	return ""
}
