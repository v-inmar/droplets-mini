package handler

import (
	"droplets_mini/pkg/utilities"
	"droplets_mini/services/todo/internal/models/dto"
	"log"
	"net/http"
)

func (handler *HandlerRepo) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := handler.service.ReadAllTasks()
	if err != nil {
		log.Printf("[Todo] error reading all tasks %v", err)
		body := dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "server error",
		}

		if err := utilities.ResponseJSON(w, http.StatusInternalServerError, body, nil); err != nil {
			log.Printf("[Todo] responding to error reading all tasks: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := utilities.ResponseJSON(w, http.StatusOK, tasks, nil); err != nil {
		log.Printf("[Todo] responding to get all tasks: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
