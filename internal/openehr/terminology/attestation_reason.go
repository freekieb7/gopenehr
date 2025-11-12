package terminology

// OpenEHR Support Terminology - Attestation Reason
// Based on: https://specifications.openehr.org/releases/TERM/latest/SupportTerminology.html#_attestation_reason
// External reference: HL7 ParticipationSignature domain
// Used in: ATTESTATION.reason

// Attestation reason constants
const (
	ATTESTATION_CODE_REASON_SIGNED    string = "240"
	ATTESTATION_CODE_REASON_WITNESSED string = "648"
)

// AttestationReasonNames provides human-readable names for attestation reason codes
var AttestationReasonNames = map[string]string{
	ATTESTATION_CODE_REASON_SIGNED:    "signed",
	ATTESTATION_CODE_REASON_WITNESSED: "witnessed",
}

// IsValidAttestationReasonCode checks if the given attestation reason code is valid
func IsValidAttestationReasonCode(code string) bool {
	_, exists := AttestationReasonNames[code]
	return exists
}

// GetAttestationReasonName returns the human-readable name for an attestation reason code
func GetAttestationReasonName(code string) string {
	if name, exists := AttestationReasonNames[code]; exists {
		return name
	}
	return "Unknown attestation reason"
}
