package errors

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errResp := HTTPError{
		Code:    code,
		Message: message,
	}

	_ = json.NewEncoder(w).Encode(errResp)
}
