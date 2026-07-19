package jquants

import (
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
