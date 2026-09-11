package jquants

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

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
			if rateLimit, ok := tc.target.(*TooManyRequests); ok && rateLimit.RetryAfter != 2*time.Second {
				t.Fatalf("RetryAfter = %v", rateLimit.RetryAfter)
			}
		})
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
