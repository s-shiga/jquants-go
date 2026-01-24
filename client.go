package jquants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

const BaseURL = "https://api.jquants.com/v2"

type Client struct {
	HttpClient    *http.Client
	BaseURL       string
	APIKey        string
	RetryInterval time.Duration
	LoopTimeout   time.Duration
}

func NewClient(httpClient *http.Client) (*Client, error) {
	APIKey, ok := os.LookupEnv("J_QUANTS_API_KEY")
	if !ok {
		return nil, errors.New("J_QUANTS_API_KEY environment variable is not set")
	}
	client := &Client{
		HttpClient:    httpClient,
		BaseURL:       BaseURL,
		APIKey:        APIKey,
		RetryInterval: 5 * time.Second,
		LoopTimeout:   20 * time.Second,
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

type BadRequest struct{ HTTPError }
type Unauthorized struct{ HTTPError }
type Forbidden struct{ HTTPError }
type PayloadTooLarge struct{ HTTPError }
type InternalServerError struct{ HTTPError }

func decodeResponse(resp *http.Response, body any) error {
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.Printf("failed to close response body: %s", closeErr.Error())
		}
	}()
	if err := json.NewDecoder(resp.Body).Decode(body); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

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

// paginatedResponse is an interface for API responses that support pagination.
type paginatedResponse[T any] interface {
	getData() []T
	getPaginationKey() *string
}

// fetchAllPages fetches all pages of a paginated API endpoint.
func fetchAllPages[T any, R paginatedResponse[T]](
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
		data = append(data, resp.getData()...)
		paginationKey = resp.getPaginationKey()
		if paginationKey == nil {
			break
		}
	}
	return data, nil
}

// fetchAllPagesWithChannel fetches all pages and sends each item to a channel.
func fetchAllPagesWithChannel[T any, R paginatedResponse[T]](
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
		for _, item := range resp.getData() {
			ch <- item
		}
		paginationKey = resp.getPaginationKey()
		if paginationKey == nil {
			break
		}
	}
	close(ch)
	return nil
}
