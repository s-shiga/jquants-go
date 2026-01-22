package jquants

import (
	"context"
	"testing"
)

func TestClient_IssueInformation(t *testing.T) {
	ctx := context.Background()
	client, err := setup()
	if err != nil {
		t.Fatalf("Failed to setup client: %v", err)
	}
	resp, err := client.IssueInformation(ctx, IssueInformationRequest{})
	if err != nil {
		t.Errorf("Failed to get issue information: %v", err)
	}
	if len(resp) == 0 {
		t.Error("Empty response")
	}
}
