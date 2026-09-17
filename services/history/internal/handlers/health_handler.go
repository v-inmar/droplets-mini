package handlers

import (
	dtomodels "droplets_mini/services/history/internal/models/dto_models"
	"droplets_mini/services/history/internal/services"
	"droplets_mini/services/history/internal/utils"
	"log"
	"net/http"
)

type HealthHandler struct {
	srvc services.HealthServiceRoot
}

func NewHealthHandler(srvc services.HealthServiceRoot) *HealthHandler {
	return &HealthHandler{
		srvc: srvc,
	}
}

func (handler *HealthHandler) GetHealthHandler(w http.ResponseWriter, r *http.Request) {
	respBody, err := handler.srvc.GetHealth()
	if err != nil {
		log.Printf("[History] error getting response body from health service: %v", err)
		errBody := dtomodels.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
		if err := utils.ResponseJSON(w, errBody, http.StatusInternalServerError, nil); err != nil {
			log.Printf("[History] responding error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := utils.ResponseJSON(w, respBody, http.StatusOK, nil); err != nil {
		log.Printf("[History] responding error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
