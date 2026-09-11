package jquants

import "testing"

func TestPageImplementsResponse(t *testing.T) {
	key := "next"
	response := Response[int](page[int]{Data: []int{1, 2}, PaginationKey: &key})
	if len(response.Items()) != 2 || response.Items()[1] != 2 {
		t.Fatalf("items = %v, want [1 2]", response.Items())
	}
	if response.NextPageKey() == nil || *response.NextPageKey() != key {
		t.Fatalf("next page key = %v, want %q", response.NextPageKey(), key)
	}
}
