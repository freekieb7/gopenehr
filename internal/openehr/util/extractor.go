package util

// Helper struct for extracting _type field
type TypeExtractor struct {
	MetaType string `json:"_type,omitzero"`
}
