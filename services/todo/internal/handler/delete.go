package handler

import (
	"database/sql"
	"droplets_mini/pkg/utilities"
	"droplets_mini/services/todo/internal/models/dto"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (handler *HandlerRepo) DeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := handler.service.DeleteTask(pidParam); err != nil {
		log.Printf("[Todo] deleting task errot: %v", err)
		errBody := dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
		if err := utilities.ResponseJSON(w, http.StatusInternalServerError, errBody, nil); err != nil {
			log.Printf("[Todo] responding to deleting task: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := utilities.ResponseJSON(w, http.StatusNoContent, nil, nil); err != nil {
		log.Printf("[Todo] responding to deleting task success: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
