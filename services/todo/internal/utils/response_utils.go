package utils

import (
	"encoding/json"
	"maps"
	"net/http"
)

// Writes data as a JSON response with the given status and headers.
func ResponseJSON(w http.ResponseWriter, data any, status int, headers http.Header) error {
	body, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	maps.Copy(w.Header(), headers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(body)
	return err

}
