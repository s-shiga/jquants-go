package jquants

import (
	"context"
	"errors"
	"testing"

	"github.com/s-shiga/jquants-go/v2/codes"
)

func TestClient_IssueInformation(t *testing.T) {
	client := setupClient(t)
	resp, err := client.IssueInformation(t.Context(), IssueInformationRequest{})
	if err != nil {
		t.Errorf("Failed to get issue information: %v", err)
	}
	if len(resp) == 0 {
		t.Error("Empty response")
	}
}

func TestClient_StockPrice(t *testing.T) {
	var code = "13010"
	client := setupClient(t)
	req := StockPriceRequest{Code: &code}
	res, err := client.StockPrice(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get stock price: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty stock price")
	}
}

func TestClient_StockPriceWithChannel(t *testing.T) {
	var code = "13010"
	client := setupClient(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	req := StockPriceRequest{Code: &code}
	ch := make(chan StockPrice)
	go func() {
		if e := client.StockPriceWithChannel(ctx, req, ch); e != nil {
			t.Errorf("Failed to get stock price: %s", e)
		}
	}()
	found := false
	for range ch {
		found = true
	}
	if !found {
		t.Error("Empty stock price")
	}
}

func TestClient_MinuteStockPrice(t *testing.T) {
	code := "86970"
	from := "2026-07-17"
	to := "2026-07-17"
	client := setupClient(t)
	req := MinuteStockPriceRequest{Code: &code, From: &from, To: &to}
	res, err := client.MinuteStockPrice(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get minute stock price: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty minute stock price")
	}
}

func TestClient_MinuteStockPriceWithChannel(t *testing.T) {
	code := "86970"
	from := "2026-07-17"
	to := "2026-07-17"
	client := setupClient(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	req := MinuteStockPriceRequest{Code: &code, From: &from, To: &to}
	ch := make(chan MinuteStockPrice)
	go func() {
		if e := client.MinuteStockPriceWithChannel(ctx, req, ch); e != nil {
			t.Errorf("Failed to get minute stock price: %s", e)
		}
	}()
	found := false
	for range ch {
		found = true
	}
	if !found {
		t.Error("Empty minute stock price")
	}
}

func TestClient_MorningSessionStockPrice(t *testing.T) {
	code := "72030"
	client := setupClient(t)
	req := MorningSessionStockPriceRequest{Code: &code}
	res, err := client.MorningSessionStockPrice(t.Context(), req)
	if err != nil {
		var noContent NoContent
		if errors.As(err, &noContent) {
			t.Skip("outside morning session window")
		}
		t.Errorf("Failed to get morning session stock price: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty morning session stock price")
	}
}

func TestClient_EarningsCalendar(t *testing.T) {
	client := setupClient(t)
	res, err := client.EarningsCalendar(t.Context(), EarningsCalendarRequest{})
	if err != nil {
		t.Errorf("Failed to get earnings calendar: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty earnings calendar")
	}
}

func TestClient_InvestorType(t *testing.T) {
	var code = codes.SectionPrime
	client := setupClient(t)
	req := InvestorTypeRequest{Section: &code}
	res, err := client.InvestorType(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get investor type: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty investor type")
	}
}
