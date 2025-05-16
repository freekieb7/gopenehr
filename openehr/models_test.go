package openehr

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	content, err := os.ReadFile("../fixture/ehr_status.json")
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]any
	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	var ehrStatus EHR_STATUS
	if err := PocUnmarshal(reflect.ValueOf(ehrStatus), data); err != nil {
		t.Fatal(err)
	}

	t.Log("\nName: " + ehrStatus.Name.Value + "\n")
}
