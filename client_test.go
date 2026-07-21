package jquants

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"
)

// transientTestItem and transientTestResponse are minimal in-package fixtures
// used to exercise the paginated fetch loop against an httptest server without
// depending on a live J-Quants endpoint.
type transientTestItem struct {
	Value string `json:"value"`
}

type transientTestResponse struct {
	Data          []transientTestItem `json:"data"`
	PaginationKey *string             `json:"pagination_key"`
}

func (r transientTestResponse) Items() []transientTestItem { return r.Data }
func (r transientTestResponse) NextPageKey() *string       { return r.PaginationKey }

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

// truncatedBody declares a large Content-Length but writes only a prefix of the
// JSON document, then closes the connection. The client's JSON decoder reads
// fewer bytes than promised and fails with io.ErrUnexpectedEOF — the same
// failure a huge option-price page produced in production.
func writeTruncated(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", "4096")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"data":[{"value":"ok`))
	// Returning without writing the remaining Content-Length bytes causes the
	// server to close the connection, truncating the body mid-stream.
}

func fetchTransientTest(ctx context.Context, c *Client) ([]transientTestItem, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (transientTestResponse, error) {
		return getJSON[transientTestResponse](ctx, c, "/test", transientTestParams{paginationKey: paginationKey})
	})
}

// TestFetch_RetriesTruncatedBodyThenSucceeds verifies that a body truncated
// mid-stream on the first request is retried and the fetch succeeds once the
// server serves a complete body.
func TestFetch_RetriesTruncatedBodyThenSucceeds(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) == 1 {
			writeTruncated(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fullBody))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key",
		WithRetryInterval(1*time.Millisecond),
		WithLoopTimeout(5*time.Second),
	)

	items, err := fetchTransientTest(t.Context(), client)
	if err != nil {
		t.Fatalf("expected fetch to succeed after retry, got error: %v", err)
	}
	if len(items) != 1 || items[0].Value != "ok" {
		t.Fatalf("unexpected items: %+v", items)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("expected exactly 2 requests (1 truncated + 1 success), got %d", got)
	}
}

// TestFetch_TruncatedBodyEveryRequestFailsWithTransientError verifies that when
// every request truncates, the fetch gives up once the bounded retry budget
// (LoopTimeout) is exhausted, and the returned error chain contains the typed
// TransientTransportError.
func TestFetch_TruncatedBodyEveryRequestFailsWithTransientError(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		writeTruncated(w)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key",
		WithRetryInterval(2*time.Millisecond),
		WithLoopTimeout(50*time.Millisecond),
	)

	items, err := fetchTransientTest(t.Context(), client)
	if err == nil {
		t.Fatalf("expected fetch to fail, got items: %+v", items)
	}
	var transient TransientTransportError
	if !errors.As(err, &transient) {
		t.Fatalf("expected TransientTransportError in chain, got: %v", err)
	}
	if !isContextError(err) {
		t.Fatalf("expected the exhausted retry budget (context deadline) in chain, got: %v", err)
	}
	if calls.Load() < 2 {
		t.Fatalf("expected multiple retry attempts, got %d", calls.Load())
	}
}

// TestFetch_ContextCancellationDuringRetrySleepAbortsPromptly verifies that
// cancelling the caller's context while the loop is sleeping between retries
// aborts promptly and stays fatal (not reclassified as transient).
func TestFetch_ContextCancellationDuringRetrySleepAbortsPromptly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeTruncated(w)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key",
		// A long retry interval guarantees the cancellation lands during the
		// sleep between retries rather than during a round trip.
		WithRetryInterval(10*time.Second),
		WithLoopTimeout(30*time.Second),
	)

	ctx, cancel := context.WithCancel(t.Context())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err := fetchTransientTest(ctx, client)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected fetch to fail after cancellation")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled in chain, got: %v", err)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("expected prompt abort, took %v", elapsed)
	}
}
