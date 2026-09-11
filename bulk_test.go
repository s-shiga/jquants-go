package jquants

import (
	"encoding/json"
	"testing"
)

func TestClient_BulkList(t *testing.T) {
	req := BulkListRequest{Endpoint: ptr("/equities/bars/daily"), From: ptr("2026-07-01"), To: ptr("2026-07-17")}
	checkEndpoint(t, "/bulk/list", "endpoint=%2Fequities%2Fbars%2Fdaily&from=2026-07-01&to=2026-07-17", `{"Key":"equities/file.csv.gz","Size":1234.0,"LastModified":"2026-07-17T12:00:00Z"}`, false, BulkFile{Key: "equities/file.csv.gz", Size: 1234, LastModified: "2026-07-17T12:00:00Z"}, func(c *Client) ([]BulkFile, error) {
		return c.BulkList(t.Context(), req)
	})
}

func TestClient_BulkGet(t *testing.T) {
	for _, tc := range []struct {
		name  string
		req   BulkGetRequest
		query string
	}{
		{"key", BulkGetRequest{Key: ptr("equities/file +.csv.gz")}, "key=equities%2Ffile+%2B.csv.gz"},
		{"endpoint and date", BulkGetRequest{Endpoint: ptr("/equities/bars/daily"), Date: ptr("2026-07-17")}, "endpoint=%2Fequities%2Fbars%2Fdaily&date=2026-07-17"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := fixtureClient(t, fixtureResponse{path: "/bulk/get", query: tc.query, body: `{"url":"https://download.invalid/file?signature=example"}`})
			got, err := c.BulkGet(t.Context(), tc.req)
			if err != nil || got != "https://download.invalid/file?signature=example" {
				t.Fatalf("BulkGet = %q, %v", got, err)
			}
		})
	}
}

func TestBulkFilePreservesIntegerPrecision(t *testing.T) {
	var file BulkFile
	if err := json.Unmarshal([]byte(`{"Size":9007199254740993}`), &file); err != nil {
		t.Fatal(err)
	}
	if file.Size != 9007199254740993 {
		t.Fatalf("Size = %d", file.Size)
	}
}

func TestBulkFileRejectsFractionalSize(t *testing.T) {
	var file BulkFile
	if err := json.Unmarshal([]byte(`{"Size":1.5}`), &file); err == nil {
		t.Fatal("json.Unmarshal accepted a fractional file size")
	}
}
