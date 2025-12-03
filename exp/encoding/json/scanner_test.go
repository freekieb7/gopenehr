package json

import (
	"bytes"
	"os"
	"testing"
)

func TestFullScanEhrStatus(t *testing.T) {
	content, err := os.ReadFile("../../../tests/fixture/ehr_status.json")
	if err != nil {
		t.Fatal(err)
	}

	scanner := NewScanner(content)
	for {
		token, err := scanner.Next()
		if err != nil {
			t.Error(err)
			break
		}

		if token.Type == TOKEN_TYPE_END_OF_DOCUMENT {
			break
		}
	}
}

func TestFullScanComposition(t *testing.T) {
	content, err := os.ReadFile("../../../tests/fixture/composition.json")
	if err != nil {
		t.Fatal(err)
	}

	scanner := NewScanner(content)
	for {
		token, err := scanner.Next()
		if err != nil {
			t.Error(err)
			break
		}

		if token.Type == TOKEN_TYPE_END_OF_DOCUMENT {
			break
		}
	}
}

func TestFullScanInstruction(t *testing.T) {
	content, err := os.ReadFile("../../../tests/fixture/instruction.json")
	if err != nil {
		t.Fatal(err)
	}

	scanner := NewScanner(content)
	for {
		token, err := scanner.Next()
		if err != nil {
			t.Error(err)
			break
		}

		if token.Type == TOKEN_TYPE_END_OF_DOCUMENT {
			break
		}
	}
}

func TestFullScanWithTokenCheckItemTree(t *testing.T) {
	content, err := os.ReadFile("../../../tests/fixture/item_tree.json")
	if err != nil {
		t.Fatal(err)
	}

	expectedTokens := []Token{
		{Type: TOKEN_TYPE_OBJECT_BEGIN, Value: []byte("{")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("_type")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("ITEM_TREE")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("name")},
		{Type: TOKEN_TYPE_OBJECT_BEGIN, Value: []byte("{")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("_type")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("DV_TEXT")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("value")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("Arbol")},
		{Type: TOKEN_TYPE_OBJECT_END, Value: []byte("}")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("archetype_node_id")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("at0003")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("items")},
		{Type: TOKEN_TYPE_ARRAY_BEGIN, Value: []byte("[")},
		{Type: TOKEN_TYPE_OBJECT_BEGIN, Value: []byte("{")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("_type")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("ELEMENT")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("name")},
		{Type: TOKEN_TYPE_OBJECT_BEGIN, Value: []byte("{")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("_type")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("DV_TEXT")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("value")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("text")},
		{Type: TOKEN_TYPE_OBJECT_END, Value: []byte("}")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("archetype_node_id")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("at0004")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("value")},
		{Type: TOKEN_TYPE_OBJECT_BEGIN, Value: []byte("{")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("_type")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("DV_TEXT")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("value")},
		{Type: TOKEN_TYPE_STRING, Value: []byte("G,GXATynzlKsNkDfktWnueORsuZDLoxXPVMVDWmVAeIhhOY,,ZFYJoRliTuBaKqzeuPSuhhNIpWmPOdfawBrVvaeUvNJgGxjLomBuTDcuQEDQMfCltwjJWGhaGXHeghOHapAzmVGtybgKISwuarLksqSOGnzUlgkFpWwJqMTQgk.hDUiQfMevlIVZpleLJTwQQZCDhRBtlUzjsjx,KBo nKIAlRqBEcmxcAXXyaVESbtSgqqnkdYqgezTawJpVAZNSQilmqJjm,rmRqFFcRlsuRwVXBrBiG,vYhsb.Cvvme ")},
		{Type: TOKEN_TYPE_OBJECT_END, Value: []byte("}")},
		{Type: TOKEN_TYPE_OBJECT_END, Value: []byte("}")},
		{Type: TOKEN_TYPE_ARRAY_END, Value: []byte("]")},
		{Type: TOKEN_TYPE_OBJECT_END, Value: []byte("}")},
	}

	scanner := NewScanner(content)
	for i, expectedToken := range expectedTokens {
		token, err := scanner.Next()
		if err != nil {
			t.Fatal(err)
		}

		if token.Type != expectedToken.Type {
			t.Errorf("at %d expected type %d, got %d", i, expectedToken.Type, token.Type)
		}

		if !bytes.Equal(token.Value, expectedToken.Value) {
			t.Errorf("at %d expected value %s, got %s", i, expectedToken.Value, token.Value)
		}
	}
}

func TestFullScanCluster(t *testing.T) {
	content, err := os.ReadFile("../../../tests/fixture/cluster.json")
	if err != nil {
		t.Fatal(err)
	}

	scanner := NewScanner(content)
	for {
		token, err := scanner.Next()
		if err != nil {
			t.Error(err)
			break
		}

		if token.Type == TOKEN_TYPE_END_OF_DOCUMENT {
			break
		}
	}
}

func BenchmarkEhrStatus(b *testing.B) {
	b.ReportAllocs()

	content, err := os.ReadFile("../../../tests/fixture/ehr_status.json")
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		scanner := NewScanner(content)

		for {
			token, err := scanner.Next()
			if err != nil {
				b.Fatal(err)
				break
			}

			if token.Type == TOKEN_TYPE_END_OF_DOCUMENT {
				break
			}
		}
	}

}
