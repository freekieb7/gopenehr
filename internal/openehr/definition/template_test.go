package definition

import (
	"encoding/xml"
	"os"
	"testing"
)

func TestTemplateUnmarshal(t *testing.T) {
	// Read the blood pressure template XML
	data, err := os.ReadFile("../../../tests/fixture/blood_pressure.template.xml")
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	var template Template
	err = xml.Unmarshal(data, &template)
	if err != nil {
		t.Fatalf("Failed to unmarshal template: %v", err)
	}

	// Basic validations
	if template.Concept != "Blutdruck" {
		t.Errorf("Expected concept 'Blutdruck', got '%s'", template.Concept)
	}

	if template.TemplateID.Value != "Blutdruck" {
		t.Errorf("Expected template_id 'Blutdruck', got '%s'", template.TemplateID.Value)
	}

	if template.Language.CodeString != "de" {
		t.Errorf("Expected language 'de', got '%s'", template.Language.CodeString)
	}

	if template.Definition.RMTypeName != "COMPOSITION" {
		t.Errorf("Expected definition rm_type_name 'COMPOSITION', got '%s'", template.Definition.RMTypeName)
	}

	if len(template.Definition.Attributes) == 0 {
		t.Error("Expected definition to have attributes")
	}

	if len(template.Annotations) == 0 {
		t.Error("Expected template to have annotations")
	}

	t.Logf("Template loaded successfully: %s", template.Concept)
	t.Logf("  UID: %s", template.UID.Value)
	t.Logf("  Lifecycle: %s", template.Description.LifecycleState)
	t.Logf("  Attributes: %d", len(template.Definition.Attributes))
	t.Logf("  Annotations: %d", len(template.Annotations))
}
