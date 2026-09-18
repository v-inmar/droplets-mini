package handlers

import (
	dtomodels "droplets_mini/services/todo/internal/models/dto_models"
	"droplets_mini/services/todo/internal/services"
	"droplets_mini/services/todo/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type TodoHandler struct {
	srvc services.TaskServiceRoot
}

func NewTodoHandler(srvc services.TaskServiceRoot) *TodoHandler {
	return &TodoHandler{
		srvc: srvc,
	}
}

func (h *TodoHandler) PostCreateTask(w http.ResponseWriter, r *http.Request) {
	var task dtomodels.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("[Todo] invalid json: %v", err)
		body := dtomodels.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid json",
		}

		if err := utils.ResponseJSON(w, body, http.StatusBadRequest, nil); err != nil {
			log.Printf("[Todo] responding to invalid json: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	task.Value = strings.TrimSpace(task.Value)
	if len(task.Value) < 1 {
		log.Print("[Todo] missing value json")
		body := dtomodels.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "missing value",
		}

		if err := utils.ResponseJSON(w, body, http.StatusBadRequest, nil); err != nil {
			log.Printf("[Todo] responding to missing value: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	taskResp, err := h.srvc.CreateNewTaskItem(r.Context(), &task)
	if err != nil {
		log.Printf("[Todo] error creating new task item: %v", err)
		body := dtomodels.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "server error",
		}

		if err := utils.ResponseJSON(w, body, http.StatusInternalServerError, nil); err != nil {
			log.Printf("[Todo] responding to error creating new task item: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := utils.ResponseJSON(w, taskResp, http.StatusCreated, nil); err != nil {
		log.Printf("[Todo] responding to created task item: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
