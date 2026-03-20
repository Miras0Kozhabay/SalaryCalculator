package handlers

import (
	"encoding/json"
	"net/http"
	"salary-calculator/internal/models"
	"salary-calculator/internal/services"
)

type SalaryHandler struct {
	Service *services.SalaryService
}

func NewSalaryHandler(s *services.SalaryService) *SalaryHandler {
	return &SalaryHandler{Service: s}
}

func (h *SalaryHandler) Calculate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	resp, err := h.Service.Calculate(&req)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
