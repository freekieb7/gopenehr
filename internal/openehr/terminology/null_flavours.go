package terminology

// Null Flavours vocabulary codes
// This vocabulary codifies 'flavours of null' for missing data items.
// Used in: ELEMENT.null_flavour

const (
	NULL_FLAVOUR_CODE_NO_INFORMATION string = "271" // no information
	NULL_FLAVOUR_CODE_UNKNOWN        string = "253" // unknown
	NULL_FLAVOUR_CODE_MASKED         string = "272" // masked
	NULL_FLAVOUR_CODE_NOT_APPLICABLE string = "273" // not applicable
)

// NullFlavourNames maps null flavour codes to their display names
var NullFlavourNames = map[string]string{
	NULL_FLAVOUR_CODE_NO_INFORMATION: "no information",
	NULL_FLAVOUR_CODE_UNKNOWN:        "unknown",
	NULL_FLAVOUR_CODE_MASKED:         "masked",
	NULL_FLAVOUR_CODE_NOT_APPLICABLE: "not applicable",
}

// IsValidNullFlavourCode checks if the provided code is a valid null flavour
func IsValidNullFlavourCode(code string) bool {
	_, exists := NullFlavourNames[code]
	return exists
}

// GetNullFlavourName returns the display name for a null flavour code
func GetNullFlavourName(code string) string {
	if name, exists := NullFlavourNames[code]; exists {
		return name
	}
	return ""
}
