package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type SchemaHandler struct {
	logger *slog.Logger
}

func NewSchemaHandler() SchemaHandler {
	return SchemaHandler{}
}

var schemas = map[string][]string{
	"myschemav1": {
		// 1. The service name must be set
		`"service.name" in resource && resource["service.name"] != ""`,

		// 2. Only allow certain environments
		`resource["deployment.environment"] in ["prod", "staging", "dev"]`,

		// 3. If environment is prod, region must be eu-west-1
		`!(resource["deployment.environment"] == "prod") || resource["cloud.region"] == "eu-west-1"`,

		// 4. The schema_url must match your known OTel schema format
		// `schema_url.matches('^https://opentelemetry.io/schemas/[0-9]+\\.[0-9]+\\.[0-9]+$')`,
	},
	// "https://my.schema.io/v2": {
	// 	// Same as v1 + stricter currency set
	// 	`currency in ["USD", "EUR", "GBP"]`,
	// 	`!(amount > 500) || approval_id != ""`,
	// },
}

func (h *SchemaHandler) GetSchemas(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(schemas)
}

func (h *SchemaHandler) GetSchema(w http.ResponseWriter, r *http.Request) error {
	schemaName := r.PathValue("name")
	if schemaName == "" {
		http.Error(w, "Missing schema name", http.StatusBadRequest)
		return nil
	}

	schema, ok := schemas[schemaName]
	if !ok {
		http.Error(w, "Schema not found", http.StatusNotFound)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schema)
	return nil
}
