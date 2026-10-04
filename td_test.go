package jquants

import (
	"encoding/json"
	"testing"
)

func TestClient_TimelyDisclosure(t *testing.T) {
	req := TimelyDisclosureRequest{Date: ptr("2026-07-17"), DiscItems: ptr("101,102")}
	checkEndpoint(t, "/td/list", "date=2026-07-17&discItems=101%2C102", `{"DiscNo":"20260717000001","RevNo":1,"DiscStatus":"revision","Docs":["g","x"],"DiscItems":["101","102"]}`, true, TimelyDisclosure{DisclosureNumber: "20260717000001", RevisionNumber: "1", DisclosureStatus: ptr("revision"), Documents: []string{"g", "x"}, DisclosureItems: []string{"101", "102"}}, func(c *Client) ([]TimelyDisclosure, error) {
		return c.TimelyDisclosure(t.Context(), req)
	})
}

func TestClient_TimelyDisclosureWithChannel(t *testing.T) {
	req := TimelyDisclosureRequest{Date: ptr("2026-07-17"), DiscItems: ptr("101,102")}
	checkEndpoint(t, "/td/list", "date=2026-07-17&discItems=101%2C102", `{"DiscNo":"20260717000001","RevNo":1,"DiscStatus":"revision","Docs":["g","x"],"DiscItems":["101","102"]}`, true, TimelyDisclosure{DisclosureNumber: "20260717000001", RevisionNumber: "1", DisclosureStatus: ptr("revision"), Documents: []string{"g", "x"}, DisclosureItems: []string{"101", "102"}}, func(c *Client) ([]TimelyDisclosure, error) {
		return collectChannel(func(ch chan<- TimelyDisclosure) error { return c.TimelyDisclosureWithChannel(t.Context(), req, ch) })
	})
}

func TestTimelyDisclosureRevisionNumber(t *testing.T) {
	for _, tc := range []struct{ name, data, want string }{
		{"number", `{"RevNo":2}`, "2"},
		{"quoted number", `{"RevNo":"2"}`, "2"},
		{"null", `{"RevNo":null}`, ""},
		{"omitted", `{}`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got TimelyDisclosure
			if err := json.Unmarshal([]byte(tc.data), &got); err != nil {
				t.Fatal(err)
			}
			if got.RevisionNumber != tc.want {
				t.Fatalf("RevisionNumber = %q, want %q", got.RevisionNumber, tc.want)
			}
		})
	}
}

func TestTimelyDisclosureParametersRejectInvalidCombinations(t *testing.T) {
	for _, req := range []TimelyDisclosureRequest{
		{},
		{Date: ptr("2026-07-17"), Code: ptr("86970")},
		{Date: ptr("2026-07-17"), From: ptr("2026-07-01"), To: ptr("2026-07-31")},
		{Code: ptr("86970"), From: ptr("2026-07-01")},
		{Code: ptr("86970"), To: ptr("2026-07-31")},
	} {
		if _, err := (timelyDisclosureParameters{TimelyDisclosureRequest: req}).values(); err == nil {
			t.Fatalf("values accepted invalid request: %#v", req)
		}
	}
}

func TestTimelyDisclosureParametersAcceptCodeRange(t *testing.T) {
	req := TimelyDisclosureRequest{Code: ptr("86970"), From: ptr("2026-07-01"), To: ptr("2026-07-31")}
	got, err := (timelyDisclosureParameters{TimelyDisclosureRequest: req}).values()
	if err != nil || got.Encode() != "code=86970&from=2026-07-01&to=2026-07-31" {
		t.Fatalf("values = %v, %v", got, err)
	}
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
