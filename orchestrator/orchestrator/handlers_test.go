package orchestrator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func float64Ptr(f float64) *float64 {
	return &f
}

func TestHandlePing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	w := httptest.NewRecorder()

	handlePing(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestHandleCalculate(t *testing.T) {
	reqBody := `{"expression": "2 + 2 * 2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	handleCalculate(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}

	var resp map[string]string
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["id"] == "" {
		t.Error("expected non-empty id in response")
	}
}

func TestHandleListExpressions(t *testing.T) {
	expressionsStore = make(map[string]*Expression)
	expressionsStore["test1"] = &Expression{
		ID:     "test1",
		Expr:   "2+2",
		Status: "done",
		Result: float64Ptr(4),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/expressions", nil)
	w := httptest.NewRecorder()

	handleListExpressions(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	var resp struct {
		Expressions []Expression `json:"expressions"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Expressions) != 1 {
		t.Errorf("expected 1 expression, got %d", len(resp.Expressions))
	}
}

func TestHandleGetExpression(t *testing.T) {
	expressionsStore = make(map[string]*Expression)
	expressionsStore["test123"] = &Expression{
		ID:     "test123",
		Expr:   "(2+2)*(3+3)",
		Status: "done",
		Result: float64Ptr(24),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/expressions/test123", nil)
	w := httptest.NewRecorder()

	handleGetExpression(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	var resp map[string]Expression
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expr, ok := resp["expression"]
	if !ok {
		t.Fatal("response JSON does not have key 'expression'")
	}

	if expr.ID != "test123" {
		t.Errorf("expected id 'test123', got %s", expr.ID)
	}
}
