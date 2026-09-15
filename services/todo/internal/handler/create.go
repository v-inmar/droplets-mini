package handler

import (
	"droplets_mini/pkg/utilities"
	"droplets_mini/services/todo/internal/models/dto"
	"encoding/json"
	"net/http"
	"strings"
)

func (handler *HandlerRepo) CreateTodo(w http.ResponseWriter, r *http.Request) {

	var todo dto.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		body := dto.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid json",
		}

		if err := utilities.ResponseJSON(w, http.StatusBadRequest, body, nil); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	todoValue := strings.TrimSpace(todo.Value)
	if len(todoValue) < 1 {
		body := dto.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "missing value",
		}

		if err := utilities.ResponseJSON(w, http.StatusBadRequest, body, nil); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

}
