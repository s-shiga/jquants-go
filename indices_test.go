package jquants

import (
	"testing"
)

func TestClient_IndexPrice(t *testing.T) {
	var indexCode = "0000"
	client := setupClient(t)
	req := IndexPriceRequest{Code: &indexCode}
	res, err := client.IndexPrice(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get index price: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty index price")
	}
}

func TestClient_TopixPrices(t *testing.T) {
	client := setupClient(t)
	res, err := client.TopixPrices(t.Context(), TopixPriceRequest{})
	if err != nil {
		t.Errorf("Failed to get topix price: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty topix price")
	}
}
