package handler

import (
	"droplets_mini/services/todo/internal/models/dto"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (handler *HandlerRepo) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	respModel := dto.HealthCheckResponse{
		Message: "todo server running",
		Status:  status,
		Time:    time.Now().UTC(),
	}

	respBody, err := json.Marshal(respModel)
	if err != nil {
		log.Printf("error while encoding response data for health check: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBody); err != nil {
		log.Printf("error write JSON response for healthcheck: %v", err)
	}
}
