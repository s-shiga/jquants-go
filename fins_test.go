package jquants

import (
	"testing"
)

func TestClient_FinancialSummary(t *testing.T) {
	req := FinancialSummaryRequest{Code: ptr("86970"), Date: ptr("2026-07-17")}
	checkEndpoint(t, "/fins/summary", "code=86970&date=2026-07-17", `{"Code":"86970","Sales":"1000000","EPS":"12.34","DEPS":"","FDiv1Q":"-","FDiv2Q":"*"}`, true, FinancialSummary{Code: "86970", NetSales: ptr(float64(1000000)), EarningsPerShare: ptr(12.34)}, func(c *Client) ([]FinancialSummary, error) {
		return c.FinancialSummary(t.Context(), req)
	})
}

func TestClient_FinancialDetails(t *testing.T) {
	req := FinancialDetailsRequest{Code: ptr("86970")}
	checkEndpoint(t, "/fins/details", "code=86970", `{"Code":"86970","FS":{"Assets":"12345","Notes":"注記"}}`, true, FinancialDetails{Code: "86970", FinancialStatement: map[string]string{"Assets": "12345", "Notes": "注記"}}, func(c *Client) ([]FinancialDetails, error) {
		return c.FinancialDetails(t.Context(), req)
	})
}

func TestClient_Dividend(t *testing.T) {
	req := DividendRequest{Code: ptr("86970")}
	checkEndpoint(t, "/fins/dividend", "code=86970", `{"Code":"86970","DivRate":12.5,"DistAmt":"-","RetEarn":null}`, true, Dividend{Code: "86970", DividendRate: ptr(12.5)}, func(c *Client) ([]Dividend, error) {
		return c.Dividend(t.Context(), req)
	})
}
