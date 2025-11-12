package terminology

// OpenEHR Support Terminology - Integrity Check Algorithms
// Based on: https://specifications.openehr.org/releases/TERM/latest/SupportTerminology.html#_integrity_check_algorithms
// External reference: HL7 IntegrityCheckAlgorithm domain
// Used in: DV_MULTIMEDIA.integrity_check_algorithm

// Integrity check algorithm constants
const (
	INTEGRITY_SHA1       string = "SHA-1"
	INTEGRITY_SHA224     string = "SHA-224"
	INTEGRITY_SHA256     string = "SHA-256"
	INTEGRITY_SHA384     string = "SHA-384"
	INTEGRITY_SHA512     string = "SHA-512"
	INTEGRITY_SHA512_224 string = "SHA-512/224"
	INTEGRITY_SHA512_256 string = "SHA-512/256"
)

// IntegrityCheckAlgorithmNames provides human-readable names for integrity check algorithm codes
var IntegrityCheckAlgorithmNames = map[string]string{
	INTEGRITY_SHA1:       "SHA-1",
	INTEGRITY_SHA224:     "SHA-224",
	INTEGRITY_SHA256:     "SHA-256",
	INTEGRITY_SHA384:     "SHA-384",
	INTEGRITY_SHA512:     "SHA-512",
	INTEGRITY_SHA512_224: "SHA-512/224",
	INTEGRITY_SHA512_256: "SHA-512/256",
}

// IsValidIntegrityCheckAlgorithm checks if the given integrity check algorithm code is valid
func IsValidIntegrityCheckAlgorithm(algorithm string) bool {
	_, exists := IntegrityCheckAlgorithmNames[algorithm]
	return exists
}

// GetIntegrityCheckAlgorithmName returns the human-readable name for an integrity check algorithm code
func GetIntegrityCheckAlgorithmName(algorithm string) string {
	if name, exists := IntegrityCheckAlgorithmNames[algorithm]; exists {
		return name
	}
	return "Unknown integrity check algorithm"
}
