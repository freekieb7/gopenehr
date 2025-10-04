package http

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freekieb7/gopenehr/internal/storage"
)

func TestExperimentHandleExecuteQuery(t *testing.T) {

	db := storage.NewDatabase()
	db.Connect(t.Context(), "postgres://postgres:example@localhost:5432/postgres?sslmode=disable")

	handler := NewQueryHandler(&slog.Logger{}, &db)

	data, err := json.Marshal(ExecuteQueryRequest{
		AQL: "SELECT e ex, COUNT(c) FROM EHR e JOIN COMPOSITION c ON e GROUP BY ex",
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/query", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	handler.ExecuteQuery(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Code)
	}
}
