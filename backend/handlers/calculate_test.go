package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status field = %q, want %q", body["status"], "ok")
	}
}

func TestCalculateBodyTooLarge(t *testing.T) {
	huge := `{"operation":"add","a":1,"b":1,"padding":"` + strings.Repeat("x", maxRequestBodyBytes) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(huge))
	rec := httptest.NewRecorder()

	Calculate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantResult float64
	}{
		{"valid add", `{"operation":"add","a":2,"b":3}`, http.StatusOK, 5},
		{"valid exponentiate", `{"operation":"exponentiate","a":2,"b":8}`, http.StatusOK, 256},
		{"valid square root", `{"operation":"square_root","a":16}`, http.StatusOK, 4},
		{"valid percentage", `{"operation":"percentage","a":25,"b":200}`, http.StatusOK, 50},
		{"divide by zero", `{"operation":"divide","a":1,"b":0}`, http.StatusBadRequest, 0},
		{"square root of negative", `{"operation":"square_root","a":-9}`, http.StatusBadRequest, 0},
		{"exponent too large", `{"operation":"exponentiate","a":10,"b":1000000}`, http.StatusBadRequest, 0},
		{"operand too large", `{"operation":"add","a":1e16,"b":1}`, http.StatusBadRequest, 0},
		{"bad json", `not-json`, http.StatusBadRequest, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			Calculate(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var resp calculateResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.Result != tt.wantResult {
					t.Fatalf("result = %v, want %v", resp.Result, tt.wantResult)
				}
			}
		})
	}
}
