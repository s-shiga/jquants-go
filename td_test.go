package jquants

import (
	"testing"
)

func TestClient_TimelyDisclosure(t *testing.T) {
	date := "2026-07-17"
	client := setupClient(t)
	req := TimelyDisclosureRequest{Date: &date}
	res, err := client.TimelyDisclosure(t.Context(), req)
	if err != nil {
		t.Errorf("Failed to get timely disclosure: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty timely disclosure")
	}
}

func TestClient_TimelyDisclosureFiles(t *testing.T) {
	date := "2026-07-17"
	client := setupClient(t)
	list, err := client.TimelyDisclosure(t.Context(), TimelyDisclosureRequest{Date: &date})
	if err != nil {
		t.Fatalf("Failed to get timely disclosure: %s", err)
	}
	if len(list) == 0 {
		t.Fatal("Empty timely disclosure")
	}
	discNo := list[0].DisclosureNumber
	res, err := client.TimelyDisclosureFiles(t.Context(), TimelyDisclosureFilesRequest{DisclosureNumber: discNo})
	if err != nil {
		t.Errorf("Failed to get timely disclosure files: %s", err)
	}
	if res.DisclosureNumber != discNo {
		t.Errorf("DisclosureNumber mismatch: got %q, want %q", res.DisclosureNumber, discNo)
	}
	if res.Files.PDF == nil && res.Files.SummaryPDF == nil && res.Files.XBRL == nil {
		t.Error("No file URLs returned")
	}
}

func TestClient_TimelyDisclosureBulk(t *testing.T) {
	client := setupClient(t)
	res, err := client.TimelyDisclosureBulk(t.Context())
	if err != nil {
		t.Errorf("Failed to get timely disclosure bulk: %s", err)
	}
	if res.LastUpdated == "" {
		t.Error("Empty last updated")
	}
	if res.URL == "" {
		t.Error("Empty URL")
	}
}
