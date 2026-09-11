package jquants

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
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
