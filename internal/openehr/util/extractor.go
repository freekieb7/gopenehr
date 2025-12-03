package util

import (
	"bytes"
	"unsafe"
)

// // Helper struct for extracting _type field
// type TypeExtractor struct {
// 	Type_ string `json:"_type,omitzero"`
// }

// Common type strings as byte slices for comparison
var (
	typeComposition     = []byte("COMPOSITION")
	typeSection         = []byte("SECTION")
	typeObservation     = []byte("OBSERVATION")
	typeEvaluation      = []byte("EVALUATION")
	typeInstruction     = []byte("INSTRUCTION")
	typeAction          = []byte("ACTION")
	typeActivity        = []byte("ACTIVITY")
	typeAdminEntry      = []byte("ADMIN_ENTRY")
	typeGenericEntry    = []byte("GENERIC_ENTRY")
	typeDvText          = []byte("DV_TEXT")
	typeDvCodedText     = []byte("DV_CODED_TEXT")
	typeDvQuantity      = []byte("DV_QUANTITY")
	typeDvCount         = []byte("DV_COUNT")
	typeDvDateTime      = []byte("DV_DATE_TIME")
	typeDvDate          = []byte("DV_DATE")
	typeDvTime          = []byte("DV_TIME")
	typeItemTree        = []byte("ITEM_TREE")
	typeItemList        = []byte("ITEM_LIST")
	typeItemTable       = []byte("ITEM_TABLE")
	typeItemSingle      = []byte("ITEM_SINGLE")
	typePartyIdentified = []byte("PARTY_IDENTIFIED")
	typePartyRelated    = []byte("PARTY_RELATED")
	typePartySelf       = []byte("PARTY_SELF")
)

// DecoderTypeField extracts the _type field from the given JSON data
// without fully unmarshaling the data.
// Returns a string using unsafe to avoid allocation for common types.
func UnsafeTypeFieldExtraction(data []byte) string {
	if len(data) == 0 || data[0] != '{' {
		return ""
	}

	i := 1     // cursor
	depth := 1 // object depth
	inString := false
	escape := false

	for i < len(data) {
		b := data[i]

		// Handle string parsing (only for VALUES, never for root keys)
		if inString {
			if escape {
				escape = false
			} else {
				switch b {
				case '\\':
					escape = true
				case '"':
					inString = false
				}
			}
			i++
			continue
		}

		switch b {
		case '"':
			// Only enter string mode if not root key
			if depth != 1 {
				inString = true
				i++
				continue
			}
			// else: root-level key → parse below

		case '{':
			depth++
			i++
			continue

		case '}':
			depth--
			i++
			if depth == 0 {
				return ""
			}
			continue

		case '[':
			// Arrays should not affect object depth
			i++
			continue

		case ']':
			i++
			continue

		default:
			if isSpace(b) {
				i++
				continue
			}
		}

		// Only parse keys when we're at root object depth
		if depth == 1 && b == '"' {
			// Parse key
			keyStart := i + 1
			j := keyStart

			for j < len(data) {
				c := data[j]
				if c == '\\' {
					j += 2
					continue
				}
				if c == '"' {
					break
				}
				j++
			}
			if j >= len(data) {
				return ""
			}

			key := data[keyStart:j]
			i = j + 1 // skip closing quote

			// Skip whitespace
			for i < len(data) && isSpace(data[i]) {
				i++
			}

			if i >= len(data) || data[i] != ':' {
				return ""
			}
			i++

			if bytes.Equal(key, []byte("_type")) {
				// Skip whitespace
				for i < len(data) && isSpace(data[i]) {
					i++
				}

				if i >= len(data) || data[i] != '"' {
					return ""
				}

				i++
				valStart := i

				for i < len(data) {
					c := data[i]
					if c == '\\' {
						i += 2
						continue
					}
					if c == '"' {
						break
					}
					i++
				}
				if i >= len(data) {
					return ""
				}

				value := data[valStart:i]

				// Fast-path matching
				switch {
				case bytes.Equal(value, typeComposition):
					return "COMPOSITION"
				case bytes.Equal(value, typeSection):
					return "SECTION"
				case bytes.Equal(value, typeObservation):
					return "OBSERVATION"
				case bytes.Equal(value, typeEvaluation):
					return "EVALUATION"
				case bytes.Equal(value, typeInstruction):
					return "INSTRUCTION"
				case bytes.Equal(value, typeAction):
					return "ACTION"
				case bytes.Equal(value, typeActivity):
					return "ACTIVITY"
				case bytes.Equal(value, typeAdminEntry):
					return "ADMIN_ENTRY"
				case bytes.Equal(value, typeGenericEntry):
					return "GENERIC_ENTRY"
				case bytes.Equal(value, typeDvText):
					return "DV_TEXT"
				case bytes.Equal(value, typeDvCodedText):
					return "DV_CODED_TEXT"
				case bytes.Equal(value, typeDvQuantity):
					return "DV_QUANTITY"
				case bytes.Equal(value, typeDvCount):
					return "DV_COUNT"
				case bytes.Equal(value, typeDvDateTime):
					return "DV_DATE_TIME"
				case bytes.Equal(value, typeDvDate):
					return "DV_DATE"
				case bytes.Equal(value, typeDvTime):
					return "DV_TIME"
				case bytes.Equal(value, typeItemTree):
					return "ITEM_TREE"
				case bytes.Equal(value, typeItemList):
					return "ITEM_LIST"
				case bytes.Equal(value, typeItemTable):
					return "ITEM_TABLE"
				case bytes.Equal(value, typeItemSingle):
					return "ITEM_SINGLE"
				case bytes.Equal(value, typePartyIdentified):
					return "PARTY_IDENTIFIED"
				case bytes.Equal(value, typePartyRelated):
					return "PARTY_RELATED"
				case bytes.Equal(value, typePartySelf):
					return "PARTY_SELF"
				default:
					if len(value) == 0 {
						return ""
					}
					return unsafe.String(unsafe.SliceData(value), len(value))
				}
			}

			// Skip non-target value
			for i < len(data) && isSpace(data[i]) {
				i++
			}
			if i >= len(data) {
				return ""
			}

			switch data[i] {
			case '"':
				inString = true
				i++
			case '{':
				// Skip object
				objDepth := 1
				i++
				for i < len(data) && objDepth > 0 {
					switch data[i] {
					case '"':
						i++
						for i < len(data) {
							if data[i] == '\\' {
								i += 2
								continue
							}
							if data[i] == '"' {
								i++
								break
							}
							i++
						}
					case '{':
						objDepth++
						i++
					case '}':
						objDepth--
						i++
					default:
						i++
					}
				}
			case '[':
				// Skip array
				arrDepth := 1
				i++
				for i < len(data) && arrDepth > 0 {
					switch data[i] {
					case '"':
						i++
						for i < len(data) {
							if data[i] == '\\' {
								i += 2
								continue
							}
							if data[i] == '"' {
								i++
								break
							}
							i++
						}
					case '[':
						arrDepth++
						i++
					case ']':
						arrDepth--
						i++
					default:
						i++
					}
				}
			default:
				// true, false, number, null
				for i < len(data) && data[i] != ',' && data[i] != '}' {
					i++
				}
			}
		}

		i++
	}

	return ""
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '\n'
}
