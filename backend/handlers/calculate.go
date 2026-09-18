package handlers

import (
	"encoding/json"
	"net/http"

	"sezzle-calculator/backend/calculator"
)

type calculateRequest struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

type calculateResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// Health responds with a simple ok status, used for readiness checks.
func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Calculate decodes a calculation request, evaluates it, and encodes the result.
func Calculate(w http.ResponseWriter, r *http.Request) {
	var req calculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	result, err := calculator.Evaluate(req.Operation, req.A, req.B)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, calculateResponse{Result: result})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
