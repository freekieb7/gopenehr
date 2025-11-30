package rm

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/pkg/utils"
	"github.com/google/uuid"
)

func TestSetModelType(t *testing.T) {
	// 3. Prepare EHR
	ehr := EHR{
		EHRID: HIER_OBJECT_ID{
			Value: uuid.NewString(),
		},
		EHRStatus: OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      VERSIONED_EHR_STATUS_TYPE,
			ID: ObjectIDFromHierObjectID(HIER_OBJECT_ID{
				Type_: utils.Some(HIER_OBJECT_ID_TYPE),
				Value: uuid.NewString(),
			}),
		},
		EHRAccess: OBJECT_REF{
			Namespace: config.NAMESPACE_LOCAL,
			Type:      VERSIONED_EHR_ACCESS_TYPE,
			ID: ObjectIDFromHierObjectID(HIER_OBJECT_ID{
				Type_: utils.Some(HIER_OBJECT_ID_TYPE),
				Value: uuid.NewString(),
			}),
		},
		TimeCreated: DV_DATE_TIME{
			Value: time.Now().UTC().Format(time.RFC3339),
		},
	}
	ehr.SetModelName()

	if !ehr.EHRID.Type_.E || ehr.EHRID.Type_.V != HIER_OBJECT_ID_TYPE {
		t.Errorf("EHRID _type not set correctly, got: %v", ehr.EHRID.Type_)
	}
	if !ehr.EHRStatus.Type_.E || ehr.EHRStatus.Type_.V != OBJECT_REF_TYPE {
		t.Errorf("EHRStatus _type not set correctly, got: %v", ehr.EHRStatus.Type_)
	}
	if !ehr.EHRAccess.Type_.E || ehr.EHRAccess.Type_.V != OBJECT_REF_TYPE {
		t.Errorf("EHRAccess _type not set correctly, got: %v", ehr.EHRAccess.Type_)
	}
}

func TestValidateValidEHR(t *testing.T) {
	ehrJSON, err := os.ReadFile("../../../tests/fixture/ehr.json")
	if err != nil {
		t.Fatalf("Failed to read EHR fixture: %v", err)
	}

	var ehr EHR
	err = json.Unmarshal(ehrJSON, &ehr)
	if err != nil {
		t.Fatalf("Failed to unmarshal EHR JSON: %v", err)
	}

	validateErr := ehr.Validate("$")
	if len(validateErr.Errs) != 0 {
		t.Errorf("Expected no validation errors, got %d", len(validateErr.Errs))
		for _, err := range validateErr.Errs {
			t.Logf("Validation error: %s", err.Message)
		}
	}
}

func TestValidateInvalidEHR(t *testing.T) {
	ehr := EHR{
		EHRID: HIER_OBJECT_ID{Value: "invalid-uuid."},
		EHRStatus: OBJECT_REF{
			Namespace: "local",
			Type:      EHR_STATUS_TYPE,
			ID:        ObjectIDFromHierObjectID(HIER_OBJECT_ID{Value: "also-invalid-uuid"}),
		},
		EHRAccess: OBJECT_REF{
			Namespace: "local",
			Type:      EHR_ACCESS_TYPE,
			ID:        ObjectIDFromHierObjectID(HIER_OBJECT_ID{Value: "yet-another-invalid-uuid"}),
		},
		TimeCreated: DV_DATE_TIME{Value: "2023-01-01T00:00:00Z"},
	}

	validateErr := ehr.Validate("$")
	if len(validateErr.Errs) != 6 {
		t.Errorf("Expected 6 validation errors, got %d", len(validateErr.Errs))
		for _, err := range validateErr.Errs {
			t.Logf("Validation error: %s", err.Message)
		}
	}
}
