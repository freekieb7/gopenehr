package terminology

// OpenEHR Support Terminology - Compression Algorithms
// Based on: https://specifications.openehr.org/releases/TERM/latest/SupportTerminology.html#_compression_algorithms
// External reference: HL7 Compression Algorithms domain
// Used in: DV_MULTIMEDIA.compression_algorithm

// Compression algorithm constants
const (
	COMPRESSION_COMPRESS string = "compress"
	COMPRESSION_DEFLATE  string = "deflate"
	COMPRESSION_GZIP     string = "gzip"
	COMPRESSION_ZLIB     string = "zlib"
	COMPRESSION_OTHER    string = "other"
)

// CompressionAlgorithmNames provides human-readable names for compression algorithm codes
var CompressionAlgorithmNames = map[string]string{
	COMPRESSION_COMPRESS: "Compress",
	COMPRESSION_DEFLATE:  "Deflate",
	COMPRESSION_GZIP:     "GZIP",
	COMPRESSION_ZLIB:     "ZLIB",
	COMPRESSION_OTHER:    "Other",
}

// IsValidCompressionAlgorithm checks if the given compression algorithm code is valid
func IsValidCompressionAlgorithm(algorithm string) bool {
	_, exists := CompressionAlgorithmNames[algorithm]
	return exists
}

// GetCompressionAlgorithmName returns the human-readable name for a compression algorithm code
func GetCompressionAlgorithmName(algorithm string) string {
	if name, exists := CompressionAlgorithmNames[algorithm]; exists {
		return name
	}
	return "Unknown compression algorithm"
}
