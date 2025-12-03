package json

// func TestEncode(t *testing.T) {
// 	content, err := os.ReadFile("../../../tests/fixture/ehr_status.json")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	var ehrStatus rm.EHR_STATUS
// 	err = UnmarshalJSON(content, &ehrStatus)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	objectContent, err := Marshal(ehrStatus)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	var fileToMap map[string]any
// 	var objectToMap map[string]any
// 	UnmarshalJSON(content, &fileToMap)
// 	UnmarshalJSON(objectContent, &objectToMap)

// 	if !reflect.DeepEqual(fileToMap, objectToMap) {
// 		t.Log(fileToMap)
// 		t.Log(objectToMap)
// 		t.Error("EHR STATUS is not decoded properly")
// 	}
// }
