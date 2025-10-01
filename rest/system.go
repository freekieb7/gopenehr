package rest

import (
	"encoding/json"
	"net/http"
)

func HandleServerInfo() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		payload := map[string]any{
			"solution":              "openEHRSys",
			"solution_version":      "v1.0",
			"vendor":                "GOpenEHR",
			"restapi_specs_version": "1.0.3",
			"conformance_profile":   "CUSTOM",
			"endpoints": []string{
				"/ehr",
				"/demographics",
				"/definition",
				"/query",
				"/admin",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(payload)
	}
}
