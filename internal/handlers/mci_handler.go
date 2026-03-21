package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *SalaryHandler) MCI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := map[string]float64{
		"mci": h.Service.Calc.MCI,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
