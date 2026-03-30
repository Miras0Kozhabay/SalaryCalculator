package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

	// Limit request body size to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var req models.CalculateRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		// Handle JSON parsing errors (400 Bad Request)
		if err == io.EOF {
			jsonError(w, "request body is empty", http.StatusBadRequest)
		} else if err.Error() == "http: request body too large" {
			jsonError(w, "request body too large (max 1MB)", http.StatusRequestEntityTooLarge)
		} else {
			var syntaxErr *json.SyntaxError
			var typeErr *json.UnmarshalTypeError

			if errors.As(err, &syntaxErr) {
				jsonError(w, fmt.Sprintf("invalid JSON at byte offset %d", syntaxErr.Offset), http.StatusBadRequest)
			} else if errors.As(err, &typeErr) {
				jsonError(w, fmt.Sprintf("invalid type for field %q: expected %s", typeErr.Field, typeErr.Type), http.StatusBadRequest)
			} else {
				jsonError(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
			}
		}
		log.Printf("error decoding request: %v", err)
		return
	}

	// Perform calculation (validation happens in service)
	resp, err := h.Service.Calculate(&req)
	if err != nil {
		// Distinguish error types via errors.Is()
		if errors.Is(err, services.ErrInvalidInput) {
			// Validation error → 400 Bad Request
			jsonError(w, "invalid input: "+err.Error(), http.StatusBadRequest)
		} else if errors.Is(err, services.ErrCalculation) {
			// Calculation or database error → 500 Internal Server Error
			jsonError(w, "calculation failed", http.StatusInternalServerError)
			log.Printf("calculation/database error: %v", err)
		} else {
			// Unknown error → 500
			jsonError(w, "internal server error", http.StatusInternalServerError)
			log.Printf("unknown error: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
