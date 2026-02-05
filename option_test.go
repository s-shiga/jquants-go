package jquants

import (
	"context"
	"testing"
)

func TestClient_IndexOptionPrice(t *testing.T) {
	date := "2025-01-06"
	client := setupClient(t)
	req := IndexOptionPriceRequest{Date: date}
	resp, err := client.IndexOptionPrice(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get index option price: %v", err)
	}
	if len(resp) == 0 {
		t.Error("Empty response")
	}
}

func TestClient_IndexOptionPriceWithChannel(t *testing.T) {
	date := "2025-01-06"
	client := setupClient(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	req := IndexOptionPriceRequest{Date: date}
	ch := make(chan IndexOptionPrice)
	go func() {
		if e := client.IndexOptionPriceWithChannel(ctx, req, ch); e != nil {
			t.Errorf("Failed to get index option price: %v", e)
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
