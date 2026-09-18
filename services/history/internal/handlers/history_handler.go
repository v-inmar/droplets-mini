package handlers

import (
	dtomodels "droplets_mini_historyservice/internal/models/dto_models"
	"droplets_mini_historyservice/internal/services"
	"droplets_mini_historyservice/internal/utils"
	"log"
	"net/http"
)

type HistoryHandler struct {
	srvc services.HistoryServiceRoot
}

func NewHistoryHandler(srvc services.HistoryServiceRoot) *HistoryHandler {
	return &HistoryHandler{
		srvc: srvc,
	}
}

func (h *HistoryHandler) GetAllHistoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[History] GetAllHistoryHandler")
	respBody, err := h.srvc.GetAllHistory(r.Context())
	if err != nil {
		log.Printf("[History] error getting response body from history service: %v", err)
		status := http.StatusInternalServerError
		errBody := dtomodels.ErrorResponse{
			Message: "internal server error",
			Status:  status,
		}

		if err := utils.ResponseJSON(w, errBody, status, nil); err != nil {
			log.Printf("[History] responding error: %v", err)
			http.Error(w, "internal server error", status)
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
