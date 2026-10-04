package jquants

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// transientTestItem and page[transientTestItem] are minimal in-package fixtures
// used to exercise the paginated fetch loop against in-memory responses without
// depending on a live J-Quants endpoint.
type transientTestItem struct {
	Value string `json:"value"`
}

type transientTestParams struct {
	paginationKey *string
}

func (p transientTestParams) values() (url.Values, error) {
	v := url.Values{}
	if p.paginationKey != nil {
		v.Add("pagination_key", *p.paginationKey)
	}
	return v, nil
}

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

// fullBody is a complete, decodable single-page response.
const fullBody = `{"data":[{"value":"ok"}]}`

// roundTripFunc lets tests simulate transport failures without opening sockets.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) Do(req *http.Request) (*http.Response, error) { return f(req) }

type failedReader struct{}

func (failedReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func truncatedResponse() *http.Response {
	return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(io.MultiReader(strings.NewReader(`{"data":[{"value":"ok`), failedReader{}))}
}

func fetchTransientTest(ctx context.Context, c *Client) ([]transientTestItem, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[transientTestItem], error) {
		return getJSON[page[transientTestItem]](ctx, c, "/test", transientTestParams{paginationKey: paginationKey})
	})
}

func TestFetch_RetriesTruncatedBodyThenSucceeds(t *testing.T) {
	var calls atomic.Int32
	client := NewClient("https://fixture.invalid", "fixture-key", WithRetryInterval(time.Millisecond), WithHTTPClient(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if calls.Add(1) == 1 {
			return truncatedResponse(), nil
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(fullBody))}, nil
	})))
	items, err := fetchTransientTest(t.Context(), client)
	if err != nil || len(items) != 1 || items[0].Value != "ok" {
		t.Fatalf("items = %#v, error = %v", items, err)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("requests = %d, want 2", got)
	}
}

func TestFetch_TruncatedBodyEveryRequestFailsWithTransientError(t *testing.T) {
	var calls atomic.Int32
	client := NewClient("https://fixture.invalid", "fixture-key", WithRetryInterval(2*time.Millisecond), WithLoopTimeout(50*time.Millisecond), WithHTTPClient(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if err := req.Context().Err(); err != nil {
			return nil, err
		}
		calls.Add(1)
		return truncatedResponse(), nil
	})))
	_, err := fetchTransientTest(t.Context(), client)
	var transient TransientTransportError
	if !errors.As(err, &transient) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected transient failure and deadline, got %v", err)
	}
	if calls.Load() < 2 {
		t.Fatalf("expected retries, got %d requests", calls.Load())
	}
}

func TestFetch_ContextCancellationDuringRetrySleepAbortsPromptly(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var calls atomic.Int32
	client := NewClient("https://fixture.invalid", "fixture-key", WithRetryInterval(10*time.Second), WithHTTPClient(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls.Add(1)
		cancel()
		return truncatedResponse(), nil
	})))
	done := make(chan error, 1)
	go func() { _, err := fetchTransientTest(ctx, client); done <- err }()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context cancellation, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("retry sleep did not stop after cancellation")
	}
	if calls.Load() != 1 {
		t.Fatalf("requests = %d, want 1", calls.Load())
	}
}

func TestGetJSON_HTTPError(t *testing.T) {
	for _, tc := range []struct {
		status int
		target any
	}{
		{210, &NoContent{}}, {400, &BadRequest{}}, {401, &Unauthorized{}},
		{403, &Forbidden{}}, {413, &PayloadTooLarge{}}, {429, &TooManyRequests{}},
		{500, &InternalServerError{}}, {502, &BadGateway{}}, {503, &ServiceUnavailable{}}, {504, &GatewayTimeout{}},
	} {
		t.Run(fmt.Sprint(tc.status), func(t *testing.T) {
			c := fixtureClient(t, fixtureResponse{path: "/test", status: tc.status, header: http.Header{"Retry-After": {"2"}}, body: `{"message":"fixture error"}`})
			_, err := getJSON[page[transientTestItem]](t.Context(), c, "/test", transientTestParams{})
			if !errors.As(err, tc.target) || !strings.Contains(err.Error(), "fixture error") {
				t.Fatalf("error = %v, want %T", err, tc.target)
			}
			var httpErr HTTPError
			if !errors.As(err, &httpErr) || httpErr.StatusCode != tc.status {
				t.Fatalf("errors.As(%T, *HTTPError) = %+v, want status %d", err, httpErr, tc.status)
			}
			if rateLimit, ok := tc.target.(*TooManyRequests); ok && rateLimit.RetryAfter != 2*time.Second {
				t.Fatalf("RetryAfter = %v", rateLimit.RetryAfter)
			}
		})
	}
}

func TestGetJSON_PreservesUnknownHTTPStatus(t *testing.T) {
	c := fixtureClient(t, fixtureResponse{path: "/test", status: http.StatusTeapot, body: `{"message":"short and stout"}`})
	_, err := getJSON[page[transientTestItem]](t.Context(), c, "/test", transientTestParams{})
	var httpErr HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusTeapot || !strings.Contains(err.Error(), "short and stout") {
		t.Fatalf("error = %v, want HTTPError with status %d", err, http.StatusTeapot)
	}
}

func TestGetJSON_Compression(t *testing.T) {
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, err := writer.Write([]byte(fullBody)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, body string
		header     http.Header
	}{
		{"plain", fullBody, nil},
		{"gzip", compressed.String(), http.Header{"Content-Encoding": {"gzip"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := fixtureClient(t, fixtureResponse{path: "/test", body: tc.body, header: tc.header})
			got, err := getJSON[page[transientTestItem]](t.Context(), c, "/test", transientTestParams{})
			if err != nil || len(got.Data) != 1 || got.Data[0].Value != "ok" {
				t.Fatalf("response = %#v, %v", got, err)
			}
		})
	}
}

func TestGetJSON_RejectsTrailingData(t *testing.T) {
	for _, suffix := range []string{" garbage", ` {"data":[]}`} {
		t.Run(suffix, func(t *testing.T) {
			c := fixtureClient(t, fixtureResponse{path: "/test", body: fullBody + suffix})
			_, err := getJSON[page[transientTestItem]](t.Context(), c, "/test", transientTestParams{})
			if err == nil || !strings.Contains(err.Error(), "trailing data") {
				t.Fatalf("error = %v, want trailing-data error", err)
			}
		})
	}
}

func TestFetch_RetrySamePage(t *testing.T) {
	for _, stream := range []bool{false, true} {
		for _, status := range []int{429, 500, 502, 503, 504} {
			t.Run(fmt.Sprintf("stream=%t/status=%d", stream, status), func(t *testing.T) {
				c := fixtureClient(t,
					fixtureResponse{path: "/test", body: `{"data":[{"value":"first"}],"pagination_key":"page2"}`},
					fixtureResponse{path: "/test", query: "pagination_key=page2", status: status, body: `{"message":"retry"}`},
					fixtureResponse{path: "/test", query: "pagination_key=page2", body: `{"data":[{"value":"second"}]}`},
				)
				fetch := func(ctx context.Context, key *string) (page[transientTestItem], error) {
					return getJSON[page[transientTestItem]](ctx, c, "/test", transientTestParams{key})
				}
				var got []transientTestItem
				var err error
				if stream {
					got, err = collectChannel(func(ch chan<- transientTestItem) error { return fetchAllPagesWithChannel(t.Context(), c, ch, fetch) })
				} else {
					got, err = fetchAllPages(t.Context(), c, fetch)
				}
				want := []transientTestItem{{Value: "first"}, {Value: "second"}}
				if err != nil || !reflect.DeepEqual(got, want) {
					t.Fatalf("items = %#v, error = %v", got, err)
				}
			})
		}
	}
}

func TestFetch_EmptyPaginationKeyEndsPagination(t *testing.T) {
	for _, stream := range []bool{false, true} {
		t.Run(fmt.Sprint(stream), func(t *testing.T) {
			c := fixtureClient(t, fixtureResponse{path: "/test", body: `{"data":[{"value":"only"}],"pagination_key":""}`})
			fetch := func(ctx context.Context, key *string) (page[transientTestItem], error) {
				return getJSON[page[transientTestItem]](ctx, c, "/test", transientTestParams{key})
			}
			var got []transientTestItem
			var err error
			if stream {
				got, err = collectChannel(func(ch chan<- transientTestItem) error {
					return fetchAllPagesWithChannel(t.Context(), c, ch, fetch)
				})
			} else {
				got, err = fetchAllPages(t.Context(), c, fetch)
			}
			want := []transientTestItem{{Value: "only"}}
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("items = %#v, error = %v; want %#v", got, err, want)
			}
		})
	}
}

func TestFetch_FatalErrors(t *testing.T) {
	for _, stream := range []bool{false, true} {
		for _, tc := range []struct {
			name   string
			status int
			body   string
		}{
			{"forbidden", 403, `{"message":"not included in plan"}`},
			{"malformed JSON", 200, `{"data":!}`},
		} {
			t.Run(fmt.Sprintf("%s/stream=%t", tc.name, stream), func(t *testing.T) {
				c := fixtureClient(t, fixtureResponse{path: "/test", status: tc.status, body: tc.body})
				fetch := func(ctx context.Context, key *string) (page[transientTestItem], error) {
					return getJSON[page[transientTestItem]](ctx, c, "/test", transientTestParams{key})
				}
				var err error
				if stream {
					_, err = collectChannel(func(ch chan<- transientTestItem) error { return fetchAllPagesWithChannel(t.Context(), c, ch, fetch) })
				} else {
					_, err = fetchAllPages(t.Context(), c, fetch)
				}
				if err == nil {
					t.Fatal("expected fatal error")
				}
				if tc.status == 403 && !errors.As(err, &Forbidden{}) {
					t.Fatalf("error = %v, want Forbidden", err)
				}
			})
		}
	}
}

func TestChannel_CancellationWhileBlocked(t *testing.T) {
	c := fixtureClient(t, fixtureResponse{path: "/test", body: fullBody})
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ch := make(chan transientTestItem)
	fetched := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- fetchAllPagesWithChannel(ctx, c, ch, func(ctx context.Context, key *string) (page[transientTestItem], error) {
			r, err := getJSON[page[transientTestItem]](ctx, c, "/test", transientTestParams{key})
			close(fetched)
			return r, err
		})
	}()
	<-fetched
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("blocked producer did not stop")
	}
	if _, ok := <-ch; ok {
		t.Fatal("channel was not closed")
	}
}

func TestCodeDateRangeValues(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		code, date, from, to *string
		want                 string
		invalid              bool
	}{
		{name: "missing filters", invalid: true},
		{name: "code", code: ptr("86970"), want: "code=86970"},
		{name: "date", date: ptr("2026-07-17"), want: "date=2026-07-17"},
		{name: "code and date", code: ptr("86970"), date: ptr("2026-07-17"), want: "code=86970&date=2026-07-17"},
		{name: "range", code: ptr("86970"), from: ptr("2026-07-01"), to: ptr("2026-07-17"), want: "code=86970&from=2026-07-01&to=2026-07-17"},
		{name: "range without code", from: ptr("2026-07-01"), to: ptr("2026-07-17"), invalid: true},
		// The API would apply from/to and ignore date, so the client refuses.
		{name: "date and from", code: ptr("86970"), date: ptr("2026-07-17"), from: ptr("2026-07-01"), invalid: true},
		{name: "date and to", date: ptr("2026-07-17"), to: ptr("2026-07-17"), invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := codeDateRangeValues(tc.code, tc.date, tc.from, tc.to, ptr("next+/="))
			if tc.invalid {
				if err == nil {
					t.Fatal("expected validation error")
				}
				return
			}
			want, parseErr := url.ParseQuery(tc.want + "&pagination_key=next%2B%2F%3D")
			if parseErr != nil {
				t.Fatal(parseErr)
			}
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("query = %v, error = %v; want %v", got, err, want)
			}
		})
	}
}

func TestStockPrice_CodeAndDate(t *testing.T) {
	for _, stream := range []bool{false, true} {
		t.Run(fmt.Sprint(stream), func(t *testing.T) {
			req := StockPriceRequest{Code: ptr("86970"), Date: ptr("2026-07-17")}
			checkEndpoint(t, "/equities/bars/daily", "code=86970&date=2026-07-17", `{"Code":"86970","UL":"0","LL":"0"}`, true, StockPrice{Code: "86970"}, func(c *Client) ([]StockPrice, error) {
				if stream {
					return collectChannel(func(ch chan<- StockPrice) error { return c.StockPriceWithChannel(t.Context(), req, ch) })
				}
				return c.StockPrice(t.Context(), req)
			})
		})
	}
}

func TestSingleResponseEndpoints_Retry(t *testing.T) {
	for _, tc := range []struct {
		path, query, body string
		fetch             func(context.Context, *Client) (any, error)
		want              any
	}{
		{"/equities/master", "", `{"data":[]}`, func(ctx context.Context, c *Client) (any, error) {
			return c.IssueInformation(ctx, IssueInformationRequest{})
		}, []IssueInformation{}},
		{"/markets/calendar", "", `{"data":[]}`, func(ctx context.Context, c *Client) (any, error) {
			return c.TradingCalendar(ctx, TradingCalendarRequest{})
		}, []TradingCalendar{}},
		{"/bulk/list", "date=2026-07-17", `{"data":[]}`, func(ctx context.Context, c *Client) (any, error) {
			return c.BulkList(ctx, BulkListRequest{Date: ptr("2026-07-17")})
		}, []BulkFile{}},
		{"/bulk/get", "key=file.csv.gz", `{"url":"https://download.invalid/file"}`, func(ctx context.Context, c *Client) (any, error) {
			return c.BulkGet(ctx, BulkGetRequest{Key: ptr("file.csv.gz")})
		}, "https://download.invalid/file"},
		{"/td/files", "discNo=123", `{"discNo":"123","files":{}}`, func(ctx context.Context, c *Client) (any, error) {
			return c.TimelyDisclosureFiles(ctx, TimelyDisclosureFilesRequest{DisclosureNumber: "123"})
		}, TimelyDisclosureFiles{DisclosureNumber: "123"}},
		{"/td/bulk", "", `{"url":"https://download.invalid/td"}`, func(ctx context.Context, c *Client) (any, error) { return c.TimelyDisclosureBulk(ctx) }, TimelyDisclosureBulk{URL: "https://download.invalid/td"}},
	} {
		for _, status := range []int{429, 500, 502, 503, 504} {
			t.Run(fmt.Sprintf("%s/%d", tc.path, status), func(t *testing.T) {
				c := fixtureClient(t,
					fixtureResponse{path: tc.path, query: tc.query, status: status, body: `{"message":"transient outage"}`},
					fixtureResponse{path: tc.path, query: tc.query, body: tc.body},
				)
				got, err := tc.fetch(t.Context(), c)
				if err != nil || !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("response = %#v, error = %v; want %#v", got, err, tc.want)
				}
			})
		}
	}
}

func TestSingleResponse_RetryAfterBoundedByTimeout(t *testing.T) {
	c := fixtureClient(t, fixtureResponse{path: "/bulk/get", query: "key=file", status: 429, header: http.Header{"Retry-After": {"3600"}}, body: `{"message":"slow down"}`})
	c.LoopTimeout = 20 * time.Millisecond
	_, err := c.BulkGet(t.Context(), BulkGetRequest{Key: ptr("file")})
	var rateLimit TooManyRequests
	if !errors.Is(err, context.DeadlineExceeded) || !errors.As(err, &rateLimit) || rateLimit.RetryAfter != time.Hour {
		t.Fatalf("expected rate limit and deadline, got %v", err)
	}
}

func TestSingleResponse_NoRetryOnFatalError(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"forbidden", 403, `{"message":"not included in plan"}`},
		{"invalid JSON", 200, `{"url":!}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := fixtureClient(t, fixtureResponse{path: "/bulk/get", query: "key=file", status: tc.status, body: tc.body})
			_, err := c.BulkGet(t.Context(), BulkGetRequest{Key: ptr("file")})
			if err == nil {
				t.Fatal("expected error")
			}
			if tc.status == 403 && !errors.As(err, &Forbidden{}) {
				t.Fatalf("expected Forbidden, got %v", err)
			}
		})
	}
}

func TestSingleResponse_TransientBodyAndCancellation(t *testing.T) {
	t.Run("truncated response", func(t *testing.T) {
		calls := 0
		c := NewClient("https://fixture.invalid", "fixture-key", WithRetryInterval(time.Millisecond), WithHTTPClient(roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return truncatedResponse(), nil
			}
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"url":"https://download.invalid/file"}`))}, nil
		})))
		got, err := c.BulkGet(t.Context(), BulkGetRequest{Key: ptr("file")})
		if err != nil || got != "https://download.invalid/file" || calls != 2 {
			t.Fatalf("response = %q, error = %v, requests = %d", got, err, calls)
		}
	})
	t.Run("cancellation after retry", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		calls := 0
		c := NewClient("https://fixture.invalid", "fixture-key", WithRetryInterval(time.Millisecond), WithHTTPClient(roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return &http.Response{StatusCode: 503, Body: io.NopCloser(strings.NewReader(`{"message":"outage"}`))}, nil
			}
			cancel()
			return nil, req.Context().Err()
		})))
		_, err := c.BulkGet(ctx, BulkGetRequest{Key: ptr("file")})
		if calls != 2 || !errors.Is(err, context.Canceled) || !errors.As(err, &ServiceUnavailable{}) {
			t.Fatalf("requests = %d, error = %v", calls, err)
		}
	})
	t.Run("already canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		c := fixtureClient(t)
		_, err := c.BulkGet(ctx, BulkGetRequest{Key: ptr("file")})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected cancellation, got %v", err)
		}
	})
}

// transportFunc lets tests drive a real *http.Client, including its timeout and
// redirect handling, without opening sockets.
type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func okResponse() *http.Response {
	return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(fullBody))}
}

// LoopTimeout budgets fetching only, so a receiver slower than the whole budget
// still gets every record.
func TestChannel_SlowReceiverNotCutOff(t *testing.T) {
	c := fixtureClient(t,
		fixtureResponse{path: "/test", body: `{"data":[{"value":"a"},{"value":"b"}],"pagination_key":"next"}`},
		fixtureResponse{path: "/test", query: "pagination_key=next", body: `{"data":[{"value":"c"}]}`},
	)
	c.LoopTimeout = 50 * time.Millisecond
	ch := make(chan transientTestItem)
	done := make(chan error, 1)
	go func() {
		done <- fetchAllPagesWithChannel(t.Context(), c, ch, func(ctx context.Context, key *string) (page[transientTestItem], error) {
			return getJSON[page[transientTestItem]](ctx, c, "/test", transientTestParams{key})
		})
	}()
	var got []string
	for item := range ch {
		time.Sleep(40 * time.Millisecond)
		got = append(got, item.Value)
	}
	if err := <-done; err != nil || !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Fatalf("records = %v, error = %v", got, err)
	}
}

// A timeout on a single attempt, such as http.Client.Timeout, is retried while
// the LoopTimeout budget remains.
func TestFetch_RetriesPerAttemptTimeout(t *testing.T) {
	var calls atomic.Int32
	hc := &http.Client{Timeout: 20 * time.Millisecond, Transport: transportFunc(func(req *http.Request) (*http.Response, error) {
		if calls.Add(1) == 1 {
			<-req.Context().Done()
			return nil, req.Context().Err()
		}
		return okResponse(), nil
	})}
	c := NewClient("https://fixture.invalid", "fixture-key", WithHTTPClient(hc), WithRetryInterval(time.Millisecond), WithLoopTimeout(time.Second))
	got, err := fetchTransientTest(t.Context(), c)
	if err != nil || len(got) != 1 || calls.Load() != 2 {
		t.Fatalf("items = %v, error = %v, requests = %d", got, err, calls.Load())
	}
}

// A server that keeps returning the same pagination key must not be followed
// until LoopTimeout; the fixture's request count proves the loop stopped.
func TestFetch_RepeatedPaginationKeyFails(t *testing.T) {
	repeating := func(t *testing.T) *Client {
		return fixtureClient(t,
			fixtureResponse{path: "/test", body: `{"data":[{"value":"a"}],"pagination_key":"same"}`},
			fixtureResponse{path: "/test", query: "pagination_key=same", body: `{"data":[{"value":"a"}],"pagination_key":"same"}`},
		)
	}
	t.Run("slice", func(t *testing.T) {
		_, err := fetchTransientTest(t.Context(), repeating(t))
		if err == nil || !strings.Contains(err.Error(), "already returned") {
			t.Fatalf("error = %v, want repeated pagination key error", err)
		}
	})
	t.Run("channel", func(t *testing.T) {
		c := repeating(t)
		_, err := collectChannel(func(ch chan<- transientTestItem) error {
			return fetchAllPagesWithChannel(t.Context(), c, ch, func(ctx context.Context, key *string) (page[transientTestItem], error) {
				return getJSON[page[transientTestItem]](ctx, c, "/test", transientTestParams{key})
			})
		})
		if err == nil || !strings.Contains(err.Error(), "already returned") {
			t.Fatalf("error = %v, want repeated pagination key error", err)
		}
	})
}

func TestParseRetryAfter(t *testing.T) {
	future := time.Now().Add(30 * time.Second).UTC().Format(http.TimeFormat)
	past := time.Now().Add(-time.Minute).UTC().Format(http.TimeFormat)
	for _, tc := range []struct {
		name, header string
		min, max     time.Duration
	}{
		{"absent", "", 0, 0},
		{"seconds", "2", 2 * time.Second, 2 * time.Second},
		{"negative", "-1", 0, 0},
		{"invalid", "soon", 0, 0},
		{"overflowing seconds", "18446744074", maxRetryAfter, maxRetryAfter},
		{"future date", future, 28 * time.Second, 30 * time.Second},
		{"past date", past, 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{Header: http.Header{}}
			if tc.header != "" {
				resp.Header.Set("Retry-After", tc.header)
			}
			if got := parseRetryAfter(resp); got < tc.min || got > tc.max {
				t.Fatalf("parseRetryAfter(%q) = %v, want between %v and %v", tc.header, got, tc.min, tc.max)
			}
		})
	}
}

func TestClient_NonPositiveSettingsUseDefaults(t *testing.T) {
	c := fixtureClient(t, fixtureResponse{path: "/test", body: fullBody})
	c.LoopTimeout, c.RetryInterval = 0, -time.Second
	if c.loopTimeout() != defaultLoopTimeout || c.retryInterval() != defaultRetryInterval {
		t.Fatalf("loopTimeout = %v, retryInterval = %v", c.loopTimeout(), c.retryInterval())
	}
	// A zero LoopTimeout used to expire before the first request was sent.
	if _, err := fetchTransientTest(t.Context(), c); err != nil {
		t.Fatal(err)
	}
	if _, ok := (&Client{}).httpClient().(*http.Client); !ok || http.DefaultClient.CheckRedirect != nil {
		t.Fatal("nil HTTPClient must use an unmodified copy of http.DefaultClient")
	}
}

func TestSendRequest_TrailingSlashBaseURL(t *testing.T) {
	c := fixtureClient(t, fixtureResponse{path: "/test", body: fullBody})
	c.BaseURL += "/"
	if _, err := fetchTransientTest(t.Context(), c); err != nil {
		t.Fatal(err)
	}
}

// net/http forwards custom headers on redirect, so the client must drop the API
// key itself when a redirect leaves the API's origin.
func TestSendRequest_RedirectAPIKey(t *testing.T) {
	for _, tc := range []struct{ name, location, wantKey string }{
		{"same origin", "https://fixture.invalid/moved", "fixture-key"},
		{"other host", "https://elsewhere.invalid/moved", ""},
		{"scheme downgrade", "http://fixture.invalid/moved", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var redirected bool
			var gotKey string
			hc := &http.Client{Transport: transportFunc(func(req *http.Request) (*http.Response, error) {
				if req.URL.Path == "/test" {
					return &http.Response{StatusCode: http.StatusFound, Header: http.Header{"Location": {tc.location}}, Body: http.NoBody}, nil
				}
				redirected, gotKey = true, req.Header.Get("x-api-key")
				return okResponse(), nil
			})}
			c := NewClient("https://fixture.invalid", "fixture-key", WithHTTPClient(hc))
			if _, err := fetchTransientTest(t.Context(), c); err != nil {
				t.Fatal(err)
			}
			if !redirected || gotKey != tc.wantKey {
				t.Fatalf("redirected = %v, x-api-key = %q; want %q", redirected, gotKey, tc.wantKey)
			}
			if hc.CheckRedirect != nil {
				t.Fatal("caller's http.Client was modified")
			}
		})
	}
}
