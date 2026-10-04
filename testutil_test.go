package jquants

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

type fixtureResponse struct {
	path, query, body string
	status            int
	header            http.Header
}

type fixtureHTTPClient struct {
	t         *testing.T
	mu        sync.Mutex
	responses []fixtureResponse
	calls     int
}

// Do serves only in-memory fixtures. It cannot fall back to the network, even
// when a real API key is present in the environment.
func (f *fixtureHTTPClient) Do(req *http.Request) (*http.Response, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := req.Context().Err(); err != nil {
		return nil, err
	}
	if f.calls >= len(f.responses) {
		f.t.Errorf("unexpected request: %s", req.URL)
		return nil, context.Canceled
	}
	want := f.responses[f.calls]
	f.calls++
	if req.Method != http.MethodGet || req.URL.Scheme != "https" || req.URL.Host != "fixture.invalid" || req.URL.Path != "/v2"+want.path {
		f.t.Errorf("unexpected request: %s %s; want GET /v2%s", req.Method, req.URL, want.path)
	}
	query, err := url.ParseQuery(want.query)
	if err != nil {
		f.t.Errorf("invalid fixture query: %v", err)
	}
	if !reflect.DeepEqual(req.URL.Query(), query) {
		f.t.Errorf("query = %v, want %v", req.URL.Query(), query)
	}
	if req.Header.Get("x-api-key") != "fixture-key" || req.Header.Get("Accept-Encoding") != "gzip" || !strings.HasPrefix(req.Header.Get("User-Agent"), "jquants-go/"+Version) {
		f.t.Error("missing or incorrect authentication, compression, or user-agent header")
	}
	status := want.status
	if status == 0 {
		status = http.StatusOK
	}
	return &http.Response{StatusCode: status, Header: want.header, Body: io.NopCloser(strings.NewReader(want.body))}, nil
}

func fixtureClient(t *testing.T, responses ...fixtureResponse) *Client {
	t.Helper()
	f := &fixtureHTTPClient{t: t, responses: responses}
	t.Cleanup(func() {
		f.mu.Lock()
		defer f.mu.Unlock()
		if f.calls != len(f.responses) {
			t.Errorf("requests = %d, want %d", f.calls, len(f.responses))
		}
	})
	return NewClient("https://fixture.invalid/v2", "fixture-key", WithHTTPClient(f), WithRetryInterval(time.Millisecond), WithLoopTimeout(time.Second))
}

func ptr[T any](v T) *T { return &v }

func checkEndpoint[T any](t *testing.T, path, query, item string, paginated bool, want T, fetch func(*Client) ([]T, error)) {
	t.Helper()
	responses := []fixtureResponse{{path: path, query: query, body: `{"data":[` + item + `]}`}}
	expected := []T{want}
	if paginated {
		responses[0].body = `{"data":[` + item + `],"pagination_key":"next+/="}`
		nextQuery := "pagination_key=next%2B%2F%3D"
		if query != "" {
			nextQuery = query + "&" + nextQuery
		}
		responses = append(responses, fixtureResponse{path: path, query: nextQuery, body: `{"data":[` + item + `],"pagination_key":null}`})
		expected = append(expected, want)
	}
	got, err := fetch(fixtureClient(t, responses...))
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("items = %#v, want %#v", got, expected)
	}
}

// Await the producer's returned error as well as channel closure so tests never
// finish with an unchecked error or a producer still using testing.T.
func collectChannel[T any](fetch func(chan<- T) error) ([]T, error) {
	ch := make(chan T)
	done := make(chan error, 1)
	go func() { done <- fetch(ch) }()
	var items []T
	for item := range ch {
		items = append(items, item)
	}
	return items, <-done
}

// rejected marks a valuesCase whose request values() must refuse.
const rejected = "<rejected>"

type valuesCase struct {
	params parameters
	want   string
}

// checkValues asserts that each request encodes to the wanted query, or that
// values() refuses it, which happens before any HTTP request is sent.
func checkValues(t *testing.T, cases []valuesCase) {
	t.Helper()
	for _, tc := range cases {
		got, err := tc.params.values()
		if tc.want == rejected {
			if err == nil {
				t.Errorf("values accepted %#v as %q", tc.params, got.Encode())
			}
			continue
		}
		if err != nil || got.Encode() != tc.want {
			t.Errorf("values(%#v) = %q, %v; want %q", tc.params, got.Encode(), err, tc.want)
		}
	}
}
