package jquants

import (
	"os"
	"testing"
)

func setupClient(t *testing.T) *Client {
	t.Helper()
	apiKey, ok := os.LookupEnv("J_QUANTS_API_KEY")
	if !ok {
		t.Fatal("J_QUANTS_API_KEY environment variable is not set")
	}
	return NewClient(BaseURL, apiKey)
}
