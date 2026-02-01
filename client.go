// Package jquants provides a Go client for the J-Quants API.
//
// J-Quants API provides access to Japanese stock market data from the
// Tokyo Stock Exchange (TSE), including stock prices, trading volumes,
// market indices, and other financial data.
//
// Basic usage:
//
//	client, err := jquants.NewClient()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	issues, err := client.IssueInformation(ctx, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
package jquants

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/time/rate"
)

// BaseURL is the default base URL for the J-Quants API v2.
const BaseURL = "https://api.jquants.com/v2"

// Plan represents a J-Quants subscription plan, which determines rate limits.
type Plan string

const (
	// Light plan allows 1 request per second.
	Light Plan = "Light"
	// Standard plan allows 2 requests per second.
	Standard Plan = "Standard"
	// Premium plan allows 8 requests per second (500 per minute).
	Premium Plan = "Premium"
)

// Client is the J-Quants API client.
// It holds the HTTP client, authentication credentials, and configuration
// for making requests to the J-Quants API.
type Client struct {
	// HttpClient is the HTTP client used for making requests.
	HttpClient *http.Client

	// BaseURL is the base URL for API requests. Defaults to BaseURL constant.
	BaseURL string

	// APIKey is the J-Quants API key for authentication.
	APIKey string

	// RetryInterval is the duration to wait before retrying after a 500 error.
	// Defaults to 5 seconds.
	RetryInterval time.Duration

	// LoopTimeout is the maximum duration for paginated requests.
	// If fetching all pages takes longer than this, the request will be cancelled.
	// Defaults to 20 seconds.
	LoopTimeout time.Duration
}

// RateLimitedTransport is an http.RoundTripper that applies rate limiting to requests.
type RateLimitedTransport struct {
	Transport http.RoundTripper
	Limiter   *rate.Limiter
}

func (t *RateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.Limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return t.Transport.RoundTrip(req)
}

type rateLimit = int

const (
	rateLimitLight    rateLimit = 1
	rateLimitStandard rateLimit = 2
	rateLimitPremium  rateLimit = 8 // 500 requests per min
)

// ClientConfig holds configuration options for creating a Client with NewClientWithConfig.
type ClientConfig struct {
	// BaseURL is the base URL for API requests. Defaults to BaseURL constant.
	BaseURL string
	// APIKey is the J-Quants API key. If empty, reads from J_QUANTS_API_KEY environment variable.
	APIKey string
	// RateLimit is the number of requests per second. Defaults to 1.
	RateLimit int
	// Timeout is the HTTP client timeout. Defaults to 10 seconds.
	Timeout time.Duration
	// RetryInterval is the duration to wait before retrying after a 500 error. Defaults to 5 seconds.
	RetryInterval time.Duration
	// LoopTimeout is the maximum duration for paginated requests. Defaults to 20 seconds.
	LoopTimeout time.Duration
}

func getAPIKey() (string, error) {
	APIKey, ok := os.LookupEnv("J_QUANTS_API_KEY")
	if !ok {
		return "", errors.New("J_QUANTS_API_KEY environment variable is not set")
	}
	return APIKey, nil
}

// NewClient creates a new J-Quants API client.
// It reads the API key from the J_QUANTS_API_KEY environment variable.
// Returns an error if the environment variable is not set.
func NewClient() (*Client, error) {
	httpClient := &http.Client{
		Timeout: 8 * time.Second,
	}
	apiKey, err := getAPIKey()
	if err != nil {
		return nil, err
	}
	client := &Client{
		HttpClient:    httpClient,
		BaseURL:       BaseURL,
		APIKey:        apiKey,
		RetryInterval: 5 * time.Second,
		LoopTimeout:   20 * time.Second,
	}
	return client, nil
}

// NewClientWithRateLimit creates a new J-Quants API client with rate limiting based on the subscription plan.
// It reads the API key from the J_QUANTS_API_KEY environment variable.
// Returns an error if the environment variable is not set.
func NewClientWithRateLimit(plan Plan) (*Client, error) {
	var limit int
	switch plan {
	case Light:
		limit = rateLimitLight
	case Standard:
		limit = rateLimitStandard
	case Premium:
		limit = rateLimitPremium
	}
	httpClient := &http.Client{
		Transport: &RateLimitedTransport{
			Transport: http.DefaultTransport,
			Limiter:   rate.NewLimiter(rate.Limit(limit), limit),
		},
		Timeout: 8 * time.Second,
	}
	apiKey, err := getAPIKey()
	if err != nil {
		return nil, err
	}
	client := &Client{
		HttpClient:    httpClient,
		BaseURL:       BaseURL,
		APIKey:        apiKey,
		RetryInterval: 5 * time.Second,
		LoopTimeout:   20 * time.Second,
	}
	return client, nil
}

// NewClientWithConfig creates a new J-Quants API client with custom configuration.
// Zero values in config are replaced with sensible defaults.
// If APIKey is empty, it reads from the J_QUANTS_API_KEY environment variable.
// Returns an error if no API key is provided and the environment variable is not set.
func NewClientWithConfig(config ClientConfig) (*Client, error) {
	// Set defaults for zero values
	if config.BaseURL == "" {
		config.BaseURL = BaseURL
	}
	if config.RateLimit == 0 {
		config.RateLimit = 1
	}
	if config.Timeout == 0 {
		config.Timeout = 8 * time.Second
	}
	if config.RetryInterval == 0 {
		config.RetryInterval = 5 * time.Second
	}
	if config.LoopTimeout == 0 {
		config.LoopTimeout = 20 * time.Second
	}

	// Get API key from config or environment
	apiKey := config.APIKey
	if apiKey == "" {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return nil, err
		}
	}

	httpClient := &http.Client{
		Transport: &RateLimitedTransport{
			Transport: http.DefaultTransport,
			Limiter:   rate.NewLimiter(rate.Limit(config.RateLimit), config.RateLimit),
		},
		Timeout: config.Timeout,
	}

	client := &Client{
		HttpClient:    httpClient,
		BaseURL:       config.BaseURL,
		APIKey:        apiKey,
		RetryInterval: config.RetryInterval,
		LoopTimeout:   config.LoopTimeout,
	}
	return client, nil
}

type parameters interface {
	values() (url.Values, error)
}

func (c *Client) sendRequest(ctx context.Context, urlPath string, param parameters) (*http.Response, error) {
	u, err := url.Parse(c.BaseURL + urlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	v, err := param.values()
	if err != nil {
		return nil, fmt.Errorf("failed to build query parameters: %w", err)
	}
	u.RawQuery = v.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// HTTPError is the base type for HTTP error responses.
type HTTPError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("%d %s: %v", e.StatusCode, e.Message, e.Err)
}

func (e HTTPError) Unwrap() error {
	return e.Err
}

// BadRequest represents an HTTP 400 error response.
type BadRequest struct{ HTTPError }

// Unauthorized represents an HTTP 401 error response.
// This typically indicates an invalid or missing API key.
type Unauthorized struct{ HTTPError }

// Forbidden represents an HTTP 403 error response.
// This typically indicates the API key does not have permission for the requested resource.
type Forbidden struct{ HTTPError }

// PayloadTooLarge represents an HTTP 413 error response.
// This occurs when the request parameters would result in too much data.
type PayloadTooLarge struct{ HTTPError }

// InternalServerError represents an HTTP 500 error response.
// The client automatically retries requests that receive this error.
type InternalServerError struct{ HTTPError }

func decodeResponse(resp *http.Response, body any) error {
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer func() {
		if clsErr := gzipReader.Close(); clsErr != nil {
			slog.Warn("failed to close response body", "error", clsErr)
		}
	}()
	if err := json.NewDecoder(gzipReader).Decode(body); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

// ErrResponse represents the error response body from the J-Quants API.
type ErrResponse struct {
	Message string `json:"message"`
}

func handleErrorResponse(resp *http.Response) error {
	err := decodeErrorResponse(resp)
	switch resp.StatusCode {
	case 400:
		return BadRequest{HTTPError{400, "bad request", err}}
	case 401:
		return Unauthorized{HTTPError{401, "unauthorized", err}}
	case 403:
		return Forbidden{HTTPError{403, "forbidden", err}}
	case 413:
		return PayloadTooLarge{HTTPError{413, "payload too large", err}}
	case 500:
		return InternalServerError{HTTPError{500, "internal server error", err}}
	default:
		return err
	}
}

func decodeErrorResponse(resp *http.Response) error {
	var errResp ErrResponse
	if err := decodeResponse(resp, &errResp); err != nil {
		return fmt.Errorf("failed to decode error response: %w", err)
	}
	return errors.New(errResp.Message)
}

// fetchAllPages fetches all pages of a paginated API endpoint.
func fetchAllPages[T any, R Response[T]](
	ctx context.Context,
	c *Client,
	fetchPage func(ctx context.Context, paginationKey *string) (R, error),
) ([]T, error) {
	data := make([]T, 0)
	var paginationKey *string
	ctx, cancel := context.WithTimeout(ctx, c.LoopTimeout)
	defer cancel()
	for {
		resp, err := fetchPage(ctx, paginationKey)
		if err != nil {
			if errors.As(err, &InternalServerError{}) {
				slog.Warn("Retrying HTTP request", "error", err.Error())
				time.Sleep(c.RetryInterval)
				continue
			}
			return nil, err
		}
		data = append(data, resp.Items()...)
		paginationKey = resp.NextPageKey()
		if paginationKey == nil {
			break
		}
	}
	return data, nil
}

// fetchAllPagesWithChannel fetches all pages and sends each item to a channel.
func fetchAllPagesWithChannel[T any, R Response[T]](
	ctx context.Context,
	c *Client,
	ch chan<- T,
	fetchPage func(ctx context.Context, paginationKey *string) (R, error),
) error {
	var paginationKey *string
	ctx, cancel := context.WithTimeout(ctx, c.LoopTimeout)
	defer cancel()
	for {
		resp, err := fetchPage(ctx, paginationKey)
		if err != nil {
			if errors.As(err, &InternalServerError{}) {
				slog.Warn("Retrying HTTP request", "error", err.Error())
				time.Sleep(c.RetryInterval)
				continue
			}
			return err
		}
		for _, item := range resp.Items() {
			ch <- item
		}
		paginationKey = resp.NextPageKey()
		if paginationKey == nil {
			break
		}
	}
	close(ch)
	return nil
}
