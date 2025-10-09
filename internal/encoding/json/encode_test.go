package json_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/freekieb7/gopenehr/internal/encoding/json"
	"github.com/freekieb7/gopenehr/internal/model"
)

func TestEncode(t *testing.T) {
	content, err := os.ReadFile("../../fixture/ehr_status.json")
	if err != nil {
		t.Fatal(err)
	}

	var ehrStatus model.EHR_STATUS
	err = json.Unmarshal(content, &ehrStatus)
	if err != nil {
		t.Fatal(err)
	}

	objectContent, err := json.Marshal(ehrStatus)
	if err != nil {
		t.Fatal(err)
	}

	var fileToMap map[string]any
	var objectToMap map[string]any
	json.Unmarshal(content, &fileToMap)
	json.Unmarshal(objectContent, &objectToMap)

	if !reflect.DeepEqual(fileToMap, objectToMap) {
		t.Log(fileToMap)
		t.Log(objectToMap)
		t.Error("EHR STATUS is not decoded properly")
	}
}
