package jquants

import (
	"context"
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
