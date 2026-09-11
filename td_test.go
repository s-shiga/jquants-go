package jquants

import (
	"testing"
)

func TestClient_TimelyDisclosure(t *testing.T) {
	req := TimelyDisclosureRequest{Date: ptr("2026-07-17"), DiscItems: ptr("101,102")}
	checkEndpoint(t, "/td/list", "date=2026-07-17&discItems=101%2C102", `{"DiscNo":"20260717000001","RevNo":"1","DiscStatus":"revision","Docs":["g","x"],"DiscItems":["101","102"]}`, true, TimelyDisclosure{DisclosureNumber: "20260717000001", RevisionNumber: "1", DisclosureStatus: ptr("revision"), Documents: []string{"g", "x"}, DisclosureItems: []string{"101", "102"}}, func(c *Client) ([]TimelyDisclosure, error) {
		return c.TimelyDisclosure(t.Context(), req)
	})
}

func TestClient_TimelyDisclosureWithChannel(t *testing.T) {
	req := TimelyDisclosureRequest{Date: ptr("2026-07-17"), DiscItems: ptr("101,102")}
	checkEndpoint(t, "/td/list", "date=2026-07-17&discItems=101%2C102", `{"DiscNo":"20260717000001","RevNo":"1","DiscStatus":"revision","Docs":["g","x"],"DiscItems":["101","102"]}`, true, TimelyDisclosure{DisclosureNumber: "20260717000001", RevisionNumber: "1", DisclosureStatus: ptr("revision"), Documents: []string{"g", "x"}, DisclosureItems: []string{"101", "102"}}, func(c *Client) ([]TimelyDisclosure, error) {
		return collectChannel(func(ch chan<- TimelyDisclosure) error { return c.TimelyDisclosureWithChannel(t.Context(), req, ch) })
	})
}

func TestClient_TimelyDisclosureFiles(t *testing.T) {
	c := fixtureClient(t, fixtureResponse{path: "/td/files", query: "discNo=20260717000001&docs=g%2Cx", body: `{"discNo":"20260717000001","files":{"pdf":"https://download.invalid/report.pdf","summaryPdf":null,"xbrl":"https://download.invalid/report.zip"}}`})
	got, err := c.TimelyDisclosureFiles(t.Context(), TimelyDisclosureFilesRequest{DisclosureNumber: "20260717000001", Docs: ptr("g,x")})
	if err != nil {
		t.Fatal(err)
	}
	if got.DisclosureNumber != "20260717000001" || got.Files.PDF == nil || *got.Files.PDF != "https://download.invalid/report.pdf" || got.Files.SummaryPDF != nil || got.Files.XBRL == nil || *got.Files.XBRL != "https://download.invalid/report.zip" {
		t.Fatalf("unexpected files: %#v", got)
	}
}

func TestClient_TimelyDisclosureBulk(t *testing.T) {
	c := fixtureClient(t, fixtureResponse{path: "/td/bulk", body: `{"lastUpdated":"2026-07-17T12:00:00Z","url":"https://download.invalid/td.csv.gz"}`})
	got, err := c.TimelyDisclosureBulk(t.Context())
	if err != nil || got.LastUpdated != "2026-07-17T12:00:00Z" || got.URL != "https://download.invalid/td.csv.gz" {
		t.Fatalf("bulk = %#v, %v", got, err)
	}
}
