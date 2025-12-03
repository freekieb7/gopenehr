package aql

// func TestExperimentHandleExecuteQuery(t *testing.T) {

// 	db := storage.NewDatabase()
// 	db.Connect(t.Context(), "postgres://postgres:example@localhost:5432/postgres?sslmode=disable")

// 	handler := NewQueryHandler(&slog.Logger{}, &db)

// 	data, err := json.Marshal(ExecuteQueryRequest{
// 		AQL: "SELECT e FROM EHR e JOIN PERSON patient ON e JOIN PARTY_RELATIONSHIP relation IN patient JOIN GROUP group ON patient WHERE relation/name/value = 'patientOf' AND group/uid/value IN (SELECT group/uid/value FROM PERSON practitioner JOIN PARTY_RELATIONSHIP relation IN practitioner JOIN GROUP group ON relation WHERE relation/name/value = 'memberOf') LIMIT 1",
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	req, err := http.NewRequest("POST", "/query", bytes.NewBuffer(data))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	res := httptest.NewRecorder()

// 	handler.ExecuteQuery(res, req)

// 	if res.Code != http.StatusOK {
// 		t.Errorf("expected status OK; got %v", res.Code)
// 	}
// }
