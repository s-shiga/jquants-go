package jquants

import (
	"context"
	"testing"

	"github.com/s-shiga/jquants-go/v2/codes"
)

func TestClient_MarginTradingOutstanding(t *testing.T) {
	var code = "13010"
	client := setupClient(t)
	req := MarginTradingOutstandingRequest{Code: &code}
	res, err := client.MarginTradingOutstanding(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get margin trading outstanding: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty margin trading outstanding")
	}
}

func TestClient_ShortSellingValue(t *testing.T) {
	var sector33Code = codes.Sector33FisheryAgricultureAndForestry
	client := setupClient(t)
	req := ShortSellingValueRequest{Sector33Code: &sector33Code}
	res, err := client.ShortSellingValue(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get short selling value: %s", err)
	}
	if len(res) == 0 {
		t.Errorf("Empty short selling value")
	}
}

func TestClient_OutstandingShortPosition(t *testing.T) {
	calcDate := "20260210"
	client := setupClient(t)
	req := OutstandingShortPositionRequest{CalculationDate: &calcDate}
	res, err := client.OutstandingShortPosition(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get outstanding short position: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty outstanding short position")
	}
}

func TestClient_MarginAlert(t *testing.T) {
	date := "20260213"
	client := setupClient(t)
	req := MarginAlertRequest{Date: &date}
	res, err := client.MarginAlert(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get margin alert: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty margin alert")
	}
}

func TestClient_BreakdownTrading(t *testing.T) {
	code := "13010"
	from := "2026-07-01"
	to := "2026-07-17"
	client := setupClient(t)
	req := BreakdownTradingRequest{Code: &code, From: &from, To: &to}
	res, err := client.BreakdownTrading(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get breakdown trading: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty breakdown trading")
	}
}

func TestClient_BreakdownTradingWithChannel(t *testing.T) {
	code := "13010"
	from := "2026-07-01"
	to := "2026-07-17"
	client := setupClient(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	req := BreakdownTradingRequest{Code: &code, From: &from, To: &to}
	ch := make(chan BreakdownTrading)
	go func() {
		if e := client.BreakdownTradingWithChannel(ctx, req, ch); e != nil {
			t.Errorf("Failed to get breakdown trading: %s", e)
		}
	}()
	found := false
	for range ch {
		found = true
	}
	if !found {
		t.Error("Empty breakdown trading")
	}
}

func TestClient_TradingCalendar(t *testing.T) {
	client := setupClient(t)
	res, err := client.TradingCalendar(t.Context(), TradingCalendarRequest{})
	if err != nil {
		t.Errorf("Failed to get trading calendar: %s", err)
	}
	if len(res) == 0 {
		t.Errorf("Empty trading calendar")
	}
}
