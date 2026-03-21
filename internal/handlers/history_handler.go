package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
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

	calcs, err := h.Service.GetHistory(limit, offset)
	if err != nil {
		jsonError(w, "failed to get history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calcs)
}
