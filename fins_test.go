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

func TestClient_FinancialDetails(t *testing.T) {
	code := "86970"
	client := setupClient(t)
	res, err := client.FinancialDetails(t.Context(), FinancialDetailsRequest{Code: &code})
	if err != nil {
		t.Errorf("Failed to get financial details: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty financial details")
	}
}

func TestClient_Dividend(t *testing.T) {
	code := "86970"
	client := setupClient(t)
	res, err := client.Dividend(t.Context(), DividendRequest{Code: &code})
	if err != nil {
		t.Errorf("Failed to get dividend: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty dividend")
	}
}
