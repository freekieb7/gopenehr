package util

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	UUIDRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	// ISO/IEC 8824 compliant version
	ISOOIDRegex = regexp.MustCompile(`^(0|1)(\.(0|[1-9][0-9]*)){0,1}(\.(0|[1-9][0-9]*))*$|^2(\.(0|[1-9][0-9]*))(\.(0|[1-9][0-9]*))*$`)
	// RFC 1034 compliant regex
	InternetIDRegex = regexp.MustCompile(`^([A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?)(?:\.([A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?))*$`)
	// RFC 2396 compliant regex for URI namespace part
	NamespaceRegex = regexp.MustCompile(`^([a-zA-Z])([a-zA-Z0-9_.:\/&?=+-])*$`)
	// RFC 3986 compliant regex for URI path part
	URIRegex = regexp.MustCompile(`^\/?([a-zA-Z0-9\-._~!$&'()*+,;=:@%]+\/?)*$`)
	// Lexical form: trunk_version [ '.' branch_number '.' branch_version ]
	VersionTreeIDRegex = regexp.MustCompile(`^([0-9]+)(\.([0-9]+)\.([0-9]+))?$`)
	// Lexical form: rm_originator '-' rm_name '-' rm_entity '.' concept_name { '-' specialisation }* '.v' number.
	ArchetypeIDRegex = regexp.MustCompile(`^([a-zA-Z0-9_]+)-([a-zA-Z0-9_]+)-([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)(-[a-zA-Z0-9_]+)*\.v([0-9]+)$`)
)

type ValidationError struct {
	Model          string `json:"model"`
	Path           string `json:"path"`
	Message        string `json:"message"`
	Recommendation string `json:"recommendation"`
}

func ValidateUID(uid string) error {
	if uid == "" {
		return errors.New("UID cannot be empty")
	}

	// Check for UUID format (8-4-4-4-12 hex digits)
	if ValidateUUID(uid) {
		return nil
	}

	// Check for ISO OID format (numbers separated by dots)
	if ValidateISOOID(uid) {
		return nil
	}

	// Check for Internet ID format (reverse domain name style)
	if ValidateInternetID(uid) {
		return nil
	}

	return fmt.Errorf("UID '%s' does not match any valid format (UUID, ISO OID, or Internet ID)", uid)
}

func ValidateUUID(uuid string) bool {
	return UUIDRegex.MatchString(uuid)
}

func ValidateISOOID(oid string) bool {
	return ISOOIDRegex.MatchString(oid)
}

func ValidateInternetID(internetID string) bool {
	// Check total length (RFC 1034: max 255 characters, but commonly 253)
	if len(internetID) == 0 || len(internetID) > 255 {
		return false
	}

	return InternetIDRegex.MatchString(internetID)
}
