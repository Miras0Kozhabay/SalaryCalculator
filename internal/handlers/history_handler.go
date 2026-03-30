package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"salary-calculator/internal/models"
)

func (h *SalaryHandler) History(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // prevent excessive queries
	}

	calcs, err := h.Service.GetHistory(limit, offset)
	if err != nil {
		log.Printf("error retrieving history: %v", err)
		jsonError(w, "failed to get history", http.StatusInternalServerError)
		return
	}

	// Return empty array instead of null if no results
	if calcs == nil {
		calcs = make([]*models.Calculation, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(calcs); err != nil {
		log.Printf("error encoding history response: %v", err)
	}
}
