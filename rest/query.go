package rest

import "net/http"

func HandleExecuteQuery() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}
