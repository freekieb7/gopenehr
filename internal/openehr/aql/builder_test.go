package aql

import "testing"

func TestToSQL(t *testing.T) {
	sql, _, err := ToSQL("SELECT e, c FROM EHR e CONTAINS COMPOSITION c", nil)
	if err != nil {
		t.Fatalf("ToSQL returned an error: %v", err)
	}

	if sql == "" {
		t.Fatalf("ToSQL returned an empty SQL query")
	}
	_ = sql // Use sql variable to avoid unused variable error
}
