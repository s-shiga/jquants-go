package jquants

import (
	"context"
	"testing"
)

func TestClient_FuturesPrice(t *testing.T) {
	date := "2026-07-17"
	category := "TOPIXF"
	client := setupClient(t)
	req := FuturesPriceRequest{Date: date, Category: &category}
	resp, err := client.FuturesPrice(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get futures price: %v", err)
	}
	if len(resp) == 0 {
		t.Error("Empty response")
	}
}

func TestClient_FuturesPriceWithChannel(t *testing.T) {
	date := "2026-07-17"
	category := "TOPIXF"
	client := setupClient(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	req := FuturesPriceRequest{Date: date, Category: &category}
	ch := make(chan FuturesPrice)
	go func() {
		if e := client.FuturesPriceWithChannel(ctx, req, ch); e != nil {
			t.Errorf("Failed to get futures price: %v", e)
		}
	}()
	found := false
	for range ch {
		found = true
	}
	if !found {
		t.Error("Empty response")
	}
}
