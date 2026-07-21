package jquants

import (
	"os"
	"testing"
)

func setupClient(t *testing.T) *Client {
	t.Helper()
	apiKey, ok := os.LookupEnv("J_QUANTS_API_KEY")
	if !ok {
		t.Skip("J_QUANTS_API_KEY environment variable is not set; skipping test against the live API")
	}
	return NewClient(BaseURL, apiKey)
}
