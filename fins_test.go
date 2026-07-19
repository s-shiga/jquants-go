package jquants

import "testing"

func TestClient_FinancialSummary(t *testing.T) {
	code := "7203"
	client := setupClient(t)
	res, err := client.FinancialSummary(t.Context(), FinancialSummaryRequest{Code: &code})
	if err != nil {
		t.Errorf("Failed to get financial summary: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty financial summary")
	}
}
