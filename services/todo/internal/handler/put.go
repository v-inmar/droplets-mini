package handler

import (
	"database/sql"
	"droplets_mini/pkg/utilities"
	"droplets_mini/services/todo/internal/models/dto"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (handler *HandlerRepo) PutHandler(w http.ResponseWriter, r *http.Request) {

	pidParam := chi.URLParam(r, "pid")

	_, err := handler.service.ReadTaskByPID(pidParam)
	if err != nil {

		log.Printf("[Todo] error reading task by pid: %v", err)

		status := http.StatusInternalServerError
		message := "internval server error"

		if errors.Is(err, sql.ErrNoRows) {
			status = http.StatusNotFound
			message = fmt.Sprintf("task with pid: %s does not exist", pidParam)
		}

		errBody := dto.ErrorResponse{
			Status:  status,
			Message: message,
		}

		if err := utilities.ResponseJSON(w, status, errBody, nil); err != nil {
			log.Printf("[Todo] responding to task not found: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	var taskRequest dto.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
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

	respBody, err := handler.service.UpdateTask(&taskRequest, pidParam)
	if err != nil {
		log.Printf("[Todo] update task service: %v", err)
		body := dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}

		if err := utilities.ResponseJSON(w, http.StatusBadRequest, body, nil); err != nil {
			log.Printf("[Todo] responding to update task service: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := utilities.ResponseJSON(w, http.StatusOK, respBody, nil); err != nil {
		log.Printf("[Todo] responding to successful update: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}
