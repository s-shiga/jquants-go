package jquants

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

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
