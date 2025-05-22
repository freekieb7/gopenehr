package aql

import "testing"

func TestScannerQueries(t *testing.T) {
	qList := []string{
		// "SELECT *",
		"SELECT * FROM A",
	}

	for _, q := range qList {
		scanner := NewScanner([]byte(q))

		for {
			tok, err := scanner.Next()
			if err != nil {
				t.Error(err)
				break
			}

			if tok.Type == TOKEN_TYPE_QUERY_END {
				break
			}
		}
	}
}
