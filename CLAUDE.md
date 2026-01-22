# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `jquants-go`, a Go client library for the J-Quants API, which provides access to Japanese stock market data from the Tokyo Stock Exchange (TSE).

## Build and Test Commands

```bash
# Run all tests (requires J_QUANTS_API_KEY environment variable)
go test ./...

# Run a specific test
go test -run TestClient_IssueInformation

# Build the package
go build ./...

# Format code
go fmt ./...

# Vet code
go vet ./...
```

**Note:** Tests make real API calls and require the `J_QUANTS_API_KEY` environment variable to be set.

## Architecture

### Client Pattern

The library uses a single `Client` struct (`client.go`) that holds HTTP client, base URL, API key, and retry/timeout settings. All API methods are methods on this `Client`.

### API Method Structure

Each API endpoint follows a consistent pattern:

1. **Response struct** - Go struct with custom `UnmarshalJSON` to handle J-Quants API quirks (e.g., numeric strings, floats that should be ints)
2. **Request struct** - Public struct with optional filter parameters (uses `*string` for optional fields)
3. **Parameters struct** - Internal struct embedding the request, adding `PaginationKey`
4. **`values()` method** - Implements the `parameters` interface to convert to URL query params
5. **`send*Request` method** - Internal method to make a single paginated request
6. **Public method** - Loops through pagination, handles 500 error retries, returns complete data

### Pagination Handling

APIs that return large datasets use pagination. The client automatically fetches all pages in a loop until `pagination_key` is nil. Some methods also offer `*WithChannel` variants for streaming results.

### Error Types

Custom error types in `client.go` wrap HTTP status codes: `BadRequest`, `Unauthorized`, `Forbidden`, `PayloadTooLarge`, `InternalServerError`. The client auto-retries on `InternalServerError`.

### Module Organization

- `client.go` - Client initialization, HTTP request handling, error types
- `equity.go` - Stock-related APIs: issue information (`/equities/master`), stock prices (`/prices/daily_quotes`)
- `markets.go` - Market data APIs: trading values, margin trading, short selling, trading calendar
- `indices.go` - Index prices including TOPIX
- `option.go` - Index options data
- `codes/codes.go` - Constants for market sections, sector codes, and index codes

### JSON Unmarshaling

The J-Quants API returns some numeric fields as strings and uses abbreviated JSON keys. Custom `UnmarshalJSON` methods translate these to proper Go types with descriptive field names.
