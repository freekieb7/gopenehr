package terminology

// Character Set constants from IANA Character Sets
// Based on OpenEHR Support Terminology specification
// Used in: ENTRY.encoding, DV_ENCAPSULATED.charset
const (
	CHARSET_TERMINOLOGY_ID_IANA string = "IANA_character-sets"
)

var CharsetTerminologyIDs = map[string]string{
	CHARSET_TERMINOLOGY_ID_IANA: "IANA Character Sets",
}

func IsValidCharsetTerminologyID(id string) bool {
	_, exists := CharsetTerminologyIDs[id]
	return exists
}

const (
	CHARSET_ISO_10646_UTF_1 string = "ISO-10646-UTF-1"
	CHARSET_ISO_8859_1_1987 string = "ISO_8859-1:1987"
	CHARSET_ISO_8859_2      string = "ISO-8859-2"
	CHARSET_ISO_8859_3_1988 string = "ISO_8859-3:1988"
	CHARSET_ISO_8859_15     string = "ISO-8859-15"
	CHARSET_US_ASCII        string = "US-ASCII"
	CHARSET_UTF_7           string = "UTF-7"
	CHARSET_UTF_8           string = "UTF-8"
	CHARSET_UTF_16          string = "UTF-16"
	CHARSET_UTF_16BE        string = "UTF-16BE"
	CHARSET_UTF_16LE        string = "UTF-16LE"
	CHARSET_UTF_32          string = "UTF-32"
	CHARSET_UTF_32BE        string = "UTF-32BE"
	CHARSET_UTF_32LE        string = "UTF-32LE"
)

// CharsetNames provides a mapping of character set codes to their names
var CharsetNames = map[string]string{
	CHARSET_ISO_10646_UTF_1: "ISO-10646-UTF-1",
	CHARSET_ISO_8859_1_1987: "ISO_8859-1:1987",
	CHARSET_ISO_8859_2:      "ISO-8859-2",
	CHARSET_ISO_8859_3_1988: "ISO_8859-3:1988",
	CHARSET_ISO_8859_15:     "ISO-8859-15",
	CHARSET_US_ASCII:        "US-ASCII",
	CHARSET_UTF_7:           "UTF-7",
	CHARSET_UTF_8:           "UTF-8",
	CHARSET_UTF_16:          "UTF-16",
	CHARSET_UTF_16BE:        "UTF-16BE",
	CHARSET_UTF_16LE:        "UTF-16LE",
	CHARSET_UTF_32:          "UTF-32",
	CHARSET_UTF_32BE:        "UTF-32BE",
	CHARSET_UTF_32LE:        "UTF-32LE",
}

// IsValidCharset checks if the given string is a valid IANA character set code
func IsValidCharset(charset string) bool {
	_, exists := CharsetNames[charset]
	return exists
}

// GetCharsetName returns the name of the character set for the given code
func GetCharsetName(charset string) (string, bool) {
	name, exists := CharsetNames[charset]
	return name, exists
}
