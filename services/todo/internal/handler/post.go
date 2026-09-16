package handler

import (
	"droplets_mini/pkg/utilities"
	"droplets_mini/services/todo/internal/models/dto"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (handler *HandlerRepo) PostHandler(w http.ResponseWriter, r *http.Request) {

	var task dto.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("[Todo] invalid json: %v", err)
		body := dto.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid json",
		}

		if err := utilities.ResponseJSON(w, http.StatusBadRequest, body, nil); err != nil {
			log.Printf("[Todo] responding to invalid json: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	taskValue := strings.TrimSpace(task.Value)
	if len(taskValue) < 1 {
		log.Print("[Todo] missing value json")
		body := dto.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "missing value",
		}

		if err := utilities.ResponseJSON(w, http.StatusBadRequest, body, nil); err != nil {
			log.Printf("[Todo] responding to missing value: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	task.Value = taskValue
	taskModel, err := handler.service.CreateNewTaskItem(&task)
	if err != nil {
		log.Printf("[Todo] error creating new task item: %v", err)
		body := dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "server error",
		}

		if err := utilities.ResponseJSON(w, http.StatusInternalServerError, body, nil); err != nil {
			log.Printf("[Todo] responding to error creating new task item: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	body := dto.CreateResponse{
		PID:       taskModel.PID,
		Value:     taskModel.Value,
		Completed: taskModel.Completed,
		CreatedAt: taskModel.CreatedAt,
	}
	if err := utilities.ResponseJSON(w, http.StatusCreated, body, nil); err != nil {
		log.Printf("[Todo] responding to created task item: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}
