package terminology

// Composition Category vocabulary constants
// This vocabulary codifies the values of the `category` attribute in Compositions.
// Used in: COMPOSITION.category
// External reference: openEHR terminology composition_category
const (
	COMPOSITION_CATEGORY_TERMINOLOGY_ID_OPENEHR string = "openehr"
)

var CompositionCategoryTerminologyIds = map[string]string{
	COMPOSITION_CATEGORY_TERMINOLOGY_ID_OPENEHR: "openEHR",
}

// IsValidCompositionCategoryTerminologyID checks if the provided terminology ID is valid for composition categories
func IsValidCompositionCategoryTerminologyID(terminologyID string) bool {
	_, exists := CompositionCategoryTerminologyIds[terminologyID]
	return exists
}

const (
	COMPOSITION_CATEGORY_CODE_PERSISTENT string = "431"
	COMPOSITION_CATEGORY_CODE_EPISODIC   string = "451"
	COMPOSITION_CATEGORY_CODE_EVENT      string = "433"
	COMPOSITION_CATEGORY_CODE_REPORT     string = "815"
)

// CompositionCategoryNames maps composition category codes to their human-readable names
var CompositionCategoryNames = map[string]string{
	COMPOSITION_CATEGORY_CODE_PERSISTENT: "persistent",
	COMPOSITION_CATEGORY_CODE_EPISODIC:   "episodic",
	COMPOSITION_CATEGORY_CODE_EVENT:      "event",
	COMPOSITION_CATEGORY_CODE_REPORT:     "report",
}

// IsValidCompositionCategoryCode checks if the given code is a valid composition category
func IsValidCompositionCategoryCode(code string) bool {
	_, exists := CompositionCategoryNames[code]
	return exists
}

// GetCompositionCategoryNameCode returns the human-readable name for the given composition category code
func GetCompositionCategoryNameCode(code string) string {
	if name, exists := CompositionCategoryNames[code]; exists {
		return name
	}
	return ""
}
