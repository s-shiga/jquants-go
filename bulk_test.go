package jquants

import (
	"strings"
	"testing"
)

func TestClient_BulkList(t *testing.T) {
	date := "2026-07-17"
	client := setupClient(t)
	req := BulkListRequest{Date: &date}
	res, err := client.BulkList(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get bulk list: %s", err)
	}
	if len(res) == 0 {
		t.Fatal("Empty bulk list")
	}
	if res[0].Key == "" {
		t.Error("Empty bulk file key")
	}
}

func TestClient_BulkGet(t *testing.T) {
	date := "2026-07-17"
	client := setupClient(t)
	list, err := client.BulkList(t.Context(), BulkListRequest{Date: &date})
	if err != nil {
		t.Fatalf("Failed to get bulk list: %s", err)
	}
	if len(list) == 0 {
		t.Fatal("Empty bulk list")
	}
	key := list[0].Key
	u, err := client.BulkGet(t.Context(), BulkGetRequest{Key: &key})
	if err != nil {
		t.Errorf("Failed to get bulk download URL: %s", err)
	}
	if !strings.HasPrefix(u, "https://") {
		t.Errorf("Unexpected bulk download URL: %s", u)
	}
}
