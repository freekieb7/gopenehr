package json_test

import (
	"os"
	"testing"

	"github.com/freekieb7/gopenehr"
	"github.com/freekieb7/gopenehr/encoding/json"
)

func BenchmarkUnmarshalEhrStatus(b *testing.B) {
	content, err := os.ReadFile("../../fixture/ehr_status.json")
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		var ehrStatus gopenehr.EHR_STATUS
		if err := json.Unmarshal(content, &ehrStatus); err != nil {
			b.Fatal(err)
		}
	}

}

func TestUnmarshalEhrStatus(t *testing.T) {
	content, err := os.ReadFile("../../fixture/ehr_status.json")
	if err != nil {
		t.Fatal(err)
	}

	var ehrStatus gopenehr.EHR_STATUS
	if err := json.Unmarshal(content, &ehrStatus); err != nil {
		t.Fatal(err)
	}

	if ehrStatus.Type_.Unwrap() != "EHR_STATUS" {
		t.Error("EHR STATUS _type is not decoded properly")
	}

	if ehrStatus.Name.Value != "EHR Status" {
		t.Error("EHR STATUS name is not decoded properly")
	}
}

func TestUnmarshalInstruction(t *testing.T) {
	content, err := os.ReadFile("../../fixture/instruction.json")
	if err != nil {
		t.Fatal(err)
	}

	var instruction gopenehr.INSTRUCTION
	if err := json.Unmarshal(content, &instruction); err != nil {
		t.Fatal(err)
	}

	if instruction.Type_.Unwrap() != "INSTRUCTION" {
		t.Error("EHR STATUS _type is not decoded properly")
	}

	if instruction.Name.Value != "Test all types" {
		t.Error("EHR STATUS name is not decoded properly")
	}
}

func TestUnmarshalComposition(t *testing.T) {
	content, err := os.ReadFile("../../fixture/composition.json")
	if err != nil {
		t.Fatal(err)
	}

	var composition gopenehr.COMPOSITION
	err = json.Unmarshal(content, &composition)
	if err != nil {
		t.Fatal(err)
	}

	if composition.Name.Value != "Test all types" {
		t.Error("COMPOSITION name is not decoded properly")
	}
}

func TestUnmarshalItemTree(t *testing.T) {
	content, err := os.ReadFile("../../fixture/item_tree.json")
	if err != nil {
		t.Fatal(err)
	}

	var itemTree gopenehr.ITEM_TREE
	err = json.Unmarshal(content, &itemTree)
	if err != nil {
		t.Fatal(err)
	}

	if itemTree.Name.Value != "Arbol" {
		t.Error("ITEM_TREE name is not decoded properly")
	}
}

func TestUnmarshalCluster(t *testing.T) {
	content, err := os.ReadFile("../../fixture/cluster.json")
	if err != nil {
		t.Fatal(err)
	}

	var cluster gopenehr.CLUSTER
	err = json.Unmarshal(content, &cluster)
	if err != nil {
		t.Fatal(err)
	}

	if cluster.Name.Value != "cluster 3" {
		t.Error("CLUSTER name is not decoded properly")
	}
}
