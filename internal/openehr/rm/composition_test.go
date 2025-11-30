package rm

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCompositionUnmarshal(t *testing.T) {
	content, err := os.ReadFile("../../../tests/fixture/composition.json")
	if err != nil {
		t.Fatal(err)
	}

	var composition COMPOSITION
	err = json.Unmarshal(content, &composition)
	if err != nil {
		t.Fatal(err)
	}

	if composition.Name.Value.(*DV_TEXT).Value != "Test all types" {
		t.Error("COMPOSITION name is not decoded properly")
	}

	if !composition.Content.E {
		t.Error("COMPOSITION content is not decoded properly")
	}

	if len(composition.Content.V) != 3 {
		t.Errorf("Expected 3 content items, got %d", len(composition.Content.V))
	}

	observation, ok := composition.Content.V[0].Value.(*OBSERVATION)
	if !ok {
		t.Error("First content item is not OBSERVATION")
	}

	observationName, ok := observation.Name.Value.(*DV_TEXT)
	if !ok || observationName.Value != "Test all types" {
		t.Error("OBSERVATION name is not decoded properly")
	}
}

func BenchmarkCompositionUnmarshal(b *testing.B) {
	content, err := os.ReadFile("../../../tests/fixture/composition.json")
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		var composition COMPOSITION
		err = json.Unmarshal(content, &composition)
		if err != nil {
			b.Fatal(err)
		}
	}
}
