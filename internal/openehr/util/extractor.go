package util

import (
	"bytes"
	"unsafe"
)

// // Helper struct for extracting _type field
// type TypeExtractor struct {
// 	Type_ string `json:"_type,omitzero"`
// }

var key []byte = []byte(`"_type":`)
var keyLen = len(key)

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
// _type must be the first field.
// Returns a string using unsafe to avoid allocation for common types.
func UnsafeTypeFieldExtraction(data []byte) string {
	// Skip initial '{'
	for i := 1; i <= len(data)-keyLen; i++ {
		b := data[i]

		if isSpace(b) {
			continue
		}

		if !bytes.Equal(data[i:i+keyLen], key) {
			return ""
		}

		i += keyLen

		// skip whitespace
		for i < len(data) && isSpace(data[i]) {
			i++
		}

		// must be quoted string
		if i < len(data) && data[i] == '"' {
			i++
			start := i
			for i < len(data) && data[i] != '"' {
				i++
			}
			if i < len(data) {
				typeBytes := data[start:i]

				// Fast path: compare against common types to return string without allocation
				switch {
				case bytes.Equal(typeBytes, typeComposition):
					return "COMPOSITION"
				case bytes.Equal(typeBytes, typeSection):
					return "SECTION"
				case bytes.Equal(typeBytes, typeObservation):
					return "OBSERVATION"
				case bytes.Equal(typeBytes, typeEvaluation):
					return "EVALUATION"
				case bytes.Equal(typeBytes, typeInstruction):
					return "INSTRUCTION"
				case bytes.Equal(typeBytes, typeAction):
					return "ACTION"
				case bytes.Equal(typeBytes, typeActivity):
					return "ACTIVITY"
				case bytes.Equal(typeBytes, typeAdminEntry):
					return "ADMIN_ENTRY"
				case bytes.Equal(typeBytes, typeGenericEntry):
					return "GENERIC_ENTRY"
				case bytes.Equal(typeBytes, typeDvText):
					return "DV_TEXT"
				case bytes.Equal(typeBytes, typeDvCodedText):
					return "DV_CODED_TEXT"
				case bytes.Equal(typeBytes, typeDvQuantity):
					return "DV_QUANTITY"
				case bytes.Equal(typeBytes, typeDvCount):
					return "DV_COUNT"
				case bytes.Equal(typeBytes, typeDvDateTime):
					return "DV_DATE_TIME"
				case bytes.Equal(typeBytes, typeDvDate):
					return "DV_DATE"
				case bytes.Equal(typeBytes, typeDvTime):
					return "DV_TIME"
				case bytes.Equal(typeBytes, typeItemTree):
					return "ITEM_TREE"
				case bytes.Equal(typeBytes, typeItemList):
					return "ITEM_LIST"
				case bytes.Equal(typeBytes, typeItemTable):
					return "ITEM_TABLE"
				case bytes.Equal(typeBytes, typeItemSingle):
					return "ITEM_SINGLE"
				case bytes.Equal(typeBytes, typePartyIdentified):
					return "PARTY_IDENTIFIED"
				case bytes.Equal(typeBytes, typePartyRelated):
					return "PARTY_RELATED"
				case bytes.Equal(typeBytes, typePartySelf):
					return "PARTY_SELF"
				default:
					// Fallback to unsafe string conversion for uncommon types
					return unsafe.String(&typeBytes[0], len(typeBytes))
				}
			}
		}
		return ""
	}
	return ""
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// func DecoderTypeField(data []byte) string {
// 	dec := json.NewDecoder(bytes.NewReader(data))

// 	tok, err := dec.Token() // must be {
// 	if err != nil || tok != json.Delim('{') {
// 		return ""
// 	}

// 	for dec.More() {
// 		t, err := dec.Token()
// 		if err != nil {
// 			return ""
// 		}

// 		key := t.(string)
// 		if key == "_type" {
// 			val, err := dec.Token()
// 			if err != nil {
// 				return ""
// 			}
// 			return val.(string)
// 		}

// 		// // skip value
// 		// if err := dec.Skip(); err != nil {
// 		// 	return "", err
// 		// }
// 	}
// 	return ""
// }
