package jquants

import (
	"context"
	"testing"

	"github.com/S-Shiga/jquants-go/v2/codes"
)

func TestClient_MarginTradingVolume(t *testing.T) {
	var code = "13010"
	ctx := context.Background()
	client, err := setup()
	if err != nil {
		t.Fatalf("Failed to setup client: %v", err)
	}
	req := MarginTradingBalanceRequest{Code: &code}
	res, err := client.MarginTradingBalance(ctx, req)
	if err != nil {
		t.Errorf("Failed to get margin trading volume: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty margin trading volume")
	}
}

func TestClient_ShortSellingValue(t *testing.T) {
	var sector33Code = codes.Sector33FisheryAgricultureAndForestry
	ctx := context.Background()
	client, err := setup()
	if err != nil {
		t.Fatalf("Failed to setup client: %v", err)
	}
	req := ShortSellingValueRequest{Sector33Code: &sector33Code}
	res, err := client.ShortSellingValue(ctx, req)
	if err != nil {
		t.Errorf("Failed to get short selling value: %s", err)
	}
	if len(res) == 0 {
		t.Errorf("Empty short selling value")
	}
}

func TestClient_TradingCalendar(t *testing.T) {
	ctx := context.Background()
	client, err := setup()
	if err != nil {
		t.Fatalf("Failed to setup client: %v", err)
	}
	res, err := client.TradingCalendar(ctx, TradingCalendarRequest{})
	if err != nil {
		t.Errorf("Failed to get trading calendar: %s", err)
	}
	if len(res) == 0 {
		t.Errorf("Empty trading calendar")
	}
}
