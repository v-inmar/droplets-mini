package handlers

import (
	"database/sql"
	dtomodels "droplets_mini_todoservice/internal/models/dto_models"
	"droplets_mini_todoservice/internal/services"
	"droplets_mini_todoservice/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type TodoHandler struct {
	srvc services.TaskServiceRoot
}

func NewTodoHandler(srvc services.TaskServiceRoot) *TodoHandler {
	return &TodoHandler{
		srvc: srvc,
	}
}

func (h *TodoHandler) PostCreateTaskHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *TodoHandler) GetAllTaskHandler(w http.ResponseWriter, r *http.Request) {

	allTasksResp, err := h.srvc.ReadAllTasks(r.Context())
	if err != nil {
		log.Printf("[Todo] error reading all tasks %v", err)
		body := dtomodels.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "server error",
		}

		if err := utils.ResponseJSON(w, body, http.StatusInternalServerError, nil); err != nil {
			log.Printf("[Todo] responding to error reading all tasks: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := utils.ResponseJSON(w, allTasksResp, http.StatusOK, nil); err != nil {
		log.Printf("[Todo] responding to get all tasks: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) PutUpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	pidParam := chi.URLParam(r, "pid")

	pidParamNum, err := strconv.ParseInt(pidParam, 10, 64)
	if err != nil {
		body := dtomodels.ErrorResponse{
			Message: fmt.Sprintf("task with pid %s not found", pidParam),
			Status:  http.StatusNotFound,
		}

		if err := utils.ResponseJSON(w, body, http.StatusNotFound, nil); err != nil {
			log.Printf("[Todo] pid conversion to int64 error response: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	task, err := h.srvc.ReadTaskByPID(r.Context(), pidParamNum)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		body := dtomodels.ErrorResponse{
			Message: fmt.Sprintf("task with pid %s not found", pidParam),
			Status:  http.StatusNotFound,
		}

		if err := utils.ResponseJSON(w, body, http.StatusNotFound, nil); err != nil {
			log.Printf("[Todo] responding to task not found: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err != nil {
		log.Printf("[Todo] failed reading task by pid: %s error: %v", pidParam, err)
		body := dtomodels.ErrorResponse{
			Message: "internal server error",
			Status:  http.StatusInternalServerError,
		}

		if err := utils.ResponseJSON(w, body, http.StatusInternalServerError, nil); err != nil {
			log.Printf("[Todo] responding internal server error while readin task by pid: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	var taskUpdate dtomodels.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&taskUpdate); err != nil {
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

	if len(strings.TrimSpace(taskUpdate.Value)) < 1 {
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

	taskUpdate.Value = strings.TrimSpace(taskUpdate.Value)

	taskResp, err := h.srvc.UpdateTask(r.Context(), &taskUpdate, task)
	if err != nil {
		log.Printf("[Todo] failed updating task by pid: %s error: %v", pidParam, err)
		body := dtomodels.ErrorResponse{
			Message: "internal server error",
			Status:  http.StatusInternalServerError,
		}

		if err := utils.ResponseJSON(w, body, http.StatusInternalServerError, nil); err != nil {
			log.Printf("[Todo] responding internal server error while reading task by pid:%s error: %v", pidParam, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := utils.ResponseJSON(w, taskResp, http.StatusOK, nil); err != nil {
		log.Printf("[Todo] error responding to successful update: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *TodoHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	pidParam := chi.URLParam(r, "pid")
	pidParamNum, err := strconv.ParseInt(pidParam, 10, 64)
	if err != nil {
		body := dtomodels.ErrorResponse{
			Message: fmt.Sprintf("task with pid %s not found", pidParam),
			Status:  http.StatusNotFound,
		}

		if err := utils.ResponseJSON(w, body, http.StatusNotFound, nil); err != nil {
			log.Printf("[Todo] pid conversion to int64 error response: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	task, err := h.srvc.ReadTaskByPIDModel(r.Context(), pidParamNum)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		body := dtomodels.ErrorResponse{
			Message: fmt.Sprintf("task with pid %s not found", pidParam),
			Status:  http.StatusNotFound,
		}

		if err := utils.ResponseJSON(w, body, http.StatusNotFound, nil); err != nil {
			log.Printf("[Todo] responding to task not found: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err != nil {
		log.Printf("[Todo] failed reading task by pid: %s error: %v", pidParam, err)
		body := dtomodels.ErrorResponse{
			Message: "internal server error",
			Status:  http.StatusInternalServerError,
		}

		if err := utils.ResponseJSON(w, body, http.StatusInternalServerError, nil); err != nil {
			log.Printf("[Todo] responding internal server error while reading task by pid: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := h.srvc.DeleteTask(r.Context(), task); err != nil {
		log.Printf("[Todo] failed service deleting task by pid: %s error: %v", pidParam, err)
		body := dtomodels.ErrorResponse{
			Message: "internal server error",
			Status:  http.StatusInternalServerError,
		}

		if err := utils.ResponseJSON(w, body, http.StatusInternalServerError, nil); err != nil {
			log.Printf("[Todo] responding internal server error while failed service deleting by pid:%s  error:%v", pidParam, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := utils.ResponseJSON(w, nil, http.StatusNoContent, nil); err != nil {
		log.Printf("[Todo] responding to deleting task success: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
