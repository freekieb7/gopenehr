package openehr

import (
	"encoding/json"
	"os"
	"testing"
)

func TestValidateValidEHR(t *testing.T) {
	ehrJSON, err := os.ReadFile("../../tests/fixture/ehr.json")
	if err != nil {
		t.Fatalf("Failed to read EHR fixture: %v", err)
	}

	var ehr EHR
	err = json.Unmarshal(ehrJSON, &ehr)
	if err != nil {
		t.Fatalf("Failed to unmarshal EHR JSON: %v", err)
	}

	errors := ehr.Validate("$")
	if len(errors) != 0 {
		t.Errorf("Expected no validation errors, got %d", len(errors))
		for _, err := range errors {
			t.Logf("Validation error: %s", err.Message)
		}
	}
}

func TestValidateInvalidEHR(t *testing.T) {
	ehr := EHR{
		EHRID: HIER_OBJECT_ID{Value: "invalid-uuid."},
		EHRStatus: OBJECT_REF{
			Namespace: "local",
			Type:      "EHR_STATUS",
			ID:        X_OBJECT_ID{Value: &HIER_OBJECT_ID{Value: "also-invalid-uuid"}},
		},
		EHRAccess: OBJECT_REF{
			Namespace: "local",
			Type:      "EHR_ACCESS",
			ID:        X_OBJECT_ID{Value: &HIER_OBJECT_ID{Value: "yet-another-invalid-uuid"}},
		},
		TimeCreated: DV_DATE_TIME{Value: "2023-01-01T00:00:00Z"},
	}

	errors := ehr.Validate("$")
	if len(errors) != 5 {
		t.Errorf("Expected 5 validation errors, got %d", len(errors))
		for _, err := range errors {
			t.Logf("Validation error: %s", err.Message)
		}
	}
}
