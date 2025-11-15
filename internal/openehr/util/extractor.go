package util

// Helper struct for extracting _type field
type TypeExtractor struct {
	Type_ string `json:"_type,omitzero"`
}
