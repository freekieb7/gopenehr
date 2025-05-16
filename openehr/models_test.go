package openehr

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	content, err := os.ReadFile("../fixture/ehr_status.json")
	if err != nil {
		t.Fatal(err)
	}

	var ehrStatus EHR_STATUS
	err = json.Unmarshal(content, &ehrStatus)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := json.Marshal(ehrStatus)
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile("../result/ehr_status.json", decoded, 0644)

	if slices.Compare(content, decoded) != 0 {
		t.Fatal(errors.New("ehr_status is corrupted"))
	}
}
